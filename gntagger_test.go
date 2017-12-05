package gntagger_test

import (
	. "github.com/gnames/gntagger"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gntagger", func() {
	Describe("Text", func() {
		Describe("NewText", func() {
			It("creates new Text object", func() {
				t := NewText(dataLong, pathLong, "abcd")
				Expect(t.Path).To(Equal(filepath.Join("testdata", "seashells_book.txt_gntagger")))
			})
		})

		Describe("Process", func() {
			It("cleans, wraps a text, and adds it to Processed field", func() {
				t := NewText(dataLong, pathLong, "abcd")
				t.Process(40)
				Expect(len(t.Processed)).To(Equal(1112333))
			})
		})

		Describe("NewNames", func() {
			It("creates new Names object", func() {
				gnt := NewGnTagger()
				gnt.Bayes = true
				t := NewText(dataLong, pathLong, "abcd")
				t.Process(80)
				n := NewNames(t, gnt)
				Expect(n.Path).To(Equal(filepath.Join("testdata", "seashells_book.txt_gntagger", "names.json")))
				Expect(n.Data.Meta.TotalNames).To(BeNumerically(">", 4000))
			})
		})

		Describe("Names", func() {
			It("Returns current name", func() {
				n := makeNames()
				cn := n.GetCurrentName()
				Expect(cn.Name).To(Equal("Phaeophleophleospora epicoccoides"))
			})
		})
	})
})

func makeNames() *Names {
	gnt := NewGnTagger()
	gnt.Bayes = true
	t := NewText(dataShort, pathShort, "abcd")
	t.Process(80)
	n := NewNames(t, gnt)
	return n
}
