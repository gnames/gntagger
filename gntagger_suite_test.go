package gntagger_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const (
	pathLong       = "./testdata/seashells_book.txt"
	pathShort      = "./testdata/short.txt"
	pathNamesAnnot = "./testdata/names_annot.json"
)

var (
	dataLong       []byte
	dataShort      []byte
	dataNamesAnnot []byte
)

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
	dataNamesAnnot, err = ioutil.ReadFile(pathNamesAnnot)
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
