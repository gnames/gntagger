package gntagger_test

import (
	"github.com/gnames/gnfinder/output"
	. "github.com/gnames/gntagger"
	"github.com/gnames/gntagger/annotation"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gntagger", func() {
	Describe("Text", func() {
		Describe("NewText", func() {
			It("creates new Text object", func() {
				t := NewText(dataLong, pathLong, "abcd")
				Expect(t.Path).To(Equal(filepath.Join("testdata",
					"seashells_book.txt_gntagger")))
			})
		})

		Describe("Process", func() {
			It("cleans, wraps a text, and adds it to Processed field", func() {
				t := NewText(dataLong, pathLong, "abcd")
				t.Process(40)
				Expect(len(t.Processed)).To(Equal(1112796))
			})
		})

		Describe("Names", func() {
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

			Describe("GetCurrentName", func() {
				It("Returns current name", func() {
					n := makeNames()
					cn := n.GetCurrentName()
					Expect(cn.Name).To(Equal("Phaeophleophleospora epicoccoides"))
				})
			})

			Describe("IsDoubtful", func() {
				It("determines if a name was soubtful", func() {
					gnt := NewGnTagger()
					ns := makeNames()
					n := ns.GetCurrentName()
					Expect(n.Odds).To(BeNumerically(">", gnt.OddsHigh))
					Expect(IsDoubtful(n, gnt)).To(BeFalse())
					n.Odds = 10.0
					Expect(n.Odds).To(BeNumerically("<", gnt.OddsHigh))
					Expect(IsDoubtful(n, gnt)).To(BeTrue())
					n.Odds = 0.0
					Expect(n.Odds).To(BeNumerically("<", gnt.OddsHigh))
					Expect(IsDoubtful(n, gnt)).To(BeFalse())
				})
			})

			Describe("UpdateAnnotations", func() {
				/* Test name data
				    n name                  odds   annot
				    0 Gastropoda            211.35
				    1 Amphineura           9327.30
				    2 Scaphopoda          10292.20
				    3 Pelecypoda          10292.20
				    4 Cephalopoda         12336.64
				    5 Octopus                37.21 Doubtful
				    6 Miiricea muricata  462914.03
				    7 Venus                  16.35 Doubtful
				    8 Venus                  16.35 Doubtful
				    9 Venus                  16.35 Doubtful
				   10 Mollusca              340.20
				   11 Gastropoda            211.00
				   When a name close to the edge changes to NotName or Accepted all
				   same names below also change. If a NotName or Accepted changes
				   to Accepted or NotName the same names below return to NotAssigned
				   or Doubtful (depending on what they had been before)
				*/
				It("updates Doubtful Venus to NotName all the way", func() {
					names := namesForAnnotations()
					ns := names.Data.Names
					gnt := NewGnTagger()
					Expect(gnt.OddsHigh).To(Equal(100.0))
					Expect(gnt.OddsLow).To(Equal(1.0))

					edge := 7
					for i, _ := range ns[0:edge] {
						v := &ns[i]
						v.Annotation = annotation.Accepted.String()
					}

					Expect(ns[8].Annotation).To(Equal(annotation.Doubtful.String()))
					Expect(ns[9].Annotation).To(Equal(annotation.Doubtful.String()))
					Expect(ns[10].Annotation).To(Equal(annotation.NotAssigned.String()))

					names.Data.Meta.CurrentName = 7
					names.UpdateAnnotations(annotation.NotName, edge, gnt)

					Expect(ns[8].Annotation).To(Equal(annotation.NotName.String()))
					Expect(ns[9].Annotation).To(Equal(annotation.NotName.String()))
					Expect(ns[10].Annotation).To(Equal(annotation.NotAssigned.String()))

					edge = 8
					names.Data.Meta.CurrentName = 8
					names.UpdateAnnotations(annotation.Accepted, edge, gnt)

					Expect(ns[8].Annotation).To(Equal(annotation.Accepted.String()))
					Expect(ns[9].Annotation).To(Equal(annotation.Doubtful.String()))
				})

				It("updates NotAssigned Gastropoda to Accepted all the way", func() {
					names := namesForAnnotations()
					ns := names.Data.Names
					gnt := NewGnTagger()

					Expect(ns[0].Annotation).To(Equal(annotation.NotAssigned.String()))
					Expect(ns[11].Annotation).To(Equal(annotation.NotAssigned.String()))

					edge := 0
					names.Data.Meta.CurrentName = 0
					names.UpdateAnnotations(annotation.Accepted, edge, gnt)

					Expect(ns[0].Annotation).To(Equal(annotation.Accepted.String()))
					Expect(ns[10].Annotation).To(Equal(annotation.NotAssigned.String()))
					Expect(ns[11].Annotation).To(Equal(annotation.Accepted.String()))

					names.UpdateAnnotations(annotation.NotName, edge, gnt)

					Expect(ns[0].Annotation).To(Equal(annotation.NotName.String()))
					Expect(ns[11].Annotation).To(Equal(annotation.NotAssigned.String()))
				})
			})
		})
	})
})

func makeNames() *Names {
	gnt := NewGnTagger()
	gnt.Bayes = true
	t := NewText(dataShort, pathShort, "v0.0.0")
	t.Process(80)
	n := NewNames(t, gnt)
	return n
}

func namesForAnnotations() *Names {
	var o output.Output
	o.FromJSON(dataNamesAnnot)
	return &Names{Path: pathNamesAnnot, Data: o}
}
