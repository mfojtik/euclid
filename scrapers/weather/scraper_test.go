package weather

import (
	"github.com/gookit/goutil/dump"
	"testing"
)

func TestScraper_Scrape(t *testing.T) {
	s := New()
	r, err := s.Scrape()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	dump.P(r)
}
