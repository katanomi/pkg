package manager

//go:generate mockgen -package=manager -destination=../mock/sigs.k8s.io/controller-runtime/pkg/manager/manager.go  sigs.k8s.io/controller-runtime/pkg/manager Manager

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	mockmgr "github.com/katanomi/pkg/mock/sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestManagerContext(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	ctx := context.TODO()

	mgr := Manager(ctx)
	g.Expect(mgr).To(BeNil())

	fakeMgr := mockmgr.NewMockManager(ctrl)
	ctx = WithManager(ctx, fakeMgr)
	g.Expect(Manager(ctx)).To(Equal(fakeMgr))
}
