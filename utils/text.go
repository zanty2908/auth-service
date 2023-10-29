package utils

func ReverseString(input string) string {
	// Convert the string to a slice of runes
	runes := []rune(input)

	// Calculate the length of the rune slice
	length := len(runes)

	// Reverse the elements in the slice
	for i := 0; i < length/2; i++ {
		runes[i], runes[length-i-1] = runes[length-i-1], runes[i]
	}

	// Convert the reversed slice back to a string
	reversedString := string(runes)

	return reversedString
}
