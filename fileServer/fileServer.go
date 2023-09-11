package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// Task description:
//POST /<file-system-path> # save a file on the file system according to the specified path, e.g. /directory1/directory2/filename.txt
//GET /<file-system-path> # serve a file
//DELETE /<file-system-path> # delete a file on the file system

var workingDirectory = "/workDir"
var port = "8080"

var absenceOfFile = "No such file or directory"

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		fmt.Printf("[%s] %s %s\n", startTime.Format("15:04:05"), r.Method, r.URL.Path)
		next(w, r)
	}
}

func main() {
	port = *flag.String("port", port, "Port to start the file server")
	workingDirectory = *flag.String("workDir", workingDirectory, "Working directory")
	flag.Parse()

	address := ":" + port
	http.ListenAndServe(address, logger(handleRequest))
}

func handleRequest(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		saveFile(rw, req)
	case http.MethodGet:
		getFile(rw, req)
	case http.MethodDelete:
		removeFile(rw, req)
	default:
		http.Error(rw, "Method isn't allowed", http.StatusMethodNotAllowed)
	}
}

func saveFile(rw http.ResponseWriter, req *http.Request) {
	filePath := path.Join(workingDirectory, req.URL.Path)
	dir := path.Dir(filePath)

	if _, err := os.Stat(filePath); err == nil {
		http.Error(rw, "File already exists", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "File %s uploaded successfully", filePath)
}

func removeFile(rw http.ResponseWriter, req *http.Request) {
	filePath := path.Join(workingDirectory, req.URL.Path)

	if _, err := os.Stat(filePath); err != nil {
		http.Error(rw, absenceOfFile, http.StatusBadRequest)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "File %s is removed", filePath)
}

func getFile(rw http.ResponseWriter, req *http.Request) {
	filePath := path.Join(workingDirectory, req.URL.Path)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(rw, absenceOfFile, http.StatusNotFound)
		return
	}
	defer file.Close()

	_, err = io.Copy(rw, file)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
