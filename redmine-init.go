package main

import (
  "os"

  . "github.com/kazu634/redmine-sync/lib"
)

func main() {
  if len(os.Args) < 1 {
    os.Exit(1)
  }

  root := os.Args[1]

  RedmineMkdir(root)
}
