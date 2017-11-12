package gntagger

type Stats struct {
	acceptedCount int
	rejectedCount int
	modifiedCount int
	addedCount int

	total int
}

func (s *Stats) acceptedRate() int {
	if s.total == 0 {
		return 0
	}
	return int(float32(s.acceptedCount) / float32(s.total) * 100)
}

func (s *Stats) rejectedRate() int {
	if s.total == 0 {
		return 0
	}
	return int(float32(s.rejectedCount) / float32(s.total) * 100)
}

func (s *Stats) modifiedRate() int {
	if s.total == 0 {
		return 0
	}
	return int(float32(s.modifiedCount) / float32(s.total) * 100)
}

func (s *Stats) addedRate() int {
	if s.total == 0 {
		return 0
	}
	return -1
}

func (s *Stats) precision() float32 {
	if s.total == 0 {
		return 0
	}
	return float32(s.rejectedCount) / float32(s.total)
}

func (s *Stats) recall() float32 {
	if s.total == 0 {
		return 0
	}
	return -1.
}
