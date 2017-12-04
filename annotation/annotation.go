package annotation

import "fmt"

type Annotation int

const (
	NotAssigned Annotation = iota
	NotName
	Accepted
	Uninomial
	Genus
	Species
	Doubtful
)

var names = []string{
	"",
	"NotName",
	"Accepted",
	"Uninomial",
	"Genus",
	"Species",
	"Doubtful",
}

var namesMap = func() map[string]Annotation {
	res := make(map[string]Annotation)
	for i, v := range names {
		res[v] = Annotation(i)
	}
	return res
}()

func NewAnnotation(s string) (Annotation, error) {
	if a, ok := namesMap[s]; ok {
		return a, nil
	} else {
		return -1, fmt.Errorf("Annotation name '%s' does not exist.", s)
	}
}

func (a Annotation) String() string {
	return names[a]
}

func (a Annotation) Format() string {
	color := a.Color()
	name := a.String()
	return fmt.Sprintf("\033[%d;40;2mAnnot: %s\033[0m", color, name)
}

func (a Annotation) Color() int {
	switch a {
	case Accepted:
		return 32 //green
	case NotName:
		return 31 //red
	case Doubtful:
		return 34 //blue
	case Species:
		return 35 //magenta
	case Genus:
		return 35 //magenta
	case Uninomial:
		return 35 //magenta
	case NotAssigned:
	default:
		return 33
	}
	return 33
}

func (a Annotation) In(as ...Annotation) bool {
	for _, v := range as {
		if a == v {
			return true
		}
	}
	return false
}
