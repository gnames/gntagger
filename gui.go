package gntagger

import (
	"fmt"
	"log"

	"strings"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
	"os"
)

type Window struct {
	x0, y0, x1, y1 int
}

func (w *Window) width() int {
	return w.x1 - w.x0
}

func (w *Window) height() int {
	return w.y1 - w.y0
}

type ViewType int

const (
	ViewText  ViewType = iota
	ViewNames
	ViewHelp
)

var (
	textData     []byte
	textDataPath string

	views 				  = map[ViewType]*Window{}
	names                 = &Names{}
	text                  = &Text{}
	saveCount             = 0
	nameViewCenterOffset  = 0
	lastReviewedNameIndex = 0
)

func initViewsMap(g *gocui.Gui) {
	maxX, maxY := g.Size()
	views[ViewNames] = &Window{-1, 3, 35, maxY - 1}
	views[ViewText] = &Window{35, 3, maxX, maxY - 1}
	views[ViewHelp] = &Window{-1, maxY - 2, maxX, maxY}
}

func InitGUI(inputData []byte, inputDataPath string) {
	var err error

	textData = inputData
	textDataPath = inputDataPath

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	for i, n := range names.Data.Names {
		lastReviewedNameIndex = i
		annotationName, err := annotationOfName(n.Annotation)
		if err != nil {
			log.Panicln(err)
		}
		if annotationName == AnnotationNotAssigned {
			break
		}
	}

	g.Cursor = true

	initViewsMap(g)
	text, names, err = prepareData(textData, textDataPath, views[ViewText].width() - 1)
	if err != nil {
		log.Panic(err)
	}
	if names.Data.Meta.TotalNames == 0 {
		fmt.Printf("\nNo names had been found in the document\n\n")
		os.Exit(0)
	}

	g.SetManagerFunc(Layout)

	if err := Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, save); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, listBack); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, listForward); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, noName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'y', gocui.ModNone, yesName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 's', gocui.ModNone, speciesName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'g', gocui.ModNone, genusName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'u', gocui.ModNone, uninomialName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'd', gocui.ModNone, doubtfulName); err != nil {
		return err
	}

	return nil
}

func Layout(g *gocui.Gui) error {
	var err error
	initViewsMap(g)

	if err = viewStats(g); err != nil {
		return err
	}

	if err = viewNames(g); err != nil {
		return err
	}

	if err = viewNames(g); err != nil {
		return err
	}

	if err = viewText(g); err != nil {
		return err
	}

	if err = viewHelp(g); err != nil {
		return err
	}
	return nil
}

func viewNames(g *gocui.Gui) error {
	vn := views[ViewNames]
	if v, err := g.SetView("names", vn.x0, vn.y0, vn.x1, vn.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Names"
		renderNamesView(g)
	}
	return nil
}

func viewText(g *gocui.Gui) error {
	vt := views[ViewText]
	if v, err := g.SetView("text", vt.x0, vt.y0, vt.x1, vt.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Text"
		renderTextView(g)
	}
	return nil
}

func viewStats(g *gocui.Gui) error {
	maxX, _ := g.Size()
	if v, err := g.SetView("stats", -1, -1, maxX, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Stats"
		renderStats(g)
	}
	return nil
}

func viewHelp(g *gocui.Gui) error {
	vh := views[ViewHelp]
	if v, err := g.SetView("help", vh.x0, vh.y0, vh.x1, vh.y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		fmt.Fprintln(v,
			"→ (yes*) next, ← back, Space no, y yes, s species, "+
				"g genus, u uninomial, d doubt, ^S save, ^C exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	save(g, v)
	return gocui.ErrQuit
}

func save(_ *gocui.Gui, _ *gocui.View) error {
	err := names.Save()
	return err
}

func speciesName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationSpecies)
	return err
}

func genusName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationGenus)
	return err
}

func uninomialName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationUninomial)
	return err
}

func doubtfulName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationDoubtful)
	return err
}

func yesName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationAccepted)
	return err
}

func noName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, AnnotationNotName)
	return err
}

func setKey(g *gocui.Gui, annotationId AnnotationId) error {
	var err error
	names.currentName().Annotation = annotationId.name()
	if names.Data.Meta.CurrentName >= lastReviewedNameIndex-3 {
		if annotationId == AnnotationNotName {
			for i := names.Data.Meta.CurrentName + 1; i < len(names.Data.Names); i++ {
				name := &names.Data.Names[i]
				annotationName, err := annotationOfName(name.Annotation)
				if err != nil {
					return nil
				}
				if names.currentName().Name == name.Name &&
					(annotationName == AnnotationNotAssigned ||
						i < lastReviewedNameIndex) {
					name.Annotation = AnnotationNotName.name()
				}
			}
		} else if annotationId != AnnotationNotAssigned {
			for i := names.Data.Meta.CurrentName + 1; i < len(names.Data.Names); i++ {
				name := &names.Data.Names[i]
				annotationName, err := annotationOfName(name.Annotation)
				if err != nil { return err }
				if names.currentName().Name == name.Name &&
					annotationName == AnnotationNotName {
					name.Annotation = AnnotationNotAssigned.name()
				}
			}
		}
	}

	if err = renderNamesView(g); err != nil {
		return err
	}
	if err = renderTextView(g); err != nil {
		return err
	}
	return err
}

