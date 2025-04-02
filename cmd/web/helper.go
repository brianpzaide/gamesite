package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/knadh/stuffbin"
)

func initFS() stuffbin.FileSystem {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path %v", err)
	}
	fs, err := stuffbin.UnStuff(exe)
	if err != nil {
		log.Fatalf("error reading the stuffed binary %v", err)
	}

	// fmt.Println("loaded files", fs.List())
	_, err = fs.Get("/white.png")
	if err != nil {
		log.Fatalf("error reading white.png: %v", err)
	}

	return fs
}

func setHomePage(fs stuffbin.FileSystem) {
	tpl, err := stuffbin.ParseTemplates(nil, fs, "/templates/wshome.tmpl")
	if err != nil {
		log.Fatalf("error parsing templates: %v", err)
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, templateData)
	if err != nil {
		log.Fatalf("error rendering templates: %v", err)
	}
	homePage = buf.String()

}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
