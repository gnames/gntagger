package gntagger

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"strings"
	"github.com/atotto/clipboard"
)

var (
	names                = &Names{}
	text                 = &Text{}
	saveCount            = 0
	nameViewCenterOffset = 0
)

func InitGUI(t *Text, n *Names) {
	text = t
	names = n
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(Layout)

	if err := Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, save);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, listBack);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, listForward);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, noName);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'y', gocui.ModNone, yesName);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", 's', gocui.ModNone, speciesName);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'g', gocui.ModNone, genusName);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'u', gocui.ModNone, uninomialName);
		err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'd', gocui.ModNone, doubtfulName);
		err != nil {
		return err
	}

	return nil
}

func Layout(g *gocui.Gui) error {
	err := viewNames(g)
	if err != nil {
		return err
	}

	err = viewText(g)
	if err != nil {
		return err
	}

	err = viewHelp(g)
	if err != nil {
		return err
	}

	return nil
}

func viewNames(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("names", -1, 3, 35, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Names"
		renderNamesView(g)
	}
	return nil
}

func viewText(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("text", 35, 3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Text"
		renderTextView(g)
	}
	return nil
}

func viewHelp(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("help", -1, maxY-2, maxX, maxY); err != nil {
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
	if annotationOfName(name.Annotation) == AnnotationNotAssigned {
		name.Annotation = AnnotationAccepted.name()
	}
	step := 1
	if names.Data.Meta.CurrentName == len(names.Data.Names)-1 {
		step = 0
	}
	names.Data.Meta.CurrentName += step
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
	for ; cursorLeft > 0 && newLinesBefore <= nameViewCenterOffset; cursorLeft-- {
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
	if annotationOfName(name.Annotation) == AnnotationNotName {
		color = AnnotationNotName.color()
	}
	_, err = fmt.Fprintf(viewText, "%s\033[40;%d;1m%s\033[0m%s",
		string(text.Original[cursorLeft+1:name.OffsetStart]),
		color,
		string(text.Original[name.OffsetStart:name.OffsetEnd]),
		string(text.Original[name.OffsetEnd:cursorRight]),
	)
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
		fmt.Fprintln(viewNames, strings.Join(nameStrings(&nm, current, i, namesTotal), "\n"))
	}
	if err = copyCurrentNameToClipboard(); err != nil {
		return err
	}
	return nil
}

func copyCurrentNameToClipboard() error {
	return clipboard.WriteAll(names.currentName().Name)
}
