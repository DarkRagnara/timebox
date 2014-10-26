package timebox

import (
	"time"
)

//Set is a container with zero or more slots in it. Slots in a set may overlap each other.
//You can use the set to calculate lanes from the slots inside.
type Set struct {
	slots []*Slot
}

//NewSet creates a new set.
func NewSet() *Set {
	return &Set{}
}

//Add adds a new slot to a set.
//The slot is passed as a pointer that is stored in the set, so be aware that later changes to the slot affect the set.
func (s *Set) Add(slot *Slot) {
	s.slots = append(s.slots, slot)
}

//Slots returns a slice of the slots in a set. Be aware that changes to the returned slots affect the set.
func (s *Set) Slots() []*Slot {
	return s.slots
}

//Find returns a slice of all slots in the set satisfying f(slot), or an empty slice if none do.
//Be aware that any changes to the returned slots affect the set.
func (s *Set) Find(f func(*Slot) bool) []*Slot {
	results := []*Slot{}
	for _, slot := range s.slots {
		if f(slot) {
			results = append(results, slot)
		}
	}
	return results
}

//Any returns whether any slot in the set satisfies f(slot).
func (s *Set) Any(f func(*Slot) bool) bool {
	for _, slot := range s.slots {
		if f(slot) {
			return true
		}
	}
	return false
}

//All returns whether any slot in the set satisfies f(slot).
func (s *Set) All(f func(*Slot) bool) bool {
	for _, slot := range s.slots {
		if !f(slot) {
			return false
		}
	}
	return true
}

//Contains returns whether any slot in the set contains a given time.
func (s *Set) Contains(t time.Time) bool {
	return s.Any(func(slot *Slot) bool {
		return slot.Contains(t)
	})
}

//Overlaps returns whether any slot in the set overlaps a given other slot.
func (s *Set) Overlaps(slot *Slot) bool {
	return s.Any(func(sl *Slot) bool {
		return sl.Overlaps(slot)
	})
}

//SplitLinear distributes the slots in the set to multiple sets in a way that ensures that no set contains multiple overlapping slots.
func (s *Set) SplitLinear() []*Set {
	result := []*Set{}
Outer:
	for _, slotToDistribute := range s.slots {
		for _, set := range result {
			if !set.Overlaps(slotToDistribute) {
				set.Add(slotToDistribute)
				continue Outer
			}
		}
		newSet := NewSet()
		newSet.Add(slotToDistribute)
		result = append(result, newSet)
	}
	return result
}
