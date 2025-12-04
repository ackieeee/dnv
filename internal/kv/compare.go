package kv

import "sort"

// Difference represents a key that exists in both files but has different values.
type Difference struct {
	Key         string
	FirstValue  string
	SecondValue string
}

// Result stores the differences discovered between two key/value sets.
type Result struct {
	MissingInFirst  []string // keys present in the second input but not the first
	MissingInSecond []string // keys present in the first input but not the second
	Differing       []Difference
}

// IsMatch returns true when no differences were found.
func (r Result) IsMatch() bool {
	return len(r.MissingInFirst) == 0 && len(r.MissingInSecond) == 0 && len(r.Differing) == 0
}

// Compare inspects the provided maps and reports missing or different entries.
func Compare(first, second map[string]string) Result {
	result := Result{}

	for key, firstVal := range first {
		secondVal, ok := second[key]
		if !ok {
			result.MissingInSecond = append(result.MissingInSecond, key)
			continue
		}
		if secondVal != firstVal {
			result.Differing = append(result.Differing, Difference{Key: key, FirstValue: firstVal, SecondValue: secondVal})
		}
	}

	for key := range second {
		if _, ok := first[key]; !ok {
			result.MissingInFirst = append(result.MissingInFirst, key)
		}
	}

	sort.Strings(result.MissingInFirst)
	sort.Strings(result.MissingInSecond)
	sort.Slice(result.Differing, func(i, j int) bool {
		return result.Differing[i].Key < result.Differing[j].Key
	})

	return result
}
