package main

import (
	"errors"

	code "github.com/hetianyi/gox"
)

type SimpleIterator interface {
	Next() (*int, error)
	HasNext() bool
}

type SimpleSetIterator struct {
	Set                *SimpleSet
	CurrentBucketIndex int
	CurrentEntry       *Entry
}

type SimpleHashSet interface {
	Add(int) bool
	Contains(int) bool
	Remove(int) bool
	Iterator() SimpleIterator
}

type Entry struct {
	Key  int
	Next *Entry
}

type SimpleSet struct {
	Size    int
	Buckets []*Entry
}

func initSimpleSet(capacity int) *SimpleSet {
	return &SimpleSet{
		Size:    0,
		Buckets: make([]*Entry, capacity),
	}
}

func (s *SimpleSet) Iterator() SimpleIterator {
	return &SimpleSetIterator{
		Set:                s,
		CurrentBucketIndex: -1,
		CurrentEntry:       nil,
	}
}

func (si *SimpleSetIterator) Next() (*int, error) {
	if si.CurrentEntry == nil || si.CurrentEntry.Next == nil {
		si.CurrentBucketIndex += 1

		for si.CurrentBucketIndex < len(si.Set.Buckets) && si.Set.Buckets[si.CurrentBucketIndex] == nil {
			si.CurrentBucketIndex += 1
		}

		if si.CurrentBucketIndex < len(si.Set.Buckets) {
			si.CurrentEntry = si.Set.Buckets[si.CurrentBucketIndex]
		} else {
			return nil, errors.New("No entry found")
		}
	} else {
		si.CurrentEntry = si.CurrentEntry.Next
	}

	return &si.CurrentEntry.Key, nil
}

func (si *SimpleSetIterator) HasNext() bool {
	if si.CurrentEntry != nil && si.CurrentEntry.Next != nil {
		return true
	}

	for i := si.CurrentBucketIndex + 1; i < len(si.Set.Buckets); i++ {
		if si.Set.Buckets[i] != nil {
			return true
		}
	}

	return false

}

func hashFunc(hashCode int32, bucketLength int) int32 {
	index := hashCode
	if index < 0 {
		index = -index
	}
	return index % int32(bucketLength)
}

func (s *SimpleSet) Add(element int) bool {
	hashCode := code.HashCode(element)
	index := hashFunc(hashCode)
	current := s.Buckets[index]
	for current != nil {
		if current.Key == element {
			return false
		}
		current = current.Next
	}
	entry := &Entry{}
	entry.Key = element
	entry.Next = s.Buckets[index]
	s.Buckets[index] = entry
	s.Size += 1
	return true
}

func (s *SimpleSet) Contains(element int) bool {
	hashCode := code.HashCode(element)
	index := hashFunc(hashCode)
	current := s.Buckets[index]
	for current != nil {
		if current.Key == element {
			return true
		}
		current = current.Next
	}
	return false
}

func (s *SimpleSet) Remove(element int) bool {
	hashCode := code.HashCode(element)
	index := hashFunc(hashCode)
	current := s.Buckets[index]
	var previous *Entry
	for current != nil {
		if current.Key == element {
			if previous == nil {
				s.Buckets[index] = current.Next
			} else {
				previous.Next = current.Next
			}
			s.Size -= 1
			return true
		}
		previous = current
		current = current.Next
	}
	return false
}

func (s *SimpleSet) ToString() string {
	var out string
	var currentEntry *Entry

	for i := 0; i < len(s.Buckets); i++ {
		currentEntry = s.Buckets[i]
		if currentEntry != nil {
			out += "[" + string(i) + "]"
			out += " " + string(currentEntry.Key)

			for currentEntry.Next != nil {
				currentEntry = currentEntry.Next
				out += " -> " + string(currentEntry.Key)
			}
			out += "\n"
		}
	}
	return out
}

func main() {

}
