package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func listOfIssueDirs(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var results []string
	for _, file := range files {
		if file.IsDir() {
			results = append(results, filepath.Join(dir, file.Name()))
			continue
		}
	}

	return results
}

func main() {
	if len(os.Args) < 1 {
		os.Exit(1)
	}

	root := os.Args[1]
	if _, err := os.Stat(root); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	issueDirs := listOfIssueDirs(root)
	for _, issueDir := range issueDirs {
		target := filepath.Join(issueDir, "*.msg")

		mails, _ := filepath.Glob(target)
		for index, src := range mails {
			target := filepath.Join(issueDir, fmt.Sprintf("%07d.msg", index))

			if err := os.Rename(src, target); err != nil {
				panic(err)
			}
		}
	}
}
