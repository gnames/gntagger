package gntagger

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Private", func() {

	Describe("wrap", func() {
		It("wraps text according to given width", func() {
			t := "This is a text that has English and Russian.\n" +
				"Например, мы вводим большое количество русских    знаков в  " +
				"Unicode and see if it works.\nIn case of very long words they " +
				"should-occupy-one-line-and-go-over-the-limit. " +
				"Lets say we have a CR \r\n " +
				"as well and the end does not have a new line........" +
				"...................."
			res := wrap([]byte(t), 35)
			output := "This is a text that has English \n" +
				"and Russian.\n" +
				"Например, мы вводим большое"
			Expect(string(res[0:96])).To(Equal(output))
		})
	})
})
