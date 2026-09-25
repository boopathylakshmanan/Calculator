package calc

import "testing"

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name string
		expr string
		want string
	}{
		{"addition", "2 + 2", "4"},
		{"precedence", "3 + 4 * 2", "11"},
		{"parentheses", "(3 + 4) * 2", "14"},
		{"nested parentheses", "2 * (3 + (4 - 1))", "12"},
		{"unary minus", "-5 + 3", "-2"},
		{"minus negative", "5 - -3", "8"},
		{"decimal division", "1 / 3", "0.3333333333"},
		{"exact division trims zeros", "10 / 4", "2.5"},
		{"decimal literal", "1.5 + 2.25", "3.75"},
		{"whitespace tolerant", "  2+2  ", "4"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Evaluate(tc.expr)
			if err != nil {
				t.Fatalf("Evaluate(%q) returned error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Fatalf("Evaluate(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	cases := []struct {
		name string
		expr string
	}{
		{"empty expression", ""},
		{"whitespace only", "   "},
		{"division by zero", "5 / 0"},
		{"trailing operator", "2 +"},
		{"invalid character", "2 + a"},
		{"unbalanced parens", "(3 + 4"},
		{"malformed number", "1..2 + 3"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Evaluate(tc.expr); err == nil {
				t.Fatalf("Evaluate(%q) expected an error, got nil", tc.expr)
			}
		})
	}
}
