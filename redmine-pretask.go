package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Redmineのissueディレクトリを列挙する
func listOfIssueDirs(dir string) []string {
	// 引数で渡されたディレクトリからファイル・ディレクトリ一覧を読み込む
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// ディレクトリであれば、resultに格納する
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

	// 引数に存在するディレクトリを指定しているかを確認する
	root := os.Args[1]
	if _, err := os.Stat(root); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	// Redmineのissueディレクトリの一覧を取得する
	issueDirs := listOfIssueDirs(root)
	// 各々のissueディレクトリに対して
	for _, issueDir := range issueDirs {
		target := filepath.Join(issueDir, "*.eml")

		// *.emlを検索して
		mails, _ := filepath.Glob(target)
		// 各々の*.emlに対して
		for index, mail := range mails {
			// ファイルを開く
			fp, err := os.Open(mail)
			if err != nil {
				panic(err)
			}

			// 一時ファイルを作成する
			fw, err := ioutil.TempFile(issueDir, "foo")
			if err != nil {
				panic(err)
			}

			// *.emlを読み込んで、処理していく
			reader := bufio.NewReader(fp)
			for {
				// 1行ずつ読み込む
				line, err := reader.ReadBytes('\n')
				if err != nil && err != io.EOF {
					fmt.Printf("Reader error: %q\n", err)
					return
				}

				// Content-Type: multipart/relatedがあると処理がうまくいかないため、
				// Content-Type: multipart/mixedに変更する
				line = bytes.Replace(line, []byte("Content-Type: multipart/related;"), []byte("Content-Type: multipart/mixed;"), -1)

				// 一時ファイルに書き込み
				_, _ = fw.Write(line)

				// 最終行まで読み込んでいれば、無限ループを抜ける
				allLinesProcessed := err == io.EOF
				if allLinesProcessed {
					break
				}
			}

			// 以降ではオリジナルの*.emlファイルを削除して、
			// 一時ファイルを該当のディレクトリにコピーしている
			src := fw.Name()
			tmp := filepath.Join(issueDir, fmt.Sprintf("%05d.eml", index))

			fp.Close()
			fw.Close()

			// 一時ファイルをオリジナルファイルのディレクトリに移動
			if err := os.Rename(src, tmp); err != nil {
				panic(err)
			}

			// オリジナルの*.emlを削除
			if err := os.Remove(mail); err != nil {
				panic(err)
			}
		}
	}
}
