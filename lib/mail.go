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
	Filename string           // ファイル名
	SentDate time.Time        // 送信日付
	Contents *enmime.Envelope // メールファイルの中身
}

// メールファイルの送信日付を time.Time に変換
func convDate(dateStr string) time.Time {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		layout := "Mon, 2 Jan 2006 15:04:05 -0700"
		t, err := time.Parse(layout, dateStr)

		if err != nil {
			panic(err)
		}
	}

	return t
}

// 新しい Eml を作成する。引数は *.eml ファイル
func NewEml(target string) *Eml {
	// メールファイルを開く
	file, err := os.Open(target)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// メールファイルを読み込む
	email := loadEmail(file)

	// Eml を返す
	return &Eml{Filename: target, SentDate: convDate(email.GetHeader("Date")), Contents: email}
}

// *.emlファイルを実際に読み込む
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

// *.emlファイルの格納ディレクトリに含まれる issue ID を取得する
func (e *Eml) RedmineId() int {
	sep := ""
	if runtime.GOOS == "windows" {
		sep = "\\"
	} else {
		sep = "/"
	}

	// メールファイルの格納ディレクトリを取得
	dir := strings.Split(filepath.Dir(filepath.Clean(e.Filename)), sep)
	idSubject := strings.Split(dir[len(dir)-1], "_")

	id, _ := strconv.Atoi(idSubject[0])

	return id
}

// 返信元のメールは削除する
func cutMailBody(body string) string {
	// メール本文でFrom:で始まる行があったら、その前までを取得する
	pos := strings.Index(body, "\nFrom:")
	if pos == -1 {
		return body
	}

	// 改行などを整える
	return convNewline(body[:pos], "\n")
}

// 改行などを整える
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

// メールの内容を文字列で返す
func (e *Eml) GenMail() string {
	return fmt.Sprintf("From: %s\nTo: %s\nCc: %s\nSubject: %s\nDate: %v\n\n%s", e.From(), e.To(), e.Cc(), e.Subject(), e.Date(), e.Body())
}
