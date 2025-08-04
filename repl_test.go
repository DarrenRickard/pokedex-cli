package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct{
		input string
		expected []string
	}{
		{
			input: "   hello   world   ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		actualLength := len(actual)
		expectedLength := len(c.expected)
		if actualLength != expectedLength {
			t.Errorf("length of input: %v does not match expected length: %v", actualLength, expectedLength)
			t.Errorf("input: %v", actual)
			t.Fail()
		}


		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word: %s does not match Expected: %s", word, expectedWord)
				t.Fail()
			}
		}
	}

}


