package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "encoding/json"
    "os"
    b64 "encoding/base64"

     "database/sql"
   _ "github.com/lib/pq"
)

type Configuration struct {
    Db struct {
        Host     string
        Port     string
        User     string
        Password string
        Database string
    }
    Listen struct {
        Port    string
    }
}

func shortenUrl(original_url string) string {
    data := []byte(original_url)
    shortened_url := b64.StdEncoding.EncodeToString(data)
        
    if(len(shortened_url) >= len(original_url)) {
        shortened_url = original_url
    }
    return shortened_url
}

func findUrlInDB(url string, configuration* Configuration) bool {
    // connect to AWS RDS that had sample malformed URLs
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
      "password=%s dbname=%s sslmode=disable",
      configuration.Db.Host, configuration.Db.Port, 
      configuration.Db.User, configuration.Db.Password, 
      configuration.Db.Database)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(fmt.Sprintf("Error while opening connection to DB: %s", err))
    }

    err = db.Ping()
    if err != nil {
        panic(fmt.Sprintf("Error while pinging DB: %s", err))
    }

    sql := "select * from malformedurl where url like url"
    data, err := db.Query(sql)
    if err != nil {
        panic(fmt.Sprintf("Error query: %s", err))
    }
    for data.Next() {
        var urlRes string
        err = data.Scan(&urlRes)
        if err != nil {
            panic(fmt.Sprintf("Error scan: %s", err))
        }
        if url == urlRes {
            return true
        }
    }
    return false
}

// need to support GET requests:
// /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
func (configuration* Configuration) getHandler(w http.ResponseWriter, r *http.Request) {

    if !strings.HasPrefix(r.URL.Path, "/urlinfo/1/")  {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }

    // get the substring after prefix
    prefix := "/urlinfo/1/"
    pos := strings.LastIndex(r.URL.Path, prefix)
    if pos == -1 {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    adjustedPosAfterPrefix := pos + len(prefix)
    if adjustedPosAfterPrefix >= len(r.URL.Path) {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    urlWithoutPrefix := r.URL.Path[adjustedPosAfterPrefix:len(r.URL.Path)]

    // check that there is at least 1 '/' after prefix
    idx := strings.IndexByte(urlWithoutPrefix, '/')
    if idx < 0 || idx >= len(urlWithoutPrefix) {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    // get the substring from after prefix until first ':'
    idx2 := strings.IndexByte(urlWithoutPrefix, ':')
    if idx2 < 0 || idx ==(idx2 + 1) {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    url := urlWithoutPrefix[:idx2]

    // get shorter url if base64 result is shorter 
    shortened_url := shortenUrl(url)

    // DB lookup
    urlFound := findUrlInDB(shortened_url, configuration)
    if urlFound {
        fmt.Fprintf(w, "invalid: %s", url)
        return
    }

    fmt.Fprintf(w, "valid: %s", url)
}

func main() {

    file, _ := os.Open("config.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err := decoder.Decode(&configuration)
    if err != nil {
        fmt.Println("Config error:", err)
    }

    // Set up Http request handler
    http.HandleFunc("/", configuration.getHandler)

    fmt.Printf("Starting server at port %s\n", configuration.Listen.Port)
    port := fmt.Sprintf(":%s", configuration.Listen.Port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
