package gntagger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

type Window struct {
	x0, y0, x1, y1 int
}

func (w *Window) width() int {
	return w.x1 - w.x0
}

type ViewType int

const (
	ViewText ViewType = iota
	ViewNames
	ViewHelp
)

var (
	views                 = map[ViewType]*Window{}
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

func InitGUI(t *Text, bayes *bool) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	initViewsMap(g)

	names = prepareFilesAndText(t, views[ViewText].width()-1, bayes)
	if names.Data.Meta.TotalNames == 0 {
		g.Close()
		fmt.Printf("\nNo names had been found in the document\n\n")
		os.Exit(0)
	}

	text = t
	g.SetManagerFunc(Layout)

	if err := Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone,
		save); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		listBack); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		listForward); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		noName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'y', gocui.ModNone,
		yesName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 's', gocui.ModNone,
		speciesName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'g', gocui.ModNone,
		genusName); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'u', gocui.ModNone,
		uninomialName); err != nil {
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
		err := renderNamesView(g)
		if err != nil {
			log.Panic(err)
		}
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
		if err := renderTextView(g); err != nil {
			log.Panic(err)
		}
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
		if err := renderStats(g); err != nil {
			log.Panic(err)
		}
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
				"g genus, u uninomial, ^S save, ^C exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	if err := save(g, v); err != nil {
		log.Panic(err)
	}
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
	names.GetCurrentName().Annotation = annotationId.name()
	if names.Data.Meta.CurrentName >= lastReviewedNameIndex-3 {
		if annotationId == AnnotationNotName {
			for i := names.Data.Meta.CurrentName + 1; i < len(names.Data.Names); i++ {
				name := &names.Data.Names[i]
				annotationName, err := annotationOfName(name.Annotation)
				if err != nil {
					return nil
				}
				if names.GetCurrentName().Name == name.Name &&
					(annotationName == AnnotationNotAssigned ||
						i < lastReviewedNameIndex) {
					name.Annotation = AnnotationNotName.name()
				}
			}
		} else if annotationId != AnnotationNotAssigned {
			for i := names.Data.Meta.CurrentName + 1; i < len(names.Data.Names); i++ {
				name := &names.Data.Names[i]
				annotationName, err := annotationOfName(name.Annotation)
				if err != nil {
					return err
				}
				if names.GetCurrentName().Name == name.Name &&
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
	name := names.GetCurrentName()
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
	vText, err := g.View("text")
	if err != nil {
		log.Panic()
	}

	_, maxY := g.Size()
	vText.Clear()

	name := names.GetCurrentName()
	cursorLeft := name.OffsetStart - 1

	newLinesBefore := 0
	for ; cursorLeft >= 0 && newLinesBefore <= nameViewCenterOffset; cursorLeft-- {
		if text.Processed[cursorLeft] == '\n' {
			newLinesBefore++
		}
	}

	newLinesAfter := 0
	cursorRight := name.OffsetEnd + 1
	for ; cursorRight < len(text.Processed)-1 && newLinesAfter < maxY/2-1; cursorRight++ {
		if text.Processed[cursorRight] == '\n' {
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
		fmt.Fprintln(vText)
	}
	_, err = fmt.Fprintf(vText, "%s\033[40;%d;1m%s\033[0m%s",
		string(text.Processed[cursorLeft+1:name.OffsetStart]),
		color,
		string(text.Processed[name.OffsetStart:name.OffsetEnd]),
		string(text.Processed[name.OffsetEnd:cursorRight]),
	)
	for i := 0; i <= newLinesAfter-nameViewCenterOffset+1; i++ {
		fmt.Fprintln(vText)
	}
	return err
}

func renderNamesView(g *gocui.Gui) error {
	viewNames, err := g.View("names")
	if err != nil {
		log.Panic(err)
	}

	saveCount++
	if saveCount >= 30 {
		if err := save(g, viewNames); err != nil {
			log.Panic(err)
		}
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
		text.AddError(fmt.Errorf("\033[31;1mCurrent names did not go to clipboard: %s\033[0m", err))
	}

	if err = renderStats(g); err != nil {
		return err
	}
	return nil
}

func renderStats(g *gocui.Gui) error {
	maxX, _ := g.Size()
	var (
		stats      Stats
		wordStates = map[string]*WordState{}
	)
	viewStats, err := g.View("stats")
	if err != nil {
		log.Panic(err)
	}

	for nameIdx := 0; nameIdx <= lastReviewedNameIndex; nameIdx++ {
		name := names.Data.Names[nameIdx]
		annotationName, err := annotationOfName(name.Annotation)
		if err != nil {
			return err
		}
		ws, ok := wordStates[name.Name]
		if !ok {
			ws = &WordState{}
			wordStates[name.Name] = ws
		}
		switch annotationName {
		case AnnotationNotName:
			ws.rejected = true
		case AnnotationAccepted:
			ws.accepted = true
		case AnnotationUninomial, AnnotationGenus, AnnotationSpecies:
			ws.modified = true
		}
	}

	for _, ws := range wordStates {
		if ws.accepted {
			stats.acceptedCount++
			stats.total++
		}
		if ws.rejected {
			stats.rejectedCount++
			stats.total++
		}
		if ws.modified {
			stats.modifiedCount++
			stats.total++
		}
		if ws.added {
			stats.addedCount++
			stats.total++
		}
	}

	viewStats.Clear()
	fmt.Fprintln(viewStats)
	fmt.Fprintln(viewStats)
	statsStrVisibleLen := 69 // The hack, since len(statsStr) >> len(statsStr_visibleChars)
	for i := 0; i < maxX-statsStrVisibleLen; i++ {
		fmt.Fprint(viewStats, " ")
	}
	fmt.Fprint(viewStats, stats.format())

	return err
}

func copyCurrentNameToClipboard() error {
	return clipboard.WriteAll(names.GetCurrentName().Name)
}
