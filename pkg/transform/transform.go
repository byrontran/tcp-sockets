package transform

func Encode(input string) string {
	result := make([]byte, len(input))

	for index := range input {
		result[index] = input[index] + 1

		// TODO @byrontran: Remove
		// fmt.Printf("Char at index %d: %c\n", index, input[index])
		// fmt.Printf("Encoded char: %c\n", result[index])
	}

	return string(result)
}

// Optional functionality
// TODO @byrontran: Come back later and write this?
func Decode(intput string) string {
	return "TODO: Implement the Decode() function in `pkg/transform/transform.go`."
}
