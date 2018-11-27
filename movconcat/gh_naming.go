package main

import (
	"regexp"
	"sort"
	"strings"
)

const ghChapNumRePtn = `GH(\d\d)(\d\d\d\d).MP4`

var (
	ghChapNumRe = regexp.MustCompile(ghChapNumRePtn)
)

// ghChapter makes filename map of file number chapters
func ghChapter(files []string) map[string][]string {
	r := make(map[string][]string)
	for _, f := range files {
		// chapterd videos
		if strings.HasPrefix(f, "GH") {
			m := ghChapNumRe.FindStringSubmatch(f)
			if len(m) == 3 {
				fn := m[2]
				r[fn] = append(r[fn], f)

			}
		}
	}

	for k := range r {
		sort.Strings(r[k])
	}

	return r
}
