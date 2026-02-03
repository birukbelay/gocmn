package util

import "github.com/lib/pq"

func ElementExists[T comparable](element T, arr ...T) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

// RemoveFirstElement Remove a specific element from the slice if it exists
func RemoveFirstElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			// Remove the element by concatenating slices (excluding the element)
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // Return original slice if element is not found
}

// RemoveAllOccurrences removes all occurrences of a specific element from the slice
func RemoveAllOccurrences[T comparable](element T, slice []T) []T {
	var newSlice []T // Initialize a new slice
	for _, v := range slice {
		if v != element {
			newSlice = append(newSlice, v) // Only add elements that don't match
		}
	}
	return newSlice
}

// AddIfNotExists element to slice if it doesn't exist
func AddIfNotExists[T comparable](element T, slice []T) []T {
	if !ElementExists(element, slice...) {
		slice = append(slice, element) // Add the element if it doesn't exist
	}
	return slice
}

// Reverse is a function to reverse a slice
func Reverse[T any](arr []T) []T {
	// Two-pointer approach to swap elements
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i] // Swap elements at i and j
	}
	return arr
}

func Merge2ArrUniqRes[T comparable](a, b []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(a)+len(b))

	// Add all elements from 'a'
	for _, v := range a {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	// Add elements from 'b' only if not already seen
	for _, v := range b {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}

func MergeArrUniqRes[T comparable](vals ...[]T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)

	// Add all elements from 'a'
	for _, arr := range vals {
		for _, v := range arr {

			if !seen[v] {
				seen[v] = true
				result = append(result, v)
			}
		}
	}
	return result
}

func MergeArrAll[T comparable](vals ...[]T) []T {
	result := make([]T, 0)
	// Add all elements from 'a'
	for _, arr := range vals {
		for _, v := range arr {
			result = append(result, v)
		}
	}
	return result
}
func ElementExistsInStringArray(email string, arr pq.StringArray) bool {
	for _, v := range arr {
		if v == email {
			return true
		}
	}
	return false
}

func UniqElements[T comparable](baseArr, arrWithMoreItems []T) []T {
	seen := make(map[T]bool, len(baseArr))
	for _, v := range baseArr {
		seen[v] = true
	}
	var added []T
	for _, v := range arrWithMoreItems {
		if _, ok := seen[v]; !ok {
			added = append(added, v)
		}
	}
	return added
}
