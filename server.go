package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"gopkg.in/russross/blackfriday.v2"
	"log"
	"fmt"
	"strings"
)

const MarkDownPages = "markdownpages"
var PageFormat string

func init() {

	byt, err := ioutil.ReadFile("page.html")
	if err != nil {
		panic(err)
	}

	PageFormat = string(byt)
}

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

	reader := strings.NewReader(
		fmt.Sprintf(PageFormat, blackfriday.Run(data)))

	return reader, nil

}

func (h *WikiHandler) post(r *http.Request) (io.Reader, error) {

	r.ParseForm()

	abs := h.absPath(r.URL.Path)
	file, err := os.Create(abs)
	if err != nil {
		// handle error
		return nil, err
	}

	fmt.Fprint(file, r.PostFormValue("post"))

	file.Close()

	data, _ := ioutil.ReadFile(abs)

	reader := bytes.NewReader(
		blackfriday.Run(data))

	return reader, nil
}

func (h *WikiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data io.Reader
	var err error

	log.Printf("%s %s\n", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet: // handle get request
		data, err = h.get(r.URL.Path)
		if err != nil {
			data, err = os.Open("form.html")
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, data)

	case http.MethodPost:
		data, err = h.post(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, data)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	handler := &WikiHandler{
		root: MarkDownPages,
	}

	http.ListenAndServe(":8080", handler)
}
