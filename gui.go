package gntagger

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
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
		v.Title = "Names"
		err := v.SetOrigin(0, 5)
		if err != nil {
			log.Panicln(err)
		}
		for i := 0; i <= maxY/2+1; i++ {
			fmt.Fprintln(v)
		}
		fmt.Fprintln(v, names.String())
		ox, oy := v.Origin()
		err = v.SetOrigin(ox, oy+(4*names.Data.Meta.CurrentName))
		if err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func viewText(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("text", 35, 3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Text"
		for i := 0; i <= maxY/2+2; i++ {
			fmt.Fprintln(v)
		}
		fmt.Fprintf(v, "%s", text.Markup(names))

		err := v.SetOrigin(0, text.OffsetY+3)
		if err != gocui.ErrUnknownView {
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
	err = updateText(g, viewText)
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
	err = updateText(g, viewText)
	return err
}

func updateText(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	v.Clear()
	for i := 0; i <= maxY/2+2; i++ {
		fmt.Fprintln(v)
	}
	fmt.Fprintln(v, text.Markup(names))
	err := v.SetOrigin(0, text.OffsetY+3)
	return err
}

func updateNamesView(g *gocui.Gui, v *gocui.View,
	increment int, annot string) error {
	saveCount++
	if saveCount == 30 {
		save(g, v)
		saveCount = 0
	}
	_, maxY := g.Size()
	name := names.currentName()
	if annot == a.Accepted() {
		if increment == 1 && name.Annotation == "" {
			name.Annotation = annot
		} else if increment == 0 {
			name.Annotation = annot
		}
	} else if annot != "" {
		name.Annotation = annot
	}
	if names.Data.Meta.CurrentName == len(names.Data.Names)-1 && increment == 1 {
		increment = 0
	}
	names.Data.Meta.CurrentName += increment
	v.Clear()
	for i := 0; i <= maxY/2+1; i++ {
		fmt.Fprintln(v)
	}
	fmt.Fprintln(v, names.String())
	ox, oy := v.Origin()
	err := v.SetOrigin(ox, oy+(4*increment))
	return err
}
