package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jhillyerd/enmime"
)

type Eml struct {
	Filename string
	SentDate time.Time
	Contents *enmime.Envelope
}

func convDate(dateStr string) time.Time {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		panic(err)
	}

	return t
}

func NewEml(target string) *Eml {
	file, err := os.Open(target)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	email := loadEmail(file)
	return &Eml{Filename: target, SentDate: convDate(email.GetHeader("Date")), Contents: email}
}

func loadEmail(reader io.Reader) *enmime.Envelope {
	email, err := enmime.ReadEnvelope(reader) // returns Email struct and error
	if err != nil {
		panic(err)
	}

	return email
}

func (e *Eml) Subject() string {
	return fmt.Sprintf("%s", e.Contents.GetHeader("Subject"))
}

func (e *Eml) From() string {
	return fmt.Sprintf("%s", e.Contents.GetHeader("From"))
}

func (e *Eml) To() string {
	return fmt.Sprintf("%s", e.Contents.GetHeader("To"))
}

func (e *Eml) Cc() string {
	return fmt.Sprintf("%s", e.Contents.GetHeader("Cc"))
}

func (e *Eml) Date() time.Time {
	return convDate(e.Contents.GetHeader("Date"))
}

func (e *Eml) Body() string {
	return cutMailBody(e.Contents.Text)
}

func (e *Eml) RedmineId() int {
	sep := ""
	if runtime.GOOS == "windows" {
		sep = "\\"
	} else {
		sep = "/"
	}

	dir := strings.Split(filepath.Dir(filepath.Clean(e.Filename)), sep)
	idSubject := strings.Split(dir[len(dir)-1], "_")

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
