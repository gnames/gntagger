package gntagger

// GnTagger keeps configuration parameters of the program
type GnTagger struct {
	// Bayes flag forces bayes name-finding even when the language of the text
	// is not supported.
	Bayes bool
	// OddsHigh marks a limit after which names are considered 'good'.
	OddsHigh float64
	// OddsLow marks a low limit for 'doubtful' names. OddsHigh is the upper
	// limit for such names.
	OddsLow float64
	// Express sets if we skip names that were already marked as 'good'
	// or 'bad'
	Express bool
}

// NewGnTagger creates a new GnTagger object
func NewGnTagger() *GnTagger {
	return &GnTagger{OddsHigh: 100.0, OddsLow: 1, Express: true}
}
