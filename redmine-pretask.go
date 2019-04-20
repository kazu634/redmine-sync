package main

import (
  "bytes"
  "bufio"
  "fmt"
  "io"
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
    target := filepath.Join(issueDir, "*.eml")

    mails, _ := filepath.Glob(target)
    for index, mail := range mails {
      fp, err := os.Open(mail)
      if err != nil {
        panic(err)
      }

      fw, err := ioutil.TempFile(issueDir, "foo")
      if err != nil {
        panic(err)
      }

      reader := bufio.NewReader(fp)
      for {
        line, err := reader.ReadBytes('\n')
        if err != nil && err != io.EOF {
          fmt.Printf("Reader error: %q\n", err)
          return
        }

        line = bytes.Replace(line, []byte("Content-Type: multipart/related;"), []byte("Content-Type: multipart/mixed;"), -1)

        _, _ = fw.Write(line)

        allLinesProcessed := err == io.EOF
        if allLinesProcessed {
            break
        }
      }

      src := fw.Name()
      tmp := filepath.Join(issueDir, fmt.Sprintf("%05d.eml", index))

      fp.Close()
      fw.Close()

      if err := os.Rename(src, tmp); err != nil {
        panic(err)
      }

      if err := os.Remove(mail); err != nil {
        panic(err)
      }
    }
  }
}


