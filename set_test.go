package timebox

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	set := NewSet()
	slot := NewSlot("Test", time.Now(), 15*time.Second)
	set.Add(slot)
	if len(set.slots) != 1 || set.slots[0] != slot {
		t.Errorf("Adding a slot to a set failed: Set %v", set)
	}
}

func TestSetContains(t *testing.T) {
	set := NewSet()
	now := time.Now()
	set.Add(NewSlot("First", now, 15*time.Second))
	set.Add(NewSlot("Second", now.Add(15*time.Second), 15*time.Second))
	set.Add(NewSlot("Third", now.Add(25*time.Second), 15*time.Second))
	set.Add(NewSlot("Fourth", now.Add(-time.Minute), 20*time.Second))

	tests := []struct {
		t        time.Time
		expected bool
	}{
		{now, true},                                                      // Start of first
		{now.Add(15 * time.Second), true},                                // End of first, but in second
		{now.Add(30 * time.Second), true},                                // End of second, but in third
		{now.Add(40 * time.Second), false},                               // End of third
		{now.Add(-time.Nanosecond), false},                               // Before first
		{now.Add(-time.Minute), true},                                    // Start of fourth
		{now.Add(-time.Minute + 20*time.Second - time.Nanosecond), true}, // Last nanosecond of fourth
		{now.Add(-time.Minute + 20*time.Second), false},                  // End of fourth
		{now.Add(-time.Minute - time.Nanosecond), false},                 // Before fourth
	}

	for _, test := range tests {
		if set.Contains(test.t) != test.expected {
			t.Errorf("set.Contains('%v') returns '%v', but '%v' was expected", test.t, !test.expected, test.expected)
		}
	}
}

func TestFind(t *testing.T) {
	set := NewSet()
	now := time.Now()
	set.Add(NewSlot("First", now, 15*time.Second))
	set.Add(NewSlot("Second", now.Add(15*time.Second), 15*time.Second))
	set.Add(NewSlot("Third", now.Add(25*time.Second), 15*time.Second))
	set.Add(NewSlot("Fourth", now.Add(-time.Minute), 20*time.Second))

	allStartingWithF := set.Find(slotStartsWithLetter('F'))
	assertEqualSlots(allStartingWithF, []*Slot{set.slots[0], set.slots[3]}, t)

	allContainingSpecificTime := set.Find(slotContainsTime(now.Add(25 * time.Second)))
	assertEqualSlots(allContainingSpecificTime, []*Slot{set.slots[1], set.slots[2]}, t)

	allStartingWithX := set.Find(slotStartsWithLetter('X'))
	assertEqualSlots(allStartingWithX, []*Slot{}, t)

}

func TestAllAndAny(t *testing.T) {
	set := NewSet()
	now := time.Now()
	set.Add(NewSlot("First", now, 15*time.Second))
	set.Add(NewSlot("Second", now.Add(15*time.Second), 15*time.Second))
	set.Add(NewSlot("Third", now.Add(25*time.Second), 15*time.Second))
	set.Add(NewSlot("Fourth", now.Add(-time.Minute), 20*time.Second))

	if !set.Any(slotStartsWithLetter('F')) {
		t.Error("No slot starts with F")
	}
	if set.All(slotStartsWithLetter('F')) {
		t.Error("All slots starts with F")
	}
	if set.Any(slotStartsWithLetter('X')) {
		t.Error("Some slot starts with X")
	}
	if !set.All(slotHasDurationOfAtLeast(time.Second)) {
		t.Error("Some slot has a duration of less than a second")
	}
}

func TestSplitLinear(t *testing.T) {
	set := NewSet()

	empty := set.SplitLinear()
	if len(empty) != 0 {
		t.Errorf("Splitting an empty set unexpectedly returns '%v'", empty)
	}

	now := time.Now()
	set.Add(NewSlot("First", now, 15*time.Second))
	set.Add(NewSlot("Second", now.Add(15*time.Second), 15*time.Second))
	set.Add(NewSlot("Third", now.Add(25*time.Second), 15*time.Second))
	set.Add(NewSlot("Fourth", now.Add(-time.Minute), 20*time.Second))

	split := set.SplitLinear()
	if len(split) != 2 {
		t.Errorf("Splitting set '%v' has wrong result '%v'", set, split)
	}
	assertEqualSlots(split[0].slots, []*Slot{set.slots[0], set.slots[1], set.slots[3]}, t)
	assertEqualSlots(split[1].slots, []*Slot{set.slots[2]}, t)
}

func assertEqualSlots(actual, expected []*Slot, t *testing.T) {
	if len(actual) != len(expected) {
		t.Errorf("Expected slots '%v', but got '%v'", expected, actual)
		return
	}
	for i, s := range actual {
		if s != expected[i] {
			t.Errorf("Expected slots '%v', but got '%v'", expected, actual)
			return
		}
	}
}

func slotStartsWithLetter(c uint8) func(*Slot) bool {
	return func(s *Slot) bool {
		return s.Name()[0] == c
	}
}

func slotContainsTime(t time.Time) func(*Slot) bool {
	return func(s *Slot) bool {
		return s.Contains(t)
	}
}

func slotHasDurationOfAtLeast(d time.Duration) func(*Slot) bool {
	return func(s *Slot) bool {
		return s.Duration() >= d
	}
}
