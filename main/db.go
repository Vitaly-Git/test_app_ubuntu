package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func add_connect_params_to_db(r *http.Request, answer string) {

	var client string = r.RemoteAddr
	var url_path string = r.URL.Path
	var params string = r.URL.RawQuery

	add_connect_to_db(client, url_path, params, answer)
}

func add_connect_to_db(client string, url string, params string, answer string) {
	db, err := sql.Open("sqlite3", "sqlite_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into connections_history (client, url, params, answer) values ('" + client + "', '" + url + "', '" + params + "', '" + answer + "')")
	if err != nil {
		panic(err)
	}
	// fmt.Println(result.LastInsertId()) // id последнего добавленного объекта
	// fmt.Println(result.RowsAffected()) // количество добавленных строк
}

func get_connections_history() string {

	column_width := "%22.22s"

	var result_answer string = "" //"client \t| url \t| params \t| answer \t| connection_time\n"

	result_answer += "" +
		" | " + fmt.Sprintf(column_width, "connection_time") +
		" | " + fmt.Sprintf(column_width, "client") +
		" | " + fmt.Sprintf(column_width, "url") +
		" | " + fmt.Sprintf(column_width, "params") +
		" | " + fmt.Sprintf(column_width, "answer") +
		" | " + "\n"

	db, err := sql.Open("sqlite3", "sqlite_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select client, url, params, answer, connection_time from connections_history order by connection_time desc limit 30")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {

		var client string
		var url string
		var params string
		var answer string
		var connection_time string

		err = rows.Scan(&client, &url, &params, &answer, &connection_time)
		if err != nil {
			panic(err)
		}

		//result_answer += "" + client + " \t| " + url + " \t| " + params + " \t| " + strings.Split(answer, "\n")[0][:10] + " \t| " + connection_time + "\n"

		result_answer += "" +
			" | " + fmt.Sprintf(column_width, connection_time) +
			" | " + fmt.Sprintf(column_width, client) +
			" | " + fmt.Sprintf(column_width, url) +
			" | " + fmt.Sprintf(column_width, params) +
			" | " + fmt.Sprintf(column_width, strings.Replace(answer, "\t", " ", -1)) +
			" | " + "\n"
	}

	return result_answer
}
