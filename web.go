package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func indexPage(w http.ResponseWriter, req *http.Request) {
	preLen := len(*rootDir)
	err := filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = path[preLen:]
		if path == "" {
			return nil
		}

		_, file := filepath.Split(path)

		if info.IsDir() {
			if file[:1] == "." {
				return filepath.SkipDir
			}
			fmt.Fprintf(w, "<b>%s</b>\n", path)
		}

		if file[:1] == "." {
			return nil
		}
		fmt.Fprintf(w, "<a href='/%[1]s/f/%[2]s' download>%[2]s  %[3]d  %[4]s</a>\n", *accessKey, path, info.Size(), info.ModTime())
		return nil
	})

	if err != nil && err != filepath.SkipDir {
		log.Fatal("Error walking dir: ", err)
	}
}
