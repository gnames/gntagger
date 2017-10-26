package gntagger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGntagger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gntagger Suite")
}
