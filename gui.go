package gntagger

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"strings"
)

var (
	names     = &Names{}
	text      = &Text{}
	a         = Annotation{}
	saveCount = 0
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

	if err := g.SetKeybinding("", 'd', gocui.ModNone,
		doubtfulName); err != nil {
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
		renderNamesView(g, v)
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
		renderTextView(g, v)
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
			"→ (yes*) next, ← back, Space no, y yes, s species, g genus, u uninomial, d doubt, ^S save, ^C exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	save(g, v)
	return gocui.ErrQuit
}

func save(g *gocui.Gui, v *gocui.View) error {
	err := names.Save()
	return err
}

func speciesName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.Species())
	return err
}

func genusName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.Genus())
	return err
}

func uninomialName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.Uninomial())
	return err
}

func doubtfulName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.Doubtful())
	return err
}

func yesName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.Accepted())
	return err
}

func noName(g *gocui.Gui, v *gocui.View) error {
	err := setKey(g, v, a.NotName())
	return err
}

func setKey(g *gocui.Gui, v *gocui.View, annot string) error {
	for _, view := range g.Views() {
		if view.Name() == "names" {
			v = view
			break
		}
	}
	err := updateNamesView(g, v, 0, annot)
	return err
}

func listForward(g *gocui.Gui, viewNames *gocui.View) error {
	var (
		viewText *gocui.View
		err      error
	)
	for _, v := range g.Views() {
		switch v.Name() {
		case "names":
			viewNames = v
		case "text":
			viewText = v
		}
	}
	err = updateNamesView(g, viewNames, 1, a.Accepted())
	if err != nil {
		return err
	}
	err = renderTextView(g, viewText)
	return err
}

func listBack(g *gocui.Gui, viewNames *gocui.View) error {
	var viewText *gocui.View
	if names.Data.Meta.CurrentName == 0 {
		return nil
	}
	for _, v := range g.Views() {
		if v.Name() == "names" {
			viewNames = v
		} else if v.Name() == "text" {
			viewText = v
		}
	}
	err := updateNamesView(g, viewNames, -1, "")
	if err != nil {
		return err
	}
	err = renderTextView(g, viewText)
	return err
}

func renderTextView(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	v.Clear()

	name := names.currentName()
	cursorLeft := name.OffsetStart
	newLinesBefore := 0
	for ; cursorLeft > 0 && newLinesBefore < maxY/2-1; cursorLeft-- {
		newLinesBefore = newLinesNum(text.Original[cursorLeft:name.OffsetStart])
	}
	newLinesAfter := 0
	cursorRight := name.OffsetEnd
	for ; cursorRight < len(text.Original)-1 && newLinesAfter < maxY/2-1; cursorRight++ {
		newLinesAfter = newLinesNum(text.Original[name.OffsetEnd:cursorRight])
	}

	for i := 0; i <= maxY/2-newLinesBefore; i++ {
		fmt.Fprintln(v)
	}
	_, err := fmt.Fprintf(v, "%s\033[40;33;1m%s\033[0m%s",
		string(text.Original[cursorLeft:name.OffsetStart]),
		string(text.Original[name.OffsetStart:name.OffsetEnd]),
		string(text.Original[name.OffsetEnd:cursorRight]),
	)
	return err
}

func updateNamesView(g *gocui.Gui, v *gocui.View, step int, annot string) error {
	saveCount++
	if saveCount >= 30 {
		save(g, v)
		saveCount = 0
	}
	name := names.currentName()
	if annot == a.Accepted() {
		if step == 1 && name.Annotation == "" {
			name.Annotation = annot
		} else if step == 0 {
			name.Annotation = annot
		}
	} else if annot != "" {
		name.Annotation = annot
	}
	if names.Data.Meta.CurrentName == len(names.Data.Names)-1 && step == 1 {
		step = 0
	}
	names.Data.Meta.CurrentName += step
	renderNamesView(g, v)
	return nil
}

func renderNamesView(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	v.Clear()
	namesTotal := len(names.Data.Names)
	namesSliceWindow := (maxY - 2) / 4 / 2
	namesSliceLeft := names.Data.Meta.CurrentName - namesSliceWindow
	if namesSliceLeft < 0 {
		namesSliceLeft = 0
	}
	namesSliceRight := names.Data.Meta.CurrentName + namesSliceWindow
	if namesSliceRight > namesTotal {
		namesSliceRight = namesTotal
	}
	fmt.Fprintln(v)
	for i := 0; i <= namesSliceWindow-names.Data.Meta.CurrentName-1; i++ {
		for j := 0; j < 4; j++ {
			fmt.Fprintln(v)
		}
	}
	for i := namesSliceLeft; i < namesSliceRight; i++ {
		current := i == names.Data.Meta.CurrentName
		nm := names.Data.Names[i]
		fmt.Fprintln(v, strings.Join(nameStrings(&nm, current, i, namesTotal), "\n"))
	}
	return nil
}
