package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	pt "github.com/monochromegane/the_platinum_searcher"
)

type custWriter struct {
	foundSomething bool
}

func (cw *custWriter) Write(p []byte) (n int, err error) {
	cw.foundSomething = true
	return 0, nil
}

func main() {
	rootPath := flag.String("rootPath", ".", "project root path")
	extension := flag.String("extension", ".", "extension filter")
	flag.Parse()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var unusedFiles []string

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), *extension) {
			fmt.Printf("searching for %s", f.Name())
			res := findStringInDir(*rootPath, strings.TrimRight(f.Name(), *extension))
			if !res {
				unusedFiles = append(unusedFiles, f.Name())
				fmt.Print(" - not found")
			}
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\nUnused files:\n")
	for _, v := range unusedFiles {
		fmt.Println(v)
	}
}

func findStringInDir(dir string, pattern string) bool {
	args := []string{"-c", "--ignore=vendor", pattern, dir}
	cw := custWriter{}
	pt := pt.PlatinumSearcher{Out: &cw, Err: os.Stderr}
	exitCode := pt.Run(args)
	if exitCode != 0 {
		log.Fatal(exitCode)
	}
	return cw.foundSomething
}
