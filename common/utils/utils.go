package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Regular expression validation
// expr: regular expression
// content: content to be verified
func CheckRegex(expr, content string) bool {

	r, err := regexp.Compile(expr)
	if err != nil {
		return false
	}

	return r.MatchString(content)
}

// Comparison tool
// Check if element item exists in slice
// If it exists, return true; if not, return false
func Contains[T comparable](slice []T, item T) bool {

	for _, value := range slice {
		if value == item {
			return true
		}
	}

	return false
}

// Filter
// If the condition function returns true, the element will be included in the result
func Filter[T interface{}](slice []T, condition func(T) bool) (result []T) {

	for _, value := range slice {
		if condition(value) {
			result = append(result, value)
		}
	}

	return result
}

// Desensitization tool
func Desensitize(content string, start, end int) string {

	if start < 0 || end < 0 || start > end {
		return content
	}

	var contentRune []rune

	for key, value := range content {
		if key >= start && key <= end {
			contentRune = append(contentRune, '*')
		} else {
			contentRune = append(contentRune, value)
		}
	}

	return string(contentRune)
}

// Convert string to int array
func StringToIntSlice(param, char string) ([]int, error) {

	intSlice := make([]int, 0)

	if param == "" {
		return intSlice, nil
	}

	stringSlice := strings.Split(param, char)

	for _, str := range stringSlice {

		num, err := strconv.Atoi(str)
		if err != nil {
			intSlice = append(intSlice, num)
			return nil, errors.New(str + " conversion failed: " + err.Error())
		}

		intSlice = append(intSlice, num)
	}

	return intSlice, nil
}
