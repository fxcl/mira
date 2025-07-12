package utils

import (
	"testing"
)

func TestCheckRegex(t *testing.T) {
	// Test case 1: Valid regex and matching content
	if !CheckRegex(`\d+`, "123") {
		t.Error("CheckRegex(`\\d+`, \"123\") should be true")
	}

	// Test case 2: Valid regex and non-matching content
	if CheckRegex(`\d+`, "abc") {
		t.Error("CheckRegex(`\\d+`, \"abc\") should be false")
	}

	// Test case 3: Invalid regex
	if CheckRegex(`[`, "abc") {
		t.Error("CheckRegex(`[`, \"abc\") should be false for invalid regex")
	}
}

func TestContains(t *testing.T) {
	// Test with integers
	intSlice := []int{1, 2, 3, 4, 5}
	if !Contains(intSlice, 3) {
		t.Error("Contains(intSlice, 3) should be true")
	}
	if Contains(intSlice, 6) {
		t.Error("Contains(intSlice, 6) should be false")
	}

	// Test with strings
	stringSlice := []string{"a", "b", "c"}
	if !Contains(stringSlice, "b") {
		t.Error("Contains(stringSlice, \"b\") should be true")
	}
	if Contains(stringSlice, "d") {
		t.Error("Contains(stringSlice, \"d\") should be false")
	}

	// Test with an empty slice
	emptySlice := []int{}
	if Contains(emptySlice, 1) {
		t.Error("Contains(emptySlice, 1) should be false")
	}
}

func TestFilter(t *testing.T) {
	// Test with integers: filter for even numbers
	intSlice := []int{1, 2, 3, 4, 5, 6}
	evenNumbers := Filter(intSlice, func(n int) bool {
		return n%2 == 0
	})
	if len(evenNumbers) != 3 || evenNumbers[0] != 2 || evenNumbers[1] != 4 || evenNumbers[2] != 6 {
		t.Errorf("Filter for even numbers failed, got: %v", evenNumbers)
	}

	// Test with strings: filter for strings with length > 3
	stringSlice := []string{"a", "bb", "ccc", "dddd"}
	longStrings := Filter(stringSlice, func(s string) bool {
		return len(s) > 3
	})
	if len(longStrings) != 1 || longStrings[0] != "dddd" {
		t.Errorf("Filter for long strings failed, got: %v", longStrings)
	}

	// Test with an empty slice
	emptySlice := []int{}
	filteredEmpty := Filter(emptySlice, func(n int) bool {
		return true // Condition doesn't matter
	})
	if len(filteredEmpty) != 0 {
		t.Error("Filter on empty slice should result in an empty slice")
	}

	// Test where no elements match
	noMatches := Filter(intSlice, func(n int) bool {
		return n > 10
	})
	if len(noMatches) != 0 {
		t.Error("Filter with no matching elements should result in an empty slice")
	}
}

func TestDesensitize(t *testing.T) {
	// Basic desensitization
	if Desensitize("123456789", 2, 5) != "12****789" {
		t.Error("Basic desensitization failed")
	}

	// start < 0
	if Desensitize("123456789", -1, 5) != "123456789" {
		t.Error("Desensitization with start < 0 failed")
	}

	// end < 0
	if Desensitize("123456789", 2, -1) != "123456789" {
		t.Error("Desensitization with end < 0 failed")
	}

	// start > end
	if Desensitize("123456789", 5, 2) != "123456789" {
		t.Error("Desensitization with start > end failed")
	}

	// Desensitize the whole string
	if Desensitize("123456789", 0, 8) != "*********" {
		t.Error("Desensitizing the whole string failed")
	}

	// Desensitize a single character
	if Desensitize("123456789", 3, 3) != "123*56789" {
		t.Error("Desensitizing a single character failed")
	}
}

func TestStringToIntSlice(t *testing.T) {
	// Basic conversion
	slice, err := StringToIntSlice("1,2,3", ",")
	if err != nil || len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("Basic conversion failed, got: %v, err: %v", slice, err)
	}

	// Empty string
	slice, err = StringToIntSlice("", ",")
	if err != nil || len(slice) != 0 {
		t.Errorf("Empty string conversion failed, got: %v, err: %v", slice, err)
	}

	// With a different delimiter
	slice, err = StringToIntSlice("4|5|6", "|")
	if err != nil || len(slice) != 3 || slice[0] != 4 || slice[1] != 5 || slice[2] != 6 {
		t.Errorf("Conversion with different delimiter failed, got: %v, err: %v", slice, err)
	}

	// Invalid number
	_, err = StringToIntSlice("1,a,3", ",")
	if err == nil {
		t.Error("Conversion with invalid number should have failed")
	}
}

func TestParseSort(t *testing.T) {
	// Ascending sort
	rule, col := ParseSort("ascending", "createTime", "id")
	if rule != "ASC" || col != "create_time" {
		t.Errorf("Ascending sort parsing failed, got rule: %s, col: %s", rule, col)
	}

	// Descending sort
	rule, col = ParseSort("descending", "updateTime", "id")
	if rule != "DESC" || col != "update_time" {
		t.Errorf("Descending sort parsing failed, got rule: %s, col: %s", rule, col)
	}

	// Default sort
	rule, col = ParseSort("descending", "", "defaultSort")
	if rule != "DESC" || col != "default_sort" {
		t.Errorf("Default sort parsing failed, got rule: %s, col: %s", rule, col)
	}

	// CamelCase to snake_case conversion
	rule, col = ParseSort("ascending", "myCustomField", "id")
	if col != "my_custom_field" {
		t.Errorf("CamelCase to snake_case conversion failed, got: %s", col)
	}
}
