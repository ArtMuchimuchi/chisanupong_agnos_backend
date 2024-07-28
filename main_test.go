package main

import "testing"

func TestCheckPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected int
	}{
		{"Less than min (need character)", "aaa", 3},
		{"Less than min (not need character)", "..!..", 3},
		{"Between range", "Aa1234a", 0},
		{"Between range (not need character)", ".!.!.!", 3},
		{"Between range (repeating)", "AAAAAA", 2},
		{"Between range (repeating with no need character)", "!!!!!!", 3},
		{"Exceed range ", "Ab212312312312312123123123", 7},
		{"Exceed range (repeating)", "11111111111111111111111111", 13},
		{"Exceed range (repeating with no need characer)", "!!!!!!!!!!!!!!!!!!!!!!!", 10},
		{"Exceed range (no need characer)", "!.!.!.!.!.!.!.!.!.!2", 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getStrongSteps(tc.password)
			if result != tc.expected {
				t.Errorf("Test %s failed. Expected %d, got %d", tc.name, tc.expected, result)
			}
		})
	}
}
