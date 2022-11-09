package filetype

import (
	"bufio"
	"strings"
	"testing"
)

// Test the shebang function
func TestDetect(t *testing.T) {
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
		if result := shebang(bufio.NewReader(strings.NewReader(tc.input))); result != tc.expected {
			t.Errorf("shebang(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

// Test the looksLikeBinary function
func TestLooksLikeBinary(t *testing.T) {
	for _, tc := range []struct {
		input    []byte
		expected bool
	}{
		{[]byte("this is just a text"), false},
		{[]byte("this is a text with a \x00 byte"), true},
		{[]byte(""), false},
	} {
		// run the test case
		if result := looksLikeBinary(bufio.NewReader(strings.NewReader(string(tc.input)))); result != tc.expected {
			t.Errorf("looksLikeBinary(%s) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
