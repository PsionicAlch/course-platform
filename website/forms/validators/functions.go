package validators

import "regexp"

func ContainsSpecialCharacters(data string) bool {
	specialChars := `[!@#$%^&*()_+\-=\[\]{};:'",.<>/?\\|]`
	matched, err := regexp.MatchString(specialChars, data)
	if err != nil {
		return false
	}
	return matched
}

func ContainsNumbers(data string) bool {
	numbers := `[0-9]`
	matched, err := regexp.MatchString(numbers, data)
	if err != nil {
		return false
	}
	return matched
}

func ContainsUppercaseCharacters(data string) bool {
	numbers := `[A-Z]`
	matched, err := regexp.MatchString(numbers, data)
	if err != nil {
		return false
	}
	return matched
}

func ContainsLowercaseCharacters(data string) bool {
	numbers := `[a-z]`
	matched, err := regexp.MatchString(numbers, data)
	if err != nil {
		return false
	}
	return matched
}
