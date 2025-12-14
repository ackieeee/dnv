package kv

import (
	"reflect"
	"testing"
)

func TestCompareMatch(t *testing.T) {
	first := map[string]string{
		"API_KEY": "secret",
		"PORT":    "8080",
	}
	second := map[string]string{
		"PORT":    "8080",
		"API_KEY": "secret",
	}

	result := Compare(first, second)
	if !result.IsMatch() {
		t.Fatalf("expected IsMatch to be true, got false: %+v", result)
	}

	if len(result.MissingInFirst) != 0 || len(result.MissingInSecond) != 0 || len(result.Differing) != 0 {
		t.Fatalf("expected no differences, got %+v", result)
	}
}

func TestCompareDifferences(t *testing.T) {
	first := map[string]string{
		"ONLY_FIRST": "foo",
		"COMMON":     "v1",
		"CHANGE":     "old",
	}
	second := map[string]string{
		"ONLY_SECOND": "bar",
		"COMMON":      "v1",
		"CHANGE":      "new",
	}

	result := Compare(first, second)

	want := Result{
		MissingInFirst:  []string{"ONLY_SECOND"},
		MissingInSecond: []string{"ONLY_FIRST"},
		Differing: []Difference{
			{Key: "CHANGE", FirstValue: "old", SecondValue: "new"},
		},
	}

	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Compare() mismatch\nwant: %+v\ngot:  %+v", want, result)
	}
}
