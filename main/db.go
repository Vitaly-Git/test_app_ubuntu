package main

import (
	"database/sql"
	"net/http"

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
