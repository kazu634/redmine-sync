package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// RedmineのIssueディレクトリを列挙する
func listOfIssueDirs(dir string) []string {
	// 引数で渡されたディレクトリを読み込み
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// 取得したディレクトリ・ファイルのリストから、
	// ディレクトリのみを result に格納する
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
	// 引数の数チェック
	if len(os.Args) < 1 {
		os.Exit(1)
	}

	// 引数が存在するディレクトリかどうかを確認する
	root := os.Args[1]
	if _, err := os.Stat(root); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	// Redmineのissueに該当するディレクトリを列挙する
	issueDirs := listOfIssueDirs(root)
	// issueディレクトリ配下の*.msgに対して
	for _, issueDir := range issueDirs {
		target := filepath.Join(issueDir, "*.msg")

		// 一括でファイル名を変換する
		mails, _ := filepath.Glob(target)
		for index, src := range mails {
			target := filepath.Join(issueDir, fmt.Sprintf("%07d.msg", index))

			if err := os.Rename(src, target); err != nil {
				panic(err)
			}
		}
	}
}
