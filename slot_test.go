package timebox

import (
	"testing"
	"time"
)

func TestEndAndDuration(t *testing.T) {
	start := time.Now()

	durations := []time.Duration{
		0 * time.Second,
		15 * time.Second,
		24 * time.Hour,
		33*time.Minute + 12*time.Second + 1754*time.Nanosecond,
	}

	for _, dur := range durations {
		end := start.Add(dur)
		s1 := NewSlotFromTimes("Test", start, end)
		assertEndTimeAndDuration(s1, end, dur, t)
		s2 := NewSlot("Test", start, dur)
		assertEndTimeAndDuration(s2, end, dur, t)
	}
}

func TestSlotContains(t *testing.T) {
	start := time.Now()
	dur := 30 * time.Second
	tests := []struct {
		t        time.Time
		expected bool
	}{
		{start, true},
		{start.Add(dur), false},
		{start.Add(-1 * time.Nanosecond), false},
		{start.Add(dur - 1*time.Nanosecond), true},
		{start.Add(dur / 2), true},
	}

	s := NewSlot("Test", start, dur)
	for _, test := range tests {
		assertContains(s, test.t, test.expected, t)
	}
}

func TestOverlaps(t *testing.T) {
	start := time.Now()
	dur := 30 * time.Second
	s := NewSlot("Test", start, dur)
	tests := []struct {
		s        *Slot
		expected bool
	}{
		{s, true},
		{NewSlot("1", start.Add(dur/2), dur), true},
		{NewSlot("2", start.Add(-dur/2), dur), true},
		{NewSlot("3", start.Add(dur/4), dur/2), true},
		{NewSlot("4", start.Add(-dur/2), dur*3/2), true},
		{NewSlot("5", start, 0*time.Nanosecond), true},
		{NewSlot("6", start.Add(dur), 0*time.Nanosecond), false},
		{NewSlot("7", start.Add(-dur), dur), false},
	}
	for _, test := range tests {
		assertOverlaps(s, test.s, test.expected, t)
	}
}

func assertContains(s *Slot, time time.Time, expected bool, t *testing.T) {
	if !s.Contains(time) && expected {
		t.Errorf("Slot '%v' does not contain time '%v', but was expected to", s, time)
	} else if s.Contains(time) && !expected {
		t.Errorf("Slot '%v' does contain time '%v', but was not expected to", s, time)
	}
}

func assertOverlaps(s, s2 *Slot, expected bool, t *testing.T) {
	if !s.Overlaps(s2) && expected {
		t.Errorf("Slot '%v' does not overlap slot '%v', but was expected to", s, s2)
	} else if s.Overlaps(s2) && !expected {
		t.Errorf("Slot '%v' does overlap slot '%v', but was not expected to", s, s2)
	}
}

func assertEndTimeAndDuration(s *Slot, expectedEnd time.Time, expectedDur time.Duration, t *testing.T) {
	realDur := s.Duration()
	if realDur != expectedDur {
		t.Errorf("Slot was expected with duration '%v' but has duration '%v", expectedDur, realDur)
	}
	realEnd := s.End()
	if realEnd != expectedEnd {
		t.Errorf("Slot was expected with end time '%v' but has end time '%v", expectedEnd, realEnd)
	}
}
