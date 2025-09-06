package main

import "testing"

type testCase struct {
	input    string
	expected []string
}

func TestCleanInput(t *testing.T) {
	cases := []testCase{
		{input: " hello  world  ", expected: []string{"hello", "world"}},
		{input: "one", expected: []string{"one"}},
		{input: "  ", expected: []string{}},
	}

	for _, tc := range cases {
		actual := cleanInput(tc.input)
		if len(actual) != len(tc.expected) {
			t.Errorf("input %q: expected len %d, got %d", tc.input, len(tc.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != tc.expected[i] {
				t.Errorf("input %q: at %d, got %q, want %q", tc.input, i, actual[i], tc.expected[i])
				break
			}
		}
	}
}
