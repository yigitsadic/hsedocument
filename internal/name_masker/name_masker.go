package name_masker

import (
	"unicode/utf8"
)

// Verilen girdiyi isim olarak maskeliyor.
func MaskFirstName(firstName string) string {
	if utf8.RuneCountInString(firstName) < 2 {
		return "*"
	}

	var maskedPart, baseString, result string

	if len(firstName) == 2 || utf8.RuneCountInString(firstName) == 3 {
		maskedPart = "*"
	} else if utf8.RuneCountInString(firstName) == 4 {
		maskedPart = "**"
	} else {
		maskedPart = "***"
	}

	if utf8.RuneCountInString(firstName) == 2 {
		baseString = string([]rune(firstName)[0:1])
	} else {
		baseString = string([]rune(firstName)[0:2])
	}

	result = baseString + maskedPart

	return result
}

// Verilen girdiyi soyisim olarak maskeliyor.
func MaskLastName(lastName string) string {
	if utf8.RuneCountInString(lastName) < 2 {
		return "*"
	}

	var maskedPart, visiblePart, result string

	if utf8.RuneCountInString(lastName) < 3 {
		maskedPart = "*"
	} else if utf8.RuneCountInString(lastName) <= 4 {
		maskedPart = "**"
	} else {
		maskedPart = "***"
	}

	if utf8.RuneCountInString(lastName) < 4 {
		visiblePart = string([]rune(lastName)[utf8.RuneCountInString(lastName)-1])
	} else {
		visiblePart = string([]rune(lastName)[utf8.RuneCountInString(lastName)-2:])
	}

	result = maskedPart + visiblePart

	return result
}
