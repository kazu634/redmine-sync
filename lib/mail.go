package util

import (
  "fmt"
  "net/mail"
  "os"
  "path/filepath"
  "runtime"
  "strconv"
  "strings"
  "time"

  "github.com/DusanKasan/parsemail"
)

type Eml struct {
  Filename string
  SentDate time.Time
  Contents *parsemail.Email
}

func NewEml(target string) *Eml {
  fp, err := os.Open(target)
  if err != nil {
    panic(err)
  }
  defer fp.Close()

  fmt.Println(target)
  email, err := parsemail.Parse(fp) // returns Email struct and error
  if err != nil {
    panic(err)
  }

  return &Eml{Filename: target, SentDate: email.Date, Contents: &email}
}

func (e *Eml) Subject() string {
  return e.Contents.Subject
}

func (e *Eml) From() string {
  return concat(e.Contents.From)
}

func (e *Eml) To() string {
  return concat(e.Contents.To)
}

func (e *Eml) Cc() string {
  return concat(e.Contents.Cc)
}

func (e *Eml) Date() time.Time {
  return e.Contents.Date
}

func (e *Eml) Body() string {
  return cutMailBody(e.Contents.TextBody)
}

func (e *Eml) RedmineId() int {
  sep := ""
  if runtime.GOOS == "windows" {
    sep = "\\"
  } else {
    sep = "/"
  }

  dir := strings.Split(filepath.Dir(filepath.Clean(e.Filename)),  sep)
  idSubject := strings.Split(dir[len(dir)-1],  "_")

  id, _ := strconv.Atoi(idSubject[0])

  return id
}

func cutMailBody(body string) string {
  pos := strings.Index(body, "\nFrom:")
  if pos == -1 {
    return body
  }

  return convNewline(body[:pos], "\n")
}

func convNewline(str, nlcode string) string {
  return strings.NewReplacer(
    "\r\n\r\n\r\n", nlcode,
    "\r\r\r", nlcode,
    "\n\n\n", nlcode,
    "\r\n\r\n", nlcode,
    "\r\r", nlcode,
    "\n\n", nlcode,
  ).Replace(str)
}

func (e *Eml) GenMail() string {
  return fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: %s\nDate: %v\n\n%s", e.From(), e.To(), e.Cc(), e.Subject(), e.Date(), e.Body())
}

func concat(targets []*mail.Address) string {
  if len(targets) == 1 {
    return fmt.Sprintf("%s <%s>", targets[0].Name, targets[0].Address)
  }

  result := ""
  for _, target := range targets {
    tmp := fmt.Sprintf("%s <%s>", target.Name, target.Address)
    result = fmt.Sprintf("%s, %s", tmp, result)
  }

  return result
}

