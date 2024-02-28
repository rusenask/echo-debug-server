package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
)

type response struct {
	ServerInfo struct {
		Port string
	}
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	RawQuery   string              `json:"raw_query"`
	Header     http.Header         `json:"header"`
	Body       string              `json:"body"`
	FormValues map[string][]string `json:"form_values"` // used if request is a form
}

func main() {
	port := getPort()

	http.HandleFunc("/multipart-form", func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header["Content-Type"][0])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if strings.HasPrefix(mediaType, "multipart/") {
			err = r.ParseMultipartForm(25)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			data := response{
				Method:     r.Method,
				Path:       r.URL.Path,
				RawQuery:   r.URL.RawQuery,
				Header:     r.Header,
				FormValues: make(map[string][]string),
			}
			data.ServerInfo.Port = port
			for k, v := range r.MultipartForm.Value {
				data.FormValues[k] = v
			}
			resp, err := json.Marshal(&data)
			if err != nil {
				log.Println("[ERROR] failed to marshal resp body, error: ", err)
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			fmt.Println(string(resp))
			w.WriteHeader(200)
			w.Write(resp)
			return
		}
		http.Error(w, "not a multipart form (header missing)", 400)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var bstr string

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("[ERROR] failed to read body, error: ", err)
		} else {
			bstr = string(body)
		}
		respData := &response{
			Method:   r.Method,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
			Header:   r.Header,
			Body:     bstr,
		}
		respData.ServerInfo.Port = port
		resp, err := json.Marshal(respData)
		if err != nil {
			log.Println("[ERROR] failed to marshal resp body, error: ", err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		if respData.RawQuery != "" {
			log.Printf("[INFO] %s %s?%s \n%s \n", respData.Method, respData.Path, respData.RawQuery, respData.Body)
		} else {
			log.Printf("[INFO] %s %s \n%s \n", respData.Method, respData.Path, respData.RawQuery, respData.Body)
		}

		w.WriteHeader(200)
		w.Write(resp)
	})

	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() string {
	if os.Getenv("PORT") != "" {
		return ":" + os.Getenv("PORT")
	}
	return ":8888"
}
