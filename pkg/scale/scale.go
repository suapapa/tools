package scale

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Scale represents a numeric scale (e.g., metric, binary, Korean)
type Scale struct {
	N    int
	Name []string
}

var (
	// Ko is the Korean scale (units of 10^4)
	Ko = Scale{
		N: 4,
		Name: []string{
			"",
			"만", // 10 ** 4
			"억", // 10 ** 8
			"조", // 10 ** 12
			"경", // 10 ** 16
			"해", // 10 ** 20
		},
	}

	// En is the standard English short scale (units of 10^3)
	En = Scale{
		N: 3,
		Name: []string{
			"",
			"Thousand",
			"Million",
			"Billion",
			"Trillion",
			"Quadrillion",
			"Quintillion",
			"Sextillion",
			"Septillion",
		},
	}

	// Bin is the binary scale (units of 10^3, used for BPS etc.)
	Bin = Scale{
		N: 3,
		Name: []string{
			"",
			"Kilo",
			"Mega",
			"Giga",
			"Tera",
			"Peta",
			"Exa",
			"Zetta",
			"Yotta",
		},
	}
)

// Format returns a human-readable string for the given number string and scale.
func Format(n string, s Scale) (string, error) {
	n = convertPowerForm(n)

	if !isNum(n) {
		return "", fmt.Errorf("invalid number: %v", n)
	}

	l := len(n)
	if l == 0 {
		return "", errors.New("empty input")
	}

	mIdx := l / s.N
	if mIdx >= len(s.Name) {
		return "", errors.New("too big for current scale")
	}

	if l%s.N == 0 {
		mIdx--
	}

	if mIdx < 0 {
		return n, nil
	}

	metric := s.Name[mIdx]
	hTo := l - (mIdx * s.N)
	head := n[:hTo]

	if metric == "" {
		return head, nil
	}

	return head + metric, nil
}

// FormatWithDecimal returns a human-readable string with one decimal point.
func FormatWithDecimal(n string, s Scale) (string, error) {
	if !isNum(n) {
		return "", fmt.Errorf("invalid number: %v", n)
	}

	l := len(n)
	if l == 0 {
		return "", errors.New("empty input")
	}

	mIdx := l / s.N
	if mIdx >= len(s.Name) {
		return n, nil // Too big, return original
	}

	if l%s.N == 0 {
		mIdx--
	}

	if mIdx <= 0 {
		return n, nil
	}

	metric := s.Name[mIdx]
	hTo := l - (mIdx * s.N)
	
	head := n[:hTo] + "." + n[hTo:hTo+1]
	return head + " " + metric, nil
}

func isNum(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !(('0' <= c) && (c <= '9')) {
			return false
		}
	}
	return true
}

func convertPowerForm(s string) string {
	numWithPower := regexp.MustCompile(`([0-9]+)\*\*([0-9]+)`)
	sm := numWithPower.FindStringSubmatch(s)

	if len(sm) != 3 {
		return s
	}

	b, _ := strconv.Atoi(sm[1])
	p, _ := strconv.Atoi(sm[2])
	r := 1
	for i := 0; i < p; i++ {
		r = r * b
	}

	return strconv.Itoa(r)
}
