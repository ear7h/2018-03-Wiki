package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"gopkg.in/russross/blackfriday.v2"
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

func (h *WikiHandler) post(r *http.Request) (io.Reader, error) {

	abs := h.absPath(r.URL.Path)
	file, err := os.Create(abs)
	if err != nil {
		// handle error
		return nil, err
	}

	io.Copy(file, r.Body)

	return nil, nil
}

func (h *WikiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data io.Reader
	var err error

	switch r.Method {
	case http.MethodGet: // handle get request
		data, err = h.get(r.URL.Path)

	case http.MethodPost:
		data, err = h.post(r)

		http.Redirect(w, r, r.URL.String(), http.StatusPermanentRedirect)
		return
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
