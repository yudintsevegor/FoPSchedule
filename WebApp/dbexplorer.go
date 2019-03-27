package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func dbExplorer(db *sql.DB, group string) [][]Subject {
	var tablesNames = make([]string, 0, 1)
	var tableName string

	//For debugging
	rowsTb, err := db.Query("SHOW TABLES")
	for rowsTb.Next() {
		err = rowsTb.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		tablesNames = append(tablesNames, tableName)
	}
	rowsTb.Close()

	var allWeek = make([][]Subject, 0, 6)
	req := fmt.Sprintf("SELECT first, second, third, fourth, fifth FROM `%v`", group)
	rows, err := db.Query(req)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var rawLes = make([]string, 5)
		err = rows.Scan(&rawLes[0], &rawLes[1], &rawLes[2], &rawLes[3], &rawLes[4])
		if err != nil {
			log.Fatal(err)
		}
		les := parsePercent(rawLes)
		allWeek = append(allWeek, les)
	}

	return allWeek
}

func parsePercent(arr []string) []Subject {
	sbj := Subject{}
	var result = make([]Subject, 0, 5)
	for _, val := range arr {
		res := rePerc.FindStringSubmatch(val)
		sbj.Name = res[1]
		sbj.Lector = res[2]
		sbj.Room = res[3]
		result = append(result, sbj)
	}

	return result
}
