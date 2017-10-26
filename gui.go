package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
)

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		quit); err != nil {
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
		o := NamesFromJSON()
		for i := 0; i <= maxY/2+3; i++ {
			fmt.Fprintln(v)
		}
		for _, n := range o.Names {
			fmt.Fprintf(v, "Kind: '%s'\n", n.Type)
			fmt.Fprintf(v, "Name: '%s'\n", n.Name)
			fmt.Fprintf(v, "\033[40;33;2mAnnot: '%s'\033[0m\n", n.Annotation)
			fmt.Fprintln(v)
		}
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
		if _, err := g.SetCurrentView("text"); err != nil {
			return err
		}
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
			"→ (yes*) next, ← back, Y yes, N no, S sp, U uninom, Ctrl-C exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
