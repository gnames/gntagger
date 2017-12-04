package gntagger

import (
	"fmt"
	"math"

	"github.com/gnames/gntagger/annotation"
)

type WordState struct {
	accepted bool
	rejected bool
	modified bool
	added    bool
}

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
	return 1. - float32(s.rejectedCount)/float32(s.total)
}

func (s *Stats) adjustPercentsTo100() {
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
		maxId := -1
		for i := 0; i < len(percents); i++ {
			if percent(rates[i])-*percents[i] < 0 {
				continue
			}
			if maxId == -1 {
				maxId = i
				continue
			}
			rateIf64 := float64(rates[i])
			rateMaxf64 := float64(rates[maxId])
			if (rateIf64 - math.Floor(rateIf64)) > (rateMaxf64 - math.Floor(rateMaxf64)) {
				maxId = i
			}
		}
		*percents[maxId] += 1
	}
}

func (stats *Stats) format() string {
	var (
		precisionStr, acceptedPercentStr       string
		rejectedPercentStr, modifiedPercentStr string
		recallStr                              = "n/a"
		addedRateStr                           = "n/a"
	)

	if stats.total == 0 {
		precisionStr = "n/a"
		acceptedPercentStr = "n/a"
		rejectedPercentStr = "n/a"
		modifiedPercentStr = "n/a"
	} else {
		stats.adjustPercentsTo100()
		precisionStr = fmt.Sprintf("%.2f", stats.precision())
		acceptedPercentStr = fmt.Sprintf("%d%%", stats.acceptedPercent)
		rejectedPercentStr = fmt.Sprintf("%d%%", stats.rejectedPercent)
		modifiedPercentStr = fmt.Sprintf("%d%%", stats.modifiedPercent)
	}

	statsStr := fmt.Sprintf(
		"Precision: %s, Recall: %s | "+
			"\033[%d;40;2mAcc. %s "+
			"\033[%d;40;2mRej. %s "+
			"\033[%d;40;2mMod. %s "+
			"\033[%d;40;2mAdd. %s \033[0m",
		precisionStr,
		recallStr,
		annotation.Accepted.Color(),
		acceptedPercentStr,
		annotation.NotName.Color(),
		rejectedPercentStr,
		annotation.Species.Color(),
		modifiedPercentStr,
		annotation.NotAssigned.Color(),
		addedRateStr,
	)
	return statsStr
}
