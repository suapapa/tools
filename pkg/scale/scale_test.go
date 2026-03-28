package scale

import "testing"

func TestIsNum(t *testing.T) {
	tt := []struct {
		n string
		b bool
	}{
		{"10", true},
		{"a", false},
		{"", false},
	}

	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			if got := isNum(tc.n); got != tc.b {
				t.Errorf("isNum(%q) = %v, want %v", tc.n, got, tc.b)
			}
		})
	}
}

func TestConvertPowerForm(t *testing.T) {
	tt := []struct {
		p   string
		exp string
	}{
		{"1**0", "1"},
		{"10**0", "1"},
		{"10**10", "10000000000"},
		{"a**b", "a**b"},
	}
	for _, tc := range tt {
		t.Run(tc.p, func(t *testing.T) {
			if got := convertPowerForm(tc.p); got != tc.exp {
				t.Errorf("convertPowerForm(%q) = %v, want %v", tc.p, got, tc.exp)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tt := []struct {
		in     string
		scale  Scale
		expect string
	}{
		{"200", Ko, "200"},
		{"30000", Ko, "3만"},
		{"1000", En, "1Thousand"},
		{"1000000", En, "1Million"},
		{"1024", Bin, "1Kilo"}, // Wait, is Format doing binary correctly?
	}

	for _, tc := range tt {
		t.Run(tc.expect, func(t *testing.T) {
			got, err := Format(tc.in, tc.scale)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.expect {
				t.Errorf("Format(%q, %v) = %q, want %q", tc.in, tc.scale, got, tc.expect)
			}
		})
	}
}
