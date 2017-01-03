package main

import (
	"regexp"
	"strings"
)

const gpFileNumRePtn = `GOPR(\d\d\d\d).MP4`
const gpChapNumRePtn = `GP(\d\d)(\d\d\d\d).MP4`

var (
	gpFileNumRe = regexp.MustCompile(gpFileNumRePtn)
	gpChapNumRe = regexp.MustCompile(gpChapNumRePtn)
)

// gpChapter makes filename map of file number chapters
func gpChapter(files []string) map[string][]string {
	r := make(map[string][]string)
	for _, f := range files {
		// single video
		if strings.HasPrefix(f, "GOPR") {
			m := gpFileNumRe.FindStringSubmatch(f)
			if len(m) == 2 {
				fn := m[1]
				r[fn] = append(r[fn], f)
			}
			continue
		}

		// chapterd videos
		if strings.HasPrefix(f, "GP") {
			m := gpChapNumRe.FindStringSubmatch(f)
			if len(m) == 3 {
				fn := m[2]
				r[fn] = append(r[fn], f)

			}

		}
	}

	return r
}
