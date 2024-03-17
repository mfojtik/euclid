package types

import (
	"testing"
	"time"
)

func TestRecordValue(t *testing.T) {
	n := RecordValue(10.0, []Value{}, 5)
	if g := len(n); g != 1 {
		t.Fatalf("expected 1 value, got %d", g)
	}

	n = RecordValue(10.0, []Value{
		{
			Timestamp: time.Now().Add(-10 * time.Minute),
			Value:     1,
		},
		{
			Timestamp: time.Now().Add(-15 * time.Minute),
			Value:     2,
		},
		{
			Timestamp: time.Now().Add(-20 * time.Minute),
			Value:     3,
		},
	}, 3)
	if g := n[len(n)-1].Value; g != 2 {
		t.Errorf("expected last recorded value to be 2, got %f", g)
	}
	if g := n[0].Value; g != 10 {
		t.Errorf("expected last recorded value to be 10, got %f", g)
	}
	if len(n) != 3 {
		t.Errorf("expected only 3 values, got %d", len(n))
	}
}

func TestSetTrend(t *testing.T) {
	n := RecordValue(6, []Value{
		{
			Timestamp: time.Now().Add(-10 * time.Minute),
			Value:     5,
		},
		{
			Timestamp: time.Now().Add(-15 * time.Minute),
			Value:     3,
		},
		{
			Timestamp: time.Now().Add(-20 * time.Minute),
			Value:     3,
		},
	}, 4)
	if r := SetTrend(n); r != "up" {
		t.Errorf("expected up trend, got %s", r)
	}

	n = RecordValue(2, []Value{
		{
			Timestamp: time.Now().Add(-10 * time.Minute),
			Value:     2,
		},
		{
			Timestamp: time.Now().Add(-15 * time.Minute),
			Value:     3,
		},
		{
			Timestamp: time.Now().Add(-20 * time.Minute),
			Value:     6,
		},
	}, 4)
	if r := SetTrend(n); r != "down" {
		t.Errorf("expected down trend, got %s", r)
	}
}
