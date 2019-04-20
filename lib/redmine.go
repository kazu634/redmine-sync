package util

import (
  "crypto/tls"
  "fmt"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strconv"

  "github.com/mattn/go-redmine"
)

type config struct {
  Endpoint string `json:"endpoint"`
  Apikey   string `json:"apikey"`
  Project  int    `json:"project"`
}

var conf config
var profile string = "profile"

func getConfig() config {
  proj, _ := strconv.Atoi(os.Getenv("REDMINEPROJECT"))

  return config{Endpoint: os.Getenv("REDMINEENDPOINT"),
              Apikey: os.Getenv("REDMINEAPIKEY"),
              Project: proj}
}

func notesIssue(id int, content string) {
  c := redmine.NewClient(conf.Endpoint, conf.Apikey)
  issue, err := c.Issue(id)
  if err != nil {
    log.Fatalf("Failed to update issue: %s\n", err)
  }

  issue.Notes = content
  issue.ProjectId = conf.Project
  err = c.UpdateIssue(*issue)
  if err != nil {
    log.Fatalf("Failed to update issue: %s\n", err)
  }
}

func UploadNote(id int, note string) {
  conf = getConfig()

  http.DefaultClient = &http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
  }

  notesIssue(id, note)
}

func RedmineMkdir(root string) {
  conf = getConfig()
  http.DefaultClient = &http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
  }

  c := redmine.NewClient(conf.Endpoint, conf.Apikey)

  issues, err := c.IssuesByFilter(nil)
  if err != nil {
    log.Fatalf("Failed to list issues: %s\n", err)
  }

  for _, i := range issues {
    dirname := fmt.Sprintf("%d", i.Id)

    target := filepath.Join(root, dirname)

    if err := os.MkdirAll(target, 0777);err != nil {
      log.Fatal(err)
    }

    fmt.Printf("Create %s.\n", target)
  }
}
