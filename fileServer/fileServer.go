package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var workingDirectory = os.Getenv("WORK_DIR")
var port = os.Getenv("PORT")

var absenceOfFile = "No such file or directory"

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		fmt.Printf("[%s] %s %s\n", startTime.Format("15:04:05"), r.Method, r.URL.Path)
		next(w, r)
	}
}

func main() {
	println(fmt.Sprintf("Server is started in %s", workingDirectory))
	http.ListenAndServe(fmt.Sprintf(":%s", port), logger(handleRequest))
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
	if req.URL.Path == "/" {
		http.Error(rw, "File name is empty", http.StatusBadRequest)
		return
	}
	dir := path.Dir(filePath)
	println(filePath, dir)
	if _, err := os.Stat(filePath); err == nil {
		http.Error(rw, "File already exists", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			println("aboab stat")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		println("create")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, req.Body)
	if err != nil {
		println("copy")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "File %s uploaded successfully", filePath)
}

func removeFile(rw http.ResponseWriter, req *http.Request) {
	filePath := path.Join(workingDirectory, req.URL.Path)
	if req.URL.Path == "/" {
		http.Error(rw, "File name is empty", http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(filePath); err != nil {
		http.Error(rw, absenceOfFile, http.StatusBadRequest)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	removeEmptyDirectories(filepath.Dir(filePath))

	fmt.Fprintf(rw, "File %s is removed", filePath)
}

func removeEmptyDirectories(dirPath string) {
	for {
		if dirPath == workingDirectory {
			break
		}
		err := os.Remove(dirPath)
		if err != nil {
			break
		}
		dirPath = filepath.Dir(dirPath)
	}
}

func getFile(rw http.ResponseWriter, req *http.Request) {
	filePath := path.Join(workingDirectory, req.URL.Path)
	if req.URL.Path == "/" {
		http.Error(rw, "File name is empty", http.StatusBadRequest)
		return
	}
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
