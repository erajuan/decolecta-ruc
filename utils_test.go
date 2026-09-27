package main

import "testing"

func TestValidRUC(t *testing.T) {
	for _, ruc := range []string{testRUC, "10123456781"} {
		if err := IsValidRuc(ruc); err != nil {
			t.Errorf("%s: %v", ruc, err)
		}
	}
}

func TestCreateRUCFromDNI(t *testing.T) {
	got, err := CreateRUCFromDNI("12345678")
	if err != nil {
		t.Fatal(err)
	}
	if got != "10123456781" {
		t.Fatalf("got %s", got)
	}
	if err := IsValidRuc(got); err != nil {
		t.Fatal(err)
	}
}

func TestLeftPadZero(t *testing.T) {
	for _, tc := range []struct {
		value  string
		length int
		want   string
	}{{"123", 6, "000123"}, {"150131", 6, "150131"}, {"1234567", 6, "1234567"}} {
		if got := LeftPadZero(tc.value, tc.length); got != tc.want {
			t.Errorf("LeftPadZero(%q, %d) = %q, want %q", tc.value, tc.length, got, tc.want)
		}
	}
}
