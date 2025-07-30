package db

import (
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/model"
)

type ScreenTimeDB struct {
	conn *DBConnection
}

func NewScreenTimeDB(conn *DBConnection) *ScreenTimeDB {
	return &ScreenTimeDB{conn: conn}
}

func (stdb *ScreenTimeDB) Insert(st model.ScreenTime) error {
	_, err := stdb.conn.db.Exec(
		"INSERT INTO screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)",
		st.Date, st.AppID, st.Title, st.Sleep,
	)
	return err
}

func (stdb *ScreenTimeDB) BulkInsert(records []model.ScreenTime) error {
	fn := "ScreenTimeDB:BulkInsert"
	tx, err := stdb.conn.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO screen_time(date, app_id, title, sleep) VALUES(?, ?, ?, ?)")
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			log.Println(fn, err)
		}
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(fn, err)
		}
	}()

	for _, st := range records {
		if _, err := stmt.Exec(st.Date, st.AppID, st.Title, st.Sleep); err != nil {
			e := tx.Rollback()
			if e != nil {
				log.Println(fn, err)
			}
			return err
		}
	}

	return tx.Commit()
}

func (stdb *ScreenTimeDB) GetByDateRange(from, to time.Time) ([]model.ScreenTime, error) {
	fn := "ScreenTimeDB:GetByDateRange"
	rows, err := stdb.conn.db.Query(
		"SELECT date, app_id, title, sleep FROM screen_time WHERE date BETWEEN ? AND ? ORDER BY date",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(fn, err)
		}
	}()

	var results []model.ScreenTime
	for rows.Next() {
		var st model.ScreenTime
		if err := rows.Scan(&st.Date, &st.AppID, &st.Title, &st.Sleep); err != nil {
			return nil, err
		}
		results = append(results, st)
	}

	return results, nil
}

func (stdb *ScreenTimeDB) GetAll() ([]model.ScreenTime, error) {
	fn := "ScreenTimeDB:GetAll"
	rows, err := stdb.conn.db.Query(
		"SELECT id, date, app_id, title, sleep FROM screen_time ORDER BY date",
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(fn, err)
		}
	}()

	var results []model.ScreenTime
	for rows.Next() {
		var st model.ScreenTime
		if err := rows.Scan(&st.ID, &st.Date, &st.AppID, &st.Title, &st.Sleep); err != nil {
			return nil, err
		}
		results = append(results, st)
	}

	return results, nil
}

func (stdb *ScreenTimeDB) GetAppUsage(from, to time.Time) (map[string]int, error) {
	fn := "ScreenTimeDB:GetAppUsage"
	rows, err := stdb.conn.db.Query(
		"SELECT app_id, SUM(sleep) FROM screen_time WHERE date BETWEEN ? AND ? GROUP BY app_id",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(fn, err)
		}
	}()

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

func (stdb *ScreenTimeDB) DeleteByID(screenTime model.ScreenTime) error {
	_, err := stdb.conn.db.Exec("DELETE FROM screen_time WHERE id = ?", screenTime.ID)
	return err
}

func (stdb *ScreenTimeDB) DeleteOldRecords(before time.Time) error {
	_, err := stdb.conn.db.Exec("DELETE FROM screen_time WHERE date < ?", before)
	return err
}
