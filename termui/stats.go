package termui

import (
	"fmt"

	"math"

	"github.com/gnames/gntagger/annotation"
)

// WordState keeps data about a name. It allows us to collect statistics about
// the quality of the name-finding algorithms.
type WordState struct {
	accepted bool
	rejected bool
	modified bool
	doubtful bool
}

// Stats is a collection of fields needed for calculating statistics.
type Stats struct {
	acceptedCount int
	rejectedCount int
	modifiedCount int
	addedCount    int
	total         int

	acceptedPercent int
	rejectedPercent int
	modifiedPercent int
	addedPercent    int
}

func (s *Stats) precision() float32 {
	tpos := float32(s.acceptedCount)
	fpos := float32(s.rejectedCount + s.modifiedCount)
	return tpos / (tpos + fpos)
}

func (s *Stats) recall() float32 {
	tpos := float32(s.acceptedCount)
	fneg := float32(s.addedCount)
	return tpos / (tpos + fneg)
}

func (s *Stats) calculatePercentage() {
	if s.total == 0 {
		return
	}

	acceptedRate := float32(s.acceptedCount) / float32(s.total)
	rejectedRate := float32(s.rejectedCount) / float32(s.total)
	modifiedRate := float32(s.modifiedCount) / float32(s.total)
	addedRate := float32(s.addedCount) / float32(s.total)

	percent := func(rate float32) int { return int(rate * 100) }

	s.acceptedPercent = percent(acceptedRate)
	s.rejectedPercent = percent(rejectedRate)
	s.modifiedPercent = percent(modifiedRate)
	s.addedPercent = percent(addedRate)

	percents := []*int{&s.acceptedPercent, &s.rejectedPercent, &s.modifiedPercent, &s.addedPercent}
	rates := []float32{acceptedRate, rejectedRate, modifiedRate, addedRate}

	totalPercent := func() int {
		return s.acceptedPercent + s.rejectedPercent + s.modifiedPercent + s.addedPercent
	}

	for totalPercent() < 100 {
		maxID := -1
		for i := range percents {
			if percent(rates[i])-*percents[i] < 0 {
				continue
			}
			if maxID == -1 {
				maxID = i
				continue
			}
			rateIf64 := float64(rates[i])
			rateMaxf64 := float64(rates[maxID])
			if (rateIf64 - math.Floor(rateIf64)) >
				(rateMaxf64 - math.Floor(rateMaxf64)) {
				maxID = i
			}
		}
		*percents[maxID]++
	}
}

func (s *Stats) format() string {
	var (
		precisionStr, recallStr, acceptedPercentStr             string
		rejectedPercentStr, modifiedPercentStr, addedPercentStr string
	)

	if s.total == 0 {
		precisionStr = "0.00"
		recallStr = "0.00"
		acceptedPercentStr = "  0%"
		rejectedPercentStr = "  0%"
		modifiedPercentStr = "  0%"
		addedPercentStr = "  0%"
	} else {
		s.calculatePercentage()
		precisionStr = fmt.Sprintf("%.2f", s.precision())
		recallStr = fmt.Sprintf("%.2f", s.recall())
		acceptedPercentStr = fmt.Sprintf("%3d%%", s.acceptedPercent)
		rejectedPercentStr = fmt.Sprintf("%3d%%", s.rejectedPercent)
		modifiedPercentStr = fmt.Sprintf("%3d%%", s.modifiedPercent)
		addedPercentStr = fmt.Sprintf("%3d%%", s.addedPercent)
	}

	statsStr := fmt.Sprintf(
		"Precision: %s, Recall: %s | "+
			"\033[%d;1mAcc. %s "+
			"\033[%d;1mRej. %s "+
			"\033[%d;1mMod. %s "+
			"\033[%d;1mAdd. %s \033[0m",
		precisionStr,
		recallStr,
		annotation.Accepted.Color(),
		acceptedPercentStr,
		annotation.NotName.Color(),
		rejectedPercentStr,
		annotation.Species.Color(),
		modifiedPercentStr,
		annotation.Doubtful.Color(),
		addedPercentStr,
	)
	return statsStr
}
