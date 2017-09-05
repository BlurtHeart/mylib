package strutils

import "testing"

func TestAtoi(t *testing.T) {
	s := "12345"
	r, err := Atoi(s)
	if r != 12345 {
		t.Errorf("want:12345, got:%d", r)
	}
	if err != nil {
		t.Errorf("12345 convert want nil error, got err:%v\n", err)
	}

	s = "-1234"
	r, err = Atoi(s)
	if r != -1234 {
		t.Errorf("want:-1234, got:%d", r)
	}
	if err != nil {
		t.Errorf("-1234 convert want nil error, got %v", err)
	}

	s = "abcd"
	_, err = Atoi(s)
	if err != InvalidSyntax {
		t.Errorf("parse %s, want invalid syntax, got %v", s, err)
	}

	s = "+2147483648"
	_, err = Atoi(s)
	if err != ValueOutOfRange {
		t.Errorf("parse %s, want value out of range, got %v", s, err)
	}

	s = ""
	_, err = Atoi(s)
	if err != InvalidSyntax {
		t.Errorf("parse %s, want invalid syntax, got %v", s, err)
	}

	s = "-2147483649"
	_, err = Atoi(s)
	if err != ValueOutOfRange {
		t.Errorf("parse %s, want value out of range, got %v", s, err)
	}

	s = "-2147483648"
	r, err = Atoi(s)
	if r != -2147483648 {
		t.Errorf("want %d, got %d", s, r)
	}
	if err != nil {
		t.Errorf("want %d, got %v", s, err)
	}
}
