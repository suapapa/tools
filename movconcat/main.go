package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"
)

const movRePtn = `(\d\d\d\d_\d\d\d\d_\d\d\d\d\d\d)_\d\d\d.MOV`
const movTimeForm = "2006_0102_150405"

var (
	flagUseDocker = flag.Bool("d", false, "use docker")
)

func main() {
	flag.Parse()
	// list up MOVs
	// root := filepath.Join(os.Getenv("HOME"), "video/blackbox/raw")
	root := ""
	log.Println("searching MOVs from", root, "...")

	files, err := filepath.Glob(filepath.Join(root, "*.MOV"))
	panicIfErr(err)

	sort.Strings(files)
	log.Println(len(files), "MOVs found.")

	timeForm := regexp.MustCompile(movRePtn)

	// group MOVs by time
	var lastT time.Time
	var lastStartTime string
	movs := make(map[string][]string)
	for _, m := range files {
		// fmt.Println(filepath.Base(movs[0]))
		matchs := timeForm.FindStringSubmatch(m)
		if len(matchs) != 2 {
			log.Println("Skip", m)
		}

		currT, err2 := time.Parse(movTimeForm, matchs[1])
		panicIfErr(err2)

		if currT.Sub(lastT) > (10*time.Minute + 5*time.Second) {
			lastStartTime = matchs[1]
			log.Println("new recoding begins from", lastStartTime)
		}

		movs[lastStartTime] = append(movs[lastStartTime], m)

		lastT = currT
	}

	// concat MOVs
	var wg sync.WaitGroup
	wg.Add(len(movs))
	for k, v := range movs {
		go func(k string, v []string) {
			defer wg.Done()

			tmp, err := os.Create(k + ".list")
			panicIfErr(err)
			defer os.Remove(tmp.Name())

			log.Println("Create file list,", tmp.Name())
			for _, f := range v {
				fmt.Fprintln(tmp, "file", f)
			}
			tmp.Close()

			o := k + ".mov"

			stdErr, err := os.Create(k + ".log")
			panicIfErr(err)
			defer stdErr.Close()

			var cmd *exec.Cmd
			if *flagUseDocker {
				// cf) https://github.com/jrottenberg/ffmpeg
				// $ sudo docker run -v $PWD:/opt/data --rm \
				// jrottenberg/ffmpeg -f concat \
				// -i /opt/data/files.list -c copy /opt/data/test.mov
				cmd = exec.Command("docker", "run", "--rm",
					"-v", os.Getenv("PWD")+":/opt/data",
					"jrottenberg/ffmpeg",
					"-f", "concat",
					"-i", filepath.Join("/opt/data", tmp.Name()),
					"-c", "copy",
					filepath.Join("/opt/data", o),
				)
			} else {
				cmd = exec.Command("ffmpeg",
					"-f", "concat",
					"-i", tmp.Name(),
					"-c", "copy",
					o,
				)
			}

			cmd.Stderr = stdErr

			err = cmd.Start()
			panicIfErr(err)

			log.Println("concatting", o, "...")

			err = cmd.Wait()
			if err != nil {
				log.Printf("failed to concat %s with error: %v\n", o, err)
			} else {
				log.Println(o, "finished.")
			}
		}(k, v)

	}

	wg.Wait()
	log.Println("all done!")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
