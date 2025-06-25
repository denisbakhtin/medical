// Package helpers contains helper utility functions, that are common to several other packages (controllers, views, models)
package helpers

import (
	"strconv"
)

// ErrorData returns template data for an error
func ErrorData(err error) map[string]any {
	return map[string]any{
		"Title": err.Error(),
		"Error": err.Error(),
	}
}

// Atoi64 converts string to int64, returns 0 if error
func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// Atouint converts string to uint, returns 0 if error
func Atouint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

// Atoi64r converts string to int64 reference
func Atoi64r(s string) *int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &i
}

// Atouintr converts string to uint reference
func Atouintr(s string) *uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	ui := uint(i)
	return &ui
}

// Atob converts string to bool
func Atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

// CountFunc counts the number of elements in slice, that satisfy filter func
func CountFunc[E any](s []E, filter func(E) bool) int {
	count := 0
	for i := range s {
		if filter(s[i]) {
			count++
		}
	}
	return count
}
