package main

import "strconv"

var (
	binScal = struct {
		n    int
		name []string
	}{
		n: 3,
		name: []string{
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

func scale(v int) string {
	n := strconv.Itoa(v)
	l := len(n)

	mIdx := l / binScal.n
	if mIdx >= len(binScal.name) {
		// too big to masure scale
		return n
	}

	if l%binScal.n == 0 {
		mIdx--
	}

	if mIdx == 0 {
		return n
	}

	metric := binScal.name[mIdx]
	hTo := l - (mIdx * binScal.n)
	head := n[:hTo] + "." + n[hTo:hTo+1]

	return head + " " + metric
}
