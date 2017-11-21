package gntagger_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const pathLong = "./testdata/seashells_book.txt"
const pathShort = "./testdata/short.txt"

var dataLong []byte
var dataShort []byte

func TestGntagger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gntagger Suite")
}

var _ = BeforeSuite(func() {
	var err error
	dataLong, err = ioutil.ReadFile(pathLong)
	Expect(err).ToNot(HaveOccurred())
	dataShort, err = ioutil.ReadFile(pathShort)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	dir := pathLong + "_gntagger"
	err := os.RemoveAll(dir)
	Expect(err).ToNot(HaveOccurred())
	dir = pathShort + "_gntagger"
	err = os.RemoveAll(dir)
	Expect(err).ToNot(HaveOccurred())
})
