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

package plugin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/component/tracing"
	"github.com/katanomi/pkg/plugin/config"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/route"
)

// Plugin
type plugin struct {
	config        *config.Config
	clients       []client.Interface
	shutdownFuncs []ShutdownFunc
}

type ShutdownFunc func() error

func NewPlugin() *plugin {
	plugin := &plugin{}

	return plugin
}

func (p *plugin) WithConfig(c *config.Config) {
	p.config = c
}

func (p *plugin) WithClient(clients ...client.Interface) *plugin {
	p.clients = append(p.clients, clients...)

	return p
}

// prepare prepare plugin component, include config, route, tracing
func (p *plugin) prepare() {
	if p.config == nil {
		p.config = config.NewConfig()
	}

	if p.config.Trace.Enable {
		closer, err := tracing.Config(&p.config.Trace)
		if err != nil {
			panic(fmt.Sprintf("add srevice error: %s", err.Error()))
		}
		if closer != nil {
			p.shutdownFuncs = append(p.shutdownFuncs, closer.Close)
		}
	}

	restful.DefaultContainer.Add(route.NewDefaultService())

	for _, each := range p.clients {
		ws, err := route.NewService(each)
		if err != nil {
			panic(fmt.Sprintf("add srevice error: %s", err.Error()))
		}

		restful.DefaultContainer.Add(ws)
	}

	restful.DefaultContainer.Add(route.NewDocService())
}

// Run plugin run and shutdown gracefully
func (p *plugin) Run() error {
	//todo watch configmap and reload config

	p.prepare()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := p.config.Server.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: restful.DefaultContainer,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("plugin server run, port:%d\n", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully.")

	for _, f := range p.shutdownFuncs {
		if err := f(); err != nil {
			log.Fatal("prepare shutdown error: ", err)
		}
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	return err
}
