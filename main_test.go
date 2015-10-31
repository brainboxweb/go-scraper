package main

import (
	"testing"
)


func TestStringToFloat(t *testing.T) {

	var tests = []struct {
		input    string
		expected float64
	}{
		{"\n£45\n\n", 45.0},
		{"rr£1.25rrr", 1.25},
		{"rr%ASD%ASD%ASD%A%SDASD%1.67", 1.67},
	}

	for _,test := range tests {
		actual := stringToFloat(test.input)
		if actual != test.expected {
			t.Errorf("stringToFloat(%s) = %f, expected %f.", test.input, actual, test.expected)
		}
	}
}
