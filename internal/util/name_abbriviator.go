package util

import "strings"

func AbbreviateName(name string) string {
	parts := strings.Fields(name)

	if len(parts) < 2 {
		return name
	}

	firstName := parts[0]
	lastName := parts[len(parts)-1]
	abbreviatedLastName := string(lastName[0]) + "."
	abbreviatedName := firstName + " " + abbreviatedLastName

	if len(parts) > 2 {
		abbreviatedName += " " + strings.Join(parts[1:len(parts)-1], " ")
	}

	return abbreviatedName
}
