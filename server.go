package main

import (
	"net/http"
	"io"
	"path"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"bytes"
)

const MarkDownPages = "markdownpages"

type WikiHandler struct {
	root string
}

func (h *WikiHandler) absPath(p string) string {
	return path.Join(h.root, p)
}

func (h *WikiHandler) get(path string) (io.Reader, error) {
	abs := h.absPath(path)

	data, err := ioutil.ReadFile(abs)
	if err != nil {
		// handle error
		return nil, err
	}

	reader := bytes.NewReader(
		blackfriday.Run(data))

	return reader, nil

}

func (h *WikiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {


	var data io.Reader
	var err error

	switch r.Method {
	case http.MethodGet: // handle get request
		data, err = h.get(r.URL.Path)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, data)
}

func main() {
	handler := &WikiHandler{
		root: MarkDownPages,
	}

	http.ListenAndServe(":8080", handler)
}
