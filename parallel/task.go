/*
Copyright 2021 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package parallel used to execute tasks in parallel
package parallel

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"reflect"

	"k8s.io/apimachinery/pkg/util/errors"
)

type Task func() (interface{}, error)

// ParallelTasks will construct a parallel tasks struct
// you could execute tasks in parallel
// eg.
// result, err :=  P("eg1", f1,f2, f3).Do().Wait()
// result, err :=  P("eg2", f1,f2, f3).Add(f4).FailFast().Do().Wait()
//
//	result, err :=  P("eg3", f1,f2, f3).Context(func()context.Context{
//			ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond) // 0.5s will timeout
//			return ctx
//	}).Do().Wait()
type ParallelTasks struct {
	name string

	tasks []Task

	// used for collect all errors
	errLock sync.Locker
	errs    []error

	// used for collect all results
	resultsLock sync.Locker
	results     []interface{}

	doneChan  chan struct{}
	doneError error
	doneOnce  sync.Once

	wg  sync.WaitGroup
	ctx context.Context

	// concurrency threshold
	threshold chan struct{}

	Options ParallelOptions

	Log *zap.SugaredLogger
}

type ParallelOptions struct {
	FailFast         bool
	ConcurrencyCount int
}

// P will construct ParallelTasks
// name will be used for log
// you must care about the variable that referenced by Closure
func P(log *zap.SugaredLogger, name string, tasks ...Task) *ParallelTasks {
	return &ParallelTasks{
		name:        name,
		tasks:       tasks,
		ctx:         context.Background(),
		resultsLock: &sync.Mutex{},
		results:     []interface{}{},
		errLock:     &sync.Mutex{},
		errs:        []error{},
		Log:         log,

		doneChan: make(chan struct{}),
	}
}

func (p *ParallelTasks) Name(name string) *ParallelTasks {
	p.name = name
	return p
}

// Add will add more tasks to ParallelTasks
// you should invoke it before invoke Do or Wait
// you must care about the variable that referenced by Closure
//
// eg1.
// pts := P("eg")
// for _, item := range itemArrar {
//  var itemTmp = item
// 	pts.Add(func()(interface{}, error){
// 		fmt.Println(itemTmp)
// 	})
// }
// -----------------------------------------
// eg2.
// func genTask(name string) Task {
//   return func()(interface{},error){
//   	  fmt.Println(name)
//      return nil,nil
//   }
// }

// pts := P("eg")
// for _, item := range itemArrar {
// 	pts.Add(genTask(item))
// }

func (p *ParallelTasks) Add(tasks ...Task) *ParallelTasks {
	p.tasks = append(p.tasks, tasks...)
	return p
}

func (p *ParallelTasks) FailFast() *ParallelTasks {
	p.Options.FailFast = true
	return p
}

func (p *ParallelTasks) SetConcurrent(count int) *ParallelTasks {
	p.Options.ConcurrencyCount = count
	return p
}

// Context will set context , up to now , task is not support to cancel
// if you cancel from context, wait will return immediately
func (p *ParallelTasks) Context(ctx context.Context) *ParallelTasks {
	p.ctx = ctx
	return p
}

// waitThreshold will wait one threshold until done
// if done ,it will return false; it will return true until got the threshold
func (p *ParallelTasks) waitThreshold() bool {
	if p.Options.ConcurrencyCount <= 0 {
		return true
	}

	if p.threshold == nil {
		p.threshold = make(chan struct{}, p.Options.ConcurrencyCount)
	}

	select {
	case p.threshold <- struct{}{}:
		return true
	case <-p.doneChan:
		return false
	}
}

func (p *ParallelTasks) releaseThreshold() {
	if p.Options.ConcurrencyCount <= 0 {
		return
	}

	<-p.threshold
}

// Do will start to execute all task in parallel
func (p *ParallelTasks) Do() *ParallelTasks {
	log := p.Log.Named(fmt.Sprintf("[ParallelTasks %s]", p.name))

	if len(p.tasks) == 0 {
		return p
	}

	go func() {
		select {
		case <-p.ctx.Done():
			p.Cancel(p.ctx.Err())
			return
		case <-p.doneChan:
			return
		}
	}()

	p.results = make([]interface{}, len(p.tasks))
	p.errs = make([]error, len(p.tasks))

	for i, task := range p.tasks {
		if !p.waitThreshold() {
			return p
		}

		p.wg.Add(1)
		go func(index int, t Task) {
			defer func() {
				p.wg.Done()
				p.releaseThreshold()
				log.Debugw("task: completed", "task-index", index+1)
			}()

			result, err := t()
			if err != nil {
				log.Errorw("task: error  to got result", "task-index", index, "result", result, "err", err.Error())
			} else {
				log.Debugw("task: got result", "task-index", index, "result", result)
			}

			// error is not nil, we should save error
			if err != nil {
				p.errLock.Lock()
				defer p.errLock.Unlock()
				p.errs[index] = err
				if p.Options.FailFast {
					log.Debugw("fail fast, will cancel", "task-index", index, "result", result)
					p.Cancel(err)
				}
				return
			}

			// error is nil, we should save result
			if !isNil(result) {
				p.results[index] = result
			}

		}(i, task)
	}

	return p
}

func (p *ParallelTasks) Cancel(cancelReason error) {
	p.done(cancelReason)
}

func (p *ParallelTasks) done(reason error) {
	p.doneOnce.Do(func() {
		p.doneError = reason
		p.doneChan <- struct{}{}
		close(p.doneChan)
	})
}

// Wait will wait all task executed, if set fail fast , it will return immediately if any task returns errors
// up to now , task is not support to cancel
// you should invoke Do() before invoke Wait()
// the result of task will be saved in []interface{}
// if you set failfast and one errors happened, it will return one error
// if you not set failfase and any errors happened, it will return []error as MultiErrors
func (p *ParallelTasks) Wait() ([]interface{}, error) {

	if len(p.tasks) == 0 {
		return nil, nil
	}

	go func() {
		p.wg.Wait()
		p.done(nil)
	}()

	p.Log.Debugw("waiting done.")
	<-p.doneChan
	p.Log.Debugw("waited done.")

	var (
		results = make([]interface{}, 0)
		errs    = make([]error, 0)
	)
	for _, result := range p.results {
		if result != nil {
			results = append(results, result)
		}
	}
	for _, err := range p.errs {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if p.doneError != nil {
		return results, p.doneError
	}

	if len(errs) > 0 {
		return results, errors.NewAggregate(errs)
	}

	return results, nil
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Func:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}
