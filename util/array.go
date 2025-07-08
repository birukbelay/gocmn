package util

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
