package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	. "github.com/kazu634/redmine-sync/lib"
)

// Redmineのissueディレクトリを列挙する
func listOfIssueDirs(dir string) []string {
	// 引数で渡されたディレクトリからファイル・ディレクトリの一覧を取得
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// ディレクトリのみを結果として返す
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
	// 引数の数をチェックする
	if len(os.Args) < 1 {
		os.Exit(1)
	}

	// 第一引数に指定されたディレクトリが実際に存在するかを確認
	root := os.Args[1]
	if _, err := os.Stat(root); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	// Redmineのissueディレクトリを列挙する
	issueDirs := listOfIssueDirs(root)
	// 各々のissueディレクトリに対して
	for _, issueDir := range issueDirs {
		// *.emlを列挙する
		target := filepath.Join(issueDir, "*.eml")
		mails, _ := filepath.Glob(target)

		// *.emlを配列に格納
		var emls []*Eml
		for _, mail := range mails {
			emls = append(emls, NewEml(mail))
		}

		// メールファイルの送信日付に基づいてソートする
		sort.Slice(emls, func(i, j int) bool {
			return emls[i].SentDate.Unix() < emls[j].SentDate.Unix()
		})

		// ソート済みのメールファイルリストに対して
		for _, eml := range emls {
			// Redmineのissue IDを取得
			id := eml.RedmineId()

			fmt.Println(eml.Filename)
			// Redmineにアップロードする
			UploadNote(id, eml.GenMail())
		}
	}
}
