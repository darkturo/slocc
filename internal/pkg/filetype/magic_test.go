package filetype

import (
	"bufio"
	"strings"
	"testing"
)

// Test the Detect function
func TestDetect(t *testing.T) {
	// iterate over the test cases
	for _, tc := range []struct {
		input    string
		expected FileType
	}{
		{"#!/usr/bin/env python", Python},
		{"#!/usr/bin/python", Python},
		{"#!/usr/bin/python3", Python},
		{"#!/usr/bin/env perl", Perl},
		{"#!/usr/bin/perl", Perl},
		{"#!/usr/bin/env ruby", Ruby},
		{"#!/usr/bin/ruby", Ruby},
		{"#!/usr/bin/bash", Bash},
		{"#!/usr/bash", Bash},
		{"this is just a text", Other},
		{"", Other},
	} {
		// run the test case
		if result := Detect(bufio.NewReader(strings.NewReader(tc.input))); result != tc.expected {
			t.Errorf("Detect(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}
