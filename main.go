// Sweetdrop
// Zero configuration, self hosted, file sharing.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/net/http2"
)

const VERSION = "0.1.0"

var (
	help = flag.Bool("h", false, "Show usage")

	h1Port  = flag.String("h1", ":53370", "HTTP port, set blank '' to disable")
	h2Port  = flag.String("h2", ":53371", "HTTPS & HTTP2 port, set blank '' to disable")
	rootDir = flag.String("root", "", "Share directory root, defaults to current")

	accessKey = flag.String("key", "", "Sets access key, defaults to random 32 char")
)

func main() {
	setup()
	startWebServer()
	select {}
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func startWebServer() {
	key := "/" + *accessKey + "/"
	mux := http.NewServeMux()
	mux.HandleFunc(key+"upload", upload)
	mux.Handle(key, http.StripPrefix(key, http.FileServer(assetFS())))
	mux.Handle(key+"f/", http.StripPrefix(key+"f/", http.FileServer(http.Dir(*rootDir))))

	if *h1Port != "" {
		log.Println("http://localhost" + *h1Port + key)
		s := &http.Server{
			Addr:    *h1Port,
			Handler: mux,
		}

		http2.ConfigureServer(s, nil)

		go func() {
			err := s.ListenAndServe()
			if err != nil {
				log.Fatal("HTTP server error:", err)
			}
		}()
	}

	if *h2Port != "" {
		log.Println("https://localhost" + *h2Port + key)

		s := &http.Server{
			Addr:      *h2Port,
			Handler:   mux,
			TLSConfig: CreateTLS(),
		}

		http2.ConfigureServer(s, nil)

		go func() {
			err := s.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatal("HTTPS/2 server error:", err)
			}
		}()
	}
}

func setup() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = func() {
		fmt.Println("Sweetdrop ", VERSION)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *rootDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		rootDir = &cwd
	} else {
		*rootDir = filepath.Clean(*rootDir)
		_, err := os.Stat(*rootDir)
		if err != nil {
			log.Fatal("Root directory not valid")
		}
	}
	log.Println("Sweetdrop "+VERSION+" sharing", *rootDir)

	if *accessKey == "" {
		*accessKey = randSeq(32)
	}
}
