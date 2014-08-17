// Copy your chrome History from "%appdata%\..\Local\Google\Chrome\User Data\Default"
// to this directory. Then run

package main

import ("fmt"
        "code.google.com/p/go-sqlite/go1/sqlite3"
        "strings"
        "net/http"
        "io/ioutil"
        )

func main() {
    // Select all youtube videos from history and order by visit_count
    sql := `SELECT * FROM urls
            WHERE url LIKE "https://www.youtube.com/watch?v=%"
            ORDER BY visit_count DESC
            LIMIT 100`

    // Open sqlite3 database stored in History file
    c, err := sqlite3.Open("Archived History")
    if err!=nil { panic("Cannot connect to History database") }
    row := make(sqlite3.RowMap)


    output := ""
    // Loop through the rows of the sql query
    for stmt, err := c.Query(sql); err == nil; err = stmt.Next() {
        
        var rowid int64
        stmt.Scan(&rowid, row)
        
        // For every video add a line with the visit_count and title. We put
        // an M in front if the video html contains the substring "Music"
        line := ""
        if containsMusic(row["url"].(string)) {
            line += "M "
        }
        line += fmt.Sprintln(row["visit_count"], row["title"])
        output += line

        // Print to show progress
        fmt.Println(line)
    }

    ioutil.WriteFile("MusicPreferences.txt", []byte(output), 0644)
}

// Open link with http get and looks for "Music" substring
func containsMusic(link string) bool {
    response, err := http.Get(link)
    if err != nil { panic("Http get error") }
    defer response.Body.Close()
    buff, err := ioutil.ReadAll(response.Body)
    if err != nil { panic("Cannot read http response body") }
    body := string(buff)
    return strings.Contains(body, "Muziek")
}

