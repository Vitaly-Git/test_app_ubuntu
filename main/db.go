package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

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

func get_connections_history_by_hours() ([]time.Time, []float64) {

	var xvalues []time.Time
	var yvalues []float64

	db, err := sql.Open("sqlite3", "sqlite_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`select 
			strftime("%Y", connection_time) as year,
			strftime("%m", connection_time) as month,
			strftime("%d", connection_time) as day,
			count(*) as connections_count
		from 
			connections_history 
		group by 
			strftime("%Y", connection_time),
			strftime("%m", connection_time),
			strftime("%d", connection_time)
		order by 
			year, month, day asc`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {

		var year_str string
		var month_str string
		var day_str string
		var hour_str string
		var connections_count_str string
		err = rows.Scan(&year_str, &month_str, &day_str, &connections_count_str)
		if err != nil {
			panic(err)
		}

		year := parseInt(year_str)
		month := parseInt(month_str)
		day := parseInt(day_str)
		hour := parseInt(hour_str)
		connections_count := parseFloat64(connections_count_str)
		xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC))
		yvalues = append(yvalues, connections_count)
	}

	return xvalues, yvalues

}
