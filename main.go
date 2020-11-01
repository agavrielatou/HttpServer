package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "encoding/json"
    "os"

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

    prefix := "/urlinfo/1/"
    pos := strings.LastIndex(r.URL.Path, prefix)
    if pos == -1 {
        return
    }
    adjustedPos := pos + len(prefix)
    if adjustedPos >= len(r.URL.Path) {
        return
    }
    url := r.URL.Path[adjustedPos:len(r.URL.Path)]

    urlFound := findUrlInDB(url, configuration)
    if urlFound {
        fmt.Fprintf(w, "Invalid URL: %s", url)
        return
    }

    fmt.Fprintf(w, "Valid myURL: %s", url)
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
