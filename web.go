package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming from", r.RemoteAddr)
	r.ParseMultipartForm(32 << 20)

	formfile, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer formfile.Close()

	// sanitize before saving
	// TODO strip escape codes and other dangerous char code points
	name := strings.Replace(header.Filename, string(os.PathSeparator), "-", -1)

	f, err := os.OpenFile(*rootDir+"/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, formfile)
}
