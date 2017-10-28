package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	names = NamesFromJSON()
)

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		quit); err != nil {
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
	return nil
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	err := viewNames(g, maxX, maxY)
	if err != nil {
		return err
	}

	err = viewText(g, maxX, maxY)
	if err != nil {
		return err
	}

	err = viewHelp(g, maxX, maxY)
	if err != nil {
		return err
	}

	return nil
}

func viewNames(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("names", -1, 3, 35, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.Title = "Names"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		err := v.SetOrigin(0, 5)
		if err != nil {
			log.Panicln(err)
		}
		for i := 0; i <= maxY/2+3; i++ {
			fmt.Fprintln(v)
		}
		fmt.Fprintln(v, names.String())
	}
	return nil
}

func viewText(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("text", 35, 3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		b, err := ioutil.ReadFile("../../testdata/seashells_book.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		v.Editable = false
		v.Title = "Text"
		v.Wrap = true
	}
	return nil
}

func viewHelp(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("help", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		fmt.Fprintln(v,
			"→ (yes*) next, ← back, Y yes, N no, S sp, U uninom, ? unknown, Ctrl-C exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func listForward(g *gocui.Gui, view *gocui.View) error {
	if names.Current == len(names.Data.Names) {
		names.Current--
		return nil
	}
	for _, v := range g.Views() {
		if v.Name() == "names" {
			view = v
			break
		}
	}
	err := updateNamesView(g, view, 1)
	return err
}

func listBack(g *gocui.Gui, view *gocui.View) error {
	if names.Current == 0 {
		return nil
	}
	for _, v := range g.Views() {
		if v.Name() == "names" {
			view = v
			break
		}
	}
	err := updateNamesView(g, view, -1)
	return err
}

func updateNamesView(g *gocui.Gui, v *gocui.View, increment int) error {
	_, maxY := g.Size()
	annot := Annotation{}
	name := &names.Data.Names[names.Current]
	if name.Annotation == "" {
		name.Annotation = annot.Accepted()
	}
	if increment == 1 {
		names.Current++
	} else {
		names.Current--
	}
	v.Clear()
	for i := 0; i <= maxY/2+3; i++ {
		fmt.Fprintln(v)
	}
	fmt.Fprintln(v, names.String())
	ox, oy := v.Origin()
	if err := v.SetOrigin(ox, oy+(4*increment)); err != nil {
		return err
	}
	return nil
}
