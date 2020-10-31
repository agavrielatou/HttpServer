package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
)

// need to support GET requests:
// /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
func getHandler(w http.ResponseWriter, r *http.Request) {

    if !strings.HasPrefix(r.URL.Path, "/urlinfo/1/")  {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }

    url := strings.Split(r.URL.Path, "/urlinfo/1/")
    if !isValid(url) {
      fmt.Fprintf(w, "Not valid URL")
      return
    }

    fmt.Fprintf(w, "Valid URL")
}

func main() {
    http.HandleFunc("/", getHandler)

    fmt.Printf("Starting server at port 8000\n")
    if err := http.ListenAndServe(":8000", nil); err != nil {
        log.Fatal(err)
    }
}
