package name_masker

import "testing"

func TestMaskFirstName(t *testing.T) {
	cases := []struct {
		Name     string
		Given    string
		Expected string
	}{
		{Name: "it should mask single character names", Given: "Y", Expected: "*"},
		{Name: "it should mask very short names", Given: "Yi", Expected: "Y*"},
		{Name: "it should mask short names", Given: "Can", Expected: "Ca*"},
		{Name: "it should mask short names", Given: "Kaan", Expected: "Ka**"},
		{Name: "it should mask medium names", Given: "Yiğit", Expected: "Yi***"},
		{Name: "it should mask upper case names", Given: "YİĞİT", Expected: "Yİ***"},
		{Name: "it should mask russian names", Given: "Володимѣръ", Expected: "Во***"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := MaskFirstName(c.Given)

			if got != c.Expected {
				t.Errorf("input=%s expected=%s but got=%s", c.Given, c.Expected, got)
			}
		})
	}
}

func TestMaskLastName(t *testing.T) {
	cases := []struct {
		Name     string
		Given    string
		Expected string
	}{
		{Name: "it should mask single character last names", Given: "Y", Expected: "*"},
		{Name: "it should mask two character last names", Given: "Yi", Expected: "*i"},
		{Name: "it should mask very short names", Given: "Can", Expected: "**n"},
		{Name: "it should mask short names", Given: "Kaan", Expected: "**an"},
		{Name: "it should mask medium names", Given: "Sadıç", Expected: "***ıç"},
		{Name: "it should mask upper case names", Given: "SADIÇ", Expected: "***IÇ"},
		{Name: "it should mask russian names", Given: "Володимѣръ", Expected: "***ръ"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := MaskLastName(c.Given)

			if got != c.Expected {
				t.Errorf("input=%s expected=%s but got=%s", c.Given, c.Expected, got)
			}
		})
	}
}
