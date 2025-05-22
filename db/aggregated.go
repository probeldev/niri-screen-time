package db

import (
	"time"

	"github.com/probeldev/niri-screen-time/model"
)

type AggregatedScreenTimeDB struct {
	conn *DBConnection
}

func NewAggregatedScreenTimeDB(conn *DBConnection) *AggregatedScreenTimeDB {
	return &AggregatedScreenTimeDB{conn: conn}
}

func (astdb *AggregatedScreenTimeDB) Insert(ast model.AggregatedScreenTime) error {
	_, err := astdb.conn.db.Exec(
		"INSERT INTO aggregated_screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)",
		ast.Date, ast.AppID, ast.Title, ast.Sleep,
	)
	return err
}

func (astdb *AggregatedScreenTimeDB) BulkInsert(records []model.AggregatedScreenTime) error {
	tx, err := astdb.conn.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO aggregated_screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, st := range records {
		if _, err := stmt.Exec(st.Date, st.AppID, st.Title, st.Sleep); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (astdb *AggregatedScreenTimeDB) GetByDateRange(
	from,
	to time.Time,
) (
	[]model.AggregatedScreenTime,
	error,
) {
	rows, err := astdb.conn.db.Query(
		"SELECT date, app_id, title, sleep FROM aggregated_screen_time WHERE date BETWEEN ? AND ? ORDER BY date",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.AggregatedScreenTime
	for rows.Next() {
		var st model.AggregatedScreenTime
		if err := rows.Scan(&st.Date, &st.AppID, &st.Title, &st.Sleep); err != nil {
			return nil, err
		}
		results = append(results, st)
	}

	return results, nil
}

func (astdb *AggregatedScreenTimeDB) GetAppUsage(from, to time.Time) (map[string]int, error) {
	rows, err := astdb.conn.db.Query(
		"SELECT app_id, SUM(sleep) FROM aggregated_screen_time WHERE date BETWEEN ? AND ? GROUP BY app_id",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var appID string
		var total int
		if err := rows.Scan(&appID, &total); err != nil {
			return nil, err
		}
		result[appID] = total
	}

	return result, nil
}

func (astdb *AggregatedScreenTimeDB) DeleteOldRecords(before time.Time) error {
	_, err := astdb.conn.db.Exec("DELETE FROM aggregated_screen_time WHERE date < ?", before)
	return err
}
