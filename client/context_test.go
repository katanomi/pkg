package client

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClientContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := Client(ctx)
	g.Expect(clt).To(BeNil())

	fakeClt := fake.NewFakeClient()
	ctx = WithClient(ctx, fakeClt)
	g.Expect(Client(ctx)).To(Equal(fakeClt))
}
