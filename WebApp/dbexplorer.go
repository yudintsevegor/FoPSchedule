package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func dbExplorer(db *sql.DB, group string) ([][]Subject, error) {
	allWeek := make([][]Subject, 0, 6)
	req := fmt.Sprintf("SELECT first, second, third, fourth, fifth FROM `%v`", group)
	rows, err := db.Query(req)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rawLes := make([]string, 5)
		if err = rows.Scan(&rawLes[0], &rawLes[1], &rawLes[2], &rawLes[3], &rawLes[4]); err != nil {
			return nil, err
		}

		// les := parsePercent(rawLes)
		allWeek = append(allWeek, parsePercent(rawLes))
	}

	return allWeek, nil
}

func parsePercent(arr []string) []Subject {
	// sbj := Subject{}
	result := make([]Subject, 0, 5)
	for _, val := range arr {
		res := strings.Split(val, "%")
		/*
			res := rePerc.FindStringSubmatch(val)
			sbj.Name = res[1]
			sbj.Lector = res[2]
			sbj.Room = res[3]
		*/

		result = append(result, Subject{
			Name:   res[0],
			Lector: res[1],
			Room:   res[2],
		})
	}

	return result
}