func listForward(g *gocui.Gui, _ *gocui.View) error {
	var err error
	name := names.currentName()
	annotationName, err := annotationOfName(name.Annotation)
	if err != nil {
		return nil
	}
	if annotationName == AnnotationNotAssigned {
		name.Annotation = AnnotationAccepted.name()
	}
	step := 1
	if names.Data.Meta.CurrentName == len(names.Data.Names)-1 {
		step = 0
	}
	names.Data.Meta.CurrentName += step
	if names.Data.Meta.CurrentName > lastReviewedNameIndex {
		lastReviewedNameIndex = names.Data.Meta.CurrentName
	}
	if err = renderNamesView(g); err != nil {
		return err
	}
	if err = renderTextView(g); err != nil {
		return err
	}
	return err
}

func listBack(g *gocui.Gui, _ *gocui.View) error {
	var err error
	if names.Data.Meta.CurrentName == 0 {
		return nil
	}
	names.Data.Meta.CurrentName -= 1
	if err = renderNamesView(g); err != nil {
		return err
	}
	if err = renderTextView(g); err != nil {
		return err
	}
	return err
}

func renderTextView(g *gocui.Gui) error {
	var (
		err      error
		viewText *gocui.View
	)
	for _, view := range g.Views() {
		if view.Name() == "text" {
			viewText = view
			break
		}
	}

	_, maxY := g.Size()
	viewText.Clear()

	name := names.currentName()
	cursorLeft := name.OffsetStart - 1
	newLinesBefore := 0
	for ; cursorLeft >= 0 && newLinesBefore <= nameViewCenterOffset; cursorLeft-- {
		if text.Original[cursorLeft] == '\n' {
			newLinesBefore++
		}
	}
	newLinesAfter := 0
	cursorRight := name.OffsetEnd + 1
	for ; cursorRight < len(text.Original)-1 && newLinesAfter < maxY/2-1; cursorRight++ {
		if text.Original[cursorRight] == '\n' {
			newLinesAfter++
		}
	}
	color := AnnotationNotAssigned.color()
	annotationName, err := annotationOfName(name.Annotation)
	if err != nil {
		return nil
	}
	if annotationName == AnnotationNotName {
		color = AnnotationNotName.color()
	}
	for i := 0; i <= nameViewCenterOffset-newLinesBefore; i++ {
		fmt.Fprintln(viewText)
	}
	_, err = fmt.Fprintf(viewText, "%s\033[40;%d;1m%s\033[0m%s",
		string(text.Original[cursorLeft+1:name.OffsetStart]),
		color,
		string(text.Original[name.OffsetStart:name.OffsetEnd]),
		string(text.Original[name.OffsetEnd:cursorRight]),
	)
	for i := 0; i <= newLinesAfter-nameViewCenterOffset+1; i++ {
		fmt.Fprintln(viewText)
	}
	return err
}

func renderNamesView(g *gocui.Gui) error {
	var (
		err       error
		viewNames *gocui.View
	)
	for _, view := range g.Views() {
		if view.Name() == "names" {
			viewNames = view
			break
		}
	}
	saveCount++
	if saveCount >= 30 {
		save(g, viewNames)
		saveCount = 0
	}
	_, maxY := g.Size()
	viewNames.Clear()
	namesTotal := len(names.Data.Names)
	namesSliceWindow := (maxY - 2) / 4 / 2
	nameViewCenterOffset = (namesSliceWindow+1)*4 - 2

	namesSliceLeft := names.Data.Meta.CurrentName - namesSliceWindow
	if namesSliceLeft < 0 {
		namesSliceLeft = 0
	}
	namesSliceRight := names.Data.Meta.CurrentName + namesSliceWindow + 1
	if namesSliceRight > namesTotal {
		namesSliceRight = namesTotal
	}
	fmt.Fprintln(viewNames)
	for i := 0; i <= namesSliceWindow-names.Data.Meta.CurrentName-1; i++ {
		for j := 0; j < 4; j++ {
			fmt.Fprintln(viewNames)
		}
	}
	for i := namesSliceLeft; i < namesSliceRight; i++ {
		current := i == names.Data.Meta.CurrentName
		nm := names.Data.Names[i]
		nameStrs, err := nameStrings(&nm, current, i, namesTotal)
		if err != nil {
			return err
		}
		fmt.Fprintln(viewNames, strings.Join(nameStrs, "\n"))
	}
	if err = copyCurrentNameToClipboard(); err != nil {
		return err
	}
	if err = renderStats(g); err != nil {
		return err
	}
	return nil
}

func renderStats(g *gocui.Gui) error {
	maxX, _ := g.Size()
	var (
		err       error
		viewStats *gocui.View
		stats     Stats
	)
	for _, view := range g.Views() {
		if view.Name() == "stats" {
			viewStats = view
			break
		}
	}

	for nameIdx := 0; nameIdx <= lastReviewedNameIndex; nameIdx++ {
		name := names.Data.Names[nameIdx]
		annotationName, err := annotationOfName(name.Annotation)
		if err != nil {
			return err
		}
		if annotationName != AnnotationNotAssigned {
			stats.total++
		}
		switch annotationName {
		case AnnotationNotName:
			stats.rejectedCount++
		case AnnotationAccepted:
			stats.acceptedCount++
		case AnnotationUninomial, AnnotationGenus, AnnotationSpecies:
			stats.modifiedCount++
		}
	}

	viewStats.Clear()
	fmt.Fprintln(viewStats)
	fmt.Fprintln(viewStats)
	statsStrVisibleLen := 69 // The hack, since len(statsStr) >> len(statsStr_visibleChars)
	for i := 0; i < maxX - statsStrVisibleLen; i++ {
		fmt.Fprint(viewStats, " ")
	}
	fmt.Fprint(viewStats, stats.format())

	return err
}

func copyCurrentNameToClipboard() error {
	return clipboard.WriteAll(names.currentName().Name)
}
