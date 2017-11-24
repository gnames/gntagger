package gntagger_test

import (
	. "github.com/gnames/gntagger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gntagger", func() {
	Describe("Text", func() {
		Describe("NewText", func() {
			It("creates new Text object", func() {
				t := NewText(dataLong, pathLong, "abcd")
				Expect(t.Path).To(Equal("testdata/seashells_book.txt_gntagger"))
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
				var bayes bool
				t := NewText(dataLong, pathLong, "abcd")
				t.Process(80)
				n := NewNames(t, &bayes)
				Expect(n.Path).To(Equal("testdata/seashells_book.txt_gntagger/names.json"))
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

		Describe("Wrap", func() {
			It("wraps text according to given width", func() {
				t := "This is a text that has English and Russian.\n" +
					"Например, мы вводим большое количество русских    знаков в  " +
					"Unicode and see if it works.\nIn case of very long words they " +
					"should-occupy-one-line-and-go-over-the-limit. " +
					"Lets say we have a CR \r\n " +
					"as well and the end does not have a new line........" +
					"...................."
				res := Wrap([]byte(t), 35)
				output := "This is a text that has English \n" +
					"and Russian.\n" +
					"Например, мы вводим большое"
				Expect(string(res[0:96])).To(Equal(output))
			})
		})
	})
})

func makeNames() *Names {
	var bayes bool
	t := NewText(dataShort, pathShort, "abcd")
	t.Process(80)
	n := NewNames(t, &bayes)
	return n
}
