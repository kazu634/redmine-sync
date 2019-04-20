package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "sort"

  . "github.com/kazu634/redmine-sync/lib"
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
    target := filepath.Join(issueDir, "*.eml")

    mails, _ := filepath.Glob(target)

    var emls []*Eml
    for _, mail := range mails {
      emls = append(emls, NewEml(mail))
    }

    sort.Slice(emls, func(i, j int) bool {
      return emls[i].SentDate.Unix() < emls[j].SentDate.Unix()
    })

    for _, eml := range emls {
      id := eml.RedmineId()

      fmt.Println(eml.Filename)
      UploadNote(id, eml.GenMail())
    }
  }
}
