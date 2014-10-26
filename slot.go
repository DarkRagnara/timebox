package timebox

import (
	"fmt"
	"time"
)

//Slot is a tuple of a name, a start time and a duration.
type Slot struct {
	name  string
	start time.Time
	dur   time.Duration
}

//NewSlot creates a new slot from a name, a start time and a duration.
func NewSlot(name string, start time.Time, duration time.Duration) *Slot {
	return &Slot{name, start, duration}
}

//NewSlotFromTimes creates a new slot from a name, a start time and an end time. The end time is not contained in the slot.
func NewSlotFromTimes(name string, start, end time.Time) *Slot {
	return &Slot{name, start, end.Sub(start)}
}

//Name returns the name given to a slot.
func (s *Slot) Name() string {
	return s.name
}

//Start returns the start time of a slot.
func (s *Slot) Start() time.Time {
	return s.start
}

//Duration returns the duration of a slot.
func (s *Slot) Duration() time.Duration {
	return s.dur
}

//End returns the end time of a slot.
func (s *Slot) End() time.Time {
	return s.start.Add(s.dur)
}

//Contains returns whether a time is in a slot. Note that the start time of a slot is contained by it, but not the end time.
func (s *Slot) Contains(t time.Time) bool {
	if t == s.start {
		return true
	}
	return s.start.Before(t) && s.End().After(t)
}

//Overlaps returns whether two slot overlaps each other, i.e. whether the start time of one slot is contained in the other one.
//Both directions are tested.
func (s *Slot) Overlaps(s2 *Slot) bool {
	return s.Contains(s2.start) || s2.Contains(s.start)
}

//String implements fmt.Stringer.
func (s *Slot) String() string {
	return fmt.Sprintf("Slot{%v, %v, %v}", s.name, s.start, s.dur)
}
