package gntagger_test

import (
	. "github.com/gnames/gntagger/annotation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Annotation", func() {
	Describe("NewAnnotation", func() {
		It("creates a new annotation", func() {
			a, err := NewAnnotation("")
			Expect(err).ToNot(HaveOccurred())
			Expect(a).To(Equal(NotAssigned))
		})

		It("break on unknown annotation string", func() {
			_, err := NewAnnotation("whaa?")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Annotation name 'whaa?' does not exist."))
		})
	})

	Describe("String", func() {
		It("returns a string respresentation of annotation", func() {
			Expect(NotName.String()).To(Equal("NotName"))
		})
	})

	Describe("Color", func() {
		It("returns terminal color for the annotation", func() {
			Expect(Species.Color()).To(Equal(35))
		})
	})

	Describe("Format", func() {
		It("returns a string colored for terminal", func() {
			res := "\033[34;40;2mAnnot: Doubtful\033[0m"
			Expect(Doubtful.Format()).To(Equal(res))
		})
	})

	Describe("In", func() {
		It("confirms if a list of annotations contains an annotation", func() {
			a := NotName
			b := a.In(Species, Accepted)
			Expect(b).To(BeFalse())
			b = a.In(Species, Accepted, NotName)
			Expect(b).To(BeTrue())
		})
	})
})
