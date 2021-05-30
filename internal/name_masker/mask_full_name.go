package name_masker

import "strings"

func MaskFullName(name string) string {
	var build strings.Builder

	s := strings.TrimSpace(name)
	split := strings.Fields(s)

	lastName := split[len(split)-1]
	rest := split[0 : len(split)-1]

	for _, n := range rest {
		build.WriteString(MaskFirstName(n))
		build.WriteString(" ")
	}

	build.WriteString(MaskLastName(lastName))

	return build.String()
}
