package controllers

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type mockChecker struct {
	name string
	err  error
}

func (m *mockChecker) Name() string {
	return m.name
}

func (m *mockChecker) CheckSetup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	return m.err
}

func (m *mockChecker) Setup(ctx context.Context, mgr manager.Manager, logger *zap.SugaredLogger) error {
	return m.err
}

func (m *mockChecker) DependentCrdInstalled(ctx context.Context, logger *zap.SugaredLogger) (bool, error) {
	return true, nil
}

func TestControllerLazyLoader(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	// cancel, _ := context.WithCancel(ctx)
	interval := 100 * time.Millisecond
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	// create a new lazy loader
	loader := NewLazyLoader(ctx, interval).(*controllerLazyLoader)

	// add a mock checker that will always fail
	checker1 := &mockChecker{name: "checker1", err: errors.New("failed to setup")}
	err := loader.LazyLoad(ctx, nil, sugar, checker1)
	g.Expect(err).To(BeNil())

	// add a mock checker that will always succeed
	checker2 := &mockChecker{name: "checker2", err: nil}
	err = loader.LazyLoad(ctx, nil, sugar, checker2)
	g.Expect(err).To(BeNil())

	// start the lazy loader
	done := make(chan struct{})
	go func() {
		err := loader.Start(ctx)
		g.Expect(err).To(BeNil())
		close(done)
	}()

	// wait for the first check to complete
	time.Sleep(2 * interval)

	g.Expect(len(loader.pending)).To(Equal(1))
	g.Expect(len(loader.done)).To(Equal(1))

	// wait for the second check to complete
	time.Sleep(2 * interval)

	g.Expect(len(loader.pending)).To(Equal(1))
	g.Expect(len(loader.done)).To(Equal(1))

	go func() {
		time.Sleep(10 * interval)
		<-ctx.Done()
	}()

}
