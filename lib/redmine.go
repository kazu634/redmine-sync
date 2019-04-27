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

// クライアント生成に必要な設定
type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
	Project  int    `json:"project"`
}

var conf config
var profile string = "profile"

// 設定の取得
func getConfig() config {
	// 環境変数 REDMINEPROJECT の内容を取得
	proj, _ := strconv.Atoi(os.Getenv("REDMINEPROJECT"))

	// config を返す (環境変数REDMINEENDPOINT, REDMINEAPIKEYを参照する)
	return config{Endpoint: os.Getenv("REDMINEENDPOINT"),
		Apikey:  os.Getenv("REDMINEAPIKEY"),
		Project: proj}
}

// Issueにノートを追加する
func notesIssue(id int, content string) {
	// APIクライアントの作成
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)

	// issueの取得
	issue, err := c.Issue(id)
	if err != nil {
		log.Fatalf("Failed to update issue: %s\n", err)
	}

	// 引数で指定された文字列をノートに指定
	issue.Notes = content
	// Project IDを指定して
	issue.ProjectId = conf.Project
	// Issueをアップデートする
	err = c.UpdateIssue(*issue)
	if err != nil {
		log.Fatalf("Failed to update issue: %s\n", err)
	}
}

// ノートを追加する
func UploadNote(id int, note string) {
	// APIクライアント生成に必要な設定を取得
	conf = getConfig()

	// 念のためオレオレ証明書も許可しておく
	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 指定された issue IDに対して、noteで指定された文字列をアップロードする
	notesIssue(id, note)
}

// RedmineからIssue IDを取得して、ディレクトリを作成する
func RedmineMkdir(root string) {
	// APIクライアント取得のための設定取得
	conf = getConfig()

	// オレオレ証明書を許可しておく
	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// APIクライアントの作成
	c := redmine.NewClient(conf.Endpoint, conf.Apikey)

	// Issue一覧の取得
	issues, err := c.IssuesByFilter(nil)
	if err != nil {
		log.Fatalf("Failed to list issues: %s\n", err)
	}

	// 各issueにつき
	for _, i := range issues {
		// ディレクトリ名の生成
		dirname := fmt.Sprintf("%d", i.Id)

		target := filepath.Join(root, dirname)

		// ディレクトリ作成
		if err := os.MkdirAll(target, 0777); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Create %s.\n", target)
	}
}
