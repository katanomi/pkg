package pointer

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBool(t *testing.T) {
	g := NewGomegaWithT(t)
	input := true
	g.Expect(Bool(input)).To(Equal(&input))
}

func TestInt(t *testing.T) {
	g := NewGomegaWithT(t)
	input := int(0)
	g.Expect(Int(input)).To(Equal(&input))
}

func TestInt64(t *testing.T) {
	g := NewGomegaWithT(t)
	input := int64(0)
	g.Expect(Int64(input)).To(Equal(&input))
}

func TestString(t *testing.T) {
	g := NewGomegaWithT(t)
	input := ""
	g.Expect(String(input)).To(Equal(&input))
}
