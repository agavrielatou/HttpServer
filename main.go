package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

     "database/sql"
   _ "github.com/lib/pq"
)

const (
//  host     = "malformedurls.cydwrqsworn7.us-east-2.rds.amazonaws.com"
  host     = "localhost"
  port     = 5432
  user     = "docker"
  password = "docker"
  dbname   = "postgres"
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
//    if !isValid(url) {
//      fmt.Fprintf(w, "Not valid URL")
//      return
//    }

    fmt.Fprintf(w, "Valid URL: %s", url)
}

func main() {
    // connect to AWS RDS that had sample malformed URLs

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
      "password=%s dbname=%s sslmode=disable",
      host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      fmt.Printf("Error while opening connection to DB: %s", err)
      panic(err)
    }

    err = db.Ping()
    if err != nil {
      fmt.Printf("Error while pinging DB: %s", err)
      panic(err)
    }
    fmt.Printf("Successfully connected!")

    sql := "select * from malformedurl"
    data, err := db.Query(sql)
    if err != nil {
        fmt.Printf("Error query: %s", err)
    }
    fmt.Printf("Successful query: %s", data)

    // Set up Http request handler
    http.HandleFunc("/", getHandler)

    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
