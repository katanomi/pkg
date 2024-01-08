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

package parallel_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/katanomi/pkg/parallel"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/util/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type executeFlag struct {
	executed bool
	sync.RWMutex
}

func (e *executeFlag) set(executed bool) {
	e.Lock()
	e.executed = executed
	e.Unlock()
}
func (e *executeFlag) get() bool {
	e.RLock()
	defer e.RUnlock()
	return e.executed
}

func TestParallelTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg/parallel")
}

func generateTask(index int, sleep time.Duration, err error, excuted *executeFlag) (task parallel.Task) {
	return func() (interface{}, error) {
		time.Sleep(sleep * time.Second)
		excuted.set(true)
		return fmt.Sprintf("task-%d", index), err
	}
}

var _ = Describe("P().Do().Wait()", func() {

	var (
		t1        parallel.Task
		t1Excuted *executeFlag
		t2        parallel.Task
		t2Excuted *executeFlag
		t3        parallel.Task
		t3Excuted *executeFlag
		t4        parallel.Task
		t4Excuted *executeFlag

		ptasks *parallel.ParallelTasks

		res     []interface{}
		errs    error
		elapsed float64

		zaplog, _ = zap.NewDevelopment()
		log       = zaplog.Sugar()
	)

	BeforeEach(func() {
		t1Excuted = &executeFlag{}
		t2Excuted = &executeFlag{}
		t3Excuted = &executeFlag{}
		t4Excuted = &executeFlag{}

		t1 = generateTask(1, 2, nil, t1Excuted)
		t2 = generateTask(2, 2, nil, t2Excuted)
		t3 = generateTask(3, 3, nil, t3Excuted)
		t4 = func() (interface{}, error) {
			t4Excuted.set(true)
			return nil, nil
		}
		ptasks = parallel.P(log, "custom case", t1, t2, t3, t4)
	})

	JustBeforeEach(func() {
		begin := time.Now()
		res, errs = ptasks.Do().Wait()
		elapsed = time.Since(begin).Seconds()
	})

	Context("when none error happened", func() {
		It("should execute task parallel and collect the results", func() {
			Expect(errs).To(BeNil())
			Expect(t1Excuted.get()).To(BeTrue())
			Expect(t2Excuted.get()).To(BeTrue())
			Expect(t3Excuted.get()).To(BeTrue())
			Expect(elapsed < 4 && elapsed > 3).To(BeTrue())
			Expect(len(res)).To(BeEquivalentTo(3))
			Expect(res).To(ContainElement("task-1"))
			Expect(res).To(ContainElement("task-2"))
			Expect(res).To(ContainElement("task-3"))
		})
	})

	Context("when some errors happened", func() {
		errT1 := errors.New("task-1 error")
		errT2 := errors.New("task-2 error")
		BeforeEach(func() {
			t1 = generateTask(1, 2, errT1, t1Excuted)
			t2 = generateTask(2, 2, errT2, t2Excuted)
			ptasks = parallel.P(log, "errors case", t1, t2, t3)
		})

		It("should return all errors happened and execute other task", func() {
			Expect(errs).ToNot(BeNil())
			multiErrs, ok := errs.(apierrors.Aggregate)
			Expect(ok).To(BeTrue())
			Expect(multiErrs.Errors()).To(ContainElement(errT1))
			Expect(multiErrs.Errors()).To(ContainElement(errT2))
			Expect(multiErrs.Error()).NotTo(BeEmpty())

			Expect(t1Excuted.get()).To(BeTrue())
			Expect(t2Excuted.get()).To(BeTrue())
			Expect(t3Excuted.get()).To(BeTrue())

			Expect(elapsed < 4 && elapsed > 3).To(BeTrue())

			Expect(len(res)).To(BeEquivalentTo(1))
			Expect(res).To(ContainElement("task-3"))
		})
	})

	Context("when set failfast and errors happened", func() {
		errT1 := errors.New("task-1 error")
		errT2 := errors.New("task-2 error")
		BeforeEach(func() {
			t1 = generateTask(1, 2, errT1, t1Excuted)
			t2 = generateTask(2, 1, errT2, t2Excuted)
			ptasks = parallel.P(log, "failfast case", t1, t2, t3).FailFast()
		})

		It("should fail immediately and return first error", func() {
			Expect(errs).ToNot(BeNil())
			Expect(errs).To(BeEquivalentTo(errT2))

			Expect(t1Excuted.get()).To(BeFalse())
			Expect(t2Excuted.get()).To(BeTrue())
			Expect(t3Excuted.get()).To(BeFalse())

			Expect(elapsed < 2 && elapsed > 1).To(BeTrue())

			Expect(len(res)).To(BeEquivalentTo(0))
		})
	})

	Context("when canceld by caller", func() {

		It("should return reason of canceld", func() {
			var cancelFunc context.CancelFunc
			var ctx context.Context
			ctx, cancelFunc = context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancelFunc()
			ptasksCancelable := parallel.P(log, "canceld case", t1, t2, t3).Context(ctx)
			begin := time.Now()
			res, errs = ptasksCancelable.Do().Wait()
			elapsed = time.Since(begin).Seconds()

			Expect(errs).ToNot(BeNil())
			Expect(errs).To(BeEquivalentTo(context.DeadlineExceeded))

			fmt.Printf("%v,%v,%v,", t1Excuted.get(), t2Excuted.get(), t3Excuted.get())
			// up to now, task is not support cancel
			// Expect(t1Excuted.executed).To(BeFalse())
			// Expect(t2Excuted.executed).To(BeFalse())
			// Expect(t3Excuted.executed).To(BeFalse())

			Expect(elapsed < 1 && elapsed > 0.1).To(BeTrue())
			Expect(len(res)).To(BeEquivalentTo(0))
		})
	})

	Context("when many task ", func() {

		flag4 := &executeFlag{}
		flag5 := &executeFlag{}
		flag6 := &executeFlag{}
		flag7 := &executeFlag{}
		flag8 := &executeFlag{}
		flag9 := &executeFlag{}
		flag10 := &executeFlag{}
		BeforeEach(func() {
			t4 := generateTask(4, 2, nil, flag4)
			t5 := generateTask(5, 1, nil, flag5)
			t6 := generateTask(6, 2, nil, flag6)
			t7 := generateTask(7, 1, nil, flag7)
			t8 := generateTask(8, 2, nil, flag8)
			t9 := generateTask(9, 3, nil, flag9)
			t10 := generateTask(10, 3, nil, flag10)
			ptasks = parallel.P(log, "many-tasks", t1, t2, t3, t4, t5, t6).Add(t7).Add(t8).Add(t9).Add(t10).Name("custom case")
		})

		It("should return results of all tasks", func() {
			Expect(errs).To(BeNil())

			Expect(t1Excuted.get()).To(BeTrue())
			Expect(t2Excuted.get()).To(BeTrue())
			Expect(t3Excuted.get()).To(BeTrue())
			Expect(flag4.get()).To(BeTrue())
			Expect(flag5.get()).To(BeTrue())
			Expect(flag6.get()).To(BeTrue())
			Expect(flag7.get()).To(BeTrue())
			Expect(flag8.get()).To(BeTrue())
			Expect(flag9.get()).To(BeTrue())
			Expect(flag10.get()).To(BeTrue())

			Expect(elapsed < 4 && elapsed > 3).To(BeTrue())
			Expect(len(res)).To(BeEquivalentTo(10))

			for i := 1; i <= 10; i++ {
				Expect(res).To(ContainElement(fmt.Sprintf("task-%d", i)))
			}
		})

		It("should return with same sort", func() {
			Expect(errs).To(BeNil())
			for i := 0; i <= 9; i++ {
				Expect(res[i]).To(Equal(fmt.Sprintf("task-%d", i+1)))
			}
		})
	})

	Context("when set conccurrent", func() {
		flag4 := &executeFlag{}
		flag5 := &executeFlag{}
		flag6 := &executeFlag{}
		flag7 := &executeFlag{}
		flag8 := &executeFlag{}
		flag9 := &executeFlag{}
		flag10 := &executeFlag{}
		BeforeEach(func() {
			t4 := generateTask(4, 2, nil, flag4)
			t5 := generateTask(5, 1, nil, flag5)
			t6 := generateTask(6, 3, nil, flag6)
			t7 := generateTask(7, 1, nil, flag7)
			t8 := generateTask(8, 2, nil, flag8)
			t9 := generateTask(9, 3, nil, flag9)
			t10 := generateTask(10, 3, nil, flag10)
			ptasks = parallel.P(log, "many-tasks", t1, t2, t3, t4, t5, t6).Add(t7).Add(t8).Add(t9).Add(t10).Name("custom case")
			ptasks.SetConcurrent(5)
		})

		It("should run tasks in conccurrent count", func() {
			Expect(errs).To(BeNil())

			Expect(t1Excuted.get()).To(BeTrue())
			Expect(t2Excuted.get()).To(BeTrue())
			Expect(t3Excuted.get()).To(BeTrue())
			Expect(flag4.get()).To(BeTrue())
			Expect(flag5.get()).To(BeTrue())
			Expect(flag6.get()).To(BeTrue())
			Expect(flag7.get()).To(BeTrue())
			Expect(flag8.get()).To(BeTrue())
			Expect(flag9.get()).To(BeTrue())
			Expect(flag10.get()).To(BeTrue())

			Expect(elapsed < 7 && elapsed > 6).To(BeTrue())
			Expect(len(res)).To(BeEquivalentTo(10))

			for i := 1; i <= 10; i++ {
				Expect(res).To(ContainElement(fmt.Sprintf("task-%d", i)))
			}
		})
	})
})

var _ = Describe("SetMaxConcurrent", func() {
	var (
		ptasks    *parallel.ParallelTasks
		zapLog, _ = zap.NewDevelopment()
		log       = zapLog.Sugar()
	)

	BeforeEach(func() {
		ptasks = parallel.P(log, "custom case")
	})

	Context("SetMaxConcurrent", func() {
		It("when max great concurrent", func() {
			ptasks.SetMaxConcurrent(5, 6)
			Expect(ptasks.Options.ConcurrencyCount).To(BeEquivalentTo(5))
		})

		It("when max less concurrent", func() {
			ptasks.SetMaxConcurrent(5, 4)
			Expect(ptasks.Options.ConcurrencyCount).To(BeEquivalentTo(4))
		})

		It("when max equal concurrent", func() {
			ptasks.SetMaxConcurrent(5, 5)
			Expect(ptasks.Options.ConcurrencyCount).To(BeEquivalentTo(5))
		})
	})
})
