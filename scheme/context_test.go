package scheme

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestSchemeContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := Scheme(ctx)
	g.Expect(clt).To(BeNil())

	ctx = WithScheme(ctx, scheme.Scheme)
	g.Expect(Scheme(ctx)).To(Equal(scheme.Scheme))
}
