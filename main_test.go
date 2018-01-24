package main

import (
	"bytes"
	"io"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	c := GenerateCode()
	for k, v := range c {
		if v > 6 || v < 1 {
			t.Errorf("wrong color generated :key '%d' and value '%d'", k, v)
		}
	}
}

func TestCompareCode(t *testing.T) {
	type compareTest struct {
		name      string
		try       Code
		solution  Code
		good      int
		misplaced int
	}

	tests := []compareTest{
		{"match", Code{1, 1, 1, 1}, Code{1, 1, 1, 1}, 4, 0},
		{"0 found", Code{1, 1, 1, 1}, Code{2, 2, 2, 2}, 0, 0},
		{"A", Code{1, 1, 1, 1}, Code{1, 1, 2, 2}, 2, 0},
		{"B", Code{1, 1, 1, 2}, Code{1, 2, 1, 1}, 2, 2},
		{"C", Code{1, 2, 2, 3}, Code{3, 1, 2, 2}, 1, 3},
		{"D", Code{1, 2, 3, 4}, Code{1, 3, 4, 2}, 1, 3},
		{"E", Code{1, 2, 3, 4}, Code{1, 3, 4, 5}, 1, 2},
		{"F", Code{1, 2, 3, 4}, Code{4, 3, 2, 1}, 0, 4},
		{"G", Code{1, 2, 3, 4}, Code{5, 3, 2, 1}, 0, 3},
		{"H", Code{2, 2, 3, 3}, Code{4, 4, 2, 2}, 0, 2},
		{"I", Code{1, 1, 2, 3}, Code{5, 4, 4, 2}, 0, 1},
		{"I", Code{1, 4, 5, 5}, Code{1, 5, 4, 6}, 1, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			good, misplaced := CompareCode(test.try, test.solution)
			if good != test.good {
				t.Errorf("wrong good : got '%d', expected '%d'", good, test.good)
			}
			if misplaced != test.misplaced {
				t.Errorf("wrong misplaced : got '%d', expected '%d'", misplaced, test.misplaced)
			}
		})
	}
}

func TestScanCode(t *testing.T) {
	type scanTest struct {
		name         string
		r            io.Reader
		expectedCode Code
		expectErr    bool
	}

	tests := []scanTest{
		{"good code", bytes.NewReader([]byte("2,2,2,2")), Code{2, 2, 2, 2}, false},
		{"not enough number", bytes.NewReader([]byte("2,2,2")), Code{2, 2, 2, 2}, true},
		{"wrong number 1", bytes.NewReader([]byte("9,2,2,2")), Code{2, 2, 2, 2}, true},
		{"wrong number 2", bytes.NewReader([]byte("12,2,2,2")), Code{2, 2, 2, 2}, true},
		{"fizz input", bytes.NewReader([]byte("dzaeazeaze")), Code{2, 2, 2, 2}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := ScanCode(test.r)
			if test.expectErr {
				if err == nil {
					t.Fatal("error expected, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("didn't expect any error, got '%s'", err.Error())
			}
			if code != test.expectedCode {
				t.Errorf("wrong code, expected %v, got %v", test.expectedCode, code)
			}
		})
	}
}
