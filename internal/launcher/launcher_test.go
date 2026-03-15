package launcher

import (
	"errors"
	"testing"
)

func TestMapToEnv(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: []string{},
		},
		{
			name:     "single value",
			input:    map[string]string{"KEY": "value"},
			expected: []string{"KEY=value"},
		},
		{
			name:     "multiple values",
			input:    map[string]string{"A": "1", "B": "2"},
			expected: []string{"A=1", "B=2"},
		},
		{
			name:     "values with equals",
			input:    map[string]string{"KEY": "a=b"},
			expected: []string{"KEY=a=b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapToEnv(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("mapToEnv() returned %d items, want %d", len(got), len(tt.expected))
				return
			}
			for i, expected := range tt.expected {
				if got[i] != expected {
					t.Errorf("mapToEnv()[%d] = %q, want %q", i, got[i], expected)
				}
			}
		})
	}
}

func TestJoinErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    []error
		expected string
	}{
		{
			name:     "empty slice",
			input:    []error{},
			expected: "launcher errors:\n",
		},
		{
			name:     "single error",
			input:    []error{errors.New("error one")},
			expected: "launcher errors:\nerror one",
		},
		{
			name:     "multiple errors",
			input:    []error{errors.New("error one"), errors.New("error two")},
			expected: "launcher errors:\nerror one\nerror two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinErrors(tt.input)
			if got == nil && tt.expected != "" {
				t.Errorf("joinErrors() = nil, want %q", tt.expected)
				return
			}
			if got != nil && got.Error() != tt.expected {
				t.Errorf("joinErrors() = %q, want %q", got.Error(), tt.expected)
			}
		})
	}
}

func TestReplaceVar(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		key      string
		val      string
		expected string
	}{
		{
			name:     "no replacement needed",
			s:        "hello world",
			key:      "FOO",
			val:      "bar",
			expected: "hello world",
		},
		{
			name:     "simple replacement",
			s:        "hello %FOO%",
			key:      "%FOO%",
			val:      "world",
			expected: "hello world",
		},
		{
			name:     "multiple replacements",
			s:        "%FOO% and %FOO%",
			key:      "%FOO%",
			val:      "X",
			expected: "X and X",
		},
		{
			name:     "no match",
			s:        "%FOO%",
			key:      "%BAR%",
			val:      "X",
			expected: "%FOO%",
		},
		{
			name:     "empty value",
			s:        "hello %FOO%",
			key:      "%FOO%",
			val:      "",
			expected: "hello ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceVar(tt.s, tt.key, tt.val)
			if got != tt.expected {
				t.Errorf("replaceVar(%q, %q, %q) = %q, want %q", tt.s, tt.key, tt.val, got, tt.expected)
			}
		})
	}
}
