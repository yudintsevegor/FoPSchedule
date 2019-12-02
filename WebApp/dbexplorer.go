package main

import (
	"database/sql"
	"fmt"
	"fopSchedule/master/common"
)

func dbExplorer(db *sql.DB, group string) ([][]common.Subject, error) {
	allWeek := make([][]common.Subject, 0, 6)
	req := fmt.Sprintf("SELECT first, second, third, fourth, fifth FROM `%v`", group)
	rows, err := db.Query(req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rawLes := make([]string, 5)
		if err = rows.Scan(&rawLes[0], &rawLes[1], &rawLes[2], &rawLes[3], &rawLes[4]); err != nil {
			return nil, err
		}

		allWeek = append(allWeek, parsePercent(rawLes))
	}

	return allWeek, nil
}
