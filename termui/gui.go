package termui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gnames/gntagger"
	"github.com/gnames/gntagger/annotation"
	"github.com/jroimartin/gocui"
)

// Window keeps information about windows parameters
type Window struct {
	x0, y0, x1, y1 int
}

func (w *Window) width() int {
	return w.x1 - w.x0
}

// ViewType describes possibl window view types
type ViewType int

const (
	// ViewText is a text view
	ViewText ViewType = iota
	// ViewNames is a names view
	ViewNames
	// ViewHelp is a help view
	ViewHelp
)

var (
	gnt                  = &gntagger.GnTagger{}
	views                = map[ViewType]*Window{}
	names                = &gntagger.Names{}
	text                 = &gntagger.Text{}
	saveCount            = 0
	nameViewCenterOffset = 0
	// lastReviewedNameIndex keeps the furtherst checked name so far. It is
	// important for knowing it in case if we moved back and want to know
	// how far from the 'edge' we are now.
	lastReviewedNameIndex = 0
)

func initViewsMap(g *gocui.Gui) {
	maxX, maxY := g.Size()
	views[ViewNames] = &Window{-1, 3, 35, maxY - 1}
	views[ViewText] = &Window{35, 3, maxX, maxY - 1}
	views[ViewHelp] = &Window{-1, maxY - 2, maxX, maxY}
}

// InitGUI initializes command line interface and sets text and names variables
func InitGUI(t *gntagger.Text, gntag *gntagger.GnTagger) {
	gnt = gntag
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	initViewsMap(g)

	names = gntagger.PrepareFilesAndText(t, views[ViewText].width()-1, gnt)
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

// Keybindings sets hotkeys for oprations on the text and names
func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone,
		express); err != nil {
		return err
	}

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

// Layout describes how different vindows are displayed on the screen
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

func express(g *gocui.Gui, _ *gocui.View) error {
	if gnt.Express {
		gnt.Express = false
	} else {
		gnt.Express = true
	}

	if err := renderStats(g); err != nil {
		return err
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
	err := setKey(g, annotation.Species)
	return err
}

func genusName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, annotation.Genus)
	return err
}

func uninomialName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, annotation.Uninomial)
	return err
}

func yesName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, annotation.Accepted)
	return err
}

func noName(g *gocui.Gui, _ *gocui.View) error {
	err := setKey(g, annotation.NotName)
	return err
}

// Changes annotation for current and, if required, the following names
func setKey(g *gocui.Gui, a annotation.Annotation) error {
	var err error

	if err = names.UpdateAnnotations(a, lastReviewedNameIndex, gnt); err != nil {
		return err
	}

	if err = renderNamesView(g); err != nil {
		return err
	}
	err = renderTextView(g)
	return err
}

func listForward(g *gocui.Gui, _ *gocui.View) error {
	var err error
	name := names.GetCurrentName()
	ann, err := annotation.NewAnnotation(name.Annotation)
	if err != nil {
		return err
	}

	if ann == annotation.NotAssigned {
		err := setKey(g, annotation.Accepted)
		if err != nil {
			return err
		}
	}

	step := 1
	if names.Data.Meta.CurrentName == len(names.Data.Names)-1 ||
		name.Annotation == annotation.Doubtful.String() {
		step = 0
	}

	names.Data.Meta.CurrentName += step
	if gnt.Express && step > 0 {
		for _, v := range names.Data.Names[names.Data.Meta.CurrentName:] {
			ann, err := annotation.NewAnnotation(v.Annotation)
			if err != nil {
				panic(fmt.Errorf("Uknown annotation %s", v.Annotation))
			}
			if ann.In(annotation.NotAssigned, annotation.Doubtful) {
				break
			} else {
				names.Data.Meta.CurrentName += 1
			}
		}
	}

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
	names.Data.Meta.CurrentName--
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

	ann, err := annotation.NewAnnotation(name.Annotation)
	if err != nil {
		return nil
	}
	color := ann.Color()
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
		nameStrs, err := gntagger.NameStrings(&nm, current, i, namesTotal)
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
	var stats Stats
	maxX, _ := g.Size()
	wordStates := make(map[string]*WordState)
	viewStats, err := g.View("stats")
	if err != nil {
		log.Panic(err)
	}

	for nameIdx := 0; nameIdx <= lastReviewedNameIndex; nameIdx++ {
		name := names.Data.Names[nameIdx]
		ann, err := annotation.NewAnnotation(name.Annotation)
		if err != nil {
			return err
		}
		ws, ok := wordStates[name.Name]
		if !ok {
			ws = &WordState{}
			wordStates[name.Name] = ws
		}
		switch ann {
		case annotation.NotName:
			ws.rejected = true
		case annotation.Accepted:
			ws.accepted = true
		case annotation.Uninomial, annotation.Genus, annotation.Species:
			ws.modified = true
		}

		if name.Odds != 0.0 && name.Odds < gnt.OddsHigh {
			ws.doubtful = true
		}
	}

	updateStats(&stats, wordStates)
	viewStats.Clear()
	fmt.Fprintln(viewStats)
	fmt.Fprintln(viewStats)
	statsStr := stats.format()
	statsStrVisibleLen := maxX - len(statsStr) + 40 // to compensate invisible chars
	for i := 0; i < statsStrVisibleLen; i++ {
		fmt.Fprint(viewStats, " ")
	}
	fmt.Fprintf(viewStats, "%s", statsStr)

	return err
}

func updateStats(s *Stats, wordStates map[string]*WordState) {
	for _, ws := range wordStates {
		if ws.doubtful {
			if ws.accepted || ws.modified {
				s.addedCount++
				s.total++
			}
		} else {
			if ws.accepted {
				s.acceptedCount++
				s.total++
			}
			if ws.rejected {
				s.rejectedCount++
				s.total++
			}
			if ws.modified {
				s.modifiedCount++
				s.total++
			}
		}
	}
}

func copyCurrentNameToClipboard() error {
	return clipboard.WriteAll(names.GetCurrentName().Name)
}
