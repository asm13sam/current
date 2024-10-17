package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func DBconnect(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	return err
}
func DbClose() error {
	return db.Close()
}

type Document struct {
	Id       int    `json:"id"`
	DocType  string `json:"doc_type"`
	IsActive bool   `json:"is_active"`
}

func DocumentGet(id int, tx *sql.Tx) (Document, error) {
	var d Document
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM document WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM document WHERE id=?", id)
	}

	err := row.Scan(
		&d.Id,
		&d.DocType,
		&d.IsActive,
	)
	return d, err
}

func DocumentGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Document, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM document"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

func DocumentCreate(d Document, tx *sql.Tx) (Document, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return d, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO document
            (doc_type, is_active)
            VALUES(?, ?);`
	res, err := tx.Exec(
		sql,
		d.DocType,
		d.IsActive,
	)
	if err != nil {
		return d, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return d, err
	}
	d.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func DocumentUpdate(d Document, tx *sql.Tx) (Document, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return d, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE document SET
                    doc_type=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		d.DocType,
		d.IsActive,
		d.Id,
	)
	if err != nil {
		return d, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func DocumentDelete(id int, tx *sql.Tx) (Document, error) {
	needCommit := false
	var err error
	var d Document
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return d, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	d, err = DocumentGet(id, tx)
	if err != nil {
		return d, err
	}

	sql := `UPDATE document SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, d.Id)
	if err != nil {
		return d, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return d, err
		}
	}
	d.IsActive = false
	return d, nil
}

func DocumentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Document, error) {

	if !DocumentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM document WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil

}

func DocumentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Document, error) {

	if !DocumentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM document WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil

}

func DocumentTestForExistingField(fieldName string) bool {
	fields := []string{"id", "doc_type", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Measure struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

func MeasureGet(id int, tx *sql.Tx) (Measure, error) {
	var m Measure
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM measure WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM measure WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.FullName,
		&m.IsActive,
	)
	return m, err
}

func MeasureGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Measure, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM measure"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Measure{}
	for rows.Next() {
		var m Measure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MeasureCreate(m Measure, tx *sql.Tx) (Measure, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO measure
            (name, full_name, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.Name,
		m.FullName,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MeasureUpdate(m Measure, tx *sql.Tx) (Measure, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE measure SET
                    name=?, full_name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.Name,
		m.FullName,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MeasureDelete(id int, tx *sql.Tx) (Measure, error) {
	needCommit := false
	var err error
	var m Measure
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MeasureGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE measure SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MeasureGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Measure, error) {

	if !MeasureTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM measure WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Measure{}
	for rows.Next() {
		var m Measure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MeasureGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Measure, error) {

	if !MeasureTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM measure WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Measure{}
	for rows.Next() {
		var m Measure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MeasureTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "full_name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type CountType struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func CountTypeGet(id int, tx *sql.Tx) (CountType, error) {
	var c CountType
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM count_type WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM count_type WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.IsActive,
	)
	return c, err
}

func CountTypeGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CountType, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM count_type"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CountType{}
	for rows.Next() {
		var c CountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CountTypeCreate(c CountType, tx *sql.Tx) (CountType, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO count_type
            (name, is_active)
            VALUES(?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CountTypeUpdate(c CountType, tx *sql.Tx) (CountType, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE count_type SET
                    name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CountTypeDelete(id int, tx *sql.Tx) (CountType, error) {
	needCommit := false
	var err error
	var c CountType
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CountTypeGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE count_type SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CountTypeGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CountType, error) {

	if !CountTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM count_type WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CountType{}
	for rows.Next() {
		var c CountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CountTypeGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CountType, error) {

	if !CountTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM count_type WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CountType{}
	for rows.Next() {
		var c CountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CountTypeTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ColorGroup struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ColorGroupId int    `json:"color_group_id"`
	IsActive     bool   `json:"is_active"`
}

func ColorGroupGet(id int, tx *sql.Tx) (ColorGroup, error) {
	var c ColorGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM color_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM color_group WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ColorGroupId,
		&c.IsActive,
	)
	return c, err
}

func ColorGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ColorGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM color_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ColorGroup{}
	for rows.Next() {
		var c ColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func ColorGroupCreate(c ColorGroup, tx *sql.Tx) (ColorGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO color_group
            (name, color_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.ColorGroupId,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ColorGroupUpdate(c ColorGroup, tx *sql.Tx) (ColorGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE color_group SET
                    name=?, color_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.ColorGroupId,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ColorGroupDelete(id int, tx *sql.Tx) (ColorGroup, error) {
	needCommit := false
	var err error
	var c ColorGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = ColorGroupGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE color_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func ColorGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ColorGroup, error) {

	if !ColorGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM color_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ColorGroup{}
	for rows.Next() {
		var c ColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ColorGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ColorGroup, error) {

	if !ColorGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM color_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ColorGroup{}
	for rows.Next() {
		var c ColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ColorGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "color_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Color struct {
	Id           int     `json:"id"`
	ColorGroupId int     `json:"color_group_id"`
	Name         string  `json:"name"`
	Total        float64 `json:"total"`
	IsActive     bool    `json:"is_active"`
}

func ColorGet(id int, tx *sql.Tx) (Color, error) {
	var c Color
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM color WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM color WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.ColorGroupId,
		&c.Name,
		&c.Total,
		&c.IsActive,
	)
	return c, err
}

func ColorGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Color, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM color"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Color{}
	for rows.Next() {
		var c Color
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func ColorCreate(c Color, tx *sql.Tx) (Color, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO color
            (color_group_id, name, total, is_active)
            VALUES(?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.ColorGroupId,
		c.Name,
		c.Total,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ColorUpdate(c Color, tx *sql.Tx) (Color, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE color SET
                    color_group_id=?, name=?, total=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.ColorGroupId,
		c.Name,
		c.Total,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ColorDelete(id int, tx *sql.Tx) (Color, error) {
	needCommit := false
	var err error
	var c Color
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = ColorGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE color SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func ColorGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Color, error) {

	if !ColorTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM color WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Color{}
	for rows.Next() {
		var c Color
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ColorGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Color, error) {

	if !ColorTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM color WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Color{}
	for rows.Next() {
		var c Color
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ColorTestForExistingField(fieldName string) bool {
	fields := []string{"id", "color_group_id", "name", "total", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type MatherialGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	MatherialGroupId int    `json:"matherial_group_id"`
	IsActive         bool   `json:"is_active"`
}

func MatherialGroupGet(id int, tx *sql.Tx) (MatherialGroup, error) {
	var m MatherialGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_group WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.MatherialGroupId,
		&m.IsActive,
	)
	return m, err
}

func MatherialGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialGroup{}
	for rows.Next() {
		var m MatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialGroupCreate(m MatherialGroup, tx *sql.Tx) (MatherialGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO matherial_group
            (name, matherial_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.Name,
		m.MatherialGroupId,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialGroupUpdate(m MatherialGroup, tx *sql.Tx) (MatherialGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial_group SET
                    name=?, matherial_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.Name,
		m.MatherialGroupId,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialGroupDelete(id int, tx *sql.Tx) (MatherialGroup, error) {
	needCommit := false
	var err error
	var m MatherialGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialGroupGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialGroup, error) {

	if !MatherialGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialGroup{}
	for rows.Next() {
		var m MatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialGroup, error) {

	if !MatherialGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialGroup{}
	for rows.Next() {
		var m MatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "matherial_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Matherial struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	MatherialGroupId int     `json:"matherial_group_id"`
	MeasureId        int     `json:"measure_id"`
	ColorGroupId     int     `json:"color_group_id"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	Total            float64 `json:"total"`
	Barcode          string  `json:"barcode"`
	CountTypeId      int     `json:"count_type_id"`
	IsActive         bool    `json:"is_active"`
}

func MatherialGet(id int, tx *sql.Tx) (Matherial, error) {
	var m Matherial
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.FullName,
		&m.MatherialGroupId,
		&m.MeasureId,
		&m.ColorGroupId,
		&m.Price,
		&m.Cost,
		&m.Total,
		&m.Barcode,
		&m.CountTypeId,
		&m.IsActive,
	)
	return m, err
}

func MatherialGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Matherial, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Matherial{}
	for rows.Next() {
		var m Matherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialCreate(m Matherial, tx *sql.Tx) (Matherial, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO matherial
            (name, full_name, matherial_group_id, measure_id, color_group_id, price, cost, total, barcode, count_type_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.Name,
		m.FullName,
		m.MatherialGroupId,
		m.MeasureId,
		m.ColorGroupId,
		m.Price,
		m.Cost,
		m.Total,
		m.Barcode,
		m.CountTypeId,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialUpdate(m Matherial, tx *sql.Tx) (Matherial, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial SET
                    name=?, full_name=?, matherial_group_id=?, measure_id=?, color_group_id=?, price=?, cost=?, total=?, barcode=?, count_type_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.Name,
		m.FullName,
		m.MatherialGroupId,
		m.MeasureId,
		m.ColorGroupId,
		m.Price,
		m.Cost,
		m.Total,
		m.Barcode,
		m.CountTypeId,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialDelete(id int, tx *sql.Tx) (Matherial, error) {
	needCommit := false
	var err error
	var m Matherial
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Matherial, error) {

	if !MatherialTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Matherial{}
	for rows.Next() {
		var m Matherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Matherial, error) {

	if !MatherialTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Matherial{}
	for rows.Next() {
		var m Matherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "full_name", "matherial_group_id", "measure_id", "color_group_id", "price", "cost", "total", "barcode", "count_type_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Cash struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Persent   float64 `json:"persent"`
	Total     float64 `json:"total"`
	Comm      string  `json:"comm"`
	IsFiscal  bool    `json:"is_fiscal"`
	IsAccount bool    `json:"is_account"`
	IsActive  bool    `json:"is_active"`
}

func CashGet(id int, tx *sql.Tx) (Cash, error) {
	var c Cash
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM cash WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM cash WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Persent,
		&c.Total,
		&c.Comm,
		&c.IsFiscal,
		&c.IsAccount,
		&c.IsActive,
	)
	return c, err
}

func CashGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Cash, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM cash"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Cash{}
	for rows.Next() {
		var c Cash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CashCreate(c Cash, tx *sql.Tx) (Cash, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO cash
            (name, persent, total, comm, is_fiscal, is_account, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.Persent,
		c.Total,
		c.Comm,
		c.IsFiscal,
		c.IsAccount,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashUpdate(c Cash, tx *sql.Tx) (Cash, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE cash SET
                    name=?, persent=?, total=?, comm=?, is_fiscal=?, is_account=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.Persent,
		c.Total,
		c.Comm,
		c.IsFiscal,
		c.IsAccount,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashDelete(id int, tx *sql.Tx) (Cash, error) {
	needCommit := false
	var err error
	var c Cash
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CashGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE cash SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CashGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Cash, error) {

	if !CashTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Cash{}
	for rows.Next() {
		var c Cash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Cash, error) {

	if !CashTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Cash{}
	for rows.Next() {
		var c Cash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "persent", "total", "comm", "is_fiscal", "is_account", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type UserGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UserGroupId int    `json:"user_group_id"`
	IsActive    bool   `json:"is_active"`
}

func UserGroupGet(id int, tx *sql.Tx) (UserGroup, error) {
	var u UserGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM user_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM user_group WHERE id=?", id)
	}

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.UserGroupId,
		&u.IsActive,
	)
	return u, err
}

func UserGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]UserGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM user_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []UserGroup{}
	for rows.Next() {
		var u UserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func UserGroupCreate(u UserGroup, tx *sql.Tx) (UserGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO user_group
            (name, user_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		u.Name,
		u.UserGroupId,
		u.IsActive,
	)
	if err != nil {
		return u, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	return u, nil
}

func UserGroupUpdate(u UserGroup, tx *sql.Tx) (UserGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE user_group SET
                    name=?, user_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		u.Name,
		u.UserGroupId,
		u.IsActive,
		u.Id,
	)
	if err != nil {
		return u, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	return u, nil
}

func UserGroupDelete(id int, tx *sql.Tx) (UserGroup, error) {
	needCommit := false
	var err error
	var u UserGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	u, err = UserGroupGet(id, tx)
	if err != nil {
		return u, err
	}

	sql := `UPDATE user_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, u.Id)
	if err != nil {
		return u, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	u.IsActive = false
	return u, nil
}

func UserGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]UserGroup, error) {

	if !UserGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM user_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []UserGroup{}
	for rows.Next() {
		var u UserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func UserGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]UserGroup, error) {

	if !UserGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM user_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []UserGroup{}
	for rows.Next() {
		var u UserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func UserGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "user_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	UserGroupId int    `json:"user_group_id"`
	CashId      int    `json:"cash_id"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Comm        string `json:"comm"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	BaseAccess  int    `json:"base_access"`
	AddAccess   int    `json:"add_access"`
	IsActive    bool   `json:"is_active"`
}

func UserGet(id int, tx *sql.Tx) (User, error) {
	var u User
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM user WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM user WHERE id=?", id)
	}

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.FullName,
		&u.UserGroupId,
		&u.CashId,
		&u.Phone,
		&u.Email,
		&u.Comm,
		&u.Login,
		&u.Password,
		&u.BaseAccess,
		&u.AddAccess,
		&u.IsActive,
	)
	return u, err
}

func UserGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]User, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM user"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func UserCreate(u User, tx *sql.Tx) (User, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO user
            (name, full_name, user_group_id, cash_id, phone, email, comm, login, password, base_access, add_access, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		u.Name,
		u.FullName,
		u.UserGroupId,
		u.CashId,
		u.Phone,
		u.Email,
		u.Comm,
		u.Login,
		u.Password,
		u.BaseAccess,
		u.AddAccess,
		u.IsActive,
	)
	if err != nil {
		return u, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	return u, nil
}

func UserUpdate(u User, tx *sql.Tx) (User, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE user SET
                    name=?, full_name=?, user_group_id=?, cash_id=?, phone=?, email=?, comm=?, login=?, password=?, base_access=?, add_access=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		u.Name,
		u.FullName,
		u.UserGroupId,
		u.CashId,
		u.Phone,
		u.Email,
		u.Comm,
		u.Login,
		u.Password,
		u.BaseAccess,
		u.AddAccess,
		u.IsActive,
		u.Id,
	)
	if err != nil {
		return u, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	return u, nil
}

func UserDelete(id int, tx *sql.Tx) (User, error) {
	needCommit := false
	var err error
	var u User
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return u, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	u, err = UserGet(id, tx)
	if err != nil {
		return u, err
	}

	sql := `UPDATE user SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, u.Id)
	if err != nil {
		return u, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return u, err
		}
	}
	u.IsActive = false
	return u, nil
}

func UserGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]User, error) {

	if !UserTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM user WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func UserGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]User, error) {

	if !UserTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM user WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func UserTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "full_name", "user_group_id", "cash_id", "phone", "email", "comm", "login", "password", "base_access", "add_access", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type EquipmentGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	EquipmentGroupId int    `json:"equipment_group_id"`
	IsActive         bool   `json:"is_active"`
}

func EquipmentGroupGet(id int, tx *sql.Tx) (EquipmentGroup, error) {
	var e EquipmentGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM equipment_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM equipment_group WHERE id=?", id)
	}

	err := row.Scan(
		&e.Id,
		&e.Name,
		&e.EquipmentGroupId,
		&e.IsActive,
	)
	return e, err
}

func EquipmentGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]EquipmentGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM equipment_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []EquipmentGroup{}
	for rows.Next() {
		var e EquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func EquipmentGroupCreate(e EquipmentGroup, tx *sql.Tx) (EquipmentGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO equipment_group
            (name, equipment_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		e.Name,
		e.EquipmentGroupId,
		e.IsActive,
	)
	if err != nil {
		return e, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return e, err
	}
	e.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

func EquipmentGroupUpdate(e EquipmentGroup, tx *sql.Tx) (EquipmentGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE equipment_group SET
                    name=?, equipment_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		e.Name,
		e.EquipmentGroupId,
		e.IsActive,
		e.Id,
	)
	if err != nil {
		return e, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

func EquipmentGroupDelete(id int, tx *sql.Tx) (EquipmentGroup, error) {
	needCommit := false
	var err error
	var e EquipmentGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	e, err = EquipmentGroupGet(id, tx)
	if err != nil {
		return e, err
	}

	sql := `UPDATE equipment_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, e.Id)
	if err != nil {
		return e, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	e.IsActive = false
	return e, nil
}

func EquipmentGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]EquipmentGroup, error) {

	if !EquipmentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM equipment_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []EquipmentGroup{}
	for rows.Next() {
		var e EquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func EquipmentGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]EquipmentGroup, error) {

	if !EquipmentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM equipment_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []EquipmentGroup{}
	for rows.Next() {
		var e EquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func EquipmentGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "equipment_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Equipment struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	EquipmentGroupId int     `json:"equipment_group_id"`
	Cost             float64 `json:"cost"`
	Total            float64 `json:"total"`
	IsActive         bool    `json:"is_active"`
}

func EquipmentGet(id int, tx *sql.Tx) (Equipment, error) {
	var e Equipment
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM equipment WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM equipment WHERE id=?", id)
	}

	err := row.Scan(
		&e.Id,
		&e.Name,
		&e.FullName,
		&e.EquipmentGroupId,
		&e.Cost,
		&e.Total,
		&e.IsActive,
	)
	return e, err
}

func EquipmentGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Equipment, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM equipment"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Equipment{}
	for rows.Next() {
		var e Equipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func EquipmentCreate(e Equipment, tx *sql.Tx) (Equipment, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO equipment
            (name, full_name, equipment_group_id, cost, total, is_active)
            VALUES(?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		e.Name,
		e.FullName,
		e.EquipmentGroupId,
		e.Cost,
		e.Total,
		e.IsActive,
	)
	if err != nil {
		return e, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return e, err
	}
	e.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

func EquipmentUpdate(e Equipment, tx *sql.Tx) (Equipment, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE equipment SET
                    name=?, full_name=?, equipment_group_id=?, cost=?, total=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		e.Name,
		e.FullName,
		e.EquipmentGroupId,
		e.Cost,
		e.Total,
		e.IsActive,
		e.Id,
	)
	if err != nil {
		return e, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

func EquipmentDelete(id int, tx *sql.Tx) (Equipment, error) {
	needCommit := false
	var err error
	var e Equipment
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return e, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	e, err = EquipmentGet(id, tx)
	if err != nil {
		return e, err
	}

	sql := `UPDATE equipment SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, e.Id)
	if err != nil {
		return e, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return e, err
		}
	}
	e.IsActive = false
	return e, nil
}

func EquipmentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Equipment, error) {

	if !EquipmentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM equipment WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Equipment{}
	for rows.Next() {
		var e Equipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func EquipmentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Equipment, error) {

	if !EquipmentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM equipment WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Equipment{}
	for rows.Next() {
		var e Equipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func EquipmentTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "full_name", "equipment_group_id", "cost", "total", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type OperationGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	OperationGroupId int    `json:"operation_group_id"`
	IsActive         bool   `json:"is_active"`
}

func OperationGroupGet(id int, tx *sql.Tx) (OperationGroup, error) {
	var o OperationGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM operation_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM operation_group WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.OperationGroupId,
		&o.IsActive,
	)
	return o, err
}

func OperationGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM operation_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationGroup{}
	for rows.Next() {
		var o OperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OperationGroupCreate(o OperationGroup, tx *sql.Tx) (OperationGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO operation_group
            (name, operation_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.Name,
		o.OperationGroupId,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationGroupUpdate(o OperationGroup, tx *sql.Tx) (OperationGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE operation_group SET
                    name=?, operation_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.Name,
		o.OperationGroupId,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationGroupDelete(id int, tx *sql.Tx) (OperationGroup, error) {
	needCommit := false
	var err error
	var o OperationGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OperationGroupGet(id, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE operation_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OperationGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationGroup, error) {

	if !OperationGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationGroup{}
	for rows.Next() {
		var o OperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationGroup, error) {

	if !OperationGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationGroup{}
	for rows.Next() {
		var o OperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "operation_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Operation struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	OperationGroupId int     `json:"operation_group_id"`
	MeasureId        int     `json:"measure_id"`
	UserId           int     `json:"user_id"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	EquipmentId      int     `json:"equipment_id"`
	EquipmentPrice   float64 `json:"equipment_price"`
	Barcode          string  `json:"barcode"`
	IsActive         bool    `json:"is_active"`
}

func OperationGet(id int, tx *sql.Tx) (Operation, error) {
	var o Operation
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM operation WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM operation WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.FullName,
		&o.OperationGroupId,
		&o.MeasureId,
		&o.UserId,
		&o.Price,
		&o.Cost,
		&o.EquipmentId,
		&o.EquipmentPrice,
		&o.Barcode,
		&o.IsActive,
	)
	return o, err
}

func OperationGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Operation, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM operation"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Operation{}
	for rows.Next() {
		var o Operation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OperationCreate(o Operation, tx *sql.Tx) (Operation, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO operation
            (name, full_name, operation_group_id, measure_id, user_id, price, cost, equipment_id, equipment_price, barcode, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.Name,
		o.FullName,
		o.OperationGroupId,
		o.MeasureId,
		o.UserId,
		o.Price,
		o.Cost,
		o.EquipmentId,
		o.EquipmentPrice,
		o.Barcode,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationUpdate(o Operation, tx *sql.Tx) (Operation, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE operation SET
                    name=?, full_name=?, operation_group_id=?, measure_id=?, user_id=?, price=?, cost=?, equipment_id=?, equipment_price=?, barcode=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.Name,
		o.FullName,
		o.OperationGroupId,
		o.MeasureId,
		o.UserId,
		o.Price,
		o.Cost,
		o.EquipmentId,
		o.EquipmentPrice,
		o.Barcode,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationDelete(id int, tx *sql.Tx) (Operation, error) {
	needCommit := false
	var err error
	var o Operation
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OperationGet(id, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE operation SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OperationGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Operation, error) {

	if !OperationTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Operation{}
	for rows.Next() {
		var o Operation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Operation, error) {

	if !OperationTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Operation{}
	for rows.Next() {
		var o Operation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "full_name", "operation_group_id", "measure_id", "user_id", "price", "cost", "equipment_id", "equipment_price", "barcode", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ProductGroup struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ProductGroupId int    `json:"product_group_id"`
	IsActive       bool   `json:"is_active"`
}

func ProductGroupGet(id int, tx *sql.Tx) (ProductGroup, error) {
	var p ProductGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM product_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM product_group WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ProductGroupId,
		&p.IsActive,
	)
	return p, err
}

func ProductGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM product_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductGroup{}
	for rows.Next() {
		var p ProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProductGroupCreate(p ProductGroup, tx *sql.Tx) (ProductGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO product_group
            (name, product_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.ProductGroupId,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductGroupUpdate(p ProductGroup, tx *sql.Tx) (ProductGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE product_group SET
                    name=?, product_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.ProductGroupId,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductGroupDelete(id int, tx *sql.Tx) (ProductGroup, error) {
	needCommit := false
	var err error
	var p ProductGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProductGroupGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE product_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProductGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductGroup, error) {

	if !ProductGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductGroup{}
	for rows.Next() {
		var p ProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductGroup, error) {

	if !ProductGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductGroup{}
	for rows.Next() {
		var p ProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "product_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Product struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	ShortName      string  `json:"short_name"`
	ProductGroupId int     `json:"product_group_id"`
	MeasureId      int     `json:"measure_id"`
	Width          float64 `json:"width"`
	Length         float64 `json:"length"`
	MinCost        float64 `json:"min_cost"`
	Cost           float64 `json:"cost"`
	UserId         int     `json:"user_id"`
	Barcode        string  `json:"barcode"`
	IsActive       bool    `json:"is_active"`
}

func ProductGet(id int, tx *sql.Tx) (Product, error) {
	var p Product
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM product WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM product WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ShortName,
		&p.ProductGroupId,
		&p.MeasureId,
		&p.Width,
		&p.Length,
		&p.MinCost,
		&p.Cost,
		&p.UserId,
		&p.Barcode,
		&p.IsActive,
	)
	return p, err
}

func ProductGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Product, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM product"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProductCreate(p Product, tx *sql.Tx) (Product, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO product
            (name, short_name, product_group_id, measure_id, width, length, min_cost, cost, user_id, barcode, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.ShortName,
		p.ProductGroupId,
		p.MeasureId,
		p.Width,
		p.Length,
		p.MinCost,
		p.Cost,
		p.UserId,
		p.Barcode,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductUpdate(p Product, tx *sql.Tx) (Product, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE product SET
                    name=?, short_name=?, product_group_id=?, measure_id=?, width=?, length=?, min_cost=?, cost=?, user_id=?, barcode=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.ShortName,
		p.ProductGroupId,
		p.MeasureId,
		p.Width,
		p.Length,
		p.MinCost,
		p.Cost,
		p.UserId,
		p.Barcode,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductDelete(id int, tx *sql.Tx) (Product, error) {
	needCommit := false
	var err error
	var p Product
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProductGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE product SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Product, error) {

	if !ProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Product, error) {

	if !ProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "short_name", "product_group_id", "measure_id", "width", "length", "min_cost", "cost", "user_id", "barcode", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ContragentGroup struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	ContragentGroupId int    `json:"contragent_group_id"`
	IsActive          bool   `json:"is_active"`
}

func ContragentGroupGet(id int, tx *sql.Tx) (ContragentGroup, error) {
	var c ContragentGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM contragent_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM contragent_group WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ContragentGroupId,
		&c.IsActive,
	)
	return c, err
}

func ContragentGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ContragentGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM contragent_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ContragentGroup{}
	for rows.Next() {
		var c ContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func ContragentGroupCreate(c ContragentGroup, tx *sql.Tx) (ContragentGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO contragent_group
            (name, contragent_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.ContragentGroupId,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContragentGroupUpdate(c ContragentGroup, tx *sql.Tx) (ContragentGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE contragent_group SET
                    name=?, contragent_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.ContragentGroupId,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContragentGroupDelete(id int, tx *sql.Tx) (ContragentGroup, error) {
	needCommit := false
	var err error
	var c ContragentGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = ContragentGroupGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE contragent_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func ContragentGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ContragentGroup, error) {

	if !ContragentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contragent_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ContragentGroup{}
	for rows.Next() {
		var c ContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContragentGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ContragentGroup, error) {

	if !ContragentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contragent_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ContragentGroup{}
	for rows.Next() {
		var c ContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContragentGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "contragent_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Contragent struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	ContragentGroupId int     `json:"contragent_group_id"`
	Phone             string  `json:"phone"`
	Email             string  `json:"email"`
	Web               string  `json:"web"`
	Comm              string  `json:"comm"`
	DirName           string  `json:"dir_name"`
	Search            string  `json:"search"`
	Total             float64 `json:"total"`
	FullName          string  `json:"full_name"`
	Edrpou            string  `json:"edrpou"`
	Ipn               string  `json:"ipn"`
	Iban              string  `json:"iban"`
	Bank              string  `json:"bank"`
	Mfo               string  `json:"mfo"`
	Fop               string  `json:"fop"`
	Address           string  `json:"address"`
	IsActive          bool    `json:"is_active"`
}

func ContragentGet(id int, tx *sql.Tx) (Contragent, error) {
	var c Contragent
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM contragent WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM contragent WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ContragentGroupId,
		&c.Phone,
		&c.Email,
		&c.Web,
		&c.Comm,
		&c.DirName,
		&c.Search,
		&c.Total,
		&c.FullName,
		&c.Edrpou,
		&c.Ipn,
		&c.Iban,
		&c.Bank,
		&c.Mfo,
		&c.Fop,
		&c.Address,
		&c.IsActive,
	)
	return c, err
}

func ContragentGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contragent, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM contragent"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contragent{}
	for rows.Next() {
		var c Contragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func ContragentCreate(c Contragent, tx *sql.Tx) (Contragent, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO contragent
            (name, contragent_group_id, phone, email, web, comm, dir_name, search, total, full_name, edrpou, ipn, iban, bank, mfo, fop, address, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.ContragentGroupId,
		c.Phone,
		c.Email,
		c.Web,
		c.Comm,
		c.DirName,
		c.Search,
		c.Total,
		c.FullName,
		c.Edrpou,
		c.Ipn,
		c.Iban,
		c.Bank,
		c.Mfo,
		c.Fop,
		c.Address,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContragentUpdate(c Contragent, tx *sql.Tx) (Contragent, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE contragent SET
                    name=?, contragent_group_id=?, phone=?, email=?, web=?, comm=?, dir_name=?, search=?, total=?, full_name=?, edrpou=?, ipn=?, iban=?, bank=?, mfo=?, fop=?, address=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.ContragentGroupId,
		c.Phone,
		c.Email,
		c.Web,
		c.Comm,
		c.DirName,
		c.Search,
		c.Total,
		c.FullName,
		c.Edrpou,
		c.Ipn,
		c.Iban,
		c.Bank,
		c.Mfo,
		c.Fop,
		c.Address,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContragentDelete(id int, tx *sql.Tx) (Contragent, error) {
	needCommit := false
	var err error
	var c Contragent
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = ContragentGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE contragent SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func ContragentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contragent, error) {

	if !ContragentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contragent WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contragent{}
	for rows.Next() {
		var c Contragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContragentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contragent, error) {

	if !ContragentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contragent WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contragent{}
	for rows.Next() {
		var c Contragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContragentTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "contragent_group_id", "phone", "email", "web", "comm", "dir_name", "search", "total", "full_name", "edrpou", "ipn", "iban", "bank", "mfo", "fop", "address", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ContragentFindByContragentSearchContactSearch(fs string) ([]Contragent, error) {
	fs = "%" + fs + "%"

	query := `
    SELECT DISTINCT contragent.* FROM contragent
    JOIN contact on contragent.id = contact.contragent_id
    WHERE
    contragent.search LIKE ?
    OR contact.search LIKE ?;;`

	rows, err := db.Query(query, fs, fs)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Contragent{}
	for rows.Next() {
		var c Contragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type Contact struct {
	Id           int     `json:"id"`
	ContragentId int     `json:"contragent_id"`
	Name         string  `json:"name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Viber        string  `json:"viber"`
	Telegram     string  `json:"telegram"`
	TelegramUid  int     `json:"telegram_uid"`
	Search       string  `json:"search"`
	Total        float64 `json:"total"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
}

func ContactGet(id int, tx *sql.Tx) (Contact, error) {
	var c Contact
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM contact WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM contact WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.ContragentId,
		&c.Name,
		&c.Phone,
		&c.Email,
		&c.Viber,
		&c.Telegram,
		&c.TelegramUid,
		&c.Search,
		&c.Total,
		&c.Comm,
		&c.IsActive,
	)
	return c, err
}

func ContactGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contact, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM contact"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contact{}
	for rows.Next() {
		var c Contact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func ContactCreate(c Contact, tx *sql.Tx) (Contact, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO contact
            (contragent_id, name, phone, email, viber, telegram, telegram_uid, search, total, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.ContragentId,
		c.Name,
		c.Phone,
		c.Email,
		c.Viber,
		c.Telegram,
		c.TelegramUid,
		c.Search,
		c.Total,
		c.Comm,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContactUpdate(c Contact, tx *sql.Tx) (Contact, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE contact SET
                    contragent_id=?, name=?, phone=?, email=?, viber=?, telegram=?, telegram_uid=?, search=?, total=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.ContragentId,
		c.Name,
		c.Phone,
		c.Email,
		c.Viber,
		c.Telegram,
		c.TelegramUid,
		c.Search,
		c.Total,
		c.Comm,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func ContactDelete(id int, tx *sql.Tx) (Contact, error) {
	needCommit := false
	var err error
	var c Contact
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = ContactGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE contact SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func ContactGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contact, error) {

	if !ContactTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contact WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contact{}
	for rows.Next() {
		var c Contact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContactGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Contact, error) {

	if !ContactTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM contact WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Contact{}
	for rows.Next() {
		var c Contact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func ContactTestForExistingField(fieldName string) bool {
	fields := []string{"id", "contragent_id", "name", "phone", "email", "viber", "telegram", "telegram_uid", "search", "total", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type OrderingStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func OrderingStatusGet(id int, tx *sql.Tx) (OrderingStatus, error) {
	var o OrderingStatus
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM ordering_status WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM ordering_status WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.IsActive,
	)
	return o, err
}

func OrderingStatusGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OrderingStatus, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM ordering_status"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OrderingStatus{}
	for rows.Next() {
		var o OrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OrderingStatusCreate(o OrderingStatus, tx *sql.Tx) (OrderingStatus, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO ordering_status
            (name, is_active)
            VALUES(?, ?);`
	res, err := tx.Exec(
		sql,
		o.Name,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OrderingStatusUpdate(o OrderingStatus, tx *sql.Tx) (OrderingStatus, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE ordering_status SET
                    name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.Name,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OrderingStatusDelete(id int, tx *sql.Tx) (OrderingStatus, error) {
	needCommit := false
	var err error
	var o OrderingStatus
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OrderingStatusGet(id, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE ordering_status SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OrderingStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OrderingStatus, error) {

	if !OrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM ordering_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OrderingStatus{}
	for rows.Next() {
		var o OrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OrderingStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OrderingStatus, error) {

	if !OrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM ordering_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OrderingStatus{}
	for rows.Next() {
		var o OrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OrderingStatusTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Ordering struct {
	Id               int     `json:"id"`
	DocumentUid      int     `json:"document_uid"`
	Name             string  `json:"name"`
	CreatedAt        string  `json:"created_at"`
	DeadlineAt       string  `json:"deadline_at"`
	UserId           int     `json:"user_id"`
	ContragentId     int     `json:"contragent_id"`
	ContactId        int     `json:"contact_id"`
	Price            float64 `json:"price"`
	Persent          float64 `json:"persent"`
	Profit           float64 `json:"profit"`
	Cost             float64 `json:"cost"`
	Info             string  `json:"info"`
	OrderingStatusId int     `json:"ordering_status_id"`
	IsActive         bool    `json:"is_active"`
}

func OrderingGet(id int, tx *sql.Tx) (Ordering, error) {
	var o Ordering
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM ordering WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM ordering WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.DocumentUid,
		&o.Name,
		&o.CreatedAt,
		&o.DeadlineAt,
		&o.UserId,
		&o.ContragentId,
		&o.ContactId,
		&o.Price,
		&o.Persent,
		&o.Profit,
		&o.Cost,
		&o.Info,
		&o.OrderingStatusId,
		&o.IsActive,
	)
	return o, err
}

func OrderingGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Ordering, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM ordering"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Ordering{}
	for rows.Next() {
		var o Ordering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OrderingCreate(o Ordering, tx *sql.Tx) (Ordering, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	doc := Document{Id: 0, DocType: "ordering", IsActive: true}
	doc, err = DocumentCreate(doc, tx)
	if err != nil {
		return o, err
	}
	o.DocumentUid = doc.Id

	t := time.Now()
	o.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO ordering
            (document_uid, name, created_at, deadline_at, user_id, contragent_id, contact_id, price, persent, profit, cost, info, ordering_status_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.DocumentUid,
		o.Name,
		o.CreatedAt,
		o.DeadlineAt,
		o.UserId,
		o.ContragentId,
		o.ContactId,
		o.Price,
		o.Persent,
		o.Profit,
		o.Cost,
		o.Info,
		o.OrderingStatusId,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OrderingUpdate(o Ordering, tx *sql.Tx) (Ordering, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE ordering SET
                    document_uid=?, name=?, created_at=?, deadline_at=?, user_id=?, contragent_id=?, contact_id=?, price=?, persent=?, profit=?, cost=?, info=?, ordering_status_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.DocumentUid,
		o.Name,
		o.CreatedAt,
		o.DeadlineAt,
		o.UserId,
		o.ContragentId,
		o.ContactId,
		o.Price,
		o.Persent,
		o.Profit,
		o.Cost,
		o.Info,
		o.OrderingStatusId,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OrderingDelete(id int, tx *sql.Tx) (Ordering, error) {
	needCommit := false
	var err error
	var o Ordering
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OrderingGet(id, tx)
	if err != nil {
		return o, err
	}

	product_to_orderings, err := ProductToOrderingGetByFilterInt("ordering_id", o.Id, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, product_to_ordering := range product_to_orderings {
		_, err = ProductToOrderingDelete(product_to_ordering.Id, tx)
		if err != nil {
			return o, err
		}
	}

	operation_to_orderings, err := OperationToOrderingGetByFilterInt("ordering_id", o.Id, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, operation_to_ordering := range operation_to_orderings {
		_, err = OperationToOrderingDelete(operation_to_ordering.Id, tx)
		if err != nil {
			return o, err
		}
	}

	matherial_to_orderings, err := MatherialToOrderingGetByFilterInt("ordering_id", o.Id, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, matherial_to_ordering := range matherial_to_orderings {
		_, err = MatherialToOrderingDelete(matherial_to_ordering.Id, tx)
		if err != nil {
			return o, err
		}
	}

	invoices, err := InvoiceGetByFilterInt("ordering_id", o.Id, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, invoice := range invoices {
		_, err = InvoiceDelete(invoice.Id, tx)
		if err != nil {
			return o, err
		}
	}

	cash_outs, err := CashOutGetByFilterInt("based_on", o.DocumentUid, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, cash_out := range cash_outs {
		_, err = CashOutDelete(cash_out.Id, tx)
		if err != nil {
			return o, err
		}
	}

	cash_ins, err := CashInGetByFilterInt("based_on", o.DocumentUid, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, cash_in := range cash_ins {
		_, err = CashInDelete(cash_in.Id, tx)
		if err != nil {
			return o, err
		}
	}

	whs_outs, err := WhsOutGetByFilterInt("based_on", o.DocumentUid, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, whs_out := range whs_outs {
		_, err = WhsOutDelete(whs_out.Id, tx)
		if err != nil {
			return o, err
		}
	}

	whs_ins, err := WhsInGetByFilterInt("based_on", o.DocumentUid, false, false, tx)
	if err != nil {
		return o, err
	}
	for _, whs_in := range whs_ins {
		_, err = WhsInDelete(whs_in.Id, tx)
		if err != nil {
			return o, err
		}
	}

	sql := `UPDATE ordering SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Ordering, error) {

	if !OrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Ordering{}
	for rows.Next() {
		var o Ordering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Ordering, error) {

	if !OrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Ordering{}
	for rows.Next() {
		var o Ordering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OrderingTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "created_at", "deadline_at", "user_id", "contragent_id", "contact_id", "price", "persent", "profit", "cost", "info", "ordering_status_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func OrderingGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]Ordering, error) {
	query := "SELECT * FROM ordering WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Ordering{}
	for rows.Next() {
		var o Ordering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OrderingGetBetweenDeadlineAt(deadline_at1, deadline_at2 string, withDeleted bool, deletedOnly bool) ([]Ordering, error) {
	query := "SELECT * FROM ordering WHERE deadline_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, deadline_at1, deadline_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Ordering{}
	for rows.Next() {
		var o Ordering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OrderingCashSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cash_sum) FROM ordering WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func OrderingGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM ordering WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM ordering WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type Owner struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Web      string  `json:"web"`
	Comm     string  `json:"comm"`
	Total    float64 `json:"total"`
	FullName string  `json:"full_name"`
	Edrpou   string  `json:"edrpou"`
	Ipn      string  `json:"ipn"`
	Iban     string  `json:"iban"`
	Bank     string  `json:"bank"`
	Mfo      string  `json:"mfo"`
	Fop      string  `json:"fop"`
	Address  string  `json:"address"`
	Sign     string  `json:"sign"`
	IsActive bool    `json:"is_active"`
}

func OwnerGet(id int, tx *sql.Tx) (Owner, error) {
	var o Owner
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM owner WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM owner WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.Phone,
		&o.Email,
		&o.Web,
		&o.Comm,
		&o.Total,
		&o.FullName,
		&o.Edrpou,
		&o.Ipn,
		&o.Iban,
		&o.Bank,
		&o.Mfo,
		&o.Fop,
		&o.Address,
		&o.Sign,
		&o.IsActive,
	)
	return o, err
}

func OwnerGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Owner, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM owner"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Owner{}
	for rows.Next() {
		var o Owner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OwnerCreate(o Owner, tx *sql.Tx) (Owner, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO owner
            (name, phone, email, web, comm, total, full_name, edrpou, ipn, iban, bank, mfo, fop, address, sign, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.Name,
		o.Phone,
		o.Email,
		o.Web,
		o.Comm,
		o.Total,
		o.FullName,
		o.Edrpou,
		o.Ipn,
		o.Iban,
		o.Bank,
		o.Mfo,
		o.Fop,
		o.Address,
		o.Sign,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OwnerUpdate(o Owner, tx *sql.Tx) (Owner, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE owner SET
                    name=?, phone=?, email=?, web=?, comm=?, total=?, full_name=?, edrpou=?, ipn=?, iban=?, bank=?, mfo=?, fop=?, address=?, sign=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.Name,
		o.Phone,
		o.Email,
		o.Web,
		o.Comm,
		o.Total,
		o.FullName,
		o.Edrpou,
		o.Ipn,
		o.Iban,
		o.Bank,
		o.Mfo,
		o.Fop,
		o.Address,
		o.Sign,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OwnerDelete(id int, tx *sql.Tx) (Owner, error) {
	needCommit := false
	var err error
	var o Owner
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OwnerGet(id, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE owner SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OwnerGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Owner, error) {

	if !OwnerTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM owner WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Owner{}
	for rows.Next() {
		var o Owner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OwnerGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Owner, error) {

	if !OwnerTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM owner WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Owner{}
	for rows.Next() {
		var o Owner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OwnerTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "phone", "email", "web", "comm", "total", "full_name", "edrpou", "ipn", "iban", "bank", "mfo", "fop", "address", "sign", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Invoice struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	OrderingId   int     `json:"ordering_id"`
	BasedOn      int     `json:"based_on"`
	OwnerId      int     `json:"owner_id"`
	Name         string  `json:"name"`
	CreatedAt    string  `json:"created_at"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
}

func InvoiceGet(id int, tx *sql.Tx) (Invoice, error) {
	var i Invoice
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM invoice WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM invoice WHERE id=?", id)
	}

	err := row.Scan(
		&i.Id,
		&i.DocumentUid,
		&i.OrderingId,
		&i.BasedOn,
		&i.OwnerId,
		&i.Name,
		&i.CreatedAt,
		&i.UserId,
		&i.ContragentId,
		&i.ContactId,
		&i.CashSum,
		&i.Comm,
		&i.IsActive,
	)
	return i, err
}

func InvoiceGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Invoice, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM invoice"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func InvoiceCreate(i Invoice, tx *sql.Tx) (Invoice, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	contragent, err := ContragentGet(i.ContragentId, tx)
	if err == nil {
		contragent.Total -= i.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return i, err
		}
	}

	contact, err := ContactGet(i.ContactId, tx)
	if err == nil {
		contact.Total -= i.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return i, err
		}
	}

	t := time.Now()
	i.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO invoice
            (document_uid, ordering_id, based_on, owner_id, name, created_at, user_id, contragent_id, contact_id, cash_sum, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		i.DocumentUid,
		i.OrderingId,
		i.BasedOn,
		i.OwnerId,
		i.Name,
		i.CreatedAt,
		i.UserId,
		i.ContragentId,
		i.ContactId,
		i.CashSum,
		i.Comm,
		i.IsActive,
	)
	if err != nil {
		return i, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return i, err
	}
	i.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func InvoiceUpdate(i Invoice, tx *sql.Tx) (Invoice, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	invoice, err := InvoiceGet(i.Id, tx)
	if err != nil {
		return i, err
	}

	contragent, err := ContragentGet(invoice.ContragentId, tx)
	if err == nil {
		contragent.Total += invoice.CashSum

	}

	if invoice.ContragentId != i.ContragentId {
		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return i, err
		}
		contragent, err = ContragentGet(i.ContragentId, tx)
		if err != nil {
			return i, err
		}
	}
	contragent.Total -= i.CashSum

	_, err = ContragentUpdate(contragent, tx)
	if err != nil {
		return i, err
	}

	contact, err := ContactGet(invoice.ContactId, tx)
	if err == nil {
		contact.Total += invoice.CashSum

	}

	if invoice.ContactId != i.ContactId {
		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return i, err
		}
		contact, err = ContactGet(i.ContactId, tx)
		if err != nil {
			return i, err
		}
	}
	contact.Total -= i.CashSum

	_, err = ContactUpdate(contact, tx)
	if err != nil {
		return i, err
	}

	sql := `UPDATE invoice SET
                    document_uid=?, ordering_id=?, based_on=?, owner_id=?, name=?, created_at=?, user_id=?, contragent_id=?, contact_id=?, cash_sum=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		i.DocumentUid,
		i.OrderingId,
		i.BasedOn,
		i.OwnerId,
		i.Name,
		i.CreatedAt,
		i.UserId,
		i.ContragentId,
		i.ContactId,
		i.CashSum,
		i.Comm,
		i.IsActive,
		i.Id,
	)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func InvoiceDelete(id int, tx *sql.Tx) (Invoice, error) {
	needCommit := false
	var err error
	var i Invoice
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	i, err = InvoiceGet(id, tx)
	if err != nil {
		return i, err
	}

	contragent, err := ContragentGet(i.ContragentId, tx)
	if err == nil {
		contragent.Total += i.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return i, err
		}
	}

	contact, err := ContactGet(i.ContactId, tx)
	if err == nil {
		contact.Total += i.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return i, err
		}
	}

	item_to_invoices, err := ItemToInvoiceGetByFilterInt("invoice_id", i.Id, false, false, tx)
	if err != nil {
		return i, err
	}
	for _, item_to_invoice := range item_to_invoices {
		_, err = ItemToInvoiceDelete(item_to_invoice.Id, tx)
		if err != nil {
			return i, err
		}
	}

	sql := `UPDATE invoice SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, i.Id)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	i.IsActive = false
	return i, nil
}

func InvoiceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Invoice, error) {

	if !InvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM invoice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func InvoiceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Invoice, error) {

	if !InvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM invoice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func InvoiceTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "ordering_id", "based_on", "owner_id", "name", "created_at", "user_id", "contragent_id", "contact_id", "cash_sum", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func InvoiceGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]Invoice, error) {
	query := "SELECT * FROM invoice WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func InvoiceCashSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cash_sum) FROM invoice WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func InvoiceGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM invoice WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM invoice WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type ItemToInvoice struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	InvoiceId int     `json:"invoice_id"`
	Number    float64 `json:"number"`
	MeasureId int     `json:"measure_id"`
	Price     float64 `json:"price"`
	Cost      float64 `json:"cost"`
	IsActive  bool    `json:"is_active"`
}

func ItemToInvoiceGet(id int, tx *sql.Tx) (ItemToInvoice, error) {
	var i ItemToInvoice
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM item_to_invoice WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM item_to_invoice WHERE id=?", id)
	}

	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.InvoiceId,
		&i.Number,
		&i.MeasureId,
		&i.Price,
		&i.Cost,
		&i.IsActive,
	)
	return i, err
}

func ItemToInvoiceGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToInvoice, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM item_to_invoice"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToInvoice{}
	for rows.Next() {
		var i ItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func ItemToInvoiceCreate(i ItemToInvoice, tx *sql.Tx) (ItemToInvoice, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	invoice, err := InvoiceGet(i.InvoiceId, tx)
	if err == nil {
		invoice.CashSum += i.Cost

		_, err = InvoiceUpdate(invoice, tx)
		if err != nil {
			return i, err
		}
	}

	sql := `INSERT INTO item_to_invoice
            (name, invoice_id, number, measure_id, price, cost, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		i.Name,
		i.InvoiceId,
		i.Number,
		i.MeasureId,
		i.Price,
		i.Cost,
		i.IsActive,
	)
	if err != nil {
		return i, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return i, err
	}
	i.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func ItemToInvoiceUpdate(i ItemToInvoice, tx *sql.Tx) (ItemToInvoice, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	item_to_invoice, err := ItemToInvoiceGet(i.Id, tx)
	if err != nil {
		return i, err
	}

	invoice, err := InvoiceGet(item_to_invoice.InvoiceId, tx)
	if err == nil {
		invoice.CashSum -= item_to_invoice.Cost

	}

	if item_to_invoice.InvoiceId != i.InvoiceId {
		_, err = InvoiceUpdate(invoice, tx)
		if err != nil {
			return i, err
		}
		invoice, err = InvoiceGet(i.InvoiceId, tx)
		if err != nil {
			return i, err
		}
	}
	invoice.CashSum += i.Cost

	_, err = InvoiceUpdate(invoice, tx)
	if err != nil {
		return i, err
	}

	sql := `UPDATE item_to_invoice SET
                    name=?, invoice_id=?, number=?, measure_id=?, price=?, cost=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		i.Name,
		i.InvoiceId,
		i.Number,
		i.MeasureId,
		i.Price,
		i.Cost,
		i.IsActive,
		i.Id,
	)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func ItemToInvoiceDelete(id int, tx *sql.Tx) (ItemToInvoice, error) {
	needCommit := false
	var err error
	var i ItemToInvoice
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	i, err = ItemToInvoiceGet(id, tx)
	if err != nil {
		return i, err
	}

	invoice, err := InvoiceGet(i.InvoiceId, tx)
	if err == nil {
		invoice.CashSum -= i.Cost

		_, err = InvoiceUpdate(invoice, tx)
		if err != nil {
			return i, err
		}
	}

	sql := `UPDATE item_to_invoice SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, i.Id)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	i.IsActive = false
	return i, nil
}

func ItemToInvoiceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToInvoice, error) {

	if !ItemToInvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM item_to_invoice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToInvoice{}
	for rows.Next() {
		var i ItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func ItemToInvoiceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToInvoice, error) {

	if !ItemToInvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM item_to_invoice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToInvoice{}
	for rows.Next() {
		var i ItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func ItemToInvoiceTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "invoice_id", "number", "measure_id", "price", "cost", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ItemToInvoiceCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM item_to_invoice WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func ItemToInvoiceGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM item_to_invoice WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM item_to_invoice WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type ProductToOrderingStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func ProductToOrderingStatusGet(id int, tx *sql.Tx) (ProductToOrderingStatus, error) {
	var p ProductToOrderingStatus
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM product_to_ordering_status WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM product_to_ordering_status WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.IsActive,
	)
	return p, err
}

func ProductToOrderingStatusGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrderingStatus, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM product_to_ordering_status"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrderingStatus{}
	for rows.Next() {
		var p ProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProductToOrderingStatusCreate(p ProductToOrderingStatus, tx *sql.Tx) (ProductToOrderingStatus, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO product_to_ordering_status
            (name, is_active)
            VALUES(?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToOrderingStatusUpdate(p ProductToOrderingStatus, tx *sql.Tx) (ProductToOrderingStatus, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE product_to_ordering_status SET
                    name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToOrderingStatusDelete(id int, tx *sql.Tx) (ProductToOrderingStatus, error) {
	needCommit := false
	var err error
	var p ProductToOrderingStatus
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProductToOrderingStatusGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE product_to_ordering_status SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProductToOrderingStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrderingStatus, error) {

	if !ProductToOrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_ordering_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrderingStatus{}
	for rows.Next() {
		var p ProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToOrderingStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrderingStatus, error) {

	if !ProductToOrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_ordering_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrderingStatus{}
	for rows.Next() {
		var p ProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToOrderingStatusTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ProductToOrdering struct {
	Id                        int     `json:"id"`
	Name                      string  `json:"name"`
	OrderingId                int     `json:"ordering_id"`
	ProductId                 int     `json:"product_id"`
	UserId                    int     `json:"user_id"`
	DeadlineAt                string  `json:"deadline_at"`
	ProductToOrderingStatusId int     `json:"product_to_ordering_status_id"`
	Width                     float64 `json:"width"`
	Length                    float64 `json:"length"`
	Pieces                    int     `json:"pieces"`
	Number                    float64 `json:"number"`
	Price                     float64 `json:"price"`
	Persent                   float64 `json:"persent"`
	Profit                    float64 `json:"profit"`
	Cost                      float64 `json:"cost"`
	Info                      string  `json:"info"`
	ProductToOrderingId       int     `json:"product_to_ordering_id"`
	IsActive                  bool    `json:"is_active"`
}

func ProductToOrderingGet(id int, tx *sql.Tx) (ProductToOrdering, error) {
	var p ProductToOrdering
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM product_to_ordering WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM product_to_ordering WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.OrderingId,
		&p.ProductId,
		&p.UserId,
		&p.DeadlineAt,
		&p.ProductToOrderingStatusId,
		&p.Width,
		&p.Length,
		&p.Pieces,
		&p.Number,
		&p.Price,
		&p.Persent,
		&p.Profit,
		&p.Cost,
		&p.Info,
		&p.ProductToOrderingId,
		&p.IsActive,
	)
	return p, err
}

func ProductToOrderingGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrdering, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM product_to_ordering"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrdering{}
	for rows.Next() {
		var p ProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProductToOrderingCreate(p ProductToOrdering, tx *sql.Tx) (ProductToOrdering, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO product_to_ordering
            (name, ordering_id, product_id, user_id, deadline_at, product_to_ordering_status_id, width, length, pieces, number, price, persent, profit, cost, info, product_to_ordering_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.OrderingId,
		p.ProductId,
		p.UserId,
		p.DeadlineAt,
		p.ProductToOrderingStatusId,
		p.Width,
		p.Length,
		p.Pieces,
		p.Number,
		p.Price,
		p.Persent,
		p.Profit,
		p.Cost,
		p.Info,
		p.ProductToOrderingId,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToOrderingUpdate(p ProductToOrdering, tx *sql.Tx) (ProductToOrdering, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE product_to_ordering SET
                    name=?, ordering_id=?, product_id=?, user_id=?, deadline_at=?, product_to_ordering_status_id=?, width=?, length=?, pieces=?, number=?, price=?, persent=?, profit=?, cost=?, info=?, product_to_ordering_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.OrderingId,
		p.ProductId,
		p.UserId,
		p.DeadlineAt,
		p.ProductToOrderingStatusId,
		p.Width,
		p.Length,
		p.Pieces,
		p.Number,
		p.Price,
		p.Persent,
		p.Profit,
		p.Cost,
		p.Info,
		p.ProductToOrderingId,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToOrderingDelete(id int, tx *sql.Tx) (ProductToOrdering, error) {
	needCommit := false
	var err error
	var p ProductToOrdering
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProductToOrderingGet(id, tx)
	if err != nil {
		return p, err
	}

	product_to_orderings, err := ProductToOrderingGetByFilterInt("product_to_ordering_id", p.Id, false, false, tx)
	if err != nil {
		return p, err
	}
	for _, product_to_ordering := range product_to_orderings {
		_, err = ProductToOrderingDelete(product_to_ordering.Id, tx)
		if err != nil {
			return p, err
		}
	}

	matherial_to_orderings, err := MatherialToOrderingGetByFilterInt("product_to_ordering_id", p.Id, false, false, tx)
	if err != nil {
		return p, err
	}
	for _, matherial_to_ordering := range matherial_to_orderings {
		_, err = MatherialToOrderingDelete(matherial_to_ordering.Id, tx)
		if err != nil {
			return p, err
		}
	}

	operation_to_orderings, err := OperationToOrderingGetByFilterInt("product_to_ordering_id", p.Id, false, false, tx)
	if err != nil {
		return p, err
	}
	for _, operation_to_ordering := range operation_to_orderings {
		_, err = OperationToOrderingDelete(operation_to_ordering.Id, tx)
		if err != nil {
			return p, err
		}
	}

	sql := `UPDATE product_to_ordering SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProductToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrdering, error) {

	if !ProductToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrdering{}
	for rows.Next() {
		var p ProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToOrdering, error) {

	if !ProductToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToOrdering{}
	for rows.Next() {
		var p ProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToOrderingTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "ordering_id", "product_id", "user_id", "deadline_at", "product_to_ordering_status_id", "width", "length", "pieces", "number", "price", "persent", "profit", "cost", "info", "product_to_ordering_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ProductToOrderingCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM product_to_ordering WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func ProductToOrderingGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM product_to_ordering WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM product_to_ordering WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type MatherialToOrdering struct {
	Id                  int     `json:"id"`
	OrderingId          int     `json:"ordering_id"`
	MatherialId         int     `json:"matherial_id"`
	Width               float64 `json:"width"`
	Length              float64 `json:"length"`
	Pieces              int     `json:"pieces"`
	ColorId             int     `json:"color_id"`
	UserId              int     `json:"user_id"`
	Number              float64 `json:"number"`
	Price               float64 `json:"price"`
	Persent             float64 `json:"persent"`
	Profit              float64 `json:"profit"`
	Cost                float64 `json:"cost"`
	Comm                string  `json:"comm"`
	ProductToOrderingId int     `json:"product_to_ordering_id"`
	IsActive            bool    `json:"is_active"`
}

func MatherialToOrderingGet(id int, tx *sql.Tx) (MatherialToOrdering, error) {
	var m MatherialToOrdering
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_to_ordering WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_to_ordering WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.OrderingId,
		&m.MatherialId,
		&m.Width,
		&m.Length,
		&m.Pieces,
		&m.ColorId,
		&m.UserId,
		&m.Number,
		&m.Price,
		&m.Persent,
		&m.Profit,
		&m.Cost,
		&m.Comm,
		&m.ProductToOrderingId,
		&m.IsActive,
	)
	return m, err
}

func MatherialToOrderingGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToOrdering, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_to_ordering"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToOrdering{}
	for rows.Next() {
		var m MatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialToOrderingCreate(m MatherialToOrdering, tx *sql.Tx) (MatherialToOrdering, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO matherial_to_ordering
            (ordering_id, matherial_id, width, length, pieces, color_id, user_id, number, price, persent, profit, cost, comm, product_to_ordering_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.OrderingId,
		m.MatherialId,
		m.Width,
		m.Length,
		m.Pieces,
		m.ColorId,
		m.UserId,
		m.Number,
		m.Price,
		m.Persent,
		m.Profit,
		m.Cost,
		m.Comm,
		m.ProductToOrderingId,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToOrderingUpdate(m MatherialToOrdering, tx *sql.Tx) (MatherialToOrdering, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial_to_ordering SET
                    ordering_id=?, matherial_id=?, width=?, length=?, pieces=?, color_id=?, user_id=?, number=?, price=?, persent=?, profit=?, cost=?, comm=?, product_to_ordering_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.OrderingId,
		m.MatherialId,
		m.Width,
		m.Length,
		m.Pieces,
		m.ColorId,
		m.UserId,
		m.Number,
		m.Price,
		m.Persent,
		m.Profit,
		m.Cost,
		m.Comm,
		m.ProductToOrderingId,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToOrderingDelete(id int, tx *sql.Tx) (MatherialToOrdering, error) {
	needCommit := false
	var err error
	var m MatherialToOrdering
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialToOrderingGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_to_ordering SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToOrdering, error) {

	if !MatherialToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToOrdering{}
	for rows.Next() {
		var m MatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToOrdering, error) {

	if !MatherialToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToOrdering{}
	for rows.Next() {
		var m MatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToOrderingTestForExistingField(fieldName string) bool {
	fields := []string{"id", "ordering_id", "matherial_id", "width", "length", "pieces", "color_id", "user_id", "number", "price", "persent", "profit", "cost", "comm", "product_to_ordering_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func MatherialToOrderingCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_ordering WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func MatherialToOrderingGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_ordering WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_ordering WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type MatherialToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	MatherialId   int     `json:"matherial_id"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	Comm          string  `json:"comm"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
}

func MatherialToProductGet(id int, tx *sql.Tx) (MatherialToProduct, error) {
	var m MatherialToProduct
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_to_product WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_to_product WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.ProductId,
		&m.MatherialId,
		&m.Number,
		&m.Coeff,
		&m.Cost,
		&m.ListName,
		&m.IsMultiselect,
		&m.Comm,
		&m.IsUsed,
		&m.IsActive,
	)
	return m, err
}

func MatherialToProductGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToProduct, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_to_product"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToProduct{}
	for rows.Next() {
		var m MatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialToProductCreate(m MatherialToProduct, tx *sql.Tx) (MatherialToProduct, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO matherial_to_product
            (product_id, matherial_id, number, coeff, cost, list_name, is_multiselect, comm, is_used, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.ProductId,
		m.MatherialId,
		m.Number,
		m.Coeff,
		m.Cost,
		m.ListName,
		m.IsMultiselect,
		m.Comm,
		m.IsUsed,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToProductUpdate(m MatherialToProduct, tx *sql.Tx) (MatherialToProduct, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial_to_product SET
                    product_id=?, matherial_id=?, number=?, coeff=?, cost=?, list_name=?, is_multiselect=?, comm=?, is_used=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.ProductId,
		m.MatherialId,
		m.Number,
		m.Coeff,
		m.Cost,
		m.ListName,
		m.IsMultiselect,
		m.Comm,
		m.IsUsed,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToProductDelete(id int, tx *sql.Tx) (MatherialToProduct, error) {
	needCommit := false
	var err error
	var m MatherialToProduct
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialToProductGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_to_product SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToProduct, error) {

	if !MatherialToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToProduct{}
	for rows.Next() {
		var m MatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToProduct, error) {

	if !MatherialToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToProduct{}
	for rows.Next() {
		var m MatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToProductTestForExistingField(fieldName string) bool {
	fields := []string{"id", "product_id", "matherial_id", "number", "coeff", "cost", "list_name", "is_multiselect", "comm", "is_used", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func MatherialToProductCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_product WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func MatherialToProductGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_product WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM matherial_to_product WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type OperationToOrdering struct {
	Id                  int     `json:"id"`
	OrderingId          int     `json:"ordering_id"`
	OperationId         int     `json:"operation_id"`
	UserId              int     `json:"user_id"`
	Number              float64 `json:"number"`
	Price               float64 `json:"price"`
	UserSum             float64 `json:"user_sum"`
	Cost                float64 `json:"cost"`
	EquipmentId         int     `json:"equipment_id"`
	EquipmentCost       float64 `json:"equipment_cost"`
	Comm                string  `json:"comm"`
	ProductToOrderingId int     `json:"product_to_ordering_id"`
	IsActive            bool    `json:"is_active"`
}

func OperationToOrderingGet(id int, tx *sql.Tx) (OperationToOrdering, error) {
	var o OperationToOrdering
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM operation_to_ordering WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM operation_to_ordering WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.OrderingId,
		&o.OperationId,
		&o.UserId,
		&o.Number,
		&o.Price,
		&o.UserSum,
		&o.Cost,
		&o.EquipmentId,
		&o.EquipmentCost,
		&o.Comm,
		&o.ProductToOrderingId,
		&o.IsActive,
	)
	return o, err
}

func OperationToOrderingGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToOrdering, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM operation_to_ordering"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToOrdering{}
	for rows.Next() {
		var o OperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OperationToOrderingCreate(o OperationToOrdering, tx *sql.Tx) (OperationToOrdering, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	equipment, err := EquipmentGet(o.EquipmentId, tx)
	if err == nil {
		equipment.Total += o.EquipmentCost

		_, err = EquipmentUpdate(equipment, tx)
		if err != nil {
			return o, err
		}
	}

	sql := `INSERT INTO operation_to_ordering
            (ordering_id, operation_id, user_id, number, price, user_sum, cost, equipment_id, equipment_cost, comm, product_to_ordering_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.OrderingId,
		o.OperationId,
		o.UserId,
		o.Number,
		o.Price,
		o.UserSum,
		o.Cost,
		o.EquipmentId,
		o.EquipmentCost,
		o.Comm,
		o.ProductToOrderingId,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationToOrderingUpdate(o OperationToOrdering, tx *sql.Tx) (OperationToOrdering, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	operation_to_ordering, err := OperationToOrderingGet(o.Id, tx)
	if err != nil {
		return o, err
	}

	equipment, err := EquipmentGet(operation_to_ordering.EquipmentId, tx)
	if err == nil {
		equipment.Total -= operation_to_ordering.EquipmentCost

	}

	if operation_to_ordering.EquipmentId != o.EquipmentId {
		_, err = EquipmentUpdate(equipment, tx)
		if err != nil {
			return o, err
		}
		equipment, err = EquipmentGet(o.EquipmentId, tx)
		if err != nil {
			return o, err
		}
	}
	equipment.Total += o.EquipmentCost

	_, err = EquipmentUpdate(equipment, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE operation_to_ordering SET
                    ordering_id=?, operation_id=?, user_id=?, number=?, price=?, user_sum=?, cost=?, equipment_id=?, equipment_cost=?, comm=?, product_to_ordering_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.OrderingId,
		o.OperationId,
		o.UserId,
		o.Number,
		o.Price,
		o.UserSum,
		o.Cost,
		o.EquipmentId,
		o.EquipmentCost,
		o.Comm,
		o.ProductToOrderingId,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationToOrderingDelete(id int, tx *sql.Tx) (OperationToOrdering, error) {
	needCommit := false
	var err error
	var o OperationToOrdering
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OperationToOrderingGet(id, tx)
	if err != nil {
		return o, err
	}

	equipment, err := EquipmentGet(o.EquipmentId, tx)
	if err == nil {
		equipment.Total -= o.EquipmentCost

		_, err = EquipmentUpdate(equipment, tx)
		if err != nil {
			return o, err
		}
	}

	sql := `UPDATE operation_to_ordering SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OperationToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToOrdering, error) {

	if !OperationToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToOrdering{}
	for rows.Next() {
		var o OperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToOrdering, error) {

	if !OperationToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_to_ordering WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToOrdering{}
	for rows.Next() {
		var o OperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationToOrderingTestForExistingField(fieldName string) bool {
	fields := []string{"id", "ordering_id", "operation_id", "user_id", "number", "price", "user_sum", "cost", "equipment_id", "equipment_cost", "comm", "product_to_ordering_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func OperationToOrderingCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM operation_to_ordering WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func OperationToOrderingGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM operation_to_ordering WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM operation_to_ordering WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type OperationToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	OperationId   int     `json:"operation_id"`
	UserId        int     `json:"user_id"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	EquipmentId   int     `json:"equipment_id"`
	EquipmentCost float64 `json:"equipment_cost"`
	Comm          string  `json:"comm"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
}

func OperationToProductGet(id int, tx *sql.Tx) (OperationToProduct, error) {
	var o OperationToProduct
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM operation_to_product WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM operation_to_product WHERE id=?", id)
	}

	err := row.Scan(
		&o.Id,
		&o.ProductId,
		&o.OperationId,
		&o.UserId,
		&o.Number,
		&o.Coeff,
		&o.Cost,
		&o.ListName,
		&o.IsMultiselect,
		&o.EquipmentId,
		&o.EquipmentCost,
		&o.Comm,
		&o.IsUsed,
		&o.IsActive,
	)
	return o, err
}

func OperationToProductGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToProduct, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM operation_to_product"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToProduct{}
	for rows.Next() {
		var o OperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func OperationToProductCreate(o OperationToProduct, tx *sql.Tx) (OperationToProduct, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO operation_to_product
            (product_id, operation_id, user_id, number, coeff, cost, list_name, is_multiselect, equipment_id, equipment_cost, comm, is_used, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		o.ProductId,
		o.OperationId,
		o.UserId,
		o.Number,
		o.Coeff,
		o.Cost,
		o.ListName,
		o.IsMultiselect,
		o.EquipmentId,
		o.EquipmentCost,
		o.Comm,
		o.IsUsed,
		o.IsActive,
	)
	if err != nil {
		return o, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return o, err
	}
	o.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationToProductUpdate(o OperationToProduct, tx *sql.Tx) (OperationToProduct, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE operation_to_product SET
                    product_id=?, operation_id=?, user_id=?, number=?, coeff=?, cost=?, list_name=?, is_multiselect=?, equipment_id=?, equipment_cost=?, comm=?, is_used=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		o.ProductId,
		o.OperationId,
		o.UserId,
		o.Number,
		o.Coeff,
		o.Cost,
		o.ListName,
		o.IsMultiselect,
		o.EquipmentId,
		o.EquipmentCost,
		o.Comm,
		o.IsUsed,
		o.IsActive,
		o.Id,
	)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func OperationToProductDelete(id int, tx *sql.Tx) (OperationToProduct, error) {
	needCommit := false
	var err error
	var o OperationToProduct
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return o, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	o, err = OperationToProductGet(id, tx)
	if err != nil {
		return o, err
	}

	sql := `UPDATE operation_to_product SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, o.Id)
	if err != nil {
		return o, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return o, err
		}
	}
	o.IsActive = false
	return o, nil
}

func OperationToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToProduct, error) {

	if !OperationToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToProduct{}
	for rows.Next() {
		var o OperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]OperationToProduct, error) {

	if !OperationToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM operation_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OperationToProduct{}
	for rows.Next() {
		var o OperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func OperationToProductTestForExistingField(fieldName string) bool {
	fields := []string{"id", "product_id", "operation_id", "user_id", "number", "coeff", "cost", "list_name", "is_multiselect", "equipment_id", "equipment_cost", "comm", "is_used", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func OperationToProductCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM operation_to_product WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func OperationToProductGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM operation_to_product WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM operation_to_product WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type ProductToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	Product2Id    int     `json:"product2_id"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
}

func ProductToProductGet(id int, tx *sql.Tx) (ProductToProduct, error) {
	var p ProductToProduct
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM product_to_product WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM product_to_product WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.ProductId,
		&p.Product2Id,
		&p.Width,
		&p.Length,
		&p.Number,
		&p.Coeff,
		&p.Cost,
		&p.ListName,
		&p.IsMultiselect,
		&p.IsUsed,
		&p.IsActive,
	)
	return p, err
}

func ProductToProductGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToProduct, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM product_to_product"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToProduct{}
	for rows.Next() {
		var p ProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProductToProductCreate(p ProductToProduct, tx *sql.Tx) (ProductToProduct, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO product_to_product
            (product_id, product2_id, width, length, number, coeff, cost, list_name, is_multiselect, is_used, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.ProductId,
		p.Product2Id,
		p.Width,
		p.Length,
		p.Number,
		p.Coeff,
		p.Cost,
		p.ListName,
		p.IsMultiselect,
		p.IsUsed,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToProductUpdate(p ProductToProduct, tx *sql.Tx) (ProductToProduct, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE product_to_product SET
                    product_id=?, product2_id=?, width=?, length=?, number=?, coeff=?, cost=?, list_name=?, is_multiselect=?, is_used=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.ProductId,
		p.Product2Id,
		p.Width,
		p.Length,
		p.Number,
		p.Coeff,
		p.Cost,
		p.ListName,
		p.IsMultiselect,
		p.IsUsed,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProductToProductDelete(id int, tx *sql.Tx) (ProductToProduct, error) {
	needCommit := false
	var err error
	var p ProductToProduct
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProductToProductGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE product_to_product SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProductToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToProduct, error) {

	if !ProductToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToProduct{}
	for rows.Next() {
		var p ProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProductToProduct, error) {

	if !ProductToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM product_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProductToProduct{}
	for rows.Next() {
		var p ProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProductToProductTestForExistingField(fieldName string) bool {
	fields := []string{"id", "product_id", "product2_id", "width", "length", "number", "coeff", "cost", "list_name", "is_multiselect", "is_used", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ProductToProductCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM product_to_product WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func ProductToProductGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM product_to_product WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM product_to_product WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type CboxCheck struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	FsUid        string  `json:"fs_uid"`
	CheckboxUid  string  `json:"checkbox_uid"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	DocumentUid  int     `json:"document_uid"`
	OrderingId   int     `json:"ordering_id"`
	BasedOn      int     `json:"based_on"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Discount     float64 `json:"discount"`
	Comm         string  `json:"comm"`
	IsCash       bool    `json:"is_cash"`
	IsActive     bool    `json:"is_active"`
}

func CboxCheckGet(id int, tx *sql.Tx) (CboxCheck, error) {
	var c CboxCheck
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM cbox_check WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM cbox_check WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.FsUid,
		&c.CheckboxUid,
		&c.UserId,
		&c.ContragentId,
		&c.DocumentUid,
		&c.OrderingId,
		&c.BasedOn,
		&c.CreatedAt,
		&c.CashSum,
		&c.Discount,
		&c.Comm,
		&c.IsCash,
		&c.IsActive,
	)
	return c, err
}

func CboxCheckGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CboxCheck, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM cbox_check"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CboxCheck{}
	for rows.Next() {
		var c CboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CboxCheckCreate(c CboxCheck, tx *sql.Tx) (CboxCheck, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	t := time.Now()
	c.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO cbox_check
            (name, fs_uid, checkbox_uid, user_id, contragent_id, document_uid, ordering_id, based_on, created_at, cash_sum, discount, comm, is_cash, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.FsUid,
		c.CheckboxUid,
		c.UserId,
		c.ContragentId,
		c.DocumentUid,
		c.OrderingId,
		c.BasedOn,
		c.CreatedAt,
		c.CashSum,
		c.Discount,
		c.Comm,
		c.IsCash,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CboxCheckUpdate(c CboxCheck, tx *sql.Tx) (CboxCheck, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE cbox_check SET
                    name=?, fs_uid=?, checkbox_uid=?, user_id=?, contragent_id=?, document_uid=?, ordering_id=?, based_on=?, created_at=?, cash_sum=?, discount=?, comm=?, is_cash=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.FsUid,
		c.CheckboxUid,
		c.UserId,
		c.ContragentId,
		c.DocumentUid,
		c.OrderingId,
		c.BasedOn,
		c.CreatedAt,
		c.CashSum,
		c.Discount,
		c.Comm,
		c.IsCash,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CboxCheckDelete(id int, tx *sql.Tx) (CboxCheck, error) {
	needCommit := false
	var err error
	var c CboxCheck
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CboxCheckGet(id, tx)
	if err != nil {
		return c, err
	}

	item_to_cbox_checks, err := ItemToCboxCheckGetByFilterInt("cbox_check_id", c.Id, false, false, tx)
	if err != nil {
		return c, err
	}
	for _, item_to_cbox_check := range item_to_cbox_checks {
		_, err = ItemToCboxCheckDelete(item_to_cbox_check.Id, tx)
		if err != nil {
			return c, err
		}
	}

	sql := `UPDATE cbox_check SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CboxCheckGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CboxCheck, error) {

	if !CboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cbox_check WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CboxCheck{}
	for rows.Next() {
		var c CboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CboxCheckGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CboxCheck, error) {

	if !CboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cbox_check WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CboxCheck{}
	for rows.Next() {
		var c CboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CboxCheckTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "fs_uid", "checkbox_uid", "user_id", "contragent_id", "document_uid", "ordering_id", "based_on", "created_at", "cash_sum", "discount", "comm", "is_cash", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func CboxCheckGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]CboxCheck, error) {
	query := "SELECT * FROM cbox_check WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CboxCheck{}
	for rows.Next() {
		var c CboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type ItemToCboxCheck struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	CboxCheckId int     `json:"cbox_check_id"`
	Number      float64 `json:"number"`
	MeasureId   int     `json:"measure_id"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	Cost        float64 `json:"cost"`
	ItemCode    string  `json:"item_code"`
	IsActive    bool    `json:"is_active"`
}

func ItemToCboxCheckGet(id int, tx *sql.Tx) (ItemToCboxCheck, error) {
	var i ItemToCboxCheck
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM item_to_cbox_check WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM item_to_cbox_check WHERE id=?", id)
	}

	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.CboxCheckId,
		&i.Number,
		&i.MeasureId,
		&i.Price,
		&i.Discount,
		&i.Cost,
		&i.ItemCode,
		&i.IsActive,
	)
	return i, err
}

func ItemToCboxCheckGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToCboxCheck, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM item_to_cbox_check"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToCboxCheck{}
	for rows.Next() {
		var i ItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func ItemToCboxCheckCreate(i ItemToCboxCheck, tx *sql.Tx) (ItemToCboxCheck, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO item_to_cbox_check
            (name, cbox_check_id, number, measure_id, price, discount, cost, item_code, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		i.Name,
		i.CboxCheckId,
		i.Number,
		i.MeasureId,
		i.Price,
		i.Discount,
		i.Cost,
		i.ItemCode,
		i.IsActive,
	)
	if err != nil {
		return i, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return i, err
	}
	i.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func ItemToCboxCheckUpdate(i ItemToCboxCheck, tx *sql.Tx) (ItemToCboxCheck, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE item_to_cbox_check SET
                    name=?, cbox_check_id=?, number=?, measure_id=?, price=?, discount=?, cost=?, item_code=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		i.Name,
		i.CboxCheckId,
		i.Number,
		i.MeasureId,
		i.Price,
		i.Discount,
		i.Cost,
		i.ItemCode,
		i.IsActive,
		i.Id,
	)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func ItemToCboxCheckDelete(id int, tx *sql.Tx) (ItemToCboxCheck, error) {
	needCommit := false
	var err error
	var i ItemToCboxCheck
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return i, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	i, err = ItemToCboxCheckGet(id, tx)
	if err != nil {
		return i, err
	}

	sql := `UPDATE item_to_cbox_check SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, i.Id)
	if err != nil {
		return i, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return i, err
		}
	}
	i.IsActive = false
	return i, nil
}

func ItemToCboxCheckGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToCboxCheck, error) {

	if !ItemToCboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM item_to_cbox_check WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToCboxCheck{}
	for rows.Next() {
		var i ItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func ItemToCboxCheckGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ItemToCboxCheck, error) {

	if !ItemToCboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM item_to_cbox_check WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ItemToCboxCheck{}
	for rows.Next() {
		var i ItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func ItemToCboxCheckTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "cbox_check_id", "number", "measure_id", "price", "discount", "cost", "item_code", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ItemToCboxCheckCostGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cost) FROM item_to_cbox_check WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func ItemToCboxCheckGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cost) FROM item_to_cbox_check WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cost) FROM item_to_cbox_check WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type CashIn struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	CashId       int     `json:"cash_id"`
	UserId       int     `json:"user_id"`
	BasedOn      int     `json:"based_on"`
	CboxCheckId  int     `json:"cbox_check_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
}

func CashInGet(id int, tx *sql.Tx) (CashIn, error) {
	var c CashIn
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM cash_in WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM cash_in WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.DocumentUid,
		&c.Name,
		&c.CashId,
		&c.UserId,
		&c.BasedOn,
		&c.CboxCheckId,
		&c.ContragentId,
		&c.ContactId,
		&c.CreatedAt,
		&c.CashSum,
		&c.Comm,
		&c.IsActive,
	)
	return c, err
}

func CashInGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashIn, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM cash_in"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashIn{}
	for rows.Next() {
		var c CashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CashInCreate(c CashIn, tx *sql.Tx) (CashIn, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	cash, err := CashGet(c.CashId, tx)
	if err == nil {
		cash.Total += c.CashSum

		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
	}

	contragent, err := ContragentGet(c.ContragentId, tx)
	if err == nil {
		contragent.Total += c.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
	}

	contact, err := ContactGet(c.ContactId, tx)
	if err == nil {
		contact.Total += c.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
	}

	doc := Document{Id: 0, DocType: "cash_in", IsActive: true}
	doc, err = DocumentCreate(doc, tx)
	if err != nil {
		return c, err
	}
	c.DocumentUid = doc.Id

	t := time.Now()
	c.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO cash_in
            (document_uid, name, cash_id, user_id, based_on, cbox_check_id, contragent_id, contact_id, created_at, cash_sum, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.DocumentUid,
		c.Name,
		c.CashId,
		c.UserId,
		c.BasedOn,
		c.CboxCheckId,
		c.ContragentId,
		c.ContactId,
		c.CreatedAt,
		c.CashSum,
		c.Comm,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)
	c.Name = fmt.Sprintf("%s-%d", c.Name, c.Id)

	c, err = CashInUpdate(c, tx)
	if err != nil {
		return c, err
	}

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashInUpdate(c CashIn, tx *sql.Tx) (CashIn, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	cash_in, err := CashInGet(c.Id, tx)
	if err != nil {
		return c, err
	}

	cash, err := CashGet(cash_in.CashId, tx)
	if err == nil {
		cash.Total -= cash_in.CashSum

	}

	if cash_in.CashId != c.CashId {
		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
		cash, err = CashGet(c.CashId, tx)
		if err != nil {
			return c, err
		}
	}
	cash.Total += c.CashSum

	_, err = CashUpdate(cash, tx)
	if err != nil {
		return c, err
	}

	contragent, err := ContragentGet(cash_in.ContragentId, tx)
	if err == nil {
		contragent.Total -= cash_in.CashSum

	}

	if cash_in.ContragentId != c.ContragentId {
		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
		contragent, err = ContragentGet(c.ContragentId, tx)
		if err != nil {
			return c, err
		}
	}
	contragent.Total += c.CashSum

	_, err = ContragentUpdate(contragent, tx)
	if err != nil {
		return c, err
	}

	contact, err := ContactGet(cash_in.ContactId, tx)
	if err == nil {
		contact.Total -= cash_in.CashSum

	}

	if cash_in.ContactId != c.ContactId {
		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
		contact, err = ContactGet(c.ContactId, tx)
		if err != nil {
			return c, err
		}
	}
	contact.Total += c.CashSum

	_, err = ContactUpdate(contact, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE cash_in SET
                    document_uid=?, name=?, cash_id=?, user_id=?, based_on=?, cbox_check_id=?, contragent_id=?, contact_id=?, created_at=?, cash_sum=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.DocumentUid,
		c.Name,
		c.CashId,
		c.UserId,
		c.BasedOn,
		c.CboxCheckId,
		c.ContragentId,
		c.ContactId,
		c.CreatedAt,
		c.CashSum,
		c.Comm,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashInDelete(id int, tx *sql.Tx) (CashIn, error) {
	needCommit := false
	var err error
	var c CashIn
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CashInGet(id, tx)
	if err != nil {
		return c, err
	}

	cash, err := CashGet(c.CashId, tx)
	if err == nil {
		cash.Total -= c.CashSum

		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
	}

	contragent, err := ContragentGet(c.ContragentId, tx)
	if err == nil {
		contragent.Total -= c.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
	}

	contact, err := ContactGet(c.ContactId, tx)
	if err == nil {
		contact.Total -= c.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
	}

	sql := `UPDATE cash_in SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CashInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashIn, error) {

	if !CashInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashIn{}
	for rows.Next() {
		var c CashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashIn, error) {

	if !CashInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashIn{}
	for rows.Next() {
		var c CashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashInTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "cash_id", "user_id", "based_on", "cbox_check_id", "contragent_id", "contact_id", "created_at", "cash_sum", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func CashInGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]CashIn, error) {
	query := "SELECT * FROM cash_in WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashIn{}
	for rows.Next() {
		var c CashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CashInCashSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_in WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func CashInGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_in WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_in WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type CashOut struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	CashId       int     `json:"cash_id"`
	UserId       int     `json:"user_id"`
	BasedOn      int     `json:"based_on"`
	CboxCheckId  int     `json:"cbox_check_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
}

func CashOutGet(id int, tx *sql.Tx) (CashOut, error) {
	var c CashOut
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM cash_out WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM cash_out WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.DocumentUid,
		&c.Name,
		&c.CashId,
		&c.UserId,
		&c.BasedOn,
		&c.CboxCheckId,
		&c.ContragentId,
		&c.ContactId,
		&c.CreatedAt,
		&c.CashSum,
		&c.Comm,
		&c.IsActive,
	)
	return c, err
}

func CashOutGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashOut, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM cash_out"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashOut{}
	for rows.Next() {
		var c CashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CashOutCreate(c CashOut, tx *sql.Tx) (CashOut, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	cash, err := CashGet(c.CashId, tx)
	if err == nil {
		cash.Total -= c.CashSum

		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
	}

	contragent, err := ContragentGet(c.ContragentId, tx)
	if err == nil {
		contragent.Total -= c.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
	}

	contact, err := ContactGet(c.ContactId, tx)
	if err == nil {
		contact.Total -= c.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
	}

	doc := Document{Id: 0, DocType: "cash_out", IsActive: true}
	doc, err = DocumentCreate(doc, tx)
	if err != nil {
		return c, err
	}
	c.DocumentUid = doc.Id

	t := time.Now()
	c.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO cash_out
            (document_uid, name, cash_id, user_id, based_on, cbox_check_id, contragent_id, contact_id, created_at, cash_sum, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.DocumentUid,
		c.Name,
		c.CashId,
		c.UserId,
		c.BasedOn,
		c.CboxCheckId,
		c.ContragentId,
		c.ContactId,
		c.CreatedAt,
		c.CashSum,
		c.Comm,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)
	c.Name = fmt.Sprintf("%s-%d", c.Name, c.Id)

	c, err = CashOutUpdate(c, tx)
	if err != nil {
		return c, err
	}

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashOutUpdate(c CashOut, tx *sql.Tx) (CashOut, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	cash_out, err := CashOutGet(c.Id, tx)
	if err != nil {
		return c, err
	}

	cash, err := CashGet(cash_out.CashId, tx)
	if err == nil {
		cash.Total += cash_out.CashSum

	}

	if cash_out.CashId != c.CashId {
		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
		cash, err = CashGet(c.CashId, tx)
		if err != nil {
			return c, err
		}
	}
	cash.Total -= c.CashSum

	_, err = CashUpdate(cash, tx)
	if err != nil {
		return c, err
	}

	contragent, err := ContragentGet(cash_out.ContragentId, tx)
	if err == nil {
		contragent.Total += cash_out.CashSum

	}

	if cash_out.ContragentId != c.ContragentId {
		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
		contragent, err = ContragentGet(c.ContragentId, tx)
		if err != nil {
			return c, err
		}
	}
	contragent.Total -= c.CashSum

	_, err = ContragentUpdate(contragent, tx)
	if err != nil {
		return c, err
	}

	contact, err := ContactGet(cash_out.ContactId, tx)
	if err == nil {
		contact.Total += cash_out.CashSum

	}

	if cash_out.ContactId != c.ContactId {
		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
		contact, err = ContactGet(c.ContactId, tx)
		if err != nil {
			return c, err
		}
	}
	contact.Total -= c.CashSum

	_, err = ContactUpdate(contact, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE cash_out SET
                    document_uid=?, name=?, cash_id=?, user_id=?, based_on=?, cbox_check_id=?, contragent_id=?, contact_id=?, created_at=?, cash_sum=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.DocumentUid,
		c.Name,
		c.CashId,
		c.UserId,
		c.BasedOn,
		c.CboxCheckId,
		c.ContragentId,
		c.ContactId,
		c.CreatedAt,
		c.CashSum,
		c.Comm,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CashOutDelete(id int, tx *sql.Tx) (CashOut, error) {
	needCommit := false
	var err error
	var c CashOut
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CashOutGet(id, tx)
	if err != nil {
		return c, err
	}

	cash, err := CashGet(c.CashId, tx)
	if err == nil {
		cash.Total += c.CashSum

		_, err = CashUpdate(cash, tx)
		if err != nil {
			return c, err
		}
	}

	contragent, err := ContragentGet(c.ContragentId, tx)
	if err == nil {
		contragent.Total += c.CashSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return c, err
		}
	}

	contact, err := ContactGet(c.ContactId, tx)
	if err == nil {
		contact.Total += c.CashSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return c, err
		}
	}

	sql := `UPDATE cash_out SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CashOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashOut, error) {

	if !CashOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashOut{}
	for rows.Next() {
		var c CashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]CashOut, error) {

	if !CashOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM cash_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashOut{}
	for rows.Next() {
		var c CashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CashOutTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "cash_id", "user_id", "based_on", "cbox_check_id", "contragent_id", "contact_id", "created_at", "cash_sum", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func CashOutGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]CashOut, error) {
	query := "SELECT * FROM cash_out WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []CashOut{}
	for rows.Next() {
		var c CashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CashOutCashSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_out WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func CashOutGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_out WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(cash_sum) FROM cash_out WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type Whs struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Comm     string `json:"comm"`
	IsActive bool   `json:"is_active"`
}

func WhsGet(id int, tx *sql.Tx) (Whs, error) {
	var w Whs
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM whs WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM whs WHERE id=?", id)
	}

	err := row.Scan(
		&w.Id,
		&w.Name,
		&w.Comm,
		&w.IsActive,
	)
	return w, err
}

func WhsGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Whs, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM whs"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Whs{}
	for rows.Next() {
		var w Whs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsCreate(w Whs, tx *sql.Tx) (Whs, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO whs
            (name, comm, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		w.Name,
		w.Comm,
		w.IsActive,
	)
	if err != nil {
		return w, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return w, err
	}
	w.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsUpdate(w Whs, tx *sql.Tx) (Whs, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE whs SET
                    name=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		w.Name,
		w.Comm,
		w.IsActive,
		w.Id,
	)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsDelete(id int, tx *sql.Tx) (Whs, error) {
	needCommit := false
	var err error
	var w Whs
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	w, err = WhsGet(id, tx)
	if err != nil {
		return w, err
	}

	sql := `UPDATE whs SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, w.Id)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	w.IsActive = false
	return w, nil
}

func WhsGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Whs, error) {

	if !WhsTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Whs{}
	for rows.Next() {
		var w Whs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Whs, error) {

	if !WhsTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Whs{}
	for rows.Next() {
		var w Whs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type WhsIn struct {
	Id                  int     `json:"id"`
	DocumentUid         int     `json:"document_uid"`
	Name                string  `json:"name"`
	BasedOn             int     `json:"based_on"`
	WhsId               int     `json:"whs_id"`
	UserId              int     `json:"user_id"`
	ContragentId        int     `json:"contragent_id"`
	ContactId           int     `json:"contact_id"`
	ContragentDocUid    string  `json:"contragent_doc_uid"`
	ContragentCreatedAt string  `json:"contragent_created_at"`
	CreatedAt           string  `json:"created_at"`
	WhsSum              float64 `json:"whs_sum"`
	Delivery            float64 `json:"delivery"`
	Comm                string  `json:"comm"`
	IsActive            bool    `json:"is_active"`
}

func WhsInGet(id int, tx *sql.Tx) (WhsIn, error) {
	var w WhsIn
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM whs_in WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM whs_in WHERE id=?", id)
	}

	err := row.Scan(
		&w.Id,
		&w.DocumentUid,
		&w.Name,
		&w.BasedOn,
		&w.WhsId,
		&w.UserId,
		&w.ContragentId,
		&w.ContactId,
		&w.ContragentDocUid,
		&w.ContragentCreatedAt,
		&w.CreatedAt,
		&w.WhsSum,
		&w.Delivery,
		&w.Comm,
		&w.IsActive,
	)
	return w, err
}

func WhsInGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsIn, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM whs_in"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsIn{}
	for rows.Next() {
		var w WhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsInCreate(w WhsIn, tx *sql.Tx) (WhsIn, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	contragent, err := ContragentGet(w.ContragentId, tx)
	if err == nil {
		contragent.Total += w.WhsSum
		contragent.Total += w.Delivery

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
	}

	contact, err := ContactGet(w.ContactId, tx)
	if err == nil {
		contact.Total += w.WhsSum
		contact.Total += w.Delivery

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
	}

	doc := Document{Id: 0, DocType: "whs_in", IsActive: true}
	doc, err = DocumentCreate(doc, tx)
	if err != nil {
		return w, err
	}
	w.DocumentUid = doc.Id

	t := time.Now()
	w.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO whs_in
            (document_uid, name, based_on, whs_id, user_id, contragent_id, contact_id, contragent_doc_uid, contragent_created_at, created_at, whs_sum, delivery, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		w.DocumentUid,
		w.Name,
		w.BasedOn,
		w.WhsId,
		w.UserId,
		w.ContragentId,
		w.ContactId,
		w.ContragentDocUid,
		w.ContragentCreatedAt,
		w.CreatedAt,
		w.WhsSum,
		w.Delivery,
		w.Comm,
		w.IsActive,
	)
	if err != nil {
		return w, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return w, err
	}
	w.Id = int(last_id)
	w.Name = fmt.Sprintf("%s-%d", w.Name, w.Id)

	w, err = WhsInUpdate(w, tx)
	if err != nil {
		return w, err
	}

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsInUpdate(w WhsIn, tx *sql.Tx) (WhsIn, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	whs_in, err := WhsInGet(w.Id, tx)
	if err != nil {
		return w, err
	}

	contragent, err := ContragentGet(whs_in.ContragentId, tx)
	if err == nil {
		contragent.Total -= whs_in.WhsSum
		contragent.Total -= whs_in.Delivery

	}

	if whs_in.ContragentId != w.ContragentId {
		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
		contragent, err = ContragentGet(w.ContragentId, tx)
		if err != nil {
			return w, err
		}
	}
	contragent.Total += w.WhsSum
	contragent.Total += w.Delivery

	_, err = ContragentUpdate(contragent, tx)
	if err != nil {
		return w, err
	}

	contact, err := ContactGet(whs_in.ContactId, tx)
	if err == nil {
		contact.Total -= whs_in.WhsSum
		contact.Total -= whs_in.Delivery

	}

	if whs_in.ContactId != w.ContactId {
		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
		contact, err = ContactGet(w.ContactId, tx)
		if err != nil {
			return w, err
		}
	}
	contact.Total += w.WhsSum
	contact.Total += w.Delivery

	_, err = ContactUpdate(contact, tx)
	if err != nil {
		return w, err
	}

	sql := `UPDATE whs_in SET
                    document_uid=?, name=?, based_on=?, whs_id=?, user_id=?, contragent_id=?, contact_id=?, contragent_doc_uid=?, contragent_created_at=?, created_at=?, whs_sum=?, delivery=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		w.DocumentUid,
		w.Name,
		w.BasedOn,
		w.WhsId,
		w.UserId,
		w.ContragentId,
		w.ContactId,
		w.ContragentDocUid,
		w.ContragentCreatedAt,
		w.CreatedAt,
		w.WhsSum,
		w.Delivery,
		w.Comm,
		w.IsActive,
		w.Id,
	)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsInDelete(id int, tx *sql.Tx) (WhsIn, error) {
	needCommit := false
	var err error
	var w WhsIn
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	w, err = WhsInGet(id, tx)
	if err != nil {
		return w, err
	}

	contragent, err := ContragentGet(w.ContragentId, tx)
	if err == nil {
		contragent.Total -= w.WhsSum
		contragent.Total -= w.Delivery

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
	}

	contact, err := ContactGet(w.ContactId, tx)
	if err == nil {
		contact.Total -= w.WhsSum
		contact.Total -= w.Delivery

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
	}

	matherial_to_whs_ins, err := MatherialToWhsInGetByFilterInt("whs_in_id", w.Id, false, false, tx)
	if err != nil {
		return w, err
	}
	for _, matherial_to_whs_in := range matherial_to_whs_ins {
		_, err = MatherialToWhsInDelete(matherial_to_whs_in.Id, tx)
		if err != nil {
			return w, err
		}
	}

	cash_outs, err := CashOutGetByFilterInt("based_on", w.DocumentUid, false, false, tx)
	if err != nil {
		return w, err
	}
	for _, cash_out := range cash_outs {
		_, err = CashOutDelete(cash_out.Id, tx)
		if err != nil {
			return w, err
		}
	}

	cash_ins, err := CashInGetByFilterInt("based_on", w.DocumentUid, false, false, tx)
	if err != nil {
		return w, err
	}
	for _, cash_in := range cash_ins {
		_, err = CashInDelete(cash_in.Id, tx)
		if err != nil {
			return w, err
		}
	}

	whs_outs, err := WhsOutGetByFilterInt("based_on", w.DocumentUid, false, false, tx)
	if err != nil {
		return w, err
	}
	for _, whs_out := range whs_outs {
		_, err = WhsOutDelete(whs_out.Id, tx)
		if err != nil {
			return w, err
		}
	}

	sql := `UPDATE whs_in SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, w.Id)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	w.IsActive = false
	return w, nil
}

func WhsInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsIn, error) {

	if !WhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsIn{}
	for rows.Next() {
		var w WhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsIn, error) {

	if !WhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsIn{}
	for rows.Next() {
		var w WhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsInTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "based_on", "whs_id", "user_id", "contragent_id", "contact_id", "contragent_doc_uid", "contragent_created_at", "created_at", "whs_sum", "delivery", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func WhsInGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WhsIn, error) {
	query := "SELECT * FROM whs_in WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsIn{}
	for rows.Next() {
		var w WhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsInGetBetweenContragentCreatedAt(contragent_created_at1, contragent_created_at2 string, withDeleted bool, deletedOnly bool) ([]WhsIn, error) {
	query := "SELECT * FROM whs_in WHERE contragent_created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, contragent_created_at1, contragent_created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsIn{}
	for rows.Next() {
		var w WhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsInWhsSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_in WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func WhsInGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_in WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_in WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type WhsOut struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	BasedOn      int     `json:"based_on"`
	WhsId        int     `json:"whs_id"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	WhsSum       float64 `json:"whs_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
}

func WhsOutGet(id int, tx *sql.Tx) (WhsOut, error) {
	var w WhsOut
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM whs_out WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM whs_out WHERE id=?", id)
	}

	err := row.Scan(
		&w.Id,
		&w.DocumentUid,
		&w.Name,
		&w.BasedOn,
		&w.WhsId,
		&w.UserId,
		&w.ContragentId,
		&w.ContactId,
		&w.CreatedAt,
		&w.WhsSum,
		&w.Comm,
		&w.IsActive,
	)
	return w, err
}

func WhsOutGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsOut, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM whs_out"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsOut{}
	for rows.Next() {
		var w WhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsOutCreate(w WhsOut, tx *sql.Tx) (WhsOut, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	contragent, err := ContragentGet(w.ContragentId, tx)
	if err == nil {
		contragent.Total -= w.WhsSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
	}

	contact, err := ContactGet(w.ContactId, tx)
	if err == nil {
		contact.Total -= w.WhsSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
	}

	doc := Document{Id: 0, DocType: "whs_out", IsActive: true}
	doc, err = DocumentCreate(doc, tx)
	if err != nil {
		return w, err
	}
	w.DocumentUid = doc.Id

	t := time.Now()
	w.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO whs_out
            (document_uid, name, based_on, whs_id, user_id, contragent_id, contact_id, created_at, whs_sum, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		w.DocumentUid,
		w.Name,
		w.BasedOn,
		w.WhsId,
		w.UserId,
		w.ContragentId,
		w.ContactId,
		w.CreatedAt,
		w.WhsSum,
		w.Comm,
		w.IsActive,
	)
	if err != nil {
		return w, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return w, err
	}
	w.Id = int(last_id)
	w.Name = fmt.Sprintf("%s-%d", w.Name, w.Id)

	w, err = WhsOutUpdate(w, tx)
	if err != nil {
		return w, err
	}

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsOutUpdate(w WhsOut, tx *sql.Tx) (WhsOut, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	whs_out, err := WhsOutGet(w.Id, tx)
	if err != nil {
		return w, err
	}

	contragent, err := ContragentGet(whs_out.ContragentId, tx)
	if err == nil {
		contragent.Total += whs_out.WhsSum

	}

	if whs_out.ContragentId != w.ContragentId {
		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
		contragent, err = ContragentGet(w.ContragentId, tx)
		if err != nil {
			return w, err
		}
	}
	contragent.Total -= w.WhsSum

	_, err = ContragentUpdate(contragent, tx)
	if err != nil {
		return w, err
	}

	contact, err := ContactGet(whs_out.ContactId, tx)
	if err == nil {
		contact.Total += whs_out.WhsSum

	}

	if whs_out.ContactId != w.ContactId {
		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
		contact, err = ContactGet(w.ContactId, tx)
		if err != nil {
			return w, err
		}
	}
	contact.Total -= w.WhsSum

	_, err = ContactUpdate(contact, tx)
	if err != nil {
		return w, err
	}

	sql := `UPDATE whs_out SET
                    document_uid=?, name=?, based_on=?, whs_id=?, user_id=?, contragent_id=?, contact_id=?, created_at=?, whs_sum=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		w.DocumentUid,
		w.Name,
		w.BasedOn,
		w.WhsId,
		w.UserId,
		w.ContragentId,
		w.ContactId,
		w.CreatedAt,
		w.WhsSum,
		w.Comm,
		w.IsActive,
		w.Id,
	)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WhsOutDelete(id int, tx *sql.Tx) (WhsOut, error) {
	needCommit := false
	var err error
	var w WhsOut
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	w, err = WhsOutGet(id, tx)
	if err != nil {
		return w, err
	}

	contragent, err := ContragentGet(w.ContragentId, tx)
	if err == nil {
		contragent.Total += w.WhsSum

		_, err = ContragentUpdate(contragent, tx)
		if err != nil {
			return w, err
		}
	}

	contact, err := ContactGet(w.ContactId, tx)
	if err == nil {
		contact.Total += w.WhsSum

		_, err = ContactUpdate(contact, tx)
		if err != nil {
			return w, err
		}
	}

	matherial_to_whs_outs, err := MatherialToWhsOutGetByFilterInt("whs_out_id", w.Id, false, false, tx)
	if err != nil {
		return w, err
	}
	for _, matherial_to_whs_out := range matherial_to_whs_outs {
		_, err = MatherialToWhsOutDelete(matherial_to_whs_out.Id, tx)
		if err != nil {
			return w, err
		}
	}

	sql := `UPDATE whs_out SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, w.Id)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	w.IsActive = false
	return w, nil
}

func WhsOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsOut, error) {

	if !WhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsOut{}
	for rows.Next() {
		var w WhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WhsOut, error) {

	if !WhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM whs_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsOut{}
	for rows.Next() {
		var w WhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WhsOutTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "based_on", "whs_id", "user_id", "contragent_id", "contact_id", "created_at", "whs_sum", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func WhsOutGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WhsOut, error) {
	query := "SELECT * FROM whs_out WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WhsOut{}
	for rows.Next() {
		var w WhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WhsOutWhsSumGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_out WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func WhsOutGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_out WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(whs_sum) FROM whs_out WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type MatherialToWhsIn struct {
	Id               int     `json:"id"`
	MatherialId      int     `json:"matherial_id"`
	ContragentMatUid string  `json:"contragent_mat_uid"`
	WhsInId          int     `json:"whs_in_id"`
	Number           float64 `json:"number"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	Width            float64 `json:"width"`
	Length           float64 `json:"length"`
	ColorId          int     `json:"color_id"`
	IsActive         bool    `json:"is_active"`
}

func MatherialToWhsInGet(id int, tx *sql.Tx) (MatherialToWhsIn, error) {
	var m MatherialToWhsIn
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_to_whs_in WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_to_whs_in WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.ContragentMatUid,
		&m.WhsInId,
		&m.Number,
		&m.Price,
		&m.Cost,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.IsActive,
	)
	return m, err
}

func MatherialToWhsInGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsIn, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_to_whs_in"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsIn{}
	for rows.Next() {
		var m MatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialToWhsInCreate(m MatherialToWhsIn, tx *sql.Tx) (MatherialToWhsIn, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	err = CreateMatherialToWhsInToNumber(&m, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(m.MatherialId, tx)
	if err == nil {
		matherial.Total += m.Number

		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
	}

	color, err := ColorGet(m.ColorId, tx)
	if err == nil {
		color.Total += m.Number

		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
	}

	whs_in, err := WhsInGet(m.WhsInId, tx)
	if err == nil {
		whs_in.WhsSum += m.Cost

		_, err = WhsInUpdate(whs_in, tx)
		if err != nil {
			return m, err
		}
	}

	sql := `INSERT INTO matherial_to_whs_in
            (matherial_id, contragent_mat_uid, whs_in_id, number, price, cost, width, length, color_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.MatherialId,
		m.ContragentMatUid,
		m.WhsInId,
		m.Number,
		m.Price,
		m.Cost,
		m.Width,
		m.Length,
		m.ColorId,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToWhsInUpdate(m MatherialToWhsIn, tx *sql.Tx) (MatherialToWhsIn, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	matherial_to_whs_in, err := MatherialToWhsInGet(m.Id, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(matherial_to_whs_in.MatherialId, tx)
	if err == nil {
		matherial.Total -= matherial_to_whs_in.Number

	}

	if matherial_to_whs_in.MatherialId != m.MatherialId {
		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
		matherial, err = MatherialGet(m.MatherialId, tx)
		if err != nil {
			return m, err
		}
	}
	matherial.Total += m.Number

	_, err = MatherialUpdate(matherial, tx)
	if err != nil {
		return m, err
	}

	color, err := ColorGet(matherial_to_whs_in.ColorId, tx)
	if err == nil {
		color.Total -= matherial_to_whs_in.Number

	}

	if matherial_to_whs_in.ColorId != m.ColorId {
		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
		color, err = ColorGet(m.ColorId, tx)
		if err != nil {
			return m, err
		}
	}
	color.Total += m.Number

	_, err = ColorUpdate(color, tx)
	if err != nil {
		return m, err
	}

	whs_in, err := WhsInGet(matherial_to_whs_in.WhsInId, tx)
	if err == nil {
		whs_in.WhsSum -= matherial_to_whs_in.Cost

	}

	if matherial_to_whs_in.WhsInId != m.WhsInId {
		_, err = WhsInUpdate(whs_in, tx)
		if err != nil {
			return m, err
		}
		whs_in, err = WhsInGet(m.WhsInId, tx)
		if err != nil {
			return m, err
		}
	}
	whs_in.WhsSum += m.Cost

	_, err = WhsInUpdate(whs_in, tx)
	if err != nil {
		return m, err
	}

	err = UpdateMatherialToWhsInToNumber(&m, matherial_to_whs_in.Number, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_to_whs_in SET
                    matherial_id=?, contragent_mat_uid=?, whs_in_id=?, number=?, price=?, cost=?, width=?, length=?, color_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.MatherialId,
		m.ContragentMatUid,
		m.WhsInId,
		m.Number,
		m.Price,
		m.Cost,
		m.Width,
		m.Length,
		m.ColorId,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToWhsInDelete(id int, tx *sql.Tx) (MatherialToWhsIn, error) {
	needCommit := false
	var err error
	var m MatherialToWhsIn
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialToWhsInGet(id, tx)
	if err != nil {
		return m, err
	}

	err = DeleteMatherialToWhsInToNumber(&m, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(m.MatherialId, tx)
	if err == nil {
		matherial.Total -= m.Number

		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
	}

	color, err := ColorGet(m.ColorId, tx)
	if err == nil {
		color.Total -= m.Number

		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
	}

	whs_in, err := WhsInGet(m.WhsInId, tx)
	if err == nil {
		whs_in.WhsSum -= m.Cost

		_, err = WhsInUpdate(whs_in, tx)
		if err != nil {
			return m, err
		}
	}

	sql := `UPDATE matherial_to_whs_in SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialToWhsInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsIn, error) {

	if !MatherialToWhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_whs_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsIn{}
	for rows.Next() {
		var m MatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToWhsInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsIn, error) {

	if !MatherialToWhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_whs_in WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsIn{}
	for rows.Next() {
		var m MatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToWhsInTestForExistingField(fieldName string) bool {
	fields := []string{"id", "matherial_id", "contragent_mat_uid", "whs_in_id", "number", "price", "cost", "width", "length", "color_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type MatherialToWhsOut struct {
	Id          int     `json:"id"`
	MatherialId int     `json:"matherial_id"`
	WhsOutId    int     `json:"whs_out_id"`
	Number      float64 `json:"number"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	Width       float64 `json:"width"`
	Length      float64 `json:"length"`
	ColorId     int     `json:"color_id"`
	IsActive    bool    `json:"is_active"`
}

func MatherialToWhsOutGet(id int, tx *sql.Tx) (MatherialToWhsOut, error) {
	var m MatherialToWhsOut
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_to_whs_out WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_to_whs_out WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.WhsOutId,
		&m.Number,
		&m.Price,
		&m.Cost,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.IsActive,
	)
	return m, err
}

func MatherialToWhsOutGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsOut, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_to_whs_out"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsOut{}
	for rows.Next() {
		var m MatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialToWhsOutCreate(m MatherialToWhsOut, tx *sql.Tx) (MatherialToWhsOut, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	err = CreateMatherialToWhsOutToNumber(&m, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(m.MatherialId, tx)
	if err == nil {
		matherial.Total -= m.Number

		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
	}

	color, err := ColorGet(m.ColorId, tx)
	if err == nil {
		color.Total -= m.Number

		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
	}

	whs_out, err := WhsOutGet(m.WhsOutId, tx)
	if err == nil {
		whs_out.WhsSum += m.Cost

		_, err = WhsOutUpdate(whs_out, tx)
		if err != nil {
			return m, err
		}
	}

	sql := `INSERT INTO matherial_to_whs_out
            (matherial_id, whs_out_id, number, price, cost, width, length, color_id, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.MatherialId,
		m.WhsOutId,
		m.Number,
		m.Price,
		m.Cost,
		m.Width,
		m.Length,
		m.ColorId,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToWhsOutUpdate(m MatherialToWhsOut, tx *sql.Tx) (MatherialToWhsOut, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	matherial_to_whs_out, err := MatherialToWhsOutGet(m.Id, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(matherial_to_whs_out.MatherialId, tx)
	if err == nil {
		matherial.Total += matherial_to_whs_out.Number

	}

	if matherial_to_whs_out.MatherialId != m.MatherialId {
		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
		matherial, err = MatherialGet(m.MatherialId, tx)
		if err != nil {
			return m, err
		}
	}
	matherial.Total -= m.Number

	_, err = MatherialUpdate(matherial, tx)
	if err != nil {
		return m, err
	}

	color, err := ColorGet(matherial_to_whs_out.ColorId, tx)
	if err == nil {
		color.Total += matherial_to_whs_out.Number

	}

	if matherial_to_whs_out.ColorId != m.ColorId {
		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
		color, err = ColorGet(m.ColorId, tx)
		if err != nil {
			return m, err
		}
	}
	color.Total -= m.Number

	_, err = ColorUpdate(color, tx)
	if err != nil {
		return m, err
	}

	whs_out, err := WhsOutGet(matherial_to_whs_out.WhsOutId, tx)
	if err == nil {
		whs_out.WhsSum -= matherial_to_whs_out.Cost

	}

	if matherial_to_whs_out.WhsOutId != m.WhsOutId {
		_, err = WhsOutUpdate(whs_out, tx)
		if err != nil {
			return m, err
		}
		whs_out, err = WhsOutGet(m.WhsOutId, tx)
		if err != nil {
			return m, err
		}
	}
	whs_out.WhsSum += m.Cost

	_, err = WhsOutUpdate(whs_out, tx)
	if err != nil {
		return m, err
	}

	err = UpdateMatherialToWhsOutToNumber(&m, matherial_to_whs_out.Number, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_to_whs_out SET
                    matherial_id=?, whs_out_id=?, number=?, price=?, cost=?, width=?, length=?, color_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.MatherialId,
		m.WhsOutId,
		m.Number,
		m.Price,
		m.Cost,
		m.Width,
		m.Length,
		m.ColorId,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialToWhsOutDelete(id int, tx *sql.Tx) (MatherialToWhsOut, error) {
	needCommit := false
	var err error
	var m MatherialToWhsOut
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialToWhsOutGet(id, tx)
	if err != nil {
		return m, err
	}

	err = DeleteMatherialToWhsOutToNumber(&m, tx)
	if err != nil {
		return m, err
	}

	matherial, err := MatherialGet(m.MatherialId, tx)
	if err == nil {
		matherial.Total += m.Number

		_, err = MatherialUpdate(matherial, tx)
		if err != nil {
			return m, err
		}
	}

	color, err := ColorGet(m.ColorId, tx)
	if err == nil {
		color.Total += m.Number

		_, err = ColorUpdate(color, tx)
		if err != nil {
			return m, err
		}
	}

	whs_out, err := WhsOutGet(m.WhsOutId, tx)
	if err == nil {
		whs_out.WhsSum -= m.Cost

		_, err = WhsOutUpdate(whs_out, tx)
		if err != nil {
			return m, err
		}
	}

	sql := `UPDATE matherial_to_whs_out SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialToWhsOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsOut, error) {

	if !MatherialToWhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_whs_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsOut{}
	for rows.Next() {
		var m MatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToWhsOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialToWhsOut, error) {

	if !MatherialToWhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_to_whs_out WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialToWhsOut{}
	for rows.Next() {
		var m MatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialToWhsOutTestForExistingField(fieldName string) bool {
	fields := []string{"id", "matherial_id", "whs_out_id", "number", "price", "cost", "width", "length", "color_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type MatherialPart struct {
	Id          int     `json:"id"`
	MatherialId int     `json:"matherial_id"`
	PartUid     int     `json:"part_uid"`
	Number      float64 `json:"number"`
	Width       float64 `json:"width"`
	Length      float64 `json:"length"`
	ColorId     int     `json:"color_id"`
	UserId      int     `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	IsRecycle   bool    `json:"is_recycle"`
	IsActive    bool    `json:"is_active"`
}

func MatherialPartGet(id int, tx *sql.Tx) (MatherialPart, error) {
	var m MatherialPart
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_part WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_part WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.PartUid,
		&m.Number,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.UserId,
		&m.CreatedAt,
		&m.IsRecycle,
		&m.IsActive,
	)
	return m, err
}

func MatherialPartGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPart, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_part"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPart{}
	for rows.Next() {
		var m MatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialPartCreate(m MatherialPart, tx *sql.Tx) (MatherialPart, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	t := time.Now()
	m.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO matherial_part
            (matherial_id, part_uid, number, width, length, color_id, user_id, created_at, is_recycle, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.MatherialId,
		m.PartUid,
		m.Number,
		m.Width,
		m.Length,
		m.ColorId,
		m.UserId,
		m.CreatedAt,
		m.IsRecycle,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialPartUpdate(m MatherialPart, tx *sql.Tx) (MatherialPart, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial_part SET
                    matherial_id=?, part_uid=?, number=?, width=?, length=?, color_id=?, user_id=?, created_at=?, is_recycle=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.MatherialId,
		m.PartUid,
		m.Number,
		m.Width,
		m.Length,
		m.ColorId,
		m.UserId,
		m.CreatedAt,
		m.IsRecycle,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialPartDelete(id int, tx *sql.Tx) (MatherialPart, error) {
	needCommit := false
	var err error
	var m MatherialPart
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialPartGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_part SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialPartGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPart, error) {

	if !MatherialPartTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_part WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPart{}
	for rows.Next() {
		var m MatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialPartGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPart, error) {

	if !MatherialPartTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_part WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPart{}
	for rows.Next() {
		var m MatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialPartTestForExistingField(fieldName string) bool {
	fields := []string{"id", "matherial_id", "part_uid", "number", "width", "length", "color_id", "user_id", "created_at", "is_recycle", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func MatherialPartGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]MatherialPart, error) {
	query := "SELECT * FROM matherial_part WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPart{}
	for rows.Next() {
		var m MatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type MatherialPartSlice struct {
	Id              int     `json:"id"`
	MatherialPartId int     `json:"matherial_part_id"`
	UserId          int     `json:"user_id"`
	CreatedAt       string  `json:"created_at"`
	Number          float64 `json:"number"`
	Width           float64 `json:"width"`
	Length          float64 `json:"length"`
	Comm            string  `json:"comm"`
	IsActive        bool    `json:"is_active"`
}

func MatherialPartSliceGet(id int, tx *sql.Tx) (MatherialPartSlice, error) {
	var m MatherialPartSlice
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM matherial_part_slice WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM matherial_part_slice WHERE id=?", id)
	}

	err := row.Scan(
		&m.Id,
		&m.MatherialPartId,
		&m.UserId,
		&m.CreatedAt,
		&m.Number,
		&m.Width,
		&m.Length,
		&m.Comm,
		&m.IsActive,
	)
	return m, err
}

func MatherialPartSliceGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPartSlice, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM matherial_part_slice"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPartSlice{}
	for rows.Next() {
		var m MatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MatherialPartSliceCreate(m MatherialPartSlice, tx *sql.Tx) (MatherialPartSlice, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	t := time.Now()
	m.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO matherial_part_slice
            (matherial_part_id, user_id, created_at, number, width, length, comm, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		m.MatherialPartId,
		m.UserId,
		m.CreatedAt,
		m.Number,
		m.Width,
		m.Length,
		m.Comm,
		m.IsActive,
	)
	if err != nil {
		return m, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialPartSliceUpdate(m MatherialPartSlice, tx *sql.Tx) (MatherialPartSlice, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE matherial_part_slice SET
                    matherial_part_id=?, user_id=?, created_at=?, number=?, width=?, length=?, comm=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		m.MatherialPartId,
		m.UserId,
		m.CreatedAt,
		m.Number,
		m.Width,
		m.Length,
		m.Comm,
		m.IsActive,
		m.Id,
	)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func MatherialPartSliceDelete(id int, tx *sql.Tx) (MatherialPartSlice, error) {
	needCommit := false
	var err error
	var m MatherialPartSlice
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return m, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	m, err = MatherialPartSliceGet(id, tx)
	if err != nil {
		return m, err
	}

	sql := `UPDATE matherial_part_slice SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, m.Id)
	if err != nil {
		return m, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return m, err
		}
	}
	m.IsActive = false
	return m, nil
}

func MatherialPartSliceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPartSlice, error) {

	if !MatherialPartSliceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_part_slice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPartSlice{}
	for rows.Next() {
		var m MatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialPartSliceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]MatherialPartSlice, error) {

	if !MatherialPartSliceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM matherial_part_slice WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPartSlice{}
	for rows.Next() {
		var m MatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func MatherialPartSliceTestForExistingField(fieldName string) bool {
	fields := []string{"id", "matherial_part_id", "user_id", "created_at", "number", "width", "length", "comm", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func MatherialPartSliceGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]MatherialPartSlice, error) {
	query := "SELECT * FROM matherial_part_slice WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []MatherialPartSlice{}
	for rows.Next() {
		var m MatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type ProjectGroup struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ProjectGroupId int    `json:"project_group_id"`
	IsActive       bool   `json:"is_active"`
}

func ProjectGroupGet(id int, tx *sql.Tx) (ProjectGroup, error) {
	var p ProjectGroup
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM project_group WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM project_group WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ProjectGroupId,
		&p.IsActive,
	)
	return p, err
}

func ProjectGroupGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectGroup, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM project_group"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectGroup{}
	for rows.Next() {
		var p ProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProjectGroupCreate(p ProjectGroup, tx *sql.Tx) (ProjectGroup, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO project_group
            (name, project_group_id, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.ProjectGroupId,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectGroupUpdate(p ProjectGroup, tx *sql.Tx) (ProjectGroup, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE project_group SET
                    name=?, project_group_id=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.ProjectGroupId,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectGroupDelete(id int, tx *sql.Tx) (ProjectGroup, error) {
	needCommit := false
	var err error
	var p ProjectGroup
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProjectGroupGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE project_group SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProjectGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectGroup, error) {

	if !ProjectGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectGroup{}
	for rows.Next() {
		var p ProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectGroup, error) {

	if !ProjectGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_group WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectGroup{}
	for rows.Next() {
		var p ProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectGroupTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "project_group_id", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ProjectStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	CodeName string `json:"code_name"`
	IsActive bool   `json:"is_active"`
}

func ProjectStatusGet(id int, tx *sql.Tx) (ProjectStatus, error) {
	var p ProjectStatus
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM project_status WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM project_status WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.CodeName,
		&p.IsActive,
	)
	return p, err
}

func ProjectStatusGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectStatus, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM project_status"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectStatus{}
	for rows.Next() {
		var p ProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProjectStatusCreate(p ProjectStatus, tx *sql.Tx) (ProjectStatus, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO project_status
            (name, code_name, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.CodeName,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectStatusUpdate(p ProjectStatus, tx *sql.Tx) (ProjectStatus, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE project_status SET
                    name=?, code_name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.CodeName,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectStatusDelete(id int, tx *sql.Tx) (ProjectStatus, error) {
	needCommit := false
	var err error
	var p ProjectStatus
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProjectStatusGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE project_status SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProjectStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectStatus, error) {

	if !ProjectStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectStatus{}
	for rows.Next() {
		var p ProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectStatus, error) {

	if !ProjectStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_status WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectStatus{}
	for rows.Next() {
		var p ProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectStatusTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "code_name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type ProjectType struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	DirName  string `json:"dir_name"`
	IsActive bool   `json:"is_active"`
}

func ProjectTypeGet(id int, tx *sql.Tx) (ProjectType, error) {
	var p ProjectType
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM project_type WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM project_type WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.DirName,
		&p.IsActive,
	)
	return p, err
}

func ProjectTypeGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectType, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM project_type"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectType{}
	for rows.Next() {
		var p ProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProjectTypeCreate(p ProjectType, tx *sql.Tx) (ProjectType, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO project_type
            (name, dir_name, is_active)
            VALUES(?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.Name,
		p.DirName,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectTypeUpdate(p ProjectType, tx *sql.Tx) (ProjectType, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE project_type SET
                    name=?, dir_name=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.Name,
		p.DirName,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectTypeDelete(id int, tx *sql.Tx) (ProjectType, error) {
	needCommit := false
	var err error
	var p ProjectType
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProjectTypeGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE project_type SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProjectTypeGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectType, error) {

	if !ProjectTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_type WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectType{}
	for rows.Next() {
		var p ProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectTypeGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]ProjectType, error) {

	if !ProjectTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project_type WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []ProjectType{}
	for rows.Next() {
		var p ProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectTypeTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "dir_name", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type Project struct {
	Id              int     `json:"id"`
	DocumentUid     int     `json:"document_uid"`
	Name            string  `json:"name"`
	ProjectGroupId  int     `json:"project_group_id"`
	UserId          int     `json:"user_id"`
	ContragentId    int     `json:"contragent_id"`
	ContactId       int     `json:"contact_id"`
	Cost            float64 `json:"cost"`
	CashSum         float64 `json:"cash_sum"`
	WhsSum          float64 `json:"whs_sum"`
	ProjectTypeId   int     `json:"project_type_id"`
	TypeDir         string  `json:"type_dir"`
	ProjectStatusId int     `json:"project_status_id"`
	NumberDir       string  `json:"number_dir"`
	Info            string  `json:"info"`
	CreatedAt       string  `json:"created_at"`
	IsInWork        bool    `json:"is_in_work"`
	IsActive        bool    `json:"is_active"`
}

func ProjectGet(id int, tx *sql.Tx) (Project, error) {
	var p Project
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM project WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM project WHERE id=?", id)
	}

	err := row.Scan(
		&p.Id,
		&p.DocumentUid,
		&p.Name,
		&p.ProjectGroupId,
		&p.UserId,
		&p.ContragentId,
		&p.ContactId,
		&p.Cost,
		&p.CashSum,
		&p.WhsSum,
		&p.ProjectTypeId,
		&p.TypeDir,
		&p.ProjectStatusId,
		&p.NumberDir,
		&p.Info,
		&p.CreatedAt,
		&p.IsInWork,
		&p.IsActive,
	)
	return p, err
}

func ProjectGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Project, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM project"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProjectCreate(p Project, tx *sql.Tx) (Project, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	t := time.Now()
	p.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO project
            (document_uid, name, project_group_id, user_id, contragent_id, contact_id, cost, cash_sum, whs_sum, project_type_id, type_dir, project_status_id, number_dir, info, created_at, is_in_work, is_active)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		p.DocumentUid,
		p.Name,
		p.ProjectGroupId,
		p.UserId,
		p.ContragentId,
		p.ContactId,
		p.Cost,
		p.CashSum,
		p.WhsSum,
		p.ProjectTypeId,
		p.TypeDir,
		p.ProjectStatusId,
		p.NumberDir,
		p.Info,
		p.CreatedAt,
		p.IsInWork,
		p.IsActive,
	)
	if err != nil {
		return p, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectUpdate(p Project, tx *sql.Tx) (Project, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE project SET
                    document_uid=?, name=?, project_group_id=?, user_id=?, contragent_id=?, contact_id=?, cost=?, cash_sum=?, whs_sum=?, project_type_id=?, type_dir=?, project_status_id=?, number_dir=?, info=?, created_at=?, is_in_work=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		p.DocumentUid,
		p.Name,
		p.ProjectGroupId,
		p.UserId,
		p.ContragentId,
		p.ContactId,
		p.Cost,
		p.CashSum,
		p.WhsSum,
		p.ProjectTypeId,
		p.TypeDir,
		p.ProjectStatusId,
		p.NumberDir,
		p.Info,
		p.CreatedAt,
		p.IsInWork,
		p.IsActive,
		p.Id,
	)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

func ProjectDelete(id int, tx *sql.Tx) (Project, error) {
	needCommit := false
	var err error
	var p Project
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return p, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	p, err = ProjectGet(id, tx)
	if err != nil {
		return p, err
	}

	sql := `UPDATE project SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, p.Id)
	if err != nil {
		return p, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return p, err
		}
	}
	p.IsActive = false
	return p, nil
}

func ProjectGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Project, error) {

	if !ProjectTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Project, error) {

	if !ProjectTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM project WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func ProjectTestForExistingField(fieldName string) bool {
	fields := []string{"id", "document_uid", "name", "project_group_id", "user_id", "contragent_id", "contact_id", "cost", "cash_sum", "whs_sum", "project_type_id", "type_dir", "project_status_id", "number_dir", "info", "created_at", "is_in_work", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func ProjectFindByProjectInfoContragentNoSearchContactNoSearch(fs string) ([]Project, error) {
	fs = "%" + fs + "%"

	query := `
       SELECT project.* FROM project
        JOIN contragent on contragent.id = project.contragent_id
        JOIN contact on contact.id = project.contact_id
        WHERE project.info LIKE ?
        AND NOT (contragent.search LIKE ? OR contact.search LIKE ?);;`

	rows, err := db.Query(query, fs, fs, fs)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func ProjectGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]Project, error) {
	query := "SELECT * FROM project WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

type Counter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	EquipmentId int    `json:"equipment_id"`
	Total       int    `json:"total"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
}

func CounterGet(id int, tx *sql.Tx) (Counter, error) {
	var c Counter
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM counter WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM counter WHERE id=?", id)
	}

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.EquipmentId,
		&c.Total,
		&c.UpdatedAt,
		&c.IsActive,
	)
	return c, err
}

func CounterGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Counter, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM counter"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Counter{}
	for rows.Next() {
		var c Counter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func CounterCreate(c Counter, tx *sql.Tx) (Counter, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO counter
            (name, equipment_id, total, updated_at, is_active)
            VALUES(?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		c.Name,
		c.EquipmentId,
		c.Total,
		c.UpdatedAt,
		c.IsActive,
	)
	if err != nil {
		return c, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return c, err
	}
	c.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CounterUpdate(c Counter, tx *sql.Tx) (Counter, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	t := time.Now()
	c.UpdatedAt = t.Format("2006-01-02T15:04:05")

	sql := `UPDATE counter SET
                    name=?, equipment_id=?, total=?, updated_at=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		c.Name,
		c.EquipmentId,
		c.Total,
		c.UpdatedAt,
		c.IsActive,
		c.Id,
	)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func CounterDelete(id int, tx *sql.Tx) (Counter, error) {
	needCommit := false
	var err error
	var c Counter
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return c, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	c, err = CounterGet(id, tx)
	if err != nil {
		return c, err
	}

	sql := `UPDATE counter SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, c.Id)
	if err != nil {
		return c, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return c, err
		}
	}
	c.IsActive = false
	return c, nil
}

func CounterGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Counter, error) {

	if !CounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM counter WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Counter{}
	for rows.Next() {
		var c Counter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CounterGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]Counter, error) {

	if !CounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM counter WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Counter{}
	for rows.Next() {
		var c Counter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func CounterTestForExistingField(fieldName string) bool {
	fields := []string{"id", "name", "equipment_id", "total", "updated_at", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type RecordToCounter struct {
	Id        int    `json:"id"`
	CounterId int    `json:"counter_id"`
	CreatedAt string `json:"created_at"`
	Number    int    `json:"number"`
	IsActive  bool   `json:"is_active"`
}

func RecordToCounterGet(id int, tx *sql.Tx) (RecordToCounter, error) {
	var r RecordToCounter
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM record_to_counter WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM record_to_counter WHERE id=?", id)
	}

	err := row.Scan(
		&r.Id,
		&r.CounterId,
		&r.CreatedAt,
		&r.Number,
		&r.IsActive,
	)
	return r, err
}

func RecordToCounterGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]RecordToCounter, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM record_to_counter"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []RecordToCounter{}
	for rows.Next() {
		var r RecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func RecordToCounterCreate(r RecordToCounter, tx *sql.Tx) (RecordToCounter, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return r, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	counter, err := CounterGet(r.CounterId, tx)
	if err == nil {
		counter.Total = r.Number

		_, err = CounterUpdate(counter, tx)
		if err != nil {
			return r, err
		}
	}

	t := time.Now()
	r.CreatedAt = t.Format("2006-01-02T15:04:05")

	sql := `INSERT INTO record_to_counter
            (counter_id, created_at, number, is_active)
            VALUES(?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		r.CounterId,
		r.CreatedAt,
		r.Number,
		r.IsActive,
	)
	if err != nil {
		return r, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return r, err
	}
	r.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func RecordToCounterUpdate(r RecordToCounter, tx *sql.Tx) (RecordToCounter, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return r, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	record_to_counter, err := RecordToCounterGet(r.Id, tx)
	if err != nil {
		return r, err
	}

	counter, err := CounterGet(record_to_counter.CounterId, tx)
	if err == nil {
		counter.Total += record_to_counter.Number

	}

	if record_to_counter.CounterId != r.CounterId {
		_, err = CounterUpdate(counter, tx)
		if err != nil {
			return r, err
		}
		counter, err = CounterGet(r.CounterId, tx)
		if err != nil {
			return r, err
		}
	}
	counter.Total = r.Number

	_, err = CounterUpdate(counter, tx)
	if err != nil {
		return r, err
	}

	sql := `UPDATE record_to_counter SET
                    counter_id=?, created_at=?, number=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		r.CounterId,
		r.CreatedAt,
		r.Number,
		r.IsActive,
		r.Id,
	)
	if err != nil {
		return r, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func RecordToCounterGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]RecordToCounter, error) {

	if !RecordToCounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM record_to_counter WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []RecordToCounter{}
	for rows.Next() {
		var r RecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil

}

func RecordToCounterGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]RecordToCounter, error) {

	if !RecordToCounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM record_to_counter WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []RecordToCounter{}
	for rows.Next() {
		var r RecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil

}

func RecordToCounterTestForExistingField(fieldName string) bool {
	fields := []string{"id", "counter_id", "created_at", "number", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

func RecordToCounterGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]RecordToCounter, error) {
	query := "SELECT * FROM record_to_counter WHERE created_at BETWEEN ? and ?"
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []RecordToCounter{}
	for rows.Next() {
		var r RecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func RecordToCounterNumberGetSumBefore(field string, id int, date string) (map[string]int, error) {
	query := fmt.Sprintf("SELECT SUM(number) FROM record_to_counter WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
	var sum int
	row := db.QueryRow(query, id, date)
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

func RecordToCounterGetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {
	query := ""
	var row *sql.Row
	if field2 == "-" && id2 == 0 {
		query = fmt.Sprintf("SELECT SUM(number) FROM record_to_counter WHERE is_active = 1 AND %s = ?", field)
		row = db.QueryRow(query, id)
	} else {
		query = fmt.Sprintf("SELECT SUM(number) FROM record_to_counter WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
		row = db.QueryRow(query, id, id2)
	}
	var sum int
	err := row.Scan(&sum)
	if err != nil {
		return map[string]int{"sum": 0}, nil
	}
	return map[string]int{"sum": sum}, nil
}

type WmcNumber struct {
	Id          int     `json:"id"`
	WhsId       int     `json:"whs_id"`
	MatherialId int     `json:"matherial_id"`
	ColorId     int     `json:"color_id"`
	Total       float64 `json:"total"`
	IsActive    bool    `json:"is_active"`
}

func WmcNumberGet(id int, tx *sql.Tx) (WmcNumber, error) {
	var w WmcNumber
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM wmc_number WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM wmc_number WHERE id=?", id)
	}

	err := row.Scan(
		&w.Id,
		&w.WhsId,
		&w.MatherialId,
		&w.ColorId,
		&w.Total,
		&w.IsActive,
	)
	return w, err
}

func WmcNumberGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WmcNumber, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM wmc_number"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WmcNumber{}
	for rows.Next() {
		var w WmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WmcNumberCreate(w WmcNumber, tx *sql.Tx) (WmcNumber, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO wmc_number
            (whs_id, matherial_id, color_id, total, is_active)
            VALUES(?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		w.WhsId,
		w.MatherialId,
		w.ColorId,
		w.Total,
		w.IsActive,
	)
	if err != nil {
		return w, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return w, err
	}
	w.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WmcNumberUpdate(w WmcNumber, tx *sql.Tx) (WmcNumber, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE wmc_number SET
                    whs_id=?, matherial_id=?, color_id=?, total=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		w.WhsId,
		w.MatherialId,
		w.ColorId,
		w.Total,
		w.IsActive,
		w.Id,
	)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	return w, nil
}

func WmcNumberDelete(id int, tx *sql.Tx) (WmcNumber, error) {
	needCommit := false
	var err error
	var w WmcNumber
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return w, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	w, err = WmcNumberGet(id, tx)
	if err != nil {
		return w, err
	}

	sql := `UPDATE wmc_number SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, w.Id)
	if err != nil {
		return w, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return w, err
		}
	}
	w.IsActive = false
	return w, nil
}

func WmcNumberGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WmcNumber, error) {

	if !WmcNumberTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM wmc_number WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WmcNumber{}
	for rows.Next() {
		var w WmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WmcNumberGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]WmcNumber, error) {

	if !WmcNumberTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM wmc_number WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WmcNumber{}
	for rows.Next() {
		var w WmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WmcNumberTestForExistingField(fieldName string) bool {
	fields := []string{"id", "whs_id", "matherial_id", "color_id", "total", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type NumbersToProduct struct {
	Id        int     `json:"id"`
	ProductId int     `json:"product_id"`
	Number    float64 `json:"number"`
	Pieces    int     `json:"pieces"`
	Size      float64 `json:"size"`
	Persent   float64 `json:"persent"`
	IsActive  bool    `json:"is_active"`
}

func NumbersToProductGet(id int, tx *sql.Tx) (NumbersToProduct, error) {
	var n NumbersToProduct
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT * FROM numbers_to_product WHERE id=?", id)
	} else {
		row = db.QueryRow("SELECT * FROM numbers_to_product WHERE id=?", id)
	}

	err := row.Scan(
		&n.Id,
		&n.ProductId,
		&n.Number,
		&n.Pieces,
		&n.Size,
		&n.Persent,
		&n.IsActive,
	)
	return n, err
}

func NumbersToProductGetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]NumbersToProduct, error) {
	var rows *sql.Rows
	var err error
	query := "SELECT * FROM numbers_to_product"
	if deletedOnly {
		query += " WHERE is_active = 0"
	} else if !withDeleted {
		query += " WHERE is_active = 1"
	}

	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []NumbersToProduct{}
	for rows.Next() {
		var n NumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func NumbersToProductCreate(n NumbersToProduct, tx *sql.Tx) (NumbersToProduct, error) {
	var err error
	needCommit := false

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return n, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `INSERT INTO numbers_to_product
            (product_id, number, pieces, size, persent, is_active)
            VALUES(?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(
		sql,
		n.ProductId,
		n.Number,
		n.Pieces,
		n.Size,
		n.Persent,
		n.IsActive,
	)
	if err != nil {
		return n, err
	}
	last_id, err := res.LastInsertId()
	if err != nil {
		return n, err
	}
	n.Id = int(last_id)

	if needCommit {
		err = tx.Commit()
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func NumbersToProductUpdate(n NumbersToProduct, tx *sql.Tx) (NumbersToProduct, error) {
	var err error
	needCommit := false
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return n, err
		}
		needCommit = true
		defer tx.Rollback()
	}

	sql := `UPDATE numbers_to_product SET
                    product_id=?, number=?, pieces=?, size=?, persent=?, is_active=?
                    WHERE id=?;`

	_, err = tx.Exec(
		sql,
		n.ProductId,
		n.Number,
		n.Pieces,
		n.Size,
		n.Persent,
		n.IsActive,
		n.Id,
	)
	if err != nil {
		return n, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func NumbersToProductDelete(id int, tx *sql.Tx) (NumbersToProduct, error) {
	needCommit := false
	var err error
	var n NumbersToProduct
	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return n, err
		}
		needCommit = true
		defer tx.Rollback()
	}
	n, err = NumbersToProductGet(id, tx)
	if err != nil {
		return n, err
	}

	sql := `UPDATE numbers_to_product SET is_active=0 WHERE id=?;`

	_, err = tx.Exec(sql, n.Id)
	if err != nil {
		return n, err
	}
	if needCommit {
		err = tx.Commit()
		if err != nil {
			return n, err
		}
	}
	n.IsActive = false
	return n, nil
}

func NumbersToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]NumbersToProduct, error) {

	if !NumbersToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM numbers_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []NumbersToProduct{}
	for rows.Next() {
		var n NumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil

}

func NumbersToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]NumbersToProduct, error) {

	if !NumbersToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	var err error
	query := fmt.Sprintf("SELECT * FROM numbers_to_product WHERE %s=?", field)
	if deletedOnly {
		query += "  AND is_active = 0"
	} else if !withDeleted {
		query += "  AND is_active = 1"
	}

	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.Query(query, param)
	} else {
		rows, err = db.Query(query, param)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []NumbersToProduct{}
	for rows.Next() {
		var n NumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil

}

func NumbersToProductTestForExistingField(fieldName string) bool {
	fields := []string{"id", "product_id", "number", "pieces", "size", "persent", "is_active"}
	for _, f := range fields {
		if fieldName == f {
			return true
		}
	}
	return false
}

type WDocument struct {
	Id       int    `json:"id"`
	DocType  string `json:"doc_type"`
	IsActive bool   `json:"is_active"`
}

func WDocumentGet(id int) (WDocument, error) {
	var d WDocument
	row := db.QueryRow(`SELECT document.* FROM document WHERE document.id=?`, id)
	err := row.Scan(
		&d.Id,
		&d.DocType,
		&d.IsActive,
	)
	return d, err
}

func WDocumentGetAll(withDeleted bool, deletedOnly bool) ([]WDocument, error) {
	query := `SELECT document.* FROM document`
	if deletedOnly {
		query += "  WHERE document.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE document.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WDocument{}
	for rows.Next() {
		var d WDocument
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

func WDocumentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WDocument, error) {

	if !DocumentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT document.* FROM document WHERE document.%s=?`, field)
	if deletedOnly {
		query += "  AND document.is_active = 0"
	} else if !withDeleted {
		query += "  AND document.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WDocument{}
	for rows.Next() {
		var d WDocument
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil

}

func WDocumentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WDocument, error) {

	if !DocumentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT document.* FROM document WHERE document.%s=?`, field)
	if deletedOnly {
		query += "  AND document.is_active = 0"
	} else if !withDeleted {
		query += "  AND document.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WDocument{}
	for rows.Next() {
		var d WDocument
		if err := rows.Scan(
			&d.Id,
			&d.DocType,
			&d.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil

}

type WMeasure struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

func WMeasureGet(id int) (WMeasure, error) {
	var m WMeasure
	row := db.QueryRow(`SELECT measure.* FROM measure WHERE measure.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.FullName,
		&m.IsActive,
	)
	return m, err
}

func WMeasureGetAll(withDeleted bool, deletedOnly bool) ([]WMeasure, error) {
	query := `SELECT measure.* FROM measure`
	if deletedOnly {
		query += "  WHERE measure.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE measure.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMeasure{}
	for rows.Next() {
		var m WMeasure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMeasureGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMeasure, error) {

	if !MeasureTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT measure.* FROM measure WHERE measure.%s=?`, field)
	if deletedOnly {
		query += "  AND measure.is_active = 0"
	} else if !withDeleted {
		query += "  AND measure.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMeasure{}
	for rows.Next() {
		var m WMeasure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMeasureGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMeasure, error) {

	if !MeasureTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT measure.* FROM measure WHERE measure.%s=?`, field)
	if deletedOnly {
		query += "  AND measure.is_active = 0"
	} else if !withDeleted {
		query += "  AND measure.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMeasure{}
	for rows.Next() {
		var m WMeasure
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WCountType struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func WCountTypeGet(id int) (WCountType, error) {
	var c WCountType
	row := db.QueryRow(`SELECT count_type.* FROM count_type WHERE count_type.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.IsActive,
	)
	return c, err
}

func WCountTypeGetAll(withDeleted bool, deletedOnly bool) ([]WCountType, error) {
	query := `SELECT count_type.* FROM count_type`
	if deletedOnly {
		query += "  WHERE count_type.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE count_type.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCountType{}
	for rows.Next() {
		var c WCountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCountTypeGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCountType, error) {

	if !CountTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT count_type.* FROM count_type WHERE count_type.%s=?`, field)
	if deletedOnly {
		query += "  AND count_type.is_active = 0"
	} else if !withDeleted {
		query += "  AND count_type.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCountType{}
	for rows.Next() {
		var c WCountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCountTypeGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCountType, error) {

	if !CountTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT count_type.* FROM count_type WHERE count_type.%s=?`, field)
	if deletedOnly {
		query += "  AND count_type.is_active = 0"
	} else if !withDeleted {
		query += "  AND count_type.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCountType{}
	for rows.Next() {
		var c WCountType
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WColorGroup struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ColorGroupId int    `json:"color_group_id"`
	IsActive     bool   `json:"is_active"`
	ColorGroup   string `json:"color_group"`
}

func WColorGroupGet(id int) (WColorGroup, error) {
	var c WColorGroup
	row := db.QueryRow(`SELECT color_group.*, IFNULL(co.name, "") FROM color_group
	LEFT JOIN color_group AS co ON color_group.color_group_id = co.id WHERE color_group.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ColorGroupId,
		&c.IsActive,
		&c.ColorGroup,
	)
	return c, err
}

func WColorGroupGetAll(withDeleted bool, deletedOnly bool) ([]WColorGroup, error) {
	query := `SELECT color_group.*, IFNULL(co.name, "") FROM color_group
	LEFT JOIN color_group AS co ON color_group.color_group_id = co.id`
	if deletedOnly {
		query += "  WHERE color_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE color_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColorGroup{}
	for rows.Next() {
		var c WColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WColorGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WColorGroup, error) {

	if !ColorGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT color_group.*, IFNULL(co.name, "") FROM color_group
	LEFT JOIN color_group AS co ON color_group.color_group_id = co.id WHERE color_group.%s=?`, field)
	if deletedOnly {
		query += "  AND color_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND color_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColorGroup{}
	for rows.Next() {
		var c WColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WColorGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WColorGroup, error) {

	if !ColorGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT color_group.*, IFNULL(co.name, "") FROM color_group
	LEFT JOIN color_group AS co ON color_group.color_group_id = co.id WHERE color_group.%s=?`, field)
	if deletedOnly {
		query += "  AND color_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND color_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColorGroup{}
	for rows.Next() {
		var c WColorGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ColorGroupId,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WColor struct {
	Id           int     `json:"id"`
	ColorGroupId int     `json:"color_group_id"`
	Name         string  `json:"name"`
	Total        float64 `json:"total"`
	IsActive     bool    `json:"is_active"`
	ColorGroup   string  `json:"color_group"`
}

func WColorGet(id int) (WColor, error) {
	var c WColor
	row := db.QueryRow(`SELECT color.*, IFNULL(color_group.name, "") FROM color
	LEFT JOIN color_group ON color.color_group_id = color_group.id WHERE color.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.ColorGroupId,
		&c.Name,
		&c.Total,
		&c.IsActive,
		&c.ColorGroup,
	)
	return c, err
}

func WColorGetAll(withDeleted bool, deletedOnly bool) ([]WColor, error) {
	query := `SELECT color.*, IFNULL(color_group.name, "") FROM color
	LEFT JOIN color_group ON color.color_group_id = color_group.id`
	if deletedOnly {
		query += "  WHERE color.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE color.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColor{}
	for rows.Next() {
		var c WColor
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WColorGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WColor, error) {

	if !ColorTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT color.*, IFNULL(color_group.name, "") FROM color
	LEFT JOIN color_group ON color.color_group_id = color_group.id WHERE color.%s=?`, field)
	if deletedOnly {
		query += "  AND color.is_active = 0"
	} else if !withDeleted {
		query += "  AND color.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColor{}
	for rows.Next() {
		var c WColor
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WColorGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WColor, error) {

	if !ColorTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT color.*, IFNULL(color_group.name, "") FROM color
	LEFT JOIN color_group ON color.color_group_id = color_group.id WHERE color.%s=?`, field)
	if deletedOnly {
		query += "  AND color.is_active = 0"
	} else if !withDeleted {
		query += "  AND color.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WColor{}
	for rows.Next() {
		var c WColor
		if err := rows.Scan(
			&c.Id,
			&c.ColorGroupId,
			&c.Name,
			&c.Total,
			&c.IsActive,
			&c.ColorGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WMatherialGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	MatherialGroupId int    `json:"matherial_group_id"`
	IsActive         bool   `json:"is_active"`
	MatherialGroup   string `json:"matherial_group"`
}

func WMatherialGroupGet(id int) (WMatherialGroup, error) {
	var m WMatherialGroup
	row := db.QueryRow(`SELECT matherial_group.*, IFNULL(ma.name, "") FROM matherial_group
	LEFT JOIN matherial_group AS ma ON matherial_group.matherial_group_id = ma.id WHERE matherial_group.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.MatherialGroupId,
		&m.IsActive,
		&m.MatherialGroup,
	)
	return m, err
}

func WMatherialGroupGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialGroup, error) {
	query := `SELECT matherial_group.*, IFNULL(ma.name, "") FROM matherial_group
	LEFT JOIN matherial_group AS ma ON matherial_group.matherial_group_id = ma.id`
	if deletedOnly {
		query += "  WHERE matherial_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialGroup{}
	for rows.Next() {
		var m WMatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
			&m.MatherialGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialGroup, error) {

	if !MatherialGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_group.*, IFNULL(ma.name, "") FROM matherial_group
	LEFT JOIN matherial_group AS ma ON matherial_group.matherial_group_id = ma.id WHERE matherial_group.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialGroup{}
	for rows.Next() {
		var m WMatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
			&m.MatherialGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialGroup, error) {

	if !MatherialGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_group.*, IFNULL(ma.name, "") FROM matherial_group
	LEFT JOIN matherial_group AS ma ON matherial_group.matherial_group_id = ma.id WHERE matherial_group.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialGroup{}
	for rows.Next() {
		var m WMatherialGroup
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.MatherialGroupId,
			&m.IsActive,
			&m.MatherialGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WMatherial struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	MatherialGroupId int     `json:"matherial_group_id"`
	MeasureId        int     `json:"measure_id"`
	ColorGroupId     int     `json:"color_group_id"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	Total            float64 `json:"total"`
	Barcode          string  `json:"barcode"`
	CountTypeId      int     `json:"count_type_id"`
	IsActive         bool    `json:"is_active"`
	MatherialGroup   string  `json:"matherial_group"`
	Measure          string  `json:"measure"`
	ColorGroup       string  `json:"color_group"`
	CountType        string  `json:"count_type"`
}

func WMatherialGet(id int) (WMatherial, error) {
	var m WMatherial
	row := db.QueryRow(`SELECT matherial.*, IFNULL(matherial_group.name, ""), IFNULL(measure.name, ""), IFNULL(color_group.name, ""), IFNULL(count_type.name, "") FROM matherial
	LEFT JOIN matherial_group ON matherial.matherial_group_id = matherial_group.id
	LEFT JOIN measure ON matherial.measure_id = measure.id
	LEFT JOIN color_group ON matherial.color_group_id = color_group.id
	LEFT JOIN count_type ON matherial.count_type_id = count_type.id WHERE matherial.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.Name,
		&m.FullName,
		&m.MatherialGroupId,
		&m.MeasureId,
		&m.ColorGroupId,
		&m.Price,
		&m.Cost,
		&m.Total,
		&m.Barcode,
		&m.CountTypeId,
		&m.IsActive,
		&m.MatherialGroup,
		&m.Measure,
		&m.ColorGroup,
		&m.CountType,
	)
	return m, err
}

func WMatherialGetAll(withDeleted bool, deletedOnly bool) ([]WMatherial, error) {
	query := `SELECT matherial.*, IFNULL(matherial_group.name, ""), IFNULL(measure.name, ""), IFNULL(color_group.name, ""), IFNULL(count_type.name, "") FROM matherial
	LEFT JOIN matherial_group ON matherial.matherial_group_id = matherial_group.id
	LEFT JOIN measure ON matherial.measure_id = measure.id
	LEFT JOIN color_group ON matherial.color_group_id = color_group.id
	LEFT JOIN count_type ON matherial.count_type_id = count_type.id`
	if deletedOnly {
		query += "  WHERE matherial.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherial{}
	for rows.Next() {
		var m WMatherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
			&m.MatherialGroup,
			&m.Measure,
			&m.ColorGroup,
			&m.CountType,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherial, error) {

	if !MatherialTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial.*, IFNULL(matherial_group.name, ""), IFNULL(measure.name, ""), IFNULL(color_group.name, ""), IFNULL(count_type.name, "") FROM matherial
	LEFT JOIN matherial_group ON matherial.matherial_group_id = matherial_group.id
	LEFT JOIN measure ON matherial.measure_id = measure.id
	LEFT JOIN color_group ON matherial.color_group_id = color_group.id
	LEFT JOIN count_type ON matherial.count_type_id = count_type.id WHERE matherial.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherial{}
	for rows.Next() {
		var m WMatherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
			&m.MatherialGroup,
			&m.Measure,
			&m.ColorGroup,
			&m.CountType,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherial, error) {

	if !MatherialTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial.*, IFNULL(matherial_group.name, ""), IFNULL(measure.name, ""), IFNULL(color_group.name, ""), IFNULL(count_type.name, "") FROM matherial
	LEFT JOIN matherial_group ON matherial.matherial_group_id = matherial_group.id
	LEFT JOIN measure ON matherial.measure_id = measure.id
	LEFT JOIN color_group ON matherial.color_group_id = color_group.id
	LEFT JOIN count_type ON matherial.count_type_id = count_type.id WHERE matherial.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherial{}
	for rows.Next() {
		var m WMatherial
		if err := rows.Scan(
			&m.Id,
			&m.Name,
			&m.FullName,
			&m.MatherialGroupId,
			&m.MeasureId,
			&m.ColorGroupId,
			&m.Price,
			&m.Cost,
			&m.Total,
			&m.Barcode,
			&m.CountTypeId,
			&m.IsActive,
			&m.MatherialGroup,
			&m.Measure,
			&m.ColorGroup,
			&m.CountType,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WCash struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Persent   float64 `json:"persent"`
	Total     float64 `json:"total"`
	Comm      string  `json:"comm"`
	IsFiscal  bool    `json:"is_fiscal"`
	IsAccount bool    `json:"is_account"`
	IsActive  bool    `json:"is_active"`
}

func WCashGet(id int) (WCash, error) {
	var c WCash
	row := db.QueryRow(`SELECT cash.* FROM cash WHERE cash.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Persent,
		&c.Total,
		&c.Comm,
		&c.IsFiscal,
		&c.IsAccount,
		&c.IsActive,
	)
	return c, err
}

func WCashGetAll(withDeleted bool, deletedOnly bool) ([]WCash, error) {
	query := `SELECT cash.* FROM cash`
	if deletedOnly {
		query += "  WHERE cash.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE cash.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCash{}
	for rows.Next() {
		var c WCash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCashGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCash, error) {

	if !CashTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash.* FROM cash WHERE cash.%s=?`, field)
	if deletedOnly {
		query += "  AND cash.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCash{}
	for rows.Next() {
		var c WCash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCashGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCash, error) {

	if !CashTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash.* FROM cash WHERE cash.%s=?`, field)
	if deletedOnly {
		query += "  AND cash.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCash{}
	for rows.Next() {
		var c WCash
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Persent,
			&c.Total,
			&c.Comm,
			&c.IsFiscal,
			&c.IsAccount,
			&c.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WUserGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UserGroupId int    `json:"user_group_id"`
	IsActive    bool   `json:"is_active"`
	UserGroup   string `json:"user_group"`
}

func WUserGroupGet(id int) (WUserGroup, error) {
	var u WUserGroup
	row := db.QueryRow(`SELECT user_group.*, IFNULL(us.name, "") FROM user_group
	LEFT JOIN user_group AS us ON user_group.user_group_id = us.id WHERE user_group.id=?`, id)
	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.UserGroupId,
		&u.IsActive,
		&u.UserGroup,
	)
	return u, err
}

func WUserGroupGetAll(withDeleted bool, deletedOnly bool) ([]WUserGroup, error) {
	query := `SELECT user_group.*, IFNULL(us.name, "") FROM user_group
	LEFT JOIN user_group AS us ON user_group.user_group_id = us.id`
	if deletedOnly {
		query += "  WHERE user_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE user_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUserGroup{}
	for rows.Next() {
		var u WUserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
			&u.UserGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func WUserGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WUserGroup, error) {

	if !UserGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT user_group.*, IFNULL(us.name, "") FROM user_group
	LEFT JOIN user_group AS us ON user_group.user_group_id = us.id WHERE user_group.%s=?`, field)
	if deletedOnly {
		query += "  AND user_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND user_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUserGroup{}
	for rows.Next() {
		var u WUserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
			&u.UserGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func WUserGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WUserGroup, error) {

	if !UserGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT user_group.*, IFNULL(us.name, "") FROM user_group
	LEFT JOIN user_group AS us ON user_group.user_group_id = us.id WHERE user_group.%s=?`, field)
	if deletedOnly {
		query += "  AND user_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND user_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUserGroup{}
	for rows.Next() {
		var u WUserGroup
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.UserGroupId,
			&u.IsActive,
			&u.UserGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

type WUser struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	UserGroupId int    `json:"user_group_id"`
	CashId      int    `json:"cash_id"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Comm        string `json:"comm"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	BaseAccess  int    `json:"base_access"`
	AddAccess   int    `json:"add_access"`
	IsActive    bool   `json:"is_active"`
	UserGroup   string `json:"user_group"`
	Cash        string `json:"cash"`
}

func WUserGet(id int) (WUser, error) {
	var u WUser
	row := db.QueryRow(`SELECT user.*, IFNULL(user_group.name, ""), IFNULL(cash.name, "") FROM user
	LEFT JOIN user_group ON user.user_group_id = user_group.id
	LEFT JOIN cash ON user.cash_id = cash.id WHERE user.id=?`, id)
	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.FullName,
		&u.UserGroupId,
		&u.CashId,
		&u.Phone,
		&u.Email,
		&u.Comm,
		&u.Login,
		&u.Password,
		&u.BaseAccess,
		&u.AddAccess,
		&u.IsActive,
		&u.UserGroup,
		&u.Cash,
	)
	return u, err
}

func WUserGetAll(withDeleted bool, deletedOnly bool) ([]WUser, error) {
	query := `SELECT user.*, IFNULL(user_group.name, ""), IFNULL(cash.name, "") FROM user
	LEFT JOIN user_group ON user.user_group_id = user_group.id
	LEFT JOIN cash ON user.cash_id = cash.id`
	if deletedOnly {
		query += "  WHERE user.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE user.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUser{}
	for rows.Next() {
		var u WUser
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
			&u.UserGroup,
			&u.Cash,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func WUserGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WUser, error) {

	if !UserTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT user.*, IFNULL(user_group.name, ""), IFNULL(cash.name, "") FROM user
	LEFT JOIN user_group ON user.user_group_id = user_group.id
	LEFT JOIN cash ON user.cash_id = cash.id WHERE user.%s=?`, field)
	if deletedOnly {
		query += "  AND user.is_active = 0"
	} else if !withDeleted {
		query += "  AND user.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUser{}
	for rows.Next() {
		var u WUser
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
			&u.UserGroup,
			&u.Cash,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

func WUserGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WUser, error) {

	if !UserTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT user.*, IFNULL(user_group.name, ""), IFNULL(cash.name, "") FROM user
	LEFT JOIN user_group ON user.user_group_id = user_group.id
	LEFT JOIN cash ON user.cash_id = cash.id WHERE user.%s=?`, field)
	if deletedOnly {
		query += "  AND user.is_active = 0"
	} else if !withDeleted {
		query += "  AND user.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WUser{}
	for rows.Next() {
		var u WUser
		if err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.FullName,
			&u.UserGroupId,
			&u.CashId,
			&u.Phone,
			&u.Email,
			&u.Comm,
			&u.Login,
			&u.Password,
			&u.BaseAccess,
			&u.AddAccess,
			&u.IsActive,
			&u.UserGroup,
			&u.Cash,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil

}

type WEquipmentGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	EquipmentGroupId int    `json:"equipment_group_id"`
	IsActive         bool   `json:"is_active"`
	EquipmentGroup   string `json:"equipment_group"`
}

func WEquipmentGroupGet(id int) (WEquipmentGroup, error) {
	var e WEquipmentGroup
	row := db.QueryRow(`SELECT equipment_group.*, IFNULL(eq.name, "") FROM equipment_group
	LEFT JOIN equipment_group AS eq ON equipment_group.equipment_group_id = eq.id WHERE equipment_group.id=?`, id)
	err := row.Scan(
		&e.Id,
		&e.Name,
		&e.EquipmentGroupId,
		&e.IsActive,
		&e.EquipmentGroup,
	)
	return e, err
}

func WEquipmentGroupGetAll(withDeleted bool, deletedOnly bool) ([]WEquipmentGroup, error) {
	query := `SELECT equipment_group.*, IFNULL(eq.name, "") FROM equipment_group
	LEFT JOIN equipment_group AS eq ON equipment_group.equipment_group_id = eq.id`
	if deletedOnly {
		query += "  WHERE equipment_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE equipment_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipmentGroup{}
	for rows.Next() {
		var e WEquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func WEquipmentGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WEquipmentGroup, error) {

	if !EquipmentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT equipment_group.*, IFNULL(eq.name, "") FROM equipment_group
	LEFT JOIN equipment_group AS eq ON equipment_group.equipment_group_id = eq.id WHERE equipment_group.%s=?`, field)
	if deletedOnly {
		query += "  AND equipment_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND equipment_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipmentGroup{}
	for rows.Next() {
		var e WEquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func WEquipmentGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WEquipmentGroup, error) {

	if !EquipmentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT equipment_group.*, IFNULL(eq.name, "") FROM equipment_group
	LEFT JOIN equipment_group AS eq ON equipment_group.equipment_group_id = eq.id WHERE equipment_group.%s=?`, field)
	if deletedOnly {
		query += "  AND equipment_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND equipment_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipmentGroup{}
	for rows.Next() {
		var e WEquipmentGroup
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.EquipmentGroupId,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

type WEquipment struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	EquipmentGroupId int     `json:"equipment_group_id"`
	Cost             float64 `json:"cost"`
	Total            float64 `json:"total"`
	IsActive         bool    `json:"is_active"`
	EquipmentGroup   string  `json:"equipment_group"`
}

func WEquipmentGet(id int) (WEquipment, error) {
	var e WEquipment
	row := db.QueryRow(`SELECT equipment.*, IFNULL(equipment_group.name, "") FROM equipment
	LEFT JOIN equipment_group ON equipment.equipment_group_id = equipment_group.id WHERE equipment.id=?`, id)
	err := row.Scan(
		&e.Id,
		&e.Name,
		&e.FullName,
		&e.EquipmentGroupId,
		&e.Cost,
		&e.Total,
		&e.IsActive,
		&e.EquipmentGroup,
	)
	return e, err
}

func WEquipmentGetAll(withDeleted bool, deletedOnly bool) ([]WEquipment, error) {
	query := `SELECT equipment.*, IFNULL(equipment_group.name, "") FROM equipment
	LEFT JOIN equipment_group ON equipment.equipment_group_id = equipment_group.id`
	if deletedOnly {
		query += "  WHERE equipment.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE equipment.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipment{}
	for rows.Next() {
		var e WEquipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func WEquipmentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WEquipment, error) {

	if !EquipmentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT equipment.*, IFNULL(equipment_group.name, "") FROM equipment
	LEFT JOIN equipment_group ON equipment.equipment_group_id = equipment_group.id WHERE equipment.%s=?`, field)
	if deletedOnly {
		query += "  AND equipment.is_active = 0"
	} else if !withDeleted {
		query += "  AND equipment.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipment{}
	for rows.Next() {
		var e WEquipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

func WEquipmentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WEquipment, error) {

	if !EquipmentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT equipment.*, IFNULL(equipment_group.name, "") FROM equipment
	LEFT JOIN equipment_group ON equipment.equipment_group_id = equipment_group.id WHERE equipment.%s=?`, field)
	if deletedOnly {
		query += "  AND equipment.is_active = 0"
	} else if !withDeleted {
		query += "  AND equipment.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WEquipment{}
	for rows.Next() {
		var e WEquipment
		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.FullName,
			&e.EquipmentGroupId,
			&e.Cost,
			&e.Total,
			&e.IsActive,
			&e.EquipmentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil

}

type WOperationGroup struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	OperationGroupId int    `json:"operation_group_id"`
	IsActive         bool   `json:"is_active"`
	OperationGroup   string `json:"operation_group"`
}

func WOperationGroupGet(id int) (WOperationGroup, error) {
	var o WOperationGroup
	row := db.QueryRow(`SELECT operation_group.*, IFNULL(op.name, "") FROM operation_group
	LEFT JOIN operation_group AS op ON operation_group.operation_group_id = op.id WHERE operation_group.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.OperationGroupId,
		&o.IsActive,
		&o.OperationGroup,
	)
	return o, err
}

func WOperationGroupGetAll(withDeleted bool, deletedOnly bool) ([]WOperationGroup, error) {
	query := `SELECT operation_group.*, IFNULL(op.name, "") FROM operation_group
	LEFT JOIN operation_group AS op ON operation_group.operation_group_id = op.id`
	if deletedOnly {
		query += "  WHERE operation_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE operation_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationGroup{}
	for rows.Next() {
		var o WOperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
			&o.OperationGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOperationGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOperationGroup, error) {

	if !OperationGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_group.*, IFNULL(op.name, "") FROM operation_group
	LEFT JOIN operation_group AS op ON operation_group.operation_group_id = op.id WHERE operation_group.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationGroup{}
	for rows.Next() {
		var o WOperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
			&o.OperationGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOperationGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOperationGroup, error) {

	if !OperationGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_group.*, IFNULL(op.name, "") FROM operation_group
	LEFT JOIN operation_group AS op ON operation_group.operation_group_id = op.id WHERE operation_group.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationGroup{}
	for rows.Next() {
		var o WOperationGroup
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.OperationGroupId,
			&o.IsActive,
			&o.OperationGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

type WOperation struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	OperationGroupId int     `json:"operation_group_id"`
	MeasureId        int     `json:"measure_id"`
	UserId           int     `json:"user_id"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	EquipmentId      int     `json:"equipment_id"`
	EquipmentPrice   float64 `json:"equipment_price"`
	Barcode          string  `json:"barcode"`
	IsActive         bool    `json:"is_active"`
	OperationGroup   string  `json:"operation_group"`
	Measure          string  `json:"measure"`
	User             string  `json:"user"`
	Equipment        string  `json:"equipment"`
}

func WOperationGet(id int) (WOperation, error) {
	var o WOperation
	row := db.QueryRow(`SELECT operation.*, IFNULL(operation_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation
	LEFT JOIN operation_group ON operation.operation_group_id = operation_group.id
	LEFT JOIN measure ON operation.measure_id = measure.id
	LEFT JOIN user ON operation.user_id = user.id
	LEFT JOIN equipment ON operation.equipment_id = equipment.id WHERE operation.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.FullName,
		&o.OperationGroupId,
		&o.MeasureId,
		&o.UserId,
		&o.Price,
		&o.Cost,
		&o.EquipmentId,
		&o.EquipmentPrice,
		&o.Barcode,
		&o.IsActive,
		&o.OperationGroup,
		&o.Measure,
		&o.User,
		&o.Equipment,
	)
	return o, err
}

func WOperationGetAll(withDeleted bool, deletedOnly bool) ([]WOperation, error) {
	query := `SELECT operation.*, IFNULL(operation_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation
	LEFT JOIN operation_group ON operation.operation_group_id = operation_group.id
	LEFT JOIN measure ON operation.measure_id = measure.id
	LEFT JOIN user ON operation.user_id = user.id
	LEFT JOIN equipment ON operation.equipment_id = equipment.id`
	if deletedOnly {
		query += "  WHERE operation.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE operation.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperation{}
	for rows.Next() {
		var o WOperation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
			&o.OperationGroup,
			&o.Measure,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOperationGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOperation, error) {

	if !OperationTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation.*, IFNULL(operation_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation
	LEFT JOIN operation_group ON operation.operation_group_id = operation_group.id
	LEFT JOIN measure ON operation.measure_id = measure.id
	LEFT JOIN user ON operation.user_id = user.id
	LEFT JOIN equipment ON operation.equipment_id = equipment.id WHERE operation.%s=?`, field)
	if deletedOnly {
		query += "  AND operation.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperation{}
	for rows.Next() {
		var o WOperation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
			&o.OperationGroup,
			&o.Measure,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOperationGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOperation, error) {

	if !OperationTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation.*, IFNULL(operation_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation
	LEFT JOIN operation_group ON operation.operation_group_id = operation_group.id
	LEFT JOIN measure ON operation.measure_id = measure.id
	LEFT JOIN user ON operation.user_id = user.id
	LEFT JOIN equipment ON operation.equipment_id = equipment.id WHERE operation.%s=?`, field)
	if deletedOnly {
		query += "  AND operation.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperation{}
	for rows.Next() {
		var o WOperation
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.FullName,
			&o.OperationGroupId,
			&o.MeasureId,
			&o.UserId,
			&o.Price,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentPrice,
			&o.Barcode,
			&o.IsActive,
			&o.OperationGroup,
			&o.Measure,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

type WProductGroup struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ProductGroupId int    `json:"product_group_id"`
	IsActive       bool   `json:"is_active"`
	ProductGroup   string `json:"product_group"`
}

func WProductGroupGet(id int) (WProductGroup, error) {
	var p WProductGroup
	row := db.QueryRow(`SELECT product_group.*, IFNULL(pr.name, "") FROM product_group
	LEFT JOIN product_group AS pr ON product_group.product_group_id = pr.id WHERE product_group.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ProductGroupId,
		&p.IsActive,
		&p.ProductGroup,
	)
	return p, err
}

func WProductGroupGetAll(withDeleted bool, deletedOnly bool) ([]WProductGroup, error) {
	query := `SELECT product_group.*, IFNULL(pr.name, "") FROM product_group
	LEFT JOIN product_group AS pr ON product_group.product_group_id = pr.id`
	if deletedOnly {
		query += "  WHERE product_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE product_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductGroup{}
	for rows.Next() {
		var p WProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
			&p.ProductGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProductGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProductGroup, error) {

	if !ProductGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_group.*, IFNULL(pr.name, "") FROM product_group
	LEFT JOIN product_group AS pr ON product_group.product_group_id = pr.id WHERE product_group.%s=?`, field)
	if deletedOnly {
		query += "  AND product_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductGroup{}
	for rows.Next() {
		var p WProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
			&p.ProductGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProductGroup, error) {

	if !ProductGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_group.*, IFNULL(pr.name, "") FROM product_group
	LEFT JOIN product_group AS pr ON product_group.product_group_id = pr.id WHERE product_group.%s=?`, field)
	if deletedOnly {
		query += "  AND product_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductGroup{}
	for rows.Next() {
		var p WProductGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProductGroupId,
			&p.IsActive,
			&p.ProductGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WProduct struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	ShortName      string  `json:"short_name"`
	ProductGroupId int     `json:"product_group_id"`
	MeasureId      int     `json:"measure_id"`
	Width          float64 `json:"width"`
	Length         float64 `json:"length"`
	MinCost        float64 `json:"min_cost"`
	Cost           float64 `json:"cost"`
	UserId         int     `json:"user_id"`
	Barcode        string  `json:"barcode"`
	IsActive       bool    `json:"is_active"`
	ProductGroup   string  `json:"product_group"`
	Measure        string  `json:"measure"`
	User           string  `json:"user"`
}

func WProductGet(id int) (WProduct, error) {
	var p WProduct
	row := db.QueryRow(`SELECT product.*, IFNULL(product_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, "") FROM product
	LEFT JOIN product_group ON product.product_group_id = product_group.id
	LEFT JOIN measure ON product.measure_id = measure.id
	LEFT JOIN user ON product.user_id = user.id WHERE product.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ShortName,
		&p.ProductGroupId,
		&p.MeasureId,
		&p.Width,
		&p.Length,
		&p.MinCost,
		&p.Cost,
		&p.UserId,
		&p.Barcode,
		&p.IsActive,
		&p.ProductGroup,
		&p.Measure,
		&p.User,
	)
	return p, err
}

func WProductGetAll(withDeleted bool, deletedOnly bool) ([]WProduct, error) {
	query := `SELECT product.*, IFNULL(product_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, "") FROM product
	LEFT JOIN product_group ON product.product_group_id = product_group.id
	LEFT JOIN measure ON product.measure_id = measure.id
	LEFT JOIN user ON product.user_id = user.id`
	if deletedOnly {
		query += "  WHERE product.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE product.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProduct{}
	for rows.Next() {
		var p WProduct
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
			&p.ProductGroup,
			&p.Measure,
			&p.User,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProduct, error) {

	if !ProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product.*, IFNULL(product_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, "") FROM product
	LEFT JOIN product_group ON product.product_group_id = product_group.id
	LEFT JOIN measure ON product.measure_id = measure.id
	LEFT JOIN user ON product.user_id = user.id WHERE product.%s=?`, field)
	if deletedOnly {
		query += "  AND product.is_active = 0"
	} else if !withDeleted {
		query += "  AND product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProduct{}
	for rows.Next() {
		var p WProduct
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
			&p.ProductGroup,
			&p.Measure,
			&p.User,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProduct, error) {

	if !ProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product.*, IFNULL(product_group.name, ""), IFNULL(measure.name, ""), IFNULL(user.name, "") FROM product
	LEFT JOIN product_group ON product.product_group_id = product_group.id
	LEFT JOIN measure ON product.measure_id = measure.id
	LEFT JOIN user ON product.user_id = user.id WHERE product.%s=?`, field)
	if deletedOnly {
		query += "  AND product.is_active = 0"
	} else if !withDeleted {
		query += "  AND product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProduct{}
	for rows.Next() {
		var p WProduct
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ShortName,
			&p.ProductGroupId,
			&p.MeasureId,
			&p.Width,
			&p.Length,
			&p.MinCost,
			&p.Cost,
			&p.UserId,
			&p.Barcode,
			&p.IsActive,
			&p.ProductGroup,
			&p.Measure,
			&p.User,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WContragentGroup struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	ContragentGroupId int    `json:"contragent_group_id"`
	IsActive          bool   `json:"is_active"`
	ContragentGroup   string `json:"contragent_group"`
}

func WContragentGroupGet(id int) (WContragentGroup, error) {
	var c WContragentGroup
	row := db.QueryRow(`SELECT contragent_group.*, IFNULL(co.name, "") FROM contragent_group
	LEFT JOIN contragent_group AS co ON contragent_group.contragent_group_id = co.id WHERE contragent_group.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ContragentGroupId,
		&c.IsActive,
		&c.ContragentGroup,
	)
	return c, err
}

func WContragentGroupGetAll(withDeleted bool, deletedOnly bool) ([]WContragentGroup, error) {
	query := `SELECT contragent_group.*, IFNULL(co.name, "") FROM contragent_group
	LEFT JOIN contragent_group AS co ON contragent_group.contragent_group_id = co.id`
	if deletedOnly {
		query += "  WHERE contragent_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE contragent_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragentGroup{}
	for rows.Next() {
		var c WContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WContragentGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WContragentGroup, error) {

	if !ContragentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contragent_group.*, IFNULL(co.name, "") FROM contragent_group
	LEFT JOIN contragent_group AS co ON contragent_group.contragent_group_id = co.id WHERE contragent_group.%s=?`, field)
	if deletedOnly {
		query += "  AND contragent_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND contragent_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragentGroup{}
	for rows.Next() {
		var c WContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WContragentGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WContragentGroup, error) {

	if !ContragentGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contragent_group.*, IFNULL(co.name, "") FROM contragent_group
	LEFT JOIN contragent_group AS co ON contragent_group.contragent_group_id = co.id WHERE contragent_group.%s=?`, field)
	if deletedOnly {
		query += "  AND contragent_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND contragent_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragentGroup{}
	for rows.Next() {
		var c WContragentGroup
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WContragent struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	ContragentGroupId int     `json:"contragent_group_id"`
	Phone             string  `json:"phone"`
	Email             string  `json:"email"`
	Web               string  `json:"web"`
	Comm              string  `json:"comm"`
	DirName           string  `json:"dir_name"`
	Search            string  `json:"search"`
	Total             float64 `json:"total"`
	FullName          string  `json:"full_name"`
	Edrpou            string  `json:"edrpou"`
	Ipn               string  `json:"ipn"`
	Iban              string  `json:"iban"`
	Bank              string  `json:"bank"`
	Mfo               string  `json:"mfo"`
	Fop               string  `json:"fop"`
	Address           string  `json:"address"`
	IsActive          bool    `json:"is_active"`
	ContragentGroup   string  `json:"contragent_group"`
}

func WContragentGet(id int) (WContragent, error) {
	var c WContragent
	row := db.QueryRow(`SELECT contragent.*, IFNULL(contragent_group.name, "") FROM contragent
	LEFT JOIN contragent_group ON contragent.contragent_group_id = contragent_group.id WHERE contragent.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.ContragentGroupId,
		&c.Phone,
		&c.Email,
		&c.Web,
		&c.Comm,
		&c.DirName,
		&c.Search,
		&c.Total,
		&c.FullName,
		&c.Edrpou,
		&c.Ipn,
		&c.Iban,
		&c.Bank,
		&c.Mfo,
		&c.Fop,
		&c.Address,
		&c.IsActive,
		&c.ContragentGroup,
	)
	return c, err
}

func WContragentGetAll(withDeleted bool, deletedOnly bool) ([]WContragent, error) {
	query := `SELECT contragent.*, IFNULL(contragent_group.name, "") FROM contragent
	LEFT JOIN contragent_group ON contragent.contragent_group_id = contragent_group.id`
	if deletedOnly {
		query += "  WHERE contragent.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE contragent.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragent{}
	for rows.Next() {
		var c WContragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WContragentGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WContragent, error) {

	if !ContragentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contragent.*, IFNULL(contragent_group.name, "") FROM contragent
	LEFT JOIN contragent_group ON contragent.contragent_group_id = contragent_group.id WHERE contragent.%s=?`, field)
	if deletedOnly {
		query += "  AND contragent.is_active = 0"
	} else if !withDeleted {
		query += "  AND contragent.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragent{}
	for rows.Next() {
		var c WContragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WContragentGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WContragent, error) {

	if !ContragentTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contragent.*, IFNULL(contragent_group.name, "") FROM contragent
	LEFT JOIN contragent_group ON contragent.contragent_group_id = contragent_group.id WHERE contragent.%s=?`, field)
	if deletedOnly {
		query += "  AND contragent.is_active = 0"
	} else if !withDeleted {
		query += "  AND contragent.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContragent{}
	for rows.Next() {
		var c WContragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WContragentFindByContragentSearchContactSearch(fs string) ([]WContragent, error) {
	fs = "%" + fs + "%"

	query := `
        SELECT DISTINCT contragent.*, contragent_group.name FROM contragent
        JOIN contragent_group ON contragent.contragent_group_id =   contragent_group.id
        JOIN contact on contragent.id = contact.contragent_id
        WHERE contragent.is_active=1
        AND contact.is_active=1
        AND (contragent.search LIKE ? OR contact.search LIKE ?);`

	rows, err := db.Query(query, fs, fs)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []WContragent{}
	for rows.Next() {
		var c WContragent
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.ContragentGroupId,
			&c.Phone,
			&c.Email,
			&c.Web,
			&c.Comm,
			&c.DirName,
			&c.Search,
			&c.Total,
			&c.FullName,
			&c.Edrpou,
			&c.Ipn,
			&c.Iban,
			&c.Bank,
			&c.Mfo,
			&c.Fop,
			&c.Address,
			&c.IsActive,
			&c.ContragentGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type WContact struct {
	Id           int     `json:"id"`
	ContragentId int     `json:"contragent_id"`
	Name         string  `json:"name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Viber        string  `json:"viber"`
	Telegram     string  `json:"telegram"`
	TelegramUid  int     `json:"telegram_uid"`
	Search       string  `json:"search"`
	Total        float64 `json:"total"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
	Contragent   string  `json:"contragent"`
}

func WContactGet(id int) (WContact, error) {
	var c WContact
	row := db.QueryRow(`SELECT contact.*, IFNULL(contragent.name, "") FROM contact
	LEFT JOIN contragent ON contact.contragent_id = contragent.id WHERE contact.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.ContragentId,
		&c.Name,
		&c.Phone,
		&c.Email,
		&c.Viber,
		&c.Telegram,
		&c.TelegramUid,
		&c.Search,
		&c.Total,
		&c.Comm,
		&c.IsActive,
		&c.Contragent,
	)
	return c, err
}

func WContactGetAll(withDeleted bool, deletedOnly bool) ([]WContact, error) {
	query := `SELECT contact.*, IFNULL(contragent.name, "") FROM contact
	LEFT JOIN contragent ON contact.contragent_id = contragent.id`
	if deletedOnly {
		query += "  WHERE contact.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE contact.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContact{}
	for rows.Next() {
		var c WContact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
			&c.Contragent,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WContactGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WContact, error) {

	if !ContactTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contact.*, IFNULL(contragent.name, "") FROM contact
	LEFT JOIN contragent ON contact.contragent_id = contragent.id WHERE contact.%s=?`, field)
	if deletedOnly {
		query += "  AND contact.is_active = 0"
	} else if !withDeleted {
		query += "  AND contact.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContact{}
	for rows.Next() {
		var c WContact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
			&c.Contragent,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WContactGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WContact, error) {

	if !ContactTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT contact.*, IFNULL(contragent.name, "") FROM contact
	LEFT JOIN contragent ON contact.contragent_id = contragent.id WHERE contact.%s=?`, field)
	if deletedOnly {
		query += "  AND contact.is_active = 0"
	} else if !withDeleted {
		query += "  AND contact.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WContact{}
	for rows.Next() {
		var c WContact
		if err := rows.Scan(
			&c.Id,
			&c.ContragentId,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Viber,
			&c.Telegram,
			&c.TelegramUid,
			&c.Search,
			&c.Total,
			&c.Comm,
			&c.IsActive,
			&c.Contragent,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WOrderingStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func WOrderingStatusGet(id int) (WOrderingStatus, error) {
	var o WOrderingStatus
	row := db.QueryRow(`SELECT ordering_status.* FROM ordering_status WHERE ordering_status.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.IsActive,
	)
	return o, err
}

func WOrderingStatusGetAll(withDeleted bool, deletedOnly bool) ([]WOrderingStatus, error) {
	query := `SELECT ordering_status.* FROM ordering_status`
	if deletedOnly {
		query += "  WHERE ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE ordering_status.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrderingStatus{}
	for rows.Next() {
		var o WOrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOrderingStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOrderingStatus, error) {

	if !OrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT ordering_status.* FROM ordering_status WHERE ordering_status.%s=?`, field)
	if deletedOnly {
		query += "  AND ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrderingStatus{}
	for rows.Next() {
		var o WOrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOrderingStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOrderingStatus, error) {

	if !OrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT ordering_status.* FROM ordering_status WHERE ordering_status.%s=?`, field)
	if deletedOnly {
		query += "  AND ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrderingStatus{}
	for rows.Next() {
		var o WOrderingStatus
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

type WOrdering struct {
	Id               int     `json:"id"`
	DocumentUid      int     `json:"document_uid"`
	Name             string  `json:"name"`
	CreatedAt        string  `json:"created_at"`
	DeadlineAt       string  `json:"deadline_at"`
	UserId           int     `json:"user_id"`
	ContragentId     int     `json:"contragent_id"`
	ContactId        int     `json:"contact_id"`
	Price            float64 `json:"price"`
	Persent          float64 `json:"persent"`
	Profit           float64 `json:"profit"`
	Cost             float64 `json:"cost"`
	Info             string  `json:"info"`
	OrderingStatusId int     `json:"ordering_status_id"`
	IsActive         bool    `json:"is_active"`
	User             string  `json:"user"`
	Contragent       string  `json:"contragent"`
	Contact          string  `json:"contact"`
	OrderingStatus   string  `json:"ordering_status"`
}

func WOrderingGet(id int) (WOrdering, error) {
	var o WOrdering
	row := db.QueryRow(`SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id WHERE ordering.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.DocumentUid,
		&o.Name,
		&o.CreatedAt,
		&o.DeadlineAt,
		&o.UserId,
		&o.ContragentId,
		&o.ContactId,
		&o.Price,
		&o.Persent,
		&o.Profit,
		&o.Cost,
		&o.Info,
		&o.OrderingStatusId,
		&o.IsActive,
		&o.User,
		&o.Contragent,
		&o.Contact,
		&o.OrderingStatus,
	)
	return o, err
}

func WOrderingGetAll(withDeleted bool, deletedOnly bool) ([]WOrdering, error) {
	query := `SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id`
	if deletedOnly {
		query += "  WHERE ordering.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE ordering.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrdering{}
	for rows.Next() {
		var o WOrdering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
			&o.User,
			&o.Contragent,
			&o.Contact,
			&o.OrderingStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOrdering, error) {

	if !OrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id WHERE ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrdering{}
	for rows.Next() {
		var o WOrdering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
			&o.User,
			&o.Contragent,
			&o.Contact,
			&o.OrderingStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOrdering, error) {

	if !OrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id WHERE ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrdering{}
	for rows.Next() {
		var o WOrdering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
			&o.User,
			&o.Contragent,
			&o.Contact,
			&o.OrderingStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOrderingGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WOrdering, error) {
	query := `SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id WHERE (ordering.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrdering{}
	for rows.Next() {
		var o WOrdering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
			&o.User,
			&o.Contragent,
			&o.Contact,
			&o.OrderingStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOrderingGetBetweenDeadlineAt(deadline_at1, deadline_at2 string, withDeleted bool, deletedOnly bool) ([]WOrdering, error) {
	query := `SELECT ordering.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(ordering_status.name, "") FROM ordering
	LEFT JOIN user ON ordering.user_id = user.id
	LEFT JOIN contragent ON ordering.contragent_id = contragent.id
	LEFT JOIN contact ON ordering.contact_id = contact.id
	LEFT JOIN ordering_status ON ordering.ordering_status_id = ordering_status.id WHERE (ordering.deadline_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND ordering.is_active = 1"
	}

	rows, err := db.Query(query, deadline_at1, deadline_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOrdering{}
	for rows.Next() {
		var o WOrdering
		if err := rows.Scan(
			&o.Id,
			&o.DocumentUid,
			&o.Name,
			&o.CreatedAt,
			&o.DeadlineAt,
			&o.UserId,
			&o.ContragentId,
			&o.ContactId,
			&o.Price,
			&o.Persent,
			&o.Profit,
			&o.Cost,
			&o.Info,
			&o.OrderingStatusId,
			&o.IsActive,
			&o.User,
			&o.Contragent,
			&o.Contact,
			&o.OrderingStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

type WOwner struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Web      string  `json:"web"`
	Comm     string  `json:"comm"`
	Total    float64 `json:"total"`
	FullName string  `json:"full_name"`
	Edrpou   string  `json:"edrpou"`
	Ipn      string  `json:"ipn"`
	Iban     string  `json:"iban"`
	Bank     string  `json:"bank"`
	Mfo      string  `json:"mfo"`
	Fop      string  `json:"fop"`
	Address  string  `json:"address"`
	Sign     string  `json:"sign"`
	IsActive bool    `json:"is_active"`
}

func WOwnerGet(id int) (WOwner, error) {
	var o WOwner
	row := db.QueryRow(`SELECT owner.* FROM owner WHERE owner.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.Name,
		&o.Phone,
		&o.Email,
		&o.Web,
		&o.Comm,
		&o.Total,
		&o.FullName,
		&o.Edrpou,
		&o.Ipn,
		&o.Iban,
		&o.Bank,
		&o.Mfo,
		&o.Fop,
		&o.Address,
		&o.Sign,
		&o.IsActive,
	)
	return o, err
}

func WOwnerGetAll(withDeleted bool, deletedOnly bool) ([]WOwner, error) {
	query := `SELECT owner.* FROM owner`
	if deletedOnly {
		query += "  WHERE owner.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE owner.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOwner{}
	for rows.Next() {
		var o WOwner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOwnerGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOwner, error) {

	if !OwnerTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT owner.* FROM owner WHERE owner.%s=?`, field)
	if deletedOnly {
		query += "  AND owner.is_active = 0"
	} else if !withDeleted {
		query += "  AND owner.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOwner{}
	for rows.Next() {
		var o WOwner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOwnerGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOwner, error) {

	if !OwnerTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT owner.* FROM owner WHERE owner.%s=?`, field)
	if deletedOnly {
		query += "  AND owner.is_active = 0"
	} else if !withDeleted {
		query += "  AND owner.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOwner{}
	for rows.Next() {
		var o WOwner
		if err := rows.Scan(
			&o.Id,
			&o.Name,
			&o.Phone,
			&o.Email,
			&o.Web,
			&o.Comm,
			&o.Total,
			&o.FullName,
			&o.Edrpou,
			&o.Ipn,
			&o.Iban,
			&o.Bank,
			&o.Mfo,
			&o.Fop,
			&o.Address,
			&o.Sign,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

type WInvoice struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	OrderingId   int     `json:"ordering_id"`
	BasedOn      int     `json:"based_on"`
	OwnerId      int     `json:"owner_id"`
	Name         string  `json:"name"`
	CreatedAt    string  `json:"created_at"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
	Ordering     string  `json:"ordering"`
	Owner        string  `json:"owner"`
	User         string  `json:"user"`
	Contragent   string  `json:"contragent"`
	Contact      string  `json:"contact"`
}

func WInvoiceGet(id int) (WInvoice, error) {
	var i WInvoice
	row := db.QueryRow(`SELECT invoice.*, IFNULL(ordering.name, ""), IFNULL(owner.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM invoice
	LEFT JOIN ordering ON invoice.ordering_id = ordering.id
	LEFT JOIN owner ON invoice.owner_id = owner.id
	LEFT JOIN user ON invoice.user_id = user.id
	LEFT JOIN contragent ON invoice.contragent_id = contragent.id
	LEFT JOIN contact ON invoice.contact_id = contact.id WHERE invoice.id=?`, id)
	err := row.Scan(
		&i.Id,
		&i.DocumentUid,
		&i.OrderingId,
		&i.BasedOn,
		&i.OwnerId,
		&i.Name,
		&i.CreatedAt,
		&i.UserId,
		&i.ContragentId,
		&i.ContactId,
		&i.CashSum,
		&i.Comm,
		&i.IsActive,
		&i.Ordering,
		&i.Owner,
		&i.User,
		&i.Contragent,
		&i.Contact,
	)
	return i, err
}

func WInvoiceGetAll(withDeleted bool, deletedOnly bool) ([]WInvoice, error) {
	query := `SELECT invoice.*, IFNULL(ordering.name, ""), IFNULL(owner.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM invoice
	LEFT JOIN ordering ON invoice.ordering_id = ordering.id
	LEFT JOIN owner ON invoice.owner_id = owner.id
	LEFT JOIN user ON invoice.user_id = user.id
	LEFT JOIN contragent ON invoice.contragent_id = contragent.id
	LEFT JOIN contact ON invoice.contact_id = contact.id`
	if deletedOnly {
		query += "  WHERE invoice.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE invoice.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WInvoice{}
	for rows.Next() {
		var i WInvoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
			&i.Ordering,
			&i.Owner,
			&i.User,
			&i.Contragent,
			&i.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func WInvoiceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WInvoice, error) {

	if !InvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT invoice.*, IFNULL(ordering.name, ""), IFNULL(owner.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM invoice
	LEFT JOIN ordering ON invoice.ordering_id = ordering.id
	LEFT JOIN owner ON invoice.owner_id = owner.id
	LEFT JOIN user ON invoice.user_id = user.id
	LEFT JOIN contragent ON invoice.contragent_id = contragent.id
	LEFT JOIN contact ON invoice.contact_id = contact.id WHERE invoice.%s=?`, field)
	if deletedOnly {
		query += "  AND invoice.is_active = 0"
	} else if !withDeleted {
		query += "  AND invoice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WInvoice{}
	for rows.Next() {
		var i WInvoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
			&i.Ordering,
			&i.Owner,
			&i.User,
			&i.Contragent,
			&i.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func WInvoiceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WInvoice, error) {

	if !InvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT invoice.*, IFNULL(ordering.name, ""), IFNULL(owner.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM invoice
	LEFT JOIN ordering ON invoice.ordering_id = ordering.id
	LEFT JOIN owner ON invoice.owner_id = owner.id
	LEFT JOIN user ON invoice.user_id = user.id
	LEFT JOIN contragent ON invoice.contragent_id = contragent.id
	LEFT JOIN contact ON invoice.contact_id = contact.id WHERE invoice.%s=?`, field)
	if deletedOnly {
		query += "  AND invoice.is_active = 0"
	} else if !withDeleted {
		query += "  AND invoice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WInvoice{}
	for rows.Next() {
		var i WInvoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
			&i.Ordering,
			&i.Owner,
			&i.User,
			&i.Contragent,
			&i.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func WInvoiceGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WInvoice, error) {
	query := `SELECT invoice.*, IFNULL(ordering.name, ""), IFNULL(owner.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM invoice
	LEFT JOIN ordering ON invoice.ordering_id = ordering.id
	LEFT JOIN owner ON invoice.owner_id = owner.id
	LEFT JOIN user ON invoice.user_id = user.id
	LEFT JOIN contragent ON invoice.contragent_id = contragent.id
	LEFT JOIN contact ON invoice.contact_id = contact.id WHERE (invoice.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND invoice.is_active = 0"
	} else if !withDeleted {
		query += "  AND invoice.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WInvoice{}
	for rows.Next() {
		var i WInvoice
		if err := rows.Scan(
			&i.Id,
			&i.DocumentUid,
			&i.OrderingId,
			&i.BasedOn,
			&i.OwnerId,
			&i.Name,
			&i.CreatedAt,
			&i.UserId,
			&i.ContragentId,
			&i.ContactId,
			&i.CashSum,
			&i.Comm,
			&i.IsActive,
			&i.Ordering,
			&i.Owner,
			&i.User,
			&i.Contragent,
			&i.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

type WItemToInvoice struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	InvoiceId int     `json:"invoice_id"`
	Number    float64 `json:"number"`
	MeasureId int     `json:"measure_id"`
	Price     float64 `json:"price"`
	Cost      float64 `json:"cost"`
	IsActive  bool    `json:"is_active"`
	Invoice   string  `json:"invoice"`
	Measure   string  `json:"measure"`
}

func WItemToInvoiceGet(id int) (WItemToInvoice, error) {
	var i WItemToInvoice
	row := db.QueryRow(`SELECT item_to_invoice.*, IFNULL(invoice.name, ""), IFNULL(measure.name, "") FROM item_to_invoice
	LEFT JOIN invoice ON item_to_invoice.invoice_id = invoice.id
	LEFT JOIN measure ON item_to_invoice.measure_id = measure.id WHERE item_to_invoice.id=?`, id)
	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.InvoiceId,
		&i.Number,
		&i.MeasureId,
		&i.Price,
		&i.Cost,
		&i.IsActive,
		&i.Invoice,
		&i.Measure,
	)
	return i, err
}

func WItemToInvoiceGetAll(withDeleted bool, deletedOnly bool) ([]WItemToInvoice, error) {
	query := `SELECT item_to_invoice.*, IFNULL(invoice.name, ""), IFNULL(measure.name, "") FROM item_to_invoice
	LEFT JOIN invoice ON item_to_invoice.invoice_id = invoice.id
	LEFT JOIN measure ON item_to_invoice.measure_id = measure.id`
	if deletedOnly {
		query += "  WHERE item_to_invoice.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE item_to_invoice.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToInvoice{}
	for rows.Next() {
		var i WItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
			&i.Invoice,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func WItemToInvoiceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WItemToInvoice, error) {

	if !ItemToInvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT item_to_invoice.*, IFNULL(invoice.name, ""), IFNULL(measure.name, "") FROM item_to_invoice
	LEFT JOIN invoice ON item_to_invoice.invoice_id = invoice.id
	LEFT JOIN measure ON item_to_invoice.measure_id = measure.id WHERE item_to_invoice.%s=?`, field)
	if deletedOnly {
		query += "  AND item_to_invoice.is_active = 0"
	} else if !withDeleted {
		query += "  AND item_to_invoice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToInvoice{}
	for rows.Next() {
		var i WItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
			&i.Invoice,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func WItemToInvoiceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WItemToInvoice, error) {

	if !ItemToInvoiceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT item_to_invoice.*, IFNULL(invoice.name, ""), IFNULL(measure.name, "") FROM item_to_invoice
	LEFT JOIN invoice ON item_to_invoice.invoice_id = invoice.id
	LEFT JOIN measure ON item_to_invoice.measure_id = measure.id WHERE item_to_invoice.%s=?`, field)
	if deletedOnly {
		query += "  AND item_to_invoice.is_active = 0"
	} else if !withDeleted {
		query += "  AND item_to_invoice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToInvoice{}
	for rows.Next() {
		var i WItemToInvoice
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.InvoiceId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Cost,
			&i.IsActive,
			&i.Invoice,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

type WProductToOrderingStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func WProductToOrderingStatusGet(id int) (WProductToOrderingStatus, error) {
	var p WProductToOrderingStatus
	row := db.QueryRow(`SELECT product_to_ordering_status.* FROM product_to_ordering_status WHERE product_to_ordering_status.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.IsActive,
	)
	return p, err
}

func WProductToOrderingStatusGetAll(withDeleted bool, deletedOnly bool) ([]WProductToOrderingStatus, error) {
	query := `SELECT product_to_ordering_status.* FROM product_to_ordering_status`
	if deletedOnly {
		query += "  WHERE product_to_ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE product_to_ordering_status.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrderingStatus{}
	for rows.Next() {
		var p WProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProductToOrderingStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProductToOrderingStatus, error) {

	if !ProductToOrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_ordering_status.* FROM product_to_ordering_status WHERE product_to_ordering_status.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_ordering_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrderingStatus{}
	for rows.Next() {
		var p WProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductToOrderingStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProductToOrderingStatus, error) {

	if !ProductToOrderingStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_ordering_status.* FROM product_to_ordering_status WHERE product_to_ordering_status.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_ordering_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_ordering_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrderingStatus{}
	for rows.Next() {
		var p WProductToOrderingStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WProductToOrdering struct {
	Id                        int     `json:"id"`
	Name                      string  `json:"name"`
	OrderingId                int     `json:"ordering_id"`
	ProductId                 int     `json:"product_id"`
	UserId                    int     `json:"user_id"`
	DeadlineAt                string  `json:"deadline_at"`
	ProductToOrderingStatusId int     `json:"product_to_ordering_status_id"`
	Width                     float64 `json:"width"`
	Length                    float64 `json:"length"`
	Pieces                    int     `json:"pieces"`
	Number                    float64 `json:"number"`
	Price                     float64 `json:"price"`
	Persent                   float64 `json:"persent"`
	Profit                    float64 `json:"profit"`
	Cost                      float64 `json:"cost"`
	Info                      string  `json:"info"`
	ProductToOrderingId       int     `json:"product_to_ordering_id"`
	IsActive                  bool    `json:"is_active"`
	Ordering                  string  `json:"ordering"`
	Product                   string  `json:"product"`
	User                      string  `json:"user"`
	ProductToOrderingStatus   string  `json:"product_to_ordering_status"`
	ProductToOrdering         string  `json:"product_to_ordering"`
}

func WProductToOrderingGet(id int) (WProductToOrdering, error) {
	var p WProductToOrdering
	row := db.QueryRow(`SELECT product_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(product.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering_status.name, ""), IFNULL(pr.name, "") FROM product_to_ordering
	LEFT JOIN ordering ON product_to_ordering.ordering_id = ordering.id
	LEFT JOIN product ON product_to_ordering.product_id = product.id
	LEFT JOIN user ON product_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering_status ON product_to_ordering.product_to_ordering_status_id = product_to_ordering_status.id
	LEFT JOIN product_to_ordering AS pr ON product_to_ordering.product_to_ordering_id = pr.id WHERE product_to_ordering.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.OrderingId,
		&p.ProductId,
		&p.UserId,
		&p.DeadlineAt,
		&p.ProductToOrderingStatusId,
		&p.Width,
		&p.Length,
		&p.Pieces,
		&p.Number,
		&p.Price,
		&p.Persent,
		&p.Profit,
		&p.Cost,
		&p.Info,
		&p.ProductToOrderingId,
		&p.IsActive,
		&p.Ordering,
		&p.Product,
		&p.User,
		&p.ProductToOrderingStatus,
		&p.ProductToOrdering,
	)
	return p, err
}

func WProductToOrderingGetAll(withDeleted bool, deletedOnly bool) ([]WProductToOrdering, error) {
	query := `SELECT product_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(product.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering_status.name, ""), IFNULL(pr.name, "") FROM product_to_ordering
	LEFT JOIN ordering ON product_to_ordering.ordering_id = ordering.id
	LEFT JOIN product ON product_to_ordering.product_id = product.id
	LEFT JOIN user ON product_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering_status ON product_to_ordering.product_to_ordering_status_id = product_to_ordering_status.id
	LEFT JOIN product_to_ordering AS pr ON product_to_ordering.product_to_ordering_id = pr.id`
	if deletedOnly {
		query += "  WHERE product_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE product_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrdering{}
	for rows.Next() {
		var p WProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
			&p.Ordering,
			&p.Product,
			&p.User,
			&p.ProductToOrderingStatus,
			&p.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProductToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProductToOrdering, error) {

	if !ProductToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(product.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering_status.name, ""), IFNULL(pr.name, "") FROM product_to_ordering
	LEFT JOIN ordering ON product_to_ordering.ordering_id = ordering.id
	LEFT JOIN product ON product_to_ordering.product_id = product.id
	LEFT JOIN user ON product_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering_status ON product_to_ordering.product_to_ordering_status_id = product_to_ordering_status.id
	LEFT JOIN product_to_ordering AS pr ON product_to_ordering.product_to_ordering_id = pr.id WHERE product_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrdering{}
	for rows.Next() {
		var p WProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
			&p.Ordering,
			&p.Product,
			&p.User,
			&p.ProductToOrderingStatus,
			&p.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProductToOrdering, error) {

	if !ProductToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(product.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering_status.name, ""), IFNULL(pr.name, "") FROM product_to_ordering
	LEFT JOIN ordering ON product_to_ordering.ordering_id = ordering.id
	LEFT JOIN product ON product_to_ordering.product_id = product.id
	LEFT JOIN user ON product_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering_status ON product_to_ordering.product_to_ordering_status_id = product_to_ordering_status.id
	LEFT JOIN product_to_ordering AS pr ON product_to_ordering.product_to_ordering_id = pr.id WHERE product_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrdering{}
	for rows.Next() {
		var p WProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
			&p.Ordering,
			&p.Product,
			&p.User,
			&p.ProductToOrderingStatus,
			&p.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductToOrderingGetBetweenUpCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WProductToOrdering, error) {
	query := `SELECT product_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(product.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering_status.name, ""), IFNULL(pr.name, "") FROM product_to_ordering
	LEFT JOIN ordering ON product_to_ordering.ordering_id = ordering.id
	LEFT JOIN product ON product_to_ordering.product_id = product.id
	LEFT JOIN user ON product_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering_status ON product_to_ordering.product_to_ordering_status_id = product_to_ordering_status.id
	LEFT JOIN product_to_ordering AS pr ON product_to_ordering.product_to_ordering_id = pr.id
                WHERE (ordering.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND product_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToOrdering{}
	for rows.Next() {
		var p WProductToOrdering
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.OrderingId,
			&p.ProductId,
			&p.UserId,
			&p.DeadlineAt,
			&p.ProductToOrderingStatusId,
			&p.Width,
			&p.Length,
			&p.Pieces,
			&p.Number,
			&p.Price,
			&p.Persent,
			&p.Profit,
			&p.Cost,
			&p.Info,
			&p.ProductToOrderingId,
			&p.IsActive,
			&p.Ordering,
			&p.Product,
			&p.User,
			&p.ProductToOrderingStatus,
			&p.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

type WMatherialToOrdering struct {
	Id                  int     `json:"id"`
	OrderingId          int     `json:"ordering_id"`
	MatherialId         int     `json:"matherial_id"`
	Width               float64 `json:"width"`
	Length              float64 `json:"length"`
	Pieces              int     `json:"pieces"`
	ColorId             int     `json:"color_id"`
	UserId              int     `json:"user_id"`
	Number              float64 `json:"number"`
	Price               float64 `json:"price"`
	Persent             float64 `json:"persent"`
	Profit              float64 `json:"profit"`
	Cost                float64 `json:"cost"`
	Comm                string  `json:"comm"`
	ProductToOrderingId int     `json:"product_to_ordering_id"`
	IsActive            bool    `json:"is_active"`
	Ordering            string  `json:"ordering"`
	Matherial           string  `json:"matherial"`
	Color               string  `json:"color"`
	User                string  `json:"user"`
	ProductToOrdering   string  `json:"product_to_ordering"`
}

func WMatherialToOrderingGet(id int) (WMatherialToOrdering, error) {
	var m WMatherialToOrdering
	row := db.QueryRow(`SELECT matherial_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering.name, "") FROM matherial_to_ordering
	LEFT JOIN ordering ON matherial_to_ordering.ordering_id = ordering.id
	LEFT JOIN matherial ON matherial_to_ordering.matherial_id = matherial.id
	LEFT JOIN color ON matherial_to_ordering.color_id = color.id
	LEFT JOIN user ON matherial_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering ON matherial_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE matherial_to_ordering.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.OrderingId,
		&m.MatherialId,
		&m.Width,
		&m.Length,
		&m.Pieces,
		&m.ColorId,
		&m.UserId,
		&m.Number,
		&m.Price,
		&m.Persent,
		&m.Profit,
		&m.Cost,
		&m.Comm,
		&m.ProductToOrderingId,
		&m.IsActive,
		&m.Ordering,
		&m.Matherial,
		&m.Color,
		&m.User,
		&m.ProductToOrdering,
	)
	return m, err
}

func WMatherialToOrderingGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialToOrdering, error) {
	query := `SELECT matherial_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering.name, "") FROM matherial_to_ordering
	LEFT JOIN ordering ON matherial_to_ordering.ordering_id = ordering.id
	LEFT JOIN matherial ON matherial_to_ordering.matherial_id = matherial.id
	LEFT JOIN color ON matherial_to_ordering.color_id = color.id
	LEFT JOIN user ON matherial_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering ON matherial_to_ordering.product_to_ordering_id = product_to_ordering.id`
	if deletedOnly {
		query += "  WHERE matherial_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToOrdering{}
	for rows.Next() {
		var m WMatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
			&m.Ordering,
			&m.Matherial,
			&m.Color,
			&m.User,
			&m.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialToOrdering, error) {

	if !MatherialToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering.name, "") FROM matherial_to_ordering
	LEFT JOIN ordering ON matherial_to_ordering.ordering_id = ordering.id
	LEFT JOIN matherial ON matherial_to_ordering.matherial_id = matherial.id
	LEFT JOIN color ON matherial_to_ordering.color_id = color.id
	LEFT JOIN user ON matherial_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering ON matherial_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE matherial_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToOrdering{}
	for rows.Next() {
		var m WMatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
			&m.Ordering,
			&m.Matherial,
			&m.Color,
			&m.User,
			&m.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialToOrdering, error) {

	if !MatherialToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering.name, "") FROM matherial_to_ordering
	LEFT JOIN ordering ON matherial_to_ordering.ordering_id = ordering.id
	LEFT JOIN matherial ON matherial_to_ordering.matherial_id = matherial.id
	LEFT JOIN color ON matherial_to_ordering.color_id = color.id
	LEFT JOIN user ON matherial_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering ON matherial_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE matherial_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToOrdering{}
	for rows.Next() {
		var m WMatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
			&m.Ordering,
			&m.Matherial,
			&m.Color,
			&m.User,
			&m.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialToOrderingGetBetweenUpCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WMatherialToOrdering, error) {
	query := `SELECT matherial_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, ""), IFNULL(product_to_ordering.name, "") FROM matherial_to_ordering
	LEFT JOIN ordering ON matherial_to_ordering.ordering_id = ordering.id
	LEFT JOIN matherial ON matherial_to_ordering.matherial_id = matherial.id
	LEFT JOIN color ON matherial_to_ordering.color_id = color.id
	LEFT JOIN user ON matherial_to_ordering.user_id = user.id
	LEFT JOIN product_to_ordering ON matherial_to_ordering.product_to_ordering_id = product_to_ordering.id
                WHERE (ordering.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND matherial_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToOrdering{}
	for rows.Next() {
		var m WMatherialToOrdering
		if err := rows.Scan(
			&m.Id,
			&m.OrderingId,
			&m.MatherialId,
			&m.Width,
			&m.Length,
			&m.Pieces,
			&m.ColorId,
			&m.UserId,
			&m.Number,
			&m.Price,
			&m.Persent,
			&m.Profit,
			&m.Cost,
			&m.Comm,
			&m.ProductToOrderingId,
			&m.IsActive,
			&m.Ordering,
			&m.Matherial,
			&m.Color,
			&m.User,
			&m.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type WMatherialToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	MatherialId   int     `json:"matherial_id"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	Comm          string  `json:"comm"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
	Product       string  `json:"product"`
	Matherial     string  `json:"matherial"`
}

func WMatherialToProductGet(id int) (WMatherialToProduct, error) {
	var m WMatherialToProduct
	row := db.QueryRow(`SELECT matherial_to_product.*, IFNULL(product.name, ""), IFNULL(matherial.name, "") FROM matherial_to_product
	LEFT JOIN product ON matherial_to_product.product_id = product.id
	LEFT JOIN matherial ON matherial_to_product.matherial_id = matherial.id WHERE matherial_to_product.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.ProductId,
		&m.MatherialId,
		&m.Number,
		&m.Coeff,
		&m.Cost,
		&m.ListName,
		&m.IsMultiselect,
		&m.Comm,
		&m.IsUsed,
		&m.IsActive,
		&m.Product,
		&m.Matherial,
	)
	return m, err
}

func WMatherialToProductGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialToProduct, error) {
	query := `SELECT matherial_to_product.*, IFNULL(product.name, ""), IFNULL(matherial.name, "") FROM matherial_to_product
	LEFT JOIN product ON matherial_to_product.product_id = product.id
	LEFT JOIN matherial ON matherial_to_product.matherial_id = matherial.id`
	if deletedOnly {
		query += "  WHERE matherial_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_to_product.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToProduct{}
	for rows.Next() {
		var m WMatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
			&m.Product,
			&m.Matherial,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialToProduct, error) {

	if !MatherialToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_product.*, IFNULL(product.name, ""), IFNULL(matherial.name, "") FROM matherial_to_product
	LEFT JOIN product ON matherial_to_product.product_id = product.id
	LEFT JOIN matherial ON matherial_to_product.matherial_id = matherial.id WHERE matherial_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToProduct{}
	for rows.Next() {
		var m WMatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
			&m.Product,
			&m.Matherial,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialToProduct, error) {

	if !MatherialToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_product.*, IFNULL(product.name, ""), IFNULL(matherial.name, "") FROM matherial_to_product
	LEFT JOIN product ON matherial_to_product.product_id = product.id
	LEFT JOIN matherial ON matherial_to_product.matherial_id = matherial.id WHERE matherial_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToProduct{}
	for rows.Next() {
		var m WMatherialToProduct
		if err := rows.Scan(
			&m.Id,
			&m.ProductId,
			&m.MatherialId,
			&m.Number,
			&m.Coeff,
			&m.Cost,
			&m.ListName,
			&m.IsMultiselect,
			&m.Comm,
			&m.IsUsed,
			&m.IsActive,
			&m.Product,
			&m.Matherial,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WOperationToOrdering struct {
	Id                  int     `json:"id"`
	OrderingId          int     `json:"ordering_id"`
	OperationId         int     `json:"operation_id"`
	UserId              int     `json:"user_id"`
	Number              float64 `json:"number"`
	Price               float64 `json:"price"`
	UserSum             float64 `json:"user_sum"`
	Cost                float64 `json:"cost"`
	EquipmentId         int     `json:"equipment_id"`
	EquipmentCost       float64 `json:"equipment_cost"`
	Comm                string  `json:"comm"`
	ProductToOrderingId int     `json:"product_to_ordering_id"`
	IsActive            bool    `json:"is_active"`
	Ordering            string  `json:"ordering"`
	Operation           string  `json:"operation"`
	User                string  `json:"user"`
	Equipment           string  `json:"equipment"`
	ProductToOrdering   string  `json:"product_to_ordering"`
}

func WOperationToOrderingGet(id int) (WOperationToOrdering, error) {
	var o WOperationToOrdering
	row := db.QueryRow(`SELECT operation_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, ""), IFNULL(product_to_ordering.name, "") FROM operation_to_ordering
	LEFT JOIN ordering ON operation_to_ordering.ordering_id = ordering.id
	LEFT JOIN operation ON operation_to_ordering.operation_id = operation.id
	LEFT JOIN user ON operation_to_ordering.user_id = user.id
	LEFT JOIN equipment ON operation_to_ordering.equipment_id = equipment.id
	LEFT JOIN product_to_ordering ON operation_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE operation_to_ordering.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.OrderingId,
		&o.OperationId,
		&o.UserId,
		&o.Number,
		&o.Price,
		&o.UserSum,
		&o.Cost,
		&o.EquipmentId,
		&o.EquipmentCost,
		&o.Comm,
		&o.ProductToOrderingId,
		&o.IsActive,
		&o.Ordering,
		&o.Operation,
		&o.User,
		&o.Equipment,
		&o.ProductToOrdering,
	)
	return o, err
}

func WOperationToOrderingGetAll(withDeleted bool, deletedOnly bool) ([]WOperationToOrdering, error) {
	query := `SELECT operation_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, ""), IFNULL(product_to_ordering.name, "") FROM operation_to_ordering
	LEFT JOIN ordering ON operation_to_ordering.ordering_id = ordering.id
	LEFT JOIN operation ON operation_to_ordering.operation_id = operation.id
	LEFT JOIN user ON operation_to_ordering.user_id = user.id
	LEFT JOIN equipment ON operation_to_ordering.equipment_id = equipment.id
	LEFT JOIN product_to_ordering ON operation_to_ordering.product_to_ordering_id = product_to_ordering.id`
	if deletedOnly {
		query += "  WHERE operation_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE operation_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToOrdering{}
	for rows.Next() {
		var o WOperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
			&o.Ordering,
			&o.Operation,
			&o.User,
			&o.Equipment,
			&o.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOperationToOrderingGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOperationToOrdering, error) {

	if !OperationToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, ""), IFNULL(product_to_ordering.name, "") FROM operation_to_ordering
	LEFT JOIN ordering ON operation_to_ordering.ordering_id = ordering.id
	LEFT JOIN operation ON operation_to_ordering.operation_id = operation.id
	LEFT JOIN user ON operation_to_ordering.user_id = user.id
	LEFT JOIN equipment ON operation_to_ordering.equipment_id = equipment.id
	LEFT JOIN product_to_ordering ON operation_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE operation_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToOrdering{}
	for rows.Next() {
		var o WOperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
			&o.Ordering,
			&o.Operation,
			&o.User,
			&o.Equipment,
			&o.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOperationToOrderingGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOperationToOrdering, error) {

	if !OperationToOrderingTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, ""), IFNULL(product_to_ordering.name, "") FROM operation_to_ordering
	LEFT JOIN ordering ON operation_to_ordering.ordering_id = ordering.id
	LEFT JOIN operation ON operation_to_ordering.operation_id = operation.id
	LEFT JOIN user ON operation_to_ordering.user_id = user.id
	LEFT JOIN equipment ON operation_to_ordering.equipment_id = equipment.id
	LEFT JOIN product_to_ordering ON operation_to_ordering.product_to_ordering_id = product_to_ordering.id WHERE operation_to_ordering.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_to_ordering.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToOrdering{}
	for rows.Next() {
		var o WOperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
			&o.Ordering,
			&o.Operation,
			&o.User,
			&o.Equipment,
			&o.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOperationToOrderingGetBetweenUpCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WOperationToOrdering, error) {
	query := `SELECT operation_to_ordering.*, IFNULL(ordering.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, ""), IFNULL(product_to_ordering.name, "") FROM operation_to_ordering
	LEFT JOIN ordering ON operation_to_ordering.ordering_id = ordering.id
	LEFT JOIN operation ON operation_to_ordering.operation_id = operation.id
	LEFT JOIN user ON operation_to_ordering.user_id = user.id
	LEFT JOIN equipment ON operation_to_ordering.equipment_id = equipment.id
	LEFT JOIN product_to_ordering ON operation_to_ordering.product_to_ordering_id = product_to_ordering.id
                WHERE (ordering.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND operation_to_ordering.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_to_ordering.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToOrdering{}
	for rows.Next() {
		var o WOperationToOrdering
		if err := rows.Scan(
			&o.Id,
			&o.OrderingId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Price,
			&o.UserSum,
			&o.Cost,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.ProductToOrderingId,
			&o.IsActive,
			&o.Ordering,
			&o.Operation,
			&o.User,
			&o.Equipment,
			&o.ProductToOrdering,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

type WOperationToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	OperationId   int     `json:"operation_id"`
	UserId        int     `json:"user_id"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	EquipmentId   int     `json:"equipment_id"`
	EquipmentCost float64 `json:"equipment_cost"`
	Comm          string  `json:"comm"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
	Product       string  `json:"product"`
	Operation     string  `json:"operation"`
	User          string  `json:"user"`
	Equipment     string  `json:"equipment"`
}

func WOperationToProductGet(id int) (WOperationToProduct, error) {
	var o WOperationToProduct
	row := db.QueryRow(`SELECT operation_to_product.*, IFNULL(product.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation_to_product
	LEFT JOIN product ON operation_to_product.product_id = product.id
	LEFT JOIN operation ON operation_to_product.operation_id = operation.id
	LEFT JOIN user ON operation_to_product.user_id = user.id
	LEFT JOIN equipment ON operation_to_product.equipment_id = equipment.id WHERE operation_to_product.id=?`, id)
	err := row.Scan(
		&o.Id,
		&o.ProductId,
		&o.OperationId,
		&o.UserId,
		&o.Number,
		&o.Coeff,
		&o.Cost,
		&o.ListName,
		&o.IsMultiselect,
		&o.EquipmentId,
		&o.EquipmentCost,
		&o.Comm,
		&o.IsUsed,
		&o.IsActive,
		&o.Product,
		&o.Operation,
		&o.User,
		&o.Equipment,
	)
	return o, err
}

func WOperationToProductGetAll(withDeleted bool, deletedOnly bool) ([]WOperationToProduct, error) {
	query := `SELECT operation_to_product.*, IFNULL(product.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation_to_product
	LEFT JOIN product ON operation_to_product.product_id = product.id
	LEFT JOIN operation ON operation_to_product.operation_id = operation.id
	LEFT JOIN user ON operation_to_product.user_id = user.id
	LEFT JOIN equipment ON operation_to_product.equipment_id = equipment.id`
	if deletedOnly {
		query += "  WHERE operation_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE operation_to_product.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToProduct{}
	for rows.Next() {
		var o WOperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
			&o.Product,
			&o.Operation,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func WOperationToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WOperationToProduct, error) {

	if !OperationToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_to_product.*, IFNULL(product.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation_to_product
	LEFT JOIN product ON operation_to_product.product_id = product.id
	LEFT JOIN operation ON operation_to_product.operation_id = operation.id
	LEFT JOIN user ON operation_to_product.user_id = user.id
	LEFT JOIN equipment ON operation_to_product.equipment_id = equipment.id WHERE operation_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToProduct{}
	for rows.Next() {
		var o WOperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
			&o.Product,
			&o.Operation,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

func WOperationToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WOperationToProduct, error) {

	if !OperationToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT operation_to_product.*, IFNULL(product.name, ""), IFNULL(operation.name, ""), IFNULL(user.name, ""), IFNULL(equipment.name, "") FROM operation_to_product
	LEFT JOIN product ON operation_to_product.product_id = product.id
	LEFT JOIN operation ON operation_to_product.operation_id = operation.id
	LEFT JOIN user ON operation_to_product.user_id = user.id
	LEFT JOIN equipment ON operation_to_product.equipment_id = equipment.id WHERE operation_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND operation_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND operation_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WOperationToProduct{}
	for rows.Next() {
		var o WOperationToProduct
		if err := rows.Scan(
			&o.Id,
			&o.ProductId,
			&o.OperationId,
			&o.UserId,
			&o.Number,
			&o.Coeff,
			&o.Cost,
			&o.ListName,
			&o.IsMultiselect,
			&o.EquipmentId,
			&o.EquipmentCost,
			&o.Comm,
			&o.IsUsed,
			&o.IsActive,
			&o.Product,
			&o.Operation,
			&o.User,
			&o.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil

}

type WProductToProduct struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	Product2Id    int     `json:"product2_id"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	Number        float64 `json:"number"`
	Coeff         float64 `json:"coeff"`
	Cost          float64 `json:"cost"`
	ListName      string  `json:"list_name"`
	IsMultiselect bool    `json:"is_multiselect"`
	IsUsed        bool    `json:"is_used"`
	IsActive      bool    `json:"is_active"`
	Product       string  `json:"product"`
	Product2      string  `json:"product2"`
}

func WProductToProductGet(id int) (WProductToProduct, error) {
	var p WProductToProduct
	row := db.QueryRow(`SELECT product_to_product.*, IFNULL(product.name, ""), IFNULL(product2.name, "") FROM product_to_product
	LEFT JOIN product ON product_to_product.product_id = product.id
	LEFT JOIN product AS product2 ON product_to_product.product2_id = product2.id WHERE product_to_product.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.ProductId,
		&p.Product2Id,
		&p.Width,
		&p.Length,
		&p.Number,
		&p.Coeff,
		&p.Cost,
		&p.ListName,
		&p.IsMultiselect,
		&p.IsUsed,
		&p.IsActive,
		&p.Product,
		&p.Product2,
	)
	return p, err
}

func WProductToProductGetAll(withDeleted bool, deletedOnly bool) ([]WProductToProduct, error) {
	query := `SELECT product_to_product.*, IFNULL(product.name, ""), IFNULL(product2.name, "") FROM product_to_product
	LEFT JOIN product ON product_to_product.product_id = product.id
	LEFT JOIN product AS product2 ON product_to_product.product2_id = product2.id`
	if deletedOnly {
		query += "  WHERE product_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE product_to_product.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToProduct{}
	for rows.Next() {
		var p WProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
			&p.Product,
			&p.Product2,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProductToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProductToProduct, error) {

	if !ProductToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_product.*, IFNULL(product.name, ""), IFNULL(product2.name, "") FROM product_to_product
	LEFT JOIN product ON product_to_product.product_id = product.id
	LEFT JOIN product AS product2 ON product_to_product.product2_id = product2.id WHERE product_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToProduct{}
	for rows.Next() {
		var p WProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
			&p.Product,
			&p.Product2,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProductToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProductToProduct, error) {

	if !ProductToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT product_to_product.*, IFNULL(product.name, ""), IFNULL(product2.name, "") FROM product_to_product
	LEFT JOIN product ON product_to_product.product_id = product.id
	LEFT JOIN product AS product2 ON product_to_product.product2_id = product2.id WHERE product_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND product_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND product_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProductToProduct{}
	for rows.Next() {
		var p WProductToProduct
		if err := rows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Product2Id,
			&p.Width,
			&p.Length,
			&p.Number,
			&p.Coeff,
			&p.Cost,
			&p.ListName,
			&p.IsMultiselect,
			&p.IsUsed,
			&p.IsActive,
			&p.Product,
			&p.Product2,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WCboxCheck struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	FsUid        string  `json:"fs_uid"`
	CheckboxUid  string  `json:"checkbox_uid"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	DocumentUid  int     `json:"document_uid"`
	OrderingId   int     `json:"ordering_id"`
	BasedOn      int     `json:"based_on"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Discount     float64 `json:"discount"`
	Comm         string  `json:"comm"`
	IsCash       bool    `json:"is_cash"`
	IsActive     bool    `json:"is_active"`
	User         string  `json:"user"`
	Contragent   string  `json:"contragent"`
	Ordering     string  `json:"ordering"`
}

func WCboxCheckGet(id int) (WCboxCheck, error) {
	var c WCboxCheck
	row := db.QueryRow(`SELECT cbox_check.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(ordering.name, "") FROM cbox_check
	LEFT JOIN user ON cbox_check.user_id = user.id
	LEFT JOIN contragent ON cbox_check.contragent_id = contragent.id
	LEFT JOIN ordering ON cbox_check.ordering_id = ordering.id WHERE cbox_check.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.FsUid,
		&c.CheckboxUid,
		&c.UserId,
		&c.ContragentId,
		&c.DocumentUid,
		&c.OrderingId,
		&c.BasedOn,
		&c.CreatedAt,
		&c.CashSum,
		&c.Discount,
		&c.Comm,
		&c.IsCash,
		&c.IsActive,
		&c.User,
		&c.Contragent,
		&c.Ordering,
	)
	return c, err
}

func WCboxCheckGetAll(withDeleted bool, deletedOnly bool) ([]WCboxCheck, error) {
	query := `SELECT cbox_check.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(ordering.name, "") FROM cbox_check
	LEFT JOIN user ON cbox_check.user_id = user.id
	LEFT JOIN contragent ON cbox_check.contragent_id = contragent.id
	LEFT JOIN ordering ON cbox_check.ordering_id = ordering.id`
	if deletedOnly {
		query += "  WHERE cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE cbox_check.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCboxCheck{}
	for rows.Next() {
		var c WCboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
			&c.User,
			&c.Contragent,
			&c.Ordering,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCboxCheckGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCboxCheck, error) {

	if !CboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cbox_check.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(ordering.name, "") FROM cbox_check
	LEFT JOIN user ON cbox_check.user_id = user.id
	LEFT JOIN contragent ON cbox_check.contragent_id = contragent.id
	LEFT JOIN ordering ON cbox_check.ordering_id = ordering.id WHERE cbox_check.%s=?`, field)
	if deletedOnly {
		query += "  AND cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  AND cbox_check.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCboxCheck{}
	for rows.Next() {
		var c WCboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
			&c.User,
			&c.Contragent,
			&c.Ordering,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCboxCheckGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCboxCheck, error) {

	if !CboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cbox_check.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(ordering.name, "") FROM cbox_check
	LEFT JOIN user ON cbox_check.user_id = user.id
	LEFT JOIN contragent ON cbox_check.contragent_id = contragent.id
	LEFT JOIN ordering ON cbox_check.ordering_id = ordering.id WHERE cbox_check.%s=?`, field)
	if deletedOnly {
		query += "  AND cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  AND cbox_check.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCboxCheck{}
	for rows.Next() {
		var c WCboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
			&c.User,
			&c.Contragent,
			&c.Ordering,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCboxCheckGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WCboxCheck, error) {
	query := `SELECT cbox_check.*, IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(ordering.name, "") FROM cbox_check
	LEFT JOIN user ON cbox_check.user_id = user.id
	LEFT JOIN contragent ON cbox_check.contragent_id = contragent.id
	LEFT JOIN ordering ON cbox_check.ordering_id = ordering.id WHERE (cbox_check.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  AND cbox_check.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCboxCheck{}
	for rows.Next() {
		var c WCboxCheck
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.FsUid,
			&c.CheckboxUid,
			&c.UserId,
			&c.ContragentId,
			&c.DocumentUid,
			&c.OrderingId,
			&c.BasedOn,
			&c.CreatedAt,
			&c.CashSum,
			&c.Discount,
			&c.Comm,
			&c.IsCash,
			&c.IsActive,
			&c.User,
			&c.Contragent,
			&c.Ordering,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type WItemToCboxCheck struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	CboxCheckId int     `json:"cbox_check_id"`
	Number      float64 `json:"number"`
	MeasureId   int     `json:"measure_id"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	Cost        float64 `json:"cost"`
	ItemCode    string  `json:"item_code"`
	IsActive    bool    `json:"is_active"`
	CboxCheck   string  `json:"cbox_check"`
	Measure     string  `json:"measure"`
}

func WItemToCboxCheckGet(id int) (WItemToCboxCheck, error) {
	var i WItemToCboxCheck
	row := db.QueryRow(`SELECT item_to_cbox_check.*, IFNULL(cbox_check.name, ""), IFNULL(measure.name, "") FROM item_to_cbox_check
	LEFT JOIN cbox_check ON item_to_cbox_check.cbox_check_id = cbox_check.id
	LEFT JOIN measure ON item_to_cbox_check.measure_id = measure.id WHERE item_to_cbox_check.id=?`, id)
	err := row.Scan(
		&i.Id,
		&i.Name,
		&i.CboxCheckId,
		&i.Number,
		&i.MeasureId,
		&i.Price,
		&i.Discount,
		&i.Cost,
		&i.ItemCode,
		&i.IsActive,
		&i.CboxCheck,
		&i.Measure,
	)
	return i, err
}

func WItemToCboxCheckGetAll(withDeleted bool, deletedOnly bool) ([]WItemToCboxCheck, error) {
	query := `SELECT item_to_cbox_check.*, IFNULL(cbox_check.name, ""), IFNULL(measure.name, "") FROM item_to_cbox_check
	LEFT JOIN cbox_check ON item_to_cbox_check.cbox_check_id = cbox_check.id
	LEFT JOIN measure ON item_to_cbox_check.measure_id = measure.id`
	if deletedOnly {
		query += "  WHERE item_to_cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE item_to_cbox_check.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToCboxCheck{}
	for rows.Next() {
		var i WItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
			&i.CboxCheck,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func WItemToCboxCheckGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WItemToCboxCheck, error) {

	if !ItemToCboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT item_to_cbox_check.*, IFNULL(cbox_check.name, ""), IFNULL(measure.name, "") FROM item_to_cbox_check
	LEFT JOIN cbox_check ON item_to_cbox_check.cbox_check_id = cbox_check.id
	LEFT JOIN measure ON item_to_cbox_check.measure_id = measure.id WHERE item_to_cbox_check.%s=?`, field)
	if deletedOnly {
		query += "  AND item_to_cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  AND item_to_cbox_check.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToCboxCheck{}
	for rows.Next() {
		var i WItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
			&i.CboxCheck,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

func WItemToCboxCheckGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WItemToCboxCheck, error) {

	if !ItemToCboxCheckTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT item_to_cbox_check.*, IFNULL(cbox_check.name, ""), IFNULL(measure.name, "") FROM item_to_cbox_check
	LEFT JOIN cbox_check ON item_to_cbox_check.cbox_check_id = cbox_check.id
	LEFT JOIN measure ON item_to_cbox_check.measure_id = measure.id WHERE item_to_cbox_check.%s=?`, field)
	if deletedOnly {
		query += "  AND item_to_cbox_check.is_active = 0"
	} else if !withDeleted {
		query += "  AND item_to_cbox_check.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WItemToCboxCheck{}
	for rows.Next() {
		var i WItemToCboxCheck
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.CboxCheckId,
			&i.Number,
			&i.MeasureId,
			&i.Price,
			&i.Discount,
			&i.Cost,
			&i.ItemCode,
			&i.IsActive,
			&i.CboxCheck,
			&i.Measure,
		); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil

}

type WCashIn struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	CashId       int     `json:"cash_id"`
	UserId       int     `json:"user_id"`
	BasedOn      int     `json:"based_on"`
	CboxCheckId  int     `json:"cbox_check_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
	Cash         string  `json:"cash"`
	User         string  `json:"user"`
	CboxCheck    string  `json:"cbox_check"`
	Contragent   string  `json:"contragent"`
	Contact      string  `json:"contact"`
}

func WCashInGet(id int) (WCashIn, error) {
	var c WCashIn
	row := db.QueryRow(`SELECT cash_in.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_in
	LEFT JOIN cash ON cash_in.cash_id = cash.id
	LEFT JOIN user ON cash_in.user_id = user.id
	LEFT JOIN cbox_check ON cash_in.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_in.contragent_id = contragent.id
	LEFT JOIN contact ON cash_in.contact_id = contact.id WHERE cash_in.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.DocumentUid,
		&c.Name,
		&c.CashId,
		&c.UserId,
		&c.BasedOn,
		&c.CboxCheckId,
		&c.ContragentId,
		&c.ContactId,
		&c.CreatedAt,
		&c.CashSum,
		&c.Comm,
		&c.IsActive,
		&c.Cash,
		&c.User,
		&c.CboxCheck,
		&c.Contragent,
		&c.Contact,
	)
	return c, err
}

func WCashInGetAll(withDeleted bool, deletedOnly bool) ([]WCashIn, error) {
	query := `SELECT cash_in.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_in
	LEFT JOIN cash ON cash_in.cash_id = cash.id
	LEFT JOIN user ON cash_in.user_id = user.id
	LEFT JOIN cbox_check ON cash_in.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_in.contragent_id = contragent.id
	LEFT JOIN contact ON cash_in.contact_id = contact.id`
	if deletedOnly {
		query += "  WHERE cash_in.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE cash_in.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashIn{}
	for rows.Next() {
		var c WCashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCashInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCashIn, error) {

	if !CashInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash_in.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_in
	LEFT JOIN cash ON cash_in.cash_id = cash.id
	LEFT JOIN user ON cash_in.user_id = user.id
	LEFT JOIN cbox_check ON cash_in.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_in.contragent_id = contragent.id
	LEFT JOIN contact ON cash_in.contact_id = contact.id WHERE cash_in.%s=?`, field)
	if deletedOnly {
		query += "  AND cash_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashIn{}
	for rows.Next() {
		var c WCashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCashInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCashIn, error) {

	if !CashInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash_in.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_in
	LEFT JOIN cash ON cash_in.cash_id = cash.id
	LEFT JOIN user ON cash_in.user_id = user.id
	LEFT JOIN cbox_check ON cash_in.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_in.contragent_id = contragent.id
	LEFT JOIN contact ON cash_in.contact_id = contact.id WHERE cash_in.%s=?`, field)
	if deletedOnly {
		query += "  AND cash_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashIn{}
	for rows.Next() {
		var c WCashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCashInGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WCashIn, error) {
	query := `SELECT cash_in.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_in
	LEFT JOIN cash ON cash_in.cash_id = cash.id
	LEFT JOIN user ON cash_in.user_id = user.id
	LEFT JOIN cbox_check ON cash_in.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_in.contragent_id = contragent.id
	LEFT JOIN contact ON cash_in.contact_id = contact.id WHERE (cash_in.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND cash_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_in.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashIn{}
	for rows.Next() {
		var c WCashIn
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type WCashOut struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	CashId       int     `json:"cash_id"`
	UserId       int     `json:"user_id"`
	BasedOn      int     `json:"based_on"`
	CboxCheckId  int     `json:"cbox_check_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	CashSum      float64 `json:"cash_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
	Cash         string  `json:"cash"`
	User         string  `json:"user"`
	CboxCheck    string  `json:"cbox_check"`
	Contragent   string  `json:"contragent"`
	Contact      string  `json:"contact"`
}

func WCashOutGet(id int) (WCashOut, error) {
	var c WCashOut
	row := db.QueryRow(`SELECT cash_out.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_out
	LEFT JOIN cash ON cash_out.cash_id = cash.id
	LEFT JOIN user ON cash_out.user_id = user.id
	LEFT JOIN cbox_check ON cash_out.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_out.contragent_id = contragent.id
	LEFT JOIN contact ON cash_out.contact_id = contact.id WHERE cash_out.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.DocumentUid,
		&c.Name,
		&c.CashId,
		&c.UserId,
		&c.BasedOn,
		&c.CboxCheckId,
		&c.ContragentId,
		&c.ContactId,
		&c.CreatedAt,
		&c.CashSum,
		&c.Comm,
		&c.IsActive,
		&c.Cash,
		&c.User,
		&c.CboxCheck,
		&c.Contragent,
		&c.Contact,
	)
	return c, err
}

func WCashOutGetAll(withDeleted bool, deletedOnly bool) ([]WCashOut, error) {
	query := `SELECT cash_out.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_out
	LEFT JOIN cash ON cash_out.cash_id = cash.id
	LEFT JOIN user ON cash_out.user_id = user.id
	LEFT JOIN cbox_check ON cash_out.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_out.contragent_id = contragent.id
	LEFT JOIN contact ON cash_out.contact_id = contact.id`
	if deletedOnly {
		query += "  WHERE cash_out.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE cash_out.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashOut{}
	for rows.Next() {
		var c WCashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCashOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCashOut, error) {

	if !CashOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash_out.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_out
	LEFT JOIN cash ON cash_out.cash_id = cash.id
	LEFT JOIN user ON cash_out.user_id = user.id
	LEFT JOIN cbox_check ON cash_out.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_out.contragent_id = contragent.id
	LEFT JOIN contact ON cash_out.contact_id = contact.id WHERE cash_out.%s=?`, field)
	if deletedOnly {
		query += "  AND cash_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashOut{}
	for rows.Next() {
		var c WCashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCashOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCashOut, error) {

	if !CashOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT cash_out.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_out
	LEFT JOIN cash ON cash_out.cash_id = cash.id
	LEFT JOIN user ON cash_out.user_id = user.id
	LEFT JOIN cbox_check ON cash_out.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_out.contragent_id = contragent.id
	LEFT JOIN contact ON cash_out.contact_id = contact.id WHERE cash_out.%s=?`, field)
	if deletedOnly {
		query += "  AND cash_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashOut{}
	for rows.Next() {
		var c WCashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCashOutGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WCashOut, error) {
	query := `SELECT cash_out.*, IFNULL(cash.name, ""), IFNULL(user.name, ""), IFNULL(cbox_check.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM cash_out
	LEFT JOIN cash ON cash_out.cash_id = cash.id
	LEFT JOIN user ON cash_out.user_id = user.id
	LEFT JOIN cbox_check ON cash_out.cbox_check_id = cbox_check.id
	LEFT JOIN contragent ON cash_out.contragent_id = contragent.id
	LEFT JOIN contact ON cash_out.contact_id = contact.id WHERE (cash_out.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND cash_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND cash_out.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCashOut{}
	for rows.Next() {
		var c WCashOut
		if err := rows.Scan(
			&c.Id,
			&c.DocumentUid,
			&c.Name,
			&c.CashId,
			&c.UserId,
			&c.BasedOn,
			&c.CboxCheckId,
			&c.ContragentId,
			&c.ContactId,
			&c.CreatedAt,
			&c.CashSum,
			&c.Comm,
			&c.IsActive,
			&c.Cash,
			&c.User,
			&c.CboxCheck,
			&c.Contragent,
			&c.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

type WWhs struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Comm     string `json:"comm"`
	IsActive bool   `json:"is_active"`
}

func WWhsGet(id int) (WWhs, error) {
	var w WWhs
	row := db.QueryRow(`SELECT whs.* FROM whs WHERE whs.id=?`, id)
	err := row.Scan(
		&w.Id,
		&w.Name,
		&w.Comm,
		&w.IsActive,
	)
	return w, err
}

func WWhsGetAll(withDeleted bool, deletedOnly bool) ([]WWhs, error) {
	query := `SELECT whs.* FROM whs`
	if deletedOnly {
		query += "  WHERE whs.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE whs.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhs{}
	for rows.Next() {
		var w WWhs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WWhsGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WWhs, error) {

	if !WhsTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs.* FROM whs WHERE whs.%s=?`, field)
	if deletedOnly {
		query += "  AND whs.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhs{}
	for rows.Next() {
		var w WWhs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWhsGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WWhs, error) {

	if !WhsTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs.* FROM whs WHERE whs.%s=?`, field)
	if deletedOnly {
		query += "  AND whs.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhs{}
	for rows.Next() {
		var w WWhs
		if err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.Comm,
			&w.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

type WWhsIn struct {
	Id                  int     `json:"id"`
	DocumentUid         int     `json:"document_uid"`
	Name                string  `json:"name"`
	BasedOn             int     `json:"based_on"`
	WhsId               int     `json:"whs_id"`
	UserId              int     `json:"user_id"`
	ContragentId        int     `json:"contragent_id"`
	ContactId           int     `json:"contact_id"`
	ContragentDocUid    string  `json:"contragent_doc_uid"`
	ContragentCreatedAt string  `json:"contragent_created_at"`
	CreatedAt           string  `json:"created_at"`
	WhsSum              float64 `json:"whs_sum"`
	Delivery            float64 `json:"delivery"`
	Comm                string  `json:"comm"`
	IsActive            bool    `json:"is_active"`
	Whs                 string  `json:"whs"`
	User                string  `json:"user"`
	Contragent          string  `json:"contragent"`
	Contact             string  `json:"contact"`
}

func WWhsInGet(id int) (WWhsIn, error) {
	var w WWhsIn
	row := db.QueryRow(`SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id WHERE whs_in.id=?`, id)
	err := row.Scan(
		&w.Id,
		&w.DocumentUid,
		&w.Name,
		&w.BasedOn,
		&w.WhsId,
		&w.UserId,
		&w.ContragentId,
		&w.ContactId,
		&w.ContragentDocUid,
		&w.ContragentCreatedAt,
		&w.CreatedAt,
		&w.WhsSum,
		&w.Delivery,
		&w.Comm,
		&w.IsActive,
		&w.Whs,
		&w.User,
		&w.Contragent,
		&w.Contact,
	)
	return w, err
}

func WWhsInGetAll(withDeleted bool, deletedOnly bool) ([]WWhsIn, error) {
	query := `SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id`
	if deletedOnly {
		query += "  WHERE whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE whs_in.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsIn{}
	for rows.Next() {
		var w WWhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WWhsInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WWhsIn, error) {

	if !WhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id WHERE whs_in.%s=?`, field)
	if deletedOnly {
		query += "  AND whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsIn{}
	for rows.Next() {
		var w WWhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWhsInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WWhsIn, error) {

	if !WhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id WHERE whs_in.%s=?`, field)
	if deletedOnly {
		query += "  AND whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsIn{}
	for rows.Next() {
		var w WWhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWhsInGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WWhsIn, error) {
	query := `SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id WHERE (whs_in.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_in.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsIn{}
	for rows.Next() {
		var w WWhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WWhsInGetBetweenContragentCreatedAt(contragent_created_at1, contragent_created_at2 string, withDeleted bool, deletedOnly bool) ([]WWhsIn, error) {
	query := `SELECT whs_in.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_in
	LEFT JOIN whs ON whs_in.whs_id = whs.id
	LEFT JOIN user ON whs_in.user_id = user.id
	LEFT JOIN contragent ON whs_in.contragent_id = contragent.id
	LEFT JOIN contact ON whs_in.contact_id = contact.id WHERE (whs_in.contragent_created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_in.is_active = 1"
	}

	rows, err := db.Query(query, contragent_created_at1, contragent_created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsIn{}
	for rows.Next() {
		var w WWhsIn
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.ContragentDocUid,
			&w.ContragentCreatedAt,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Delivery,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

type WWhsOut struct {
	Id           int     `json:"id"`
	DocumentUid  int     `json:"document_uid"`
	Name         string  `json:"name"`
	BasedOn      int     `json:"based_on"`
	WhsId        int     `json:"whs_id"`
	UserId       int     `json:"user_id"`
	ContragentId int     `json:"contragent_id"`
	ContactId    int     `json:"contact_id"`
	CreatedAt    string  `json:"created_at"`
	WhsSum       float64 `json:"whs_sum"`
	Comm         string  `json:"comm"`
	IsActive     bool    `json:"is_active"`
	Whs          string  `json:"whs"`
	User         string  `json:"user"`
	Contragent   string  `json:"contragent"`
	Contact      string  `json:"contact"`
}

func WWhsOutGet(id int) (WWhsOut, error) {
	var w WWhsOut
	row := db.QueryRow(`SELECT whs_out.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_out
	LEFT JOIN whs ON whs_out.whs_id = whs.id
	LEFT JOIN user ON whs_out.user_id = user.id
	LEFT JOIN contragent ON whs_out.contragent_id = contragent.id
	LEFT JOIN contact ON whs_out.contact_id = contact.id WHERE whs_out.id=?`, id)
	err := row.Scan(
		&w.Id,
		&w.DocumentUid,
		&w.Name,
		&w.BasedOn,
		&w.WhsId,
		&w.UserId,
		&w.ContragentId,
		&w.ContactId,
		&w.CreatedAt,
		&w.WhsSum,
		&w.Comm,
		&w.IsActive,
		&w.Whs,
		&w.User,
		&w.Contragent,
		&w.Contact,
	)
	return w, err
}

func WWhsOutGetAll(withDeleted bool, deletedOnly bool) ([]WWhsOut, error) {
	query := `SELECT whs_out.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_out
	LEFT JOIN whs ON whs_out.whs_id = whs.id
	LEFT JOIN user ON whs_out.user_id = user.id
	LEFT JOIN contragent ON whs_out.contragent_id = contragent.id
	LEFT JOIN contact ON whs_out.contact_id = contact.id`
	if deletedOnly {
		query += "  WHERE whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE whs_out.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsOut{}
	for rows.Next() {
		var w WWhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WWhsOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WWhsOut, error) {

	if !WhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs_out.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_out
	LEFT JOIN whs ON whs_out.whs_id = whs.id
	LEFT JOIN user ON whs_out.user_id = user.id
	LEFT JOIN contragent ON whs_out.contragent_id = contragent.id
	LEFT JOIN contact ON whs_out.contact_id = contact.id WHERE whs_out.%s=?`, field)
	if deletedOnly {
		query += "  AND whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsOut{}
	for rows.Next() {
		var w WWhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWhsOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WWhsOut, error) {

	if !WhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT whs_out.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_out
	LEFT JOIN whs ON whs_out.whs_id = whs.id
	LEFT JOIN user ON whs_out.user_id = user.id
	LEFT JOIN contragent ON whs_out.contragent_id = contragent.id
	LEFT JOIN contact ON whs_out.contact_id = contact.id WHERE whs_out.%s=?`, field)
	if deletedOnly {
		query += "  AND whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsOut{}
	for rows.Next() {
		var w WWhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWhsOutGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WWhsOut, error) {
	query := `SELECT whs_out.*, IFNULL(whs.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, "") FROM whs_out
	LEFT JOIN whs ON whs_out.whs_id = whs.id
	LEFT JOIN user ON whs_out.user_id = user.id
	LEFT JOIN contragent ON whs_out.contragent_id = contragent.id
	LEFT JOIN contact ON whs_out.contact_id = contact.id WHERE (whs_out.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND whs_out.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWhsOut{}
	for rows.Next() {
		var w WWhsOut
		if err := rows.Scan(
			&w.Id,
			&w.DocumentUid,
			&w.Name,
			&w.BasedOn,
			&w.WhsId,
			&w.UserId,
			&w.ContragentId,
			&w.ContactId,
			&w.CreatedAt,
			&w.WhsSum,
			&w.Comm,
			&w.IsActive,
			&w.Whs,
			&w.User,
			&w.Contragent,
			&w.Contact,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

type WMatherialToWhsIn struct {
	Id               int     `json:"id"`
	MatherialId      int     `json:"matherial_id"`
	ContragentMatUid string  `json:"contragent_mat_uid"`
	WhsInId          int     `json:"whs_in_id"`
	Number           float64 `json:"number"`
	Price            float64 `json:"price"`
	Cost             float64 `json:"cost"`
	Width            float64 `json:"width"`
	Length           float64 `json:"length"`
	ColorId          int     `json:"color_id"`
	IsActive         bool    `json:"is_active"`
	Matherial        string  `json:"matherial"`
	WhsIn            string  `json:"whs_in"`
	Color            string  `json:"color"`
}

func WMatherialToWhsInGet(id int) (WMatherialToWhsIn, error) {
	var m WMatherialToWhsIn
	row := db.QueryRow(`SELECT matherial_to_whs_in.*, IFNULL(matherial.name, ""), IFNULL(whs_in.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_in
	LEFT JOIN matherial ON matherial_to_whs_in.matherial_id = matherial.id
	LEFT JOIN whs_in ON matherial_to_whs_in.whs_in_id = whs_in.id
	LEFT JOIN color ON matherial_to_whs_in.color_id = color.id WHERE matherial_to_whs_in.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.ContragentMatUid,
		&m.WhsInId,
		&m.Number,
		&m.Price,
		&m.Cost,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.IsActive,
		&m.Matherial,
		&m.WhsIn,
		&m.Color,
	)
	return m, err
}

func WMatherialToWhsInGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialToWhsIn, error) {
	query := `SELECT matherial_to_whs_in.*, IFNULL(matherial.name, ""), IFNULL(whs_in.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_in
	LEFT JOIN matherial ON matherial_to_whs_in.matherial_id = matherial.id
	LEFT JOIN whs_in ON matherial_to_whs_in.whs_in_id = whs_in.id
	LEFT JOIN color ON matherial_to_whs_in.color_id = color.id`
	if deletedOnly {
		query += "  WHERE matherial_to_whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_to_whs_in.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsIn{}
	for rows.Next() {
		var m WMatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsIn,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialToWhsInGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialToWhsIn, error) {

	if !MatherialToWhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_whs_in.*, IFNULL(matherial.name, ""), IFNULL(whs_in.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_in
	LEFT JOIN matherial ON matherial_to_whs_in.matherial_id = matherial.id
	LEFT JOIN whs_in ON matherial_to_whs_in.whs_in_id = whs_in.id
	LEFT JOIN color ON matherial_to_whs_in.color_id = color.id WHERE matherial_to_whs_in.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_whs_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsIn{}
	for rows.Next() {
		var m WMatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsIn,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialToWhsInGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialToWhsIn, error) {

	if !MatherialToWhsInTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_whs_in.*, IFNULL(matherial.name, ""), IFNULL(whs_in.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_in
	LEFT JOIN matherial ON matherial_to_whs_in.matherial_id = matherial.id
	LEFT JOIN whs_in ON matherial_to_whs_in.whs_in_id = whs_in.id
	LEFT JOIN color ON matherial_to_whs_in.color_id = color.id WHERE matherial_to_whs_in.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_whs_in.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_whs_in.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsIn{}
	for rows.Next() {
		var m WMatherialToWhsIn
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.ContragentMatUid,
			&m.WhsInId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsIn,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WMatherialToWhsOut struct {
	Id          int     `json:"id"`
	MatherialId int     `json:"matherial_id"`
	WhsOutId    int     `json:"whs_out_id"`
	Number      float64 `json:"number"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	Width       float64 `json:"width"`
	Length      float64 `json:"length"`
	ColorId     int     `json:"color_id"`
	IsActive    bool    `json:"is_active"`
	Matherial   string  `json:"matherial"`
	WhsOut      string  `json:"whs_out"`
	Color       string  `json:"color"`
}

func WMatherialToWhsOutGet(id int) (WMatherialToWhsOut, error) {
	var m WMatherialToWhsOut
	row := db.QueryRow(`SELECT matherial_to_whs_out.*, IFNULL(matherial.name, ""), IFNULL(whs_out.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_out
	LEFT JOIN matherial ON matherial_to_whs_out.matherial_id = matherial.id
	LEFT JOIN whs_out ON matherial_to_whs_out.whs_out_id = whs_out.id
	LEFT JOIN color ON matherial_to_whs_out.color_id = color.id WHERE matherial_to_whs_out.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.WhsOutId,
		&m.Number,
		&m.Price,
		&m.Cost,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.IsActive,
		&m.Matherial,
		&m.WhsOut,
		&m.Color,
	)
	return m, err
}

func WMatherialToWhsOutGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialToWhsOut, error) {
	query := `SELECT matherial_to_whs_out.*, IFNULL(matherial.name, ""), IFNULL(whs_out.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_out
	LEFT JOIN matherial ON matherial_to_whs_out.matherial_id = matherial.id
	LEFT JOIN whs_out ON matherial_to_whs_out.whs_out_id = whs_out.id
	LEFT JOIN color ON matherial_to_whs_out.color_id = color.id`
	if deletedOnly {
		query += "  WHERE matherial_to_whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_to_whs_out.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsOut{}
	for rows.Next() {
		var m WMatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsOut,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialToWhsOutGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialToWhsOut, error) {

	if !MatherialToWhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_whs_out.*, IFNULL(matherial.name, ""), IFNULL(whs_out.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_out
	LEFT JOIN matherial ON matherial_to_whs_out.matherial_id = matherial.id
	LEFT JOIN whs_out ON matherial_to_whs_out.whs_out_id = whs_out.id
	LEFT JOIN color ON matherial_to_whs_out.color_id = color.id WHERE matherial_to_whs_out.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_whs_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsOut{}
	for rows.Next() {
		var m WMatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsOut,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialToWhsOutGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialToWhsOut, error) {

	if !MatherialToWhsOutTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_to_whs_out.*, IFNULL(matherial.name, ""), IFNULL(whs_out.name, ""), IFNULL(color.name, "") FROM matherial_to_whs_out
	LEFT JOIN matherial ON matherial_to_whs_out.matherial_id = matherial.id
	LEFT JOIN whs_out ON matherial_to_whs_out.whs_out_id = whs_out.id
	LEFT JOIN color ON matherial_to_whs_out.color_id = color.id WHERE matherial_to_whs_out.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_to_whs_out.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_to_whs_out.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialToWhsOut{}
	for rows.Next() {
		var m WMatherialToWhsOut
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.WhsOutId,
			&m.Number,
			&m.Price,
			&m.Cost,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.IsActive,
			&m.Matherial,
			&m.WhsOut,
			&m.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

type WMatherialPart struct {
	Id          int     `json:"id"`
	MatherialId int     `json:"matherial_id"`
	PartUid     int     `json:"part_uid"`
	Number      float64 `json:"number"`
	Width       float64 `json:"width"`
	Length      float64 `json:"length"`
	ColorId     int     `json:"color_id"`
	UserId      int     `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	IsRecycle   bool    `json:"is_recycle"`
	IsActive    bool    `json:"is_active"`
	Matherial   string  `json:"matherial"`
	Color       string  `json:"color"`
	User        string  `json:"user"`
}

func WMatherialPartGet(id int) (WMatherialPart, error) {
	var m WMatherialPart
	row := db.QueryRow(`SELECT matherial_part.*, IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, "") FROM matherial_part
	LEFT JOIN matherial ON matherial_part.matherial_id = matherial.id
	LEFT JOIN color ON matherial_part.color_id = color.id
	LEFT JOIN user ON matherial_part.user_id = user.id WHERE matherial_part.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.MatherialId,
		&m.PartUid,
		&m.Number,
		&m.Width,
		&m.Length,
		&m.ColorId,
		&m.UserId,
		&m.CreatedAt,
		&m.IsRecycle,
		&m.IsActive,
		&m.Matherial,
		&m.Color,
		&m.User,
	)
	return m, err
}

func WMatherialPartGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialPart, error) {
	query := `SELECT matherial_part.*, IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, "") FROM matherial_part
	LEFT JOIN matherial ON matherial_part.matherial_id = matherial.id
	LEFT JOIN color ON matherial_part.color_id = color.id
	LEFT JOIN user ON matherial_part.user_id = user.id`
	if deletedOnly {
		query += "  WHERE matherial_part.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_part.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPart{}
	for rows.Next() {
		var m WMatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
			&m.Matherial,
			&m.Color,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialPartGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialPart, error) {

	if !MatherialPartTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_part.*, IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, "") FROM matherial_part
	LEFT JOIN matherial ON matherial_part.matherial_id = matherial.id
	LEFT JOIN color ON matherial_part.color_id = color.id
	LEFT JOIN user ON matherial_part.user_id = user.id WHERE matherial_part.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_part.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPart{}
	for rows.Next() {
		var m WMatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
			&m.Matherial,
			&m.Color,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialPartGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialPart, error) {

	if !MatherialPartTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_part.*, IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, "") FROM matherial_part
	LEFT JOIN matherial ON matherial_part.matherial_id = matherial.id
	LEFT JOIN color ON matherial_part.color_id = color.id
	LEFT JOIN user ON matherial_part.user_id = user.id WHERE matherial_part.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_part.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPart{}
	for rows.Next() {
		var m WMatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
			&m.Matherial,
			&m.Color,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialPartGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WMatherialPart, error) {
	query := `SELECT matherial_part.*, IFNULL(matherial.name, ""), IFNULL(color.name, ""), IFNULL(user.name, "") FROM matherial_part
	LEFT JOIN matherial ON matherial_part.matherial_id = matherial.id
	LEFT JOIN color ON matherial_part.color_id = color.id
	LEFT JOIN user ON matherial_part.user_id = user.id WHERE (matherial_part.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND matherial_part.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPart{}
	for rows.Next() {
		var m WMatherialPart
		if err := rows.Scan(
			&m.Id,
			&m.MatherialId,
			&m.PartUid,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.ColorId,
			&m.UserId,
			&m.CreatedAt,
			&m.IsRecycle,
			&m.IsActive,
			&m.Matherial,
			&m.Color,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type WMatherialPartSlice struct {
	Id              int     `json:"id"`
	MatherialPartId int     `json:"matherial_part_id"`
	UserId          int     `json:"user_id"`
	CreatedAt       string  `json:"created_at"`
	Number          float64 `json:"number"`
	Width           float64 `json:"width"`
	Length          float64 `json:"length"`
	Comm            string  `json:"comm"`
	IsActive        bool    `json:"is_active"`
	MatherialPart   string  `json:"matherial_part"`
	User            string  `json:"user"`
}

func WMatherialPartSliceGet(id int) (WMatherialPartSlice, error) {
	var m WMatherialPartSlice
	row := db.QueryRow(`SELECT matherial_part_slice.*, IFNULL(matherial_part.name, ""), IFNULL(user.name, "") FROM matherial_part_slice
	LEFT JOIN matherial_part ON matherial_part_slice.matherial_part_id = matherial_part.id
	LEFT JOIN user ON matherial_part_slice.user_id = user.id WHERE matherial_part_slice.id=?`, id)
	err := row.Scan(
		&m.Id,
		&m.MatherialPartId,
		&m.UserId,
		&m.CreatedAt,
		&m.Number,
		&m.Width,
		&m.Length,
		&m.Comm,
		&m.IsActive,
		&m.MatherialPart,
		&m.User,
	)
	return m, err
}

func WMatherialPartSliceGetAll(withDeleted bool, deletedOnly bool) ([]WMatherialPartSlice, error) {
	query := `SELECT matherial_part_slice.*, IFNULL(matherial_part.name, ""), IFNULL(user.name, "") FROM matherial_part_slice
	LEFT JOIN matherial_part ON matherial_part_slice.matherial_part_id = matherial_part.id
	LEFT JOIN user ON matherial_part_slice.user_id = user.id`
	if deletedOnly {
		query += "  WHERE matherial_part_slice.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE matherial_part_slice.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPartSlice{}
	for rows.Next() {
		var m WMatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
			&m.MatherialPart,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func WMatherialPartSliceGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WMatherialPartSlice, error) {

	if !MatherialPartSliceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_part_slice.*, IFNULL(matherial_part.name, ""), IFNULL(user.name, "") FROM matherial_part_slice
	LEFT JOIN matherial_part ON matherial_part_slice.matherial_part_id = matherial_part.id
	LEFT JOIN user ON matherial_part_slice.user_id = user.id WHERE matherial_part_slice.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_part_slice.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part_slice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPartSlice{}
	for rows.Next() {
		var m WMatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
			&m.MatherialPart,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialPartSliceGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WMatherialPartSlice, error) {

	if !MatherialPartSliceTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT matherial_part_slice.*, IFNULL(matherial_part.name, ""), IFNULL(user.name, "") FROM matherial_part_slice
	LEFT JOIN matherial_part ON matherial_part_slice.matherial_part_id = matherial_part.id
	LEFT JOIN user ON matherial_part_slice.user_id = user.id WHERE matherial_part_slice.%s=?`, field)
	if deletedOnly {
		query += "  AND matherial_part_slice.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part_slice.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPartSlice{}
	for rows.Next() {
		var m WMatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
			&m.MatherialPart,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil

}

func WMatherialPartSliceGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WMatherialPartSlice, error) {
	query := `SELECT matherial_part_slice.*, IFNULL(matherial_part.name, ""), IFNULL(user.name, "") FROM matherial_part_slice
	LEFT JOIN matherial_part ON matherial_part_slice.matherial_part_id = matherial_part.id
	LEFT JOIN user ON matherial_part_slice.user_id = user.id WHERE (matherial_part_slice.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND matherial_part_slice.is_active = 0"
	} else if !withDeleted {
		query += "  AND matherial_part_slice.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WMatherialPartSlice{}
	for rows.Next() {
		var m WMatherialPartSlice
		if err := rows.Scan(
			&m.Id,
			&m.MatherialPartId,
			&m.UserId,
			&m.CreatedAt,
			&m.Number,
			&m.Width,
			&m.Length,
			&m.Comm,
			&m.IsActive,
			&m.MatherialPart,
			&m.User,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type WProjectGroup struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ProjectGroupId int    `json:"project_group_id"`
	IsActive       bool   `json:"is_active"`
	ProjectGroup   string `json:"project_group"`
}

func WProjectGroupGet(id int) (WProjectGroup, error) {
	var p WProjectGroup
	row := db.QueryRow(`SELECT project_group.*, IFNULL(pr.name, "") FROM project_group
	LEFT JOIN project_group AS pr ON project_group.project_group_id = pr.id WHERE project_group.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.ProjectGroupId,
		&p.IsActive,
		&p.ProjectGroup,
	)
	return p, err
}

func WProjectGroupGetAll(withDeleted bool, deletedOnly bool) ([]WProjectGroup, error) {
	query := `SELECT project_group.*, IFNULL(pr.name, "") FROM project_group
	LEFT JOIN project_group AS pr ON project_group.project_group_id = pr.id`
	if deletedOnly {
		query += "  WHERE project_group.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE project_group.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectGroup{}
	for rows.Next() {
		var p WProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
			&p.ProjectGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProjectGroupGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProjectGroup, error) {

	if !ProjectGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_group.*, IFNULL(pr.name, "") FROM project_group
	LEFT JOIN project_group AS pr ON project_group.project_group_id = pr.id WHERE project_group.%s=?`, field)
	if deletedOnly {
		query += "  AND project_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectGroup{}
	for rows.Next() {
		var p WProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
			&p.ProjectGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProjectGroupGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProjectGroup, error) {

	if !ProjectGroupTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_group.*, IFNULL(pr.name, "") FROM project_group
	LEFT JOIN project_group AS pr ON project_group.project_group_id = pr.id WHERE project_group.%s=?`, field)
	if deletedOnly {
		query += "  AND project_group.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_group.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectGroup{}
	for rows.Next() {
		var p WProjectGroup
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.ProjectGroupId,
			&p.IsActive,
			&p.ProjectGroup,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WProjectStatus struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	CodeName string `json:"code_name"`
	IsActive bool   `json:"is_active"`
}

func WProjectStatusGet(id int) (WProjectStatus, error) {
	var p WProjectStatus
	row := db.QueryRow(`SELECT project_status.* FROM project_status WHERE project_status.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.CodeName,
		&p.IsActive,
	)
	return p, err
}

func WProjectStatusGetAll(withDeleted bool, deletedOnly bool) ([]WProjectStatus, error) {
	query := `SELECT project_status.* FROM project_status`
	if deletedOnly {
		query += "  WHERE project_status.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE project_status.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectStatus{}
	for rows.Next() {
		var p WProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProjectStatusGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProjectStatus, error) {

	if !ProjectStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_status.* FROM project_status WHERE project_status.%s=?`, field)
	if deletedOnly {
		query += "  AND project_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectStatus{}
	for rows.Next() {
		var p WProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProjectStatusGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProjectStatus, error) {

	if !ProjectStatusTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_status.* FROM project_status WHERE project_status.%s=?`, field)
	if deletedOnly {
		query += "  AND project_status.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_status.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectStatus{}
	for rows.Next() {
		var p WProjectStatus
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.CodeName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WProjectType struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	DirName  string `json:"dir_name"`
	IsActive bool   `json:"is_active"`
}

func WProjectTypeGet(id int) (WProjectType, error) {
	var p WProjectType
	row := db.QueryRow(`SELECT project_type.* FROM project_type WHERE project_type.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.DirName,
		&p.IsActive,
	)
	return p, err
}

func WProjectTypeGetAll(withDeleted bool, deletedOnly bool) ([]WProjectType, error) {
	query := `SELECT project_type.* FROM project_type`
	if deletedOnly {
		query += "  WHERE project_type.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE project_type.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectType{}
	for rows.Next() {
		var p WProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProjectTypeGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProjectType, error) {

	if !ProjectTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_type.* FROM project_type WHERE project_type.%s=?`, field)
	if deletedOnly {
		query += "  AND project_type.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_type.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectType{}
	for rows.Next() {
		var p WProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProjectTypeGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProjectType, error) {

	if !ProjectTypeTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project_type.* FROM project_type WHERE project_type.%s=?`, field)
	if deletedOnly {
		query += "  AND project_type.is_active = 0"
	} else if !withDeleted {
		query += "  AND project_type.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProjectType{}
	for rows.Next() {
		var p WProjectType
		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.DirName,
			&p.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

type WProject struct {
	Id              int     `json:"id"`
	DocumentUid     int     `json:"document_uid"`
	Name            string  `json:"name"`
	ProjectGroupId  int     `json:"project_group_id"`
	UserId          int     `json:"user_id"`
	ContragentId    int     `json:"contragent_id"`
	ContactId       int     `json:"contact_id"`
	Cost            float64 `json:"cost"`
	CashSum         float64 `json:"cash_sum"`
	WhsSum          float64 `json:"whs_sum"`
	ProjectTypeId   int     `json:"project_type_id"`
	TypeDir         string  `json:"type_dir"`
	ProjectStatusId int     `json:"project_status_id"`
	NumberDir       string  `json:"number_dir"`
	Info            string  `json:"info"`
	CreatedAt       string  `json:"created_at"`
	IsInWork        bool    `json:"is_in_work"`
	IsActive        bool    `json:"is_active"`
	ProjectGroup    string  `json:"project_group"`
	User            string  `json:"user"`
	Contragent      string  `json:"contragent"`
	Contact         string  `json:"contact"`
	ProjectType     string  `json:"project_type"`
	ProjectStatus   string  `json:"project_status"`
}

func WProjectGet(id int) (WProject, error) {
	var p WProject
	row := db.QueryRow(`SELECT project.*, IFNULL(project_group.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(project_type.name, ""), IFNULL(project_status.name, "") FROM project
	LEFT JOIN project_group ON project.project_group_id = project_group.id
	LEFT JOIN user ON project.user_id = user.id
	LEFT JOIN contragent ON project.contragent_id = contragent.id
	LEFT JOIN contact ON project.contact_id = contact.id
	LEFT JOIN project_type ON project.project_type_id = project_type.id
	LEFT JOIN project_status ON project.project_status_id = project_status.id WHERE project.id=?`, id)
	err := row.Scan(
		&p.Id,
		&p.DocumentUid,
		&p.Name,
		&p.ProjectGroupId,
		&p.UserId,
		&p.ContragentId,
		&p.ContactId,
		&p.Cost,
		&p.CashSum,
		&p.WhsSum,
		&p.ProjectTypeId,
		&p.TypeDir,
		&p.ProjectStatusId,
		&p.NumberDir,
		&p.Info,
		&p.CreatedAt,
		&p.IsInWork,
		&p.IsActive,
		&p.ProjectGroup,
		&p.User,
		&p.Contragent,
		&p.Contact,
		&p.ProjectType,
		&p.ProjectStatus,
	)
	return p, err
}

func WProjectGetAll(withDeleted bool, deletedOnly bool) ([]WProject, error) {
	query := `SELECT project.*, IFNULL(project_group.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(project_type.name, ""), IFNULL(project_status.name, "") FROM project
	LEFT JOIN project_group ON project.project_group_id = project_group.id
	LEFT JOIN user ON project.user_id = user.id
	LEFT JOIN contragent ON project.contragent_id = contragent.id
	LEFT JOIN contact ON project.contact_id = contact.id
	LEFT JOIN project_type ON project.project_type_id = project_type.id
	LEFT JOIN project_status ON project.project_status_id = project_status.id`
	if deletedOnly {
		query += "  WHERE project.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE project.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProject{}
	for rows.Next() {
		var p WProject
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
			&p.ProjectGroup,
			&p.User,
			&p.Contragent,
			&p.Contact,
			&p.ProjectType,
			&p.ProjectStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProjectGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WProject, error) {

	if !ProjectTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project.*, IFNULL(project_group.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(project_type.name, ""), IFNULL(project_status.name, "") FROM project
	LEFT JOIN project_group ON project.project_group_id = project_group.id
	LEFT JOIN user ON project.user_id = user.id
	LEFT JOIN contragent ON project.contragent_id = contragent.id
	LEFT JOIN contact ON project.contact_id = contact.id
	LEFT JOIN project_type ON project.project_type_id = project_type.id
	LEFT JOIN project_status ON project.project_status_id = project_status.id WHERE project.%s=?`, field)
	if deletedOnly {
		query += "  AND project.is_active = 0"
	} else if !withDeleted {
		query += "  AND project.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProject{}
	for rows.Next() {
		var p WProject
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
			&p.ProjectGroup,
			&p.User,
			&p.Contragent,
			&p.Contact,
			&p.ProjectType,
			&p.ProjectStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProjectGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WProject, error) {

	if !ProjectTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT project.*, IFNULL(project_group.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(project_type.name, ""), IFNULL(project_status.name, "") FROM project
	LEFT JOIN project_group ON project.project_group_id = project_group.id
	LEFT JOIN user ON project.user_id = user.id
	LEFT JOIN contragent ON project.contragent_id = contragent.id
	LEFT JOIN contact ON project.contact_id = contact.id
	LEFT JOIN project_type ON project.project_type_id = project_type.id
	LEFT JOIN project_status ON project.project_status_id = project_status.id WHERE project.%s=?`, field)
	if deletedOnly {
		query += "  AND project.is_active = 0"
	} else if !withDeleted {
		query += "  AND project.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProject{}
	for rows.Next() {
		var p WProject
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
			&p.ProjectGroup,
			&p.User,
			&p.Contragent,
			&p.Contact,
			&p.ProjectType,
			&p.ProjectStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil

}

func WProjectFindByProjectInfoContragentNoSearchContactNoSearch(fs string) ([]WProject, error) {
	fs = "%" + fs + "%"

	query := `
SELECT project.*, project_group.name, user.name, contragent.name, contact.name, project_type.name, project_status.name FROM project
    JOIN project_group ON project.project_group_id = project_group.id
    JOIN user ON project.user_id = user.id
    JOIN contragent on contragent.id = project.contragent_id
    JOIN contact on contact.id = project.contact_id
    JOIN project_type ON project.project_type_id = project_type.id
    JOIN project_status ON project.project_status_id = project_status.id
                WHERE project.is_active=1
                AND contragent.is_active=1
                AND contact.is_active=1
                AND (
                 project.info LIKE ?
            )
                AND NOT (
                 contragent.search LIKE ?
                OR contact.search LIKE ?
            );`

	rows, err := db.Query(query, fs, fs, fs)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []WProject{}
	for rows.Next() {
		var p WProject
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
			&p.ProjectGroup,
			&p.User,
			&p.Contragent,
			&p.Contact,
			&p.ProjectType,
			&p.ProjectStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func WProjectGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WProject, error) {
	query := `SELECT project.*, IFNULL(project_group.name, ""), IFNULL(user.name, ""), IFNULL(contragent.name, ""), IFNULL(contact.name, ""), IFNULL(project_type.name, ""), IFNULL(project_status.name, "") FROM project
	LEFT JOIN project_group ON project.project_group_id = project_group.id
	LEFT JOIN user ON project.user_id = user.id
	LEFT JOIN contragent ON project.contragent_id = contragent.id
	LEFT JOIN contact ON project.contact_id = contact.id
	LEFT JOIN project_type ON project.project_type_id = project_type.id
	LEFT JOIN project_status ON project.project_status_id = project_status.id WHERE (project.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND project.is_active = 0"
	} else if !withDeleted {
		query += "  AND project.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WProject{}
	for rows.Next() {
		var p WProject
		if err := rows.Scan(
			&p.Id,
			&p.DocumentUid,
			&p.Name,
			&p.ProjectGroupId,
			&p.UserId,
			&p.ContragentId,
			&p.ContactId,
			&p.Cost,
			&p.CashSum,
			&p.WhsSum,
			&p.ProjectTypeId,
			&p.TypeDir,
			&p.ProjectStatusId,
			&p.NumberDir,
			&p.Info,
			&p.CreatedAt,
			&p.IsInWork,
			&p.IsActive,
			&p.ProjectGroup,
			&p.User,
			&p.Contragent,
			&p.Contact,
			&p.ProjectType,
			&p.ProjectStatus,
		); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

type WCounter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	EquipmentId int    `json:"equipment_id"`
	Total       int    `json:"total"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
	Equipment   string `json:"equipment"`
}

func WCounterGet(id int) (WCounter, error) {
	var c WCounter
	row := db.QueryRow(`SELECT counter.*, IFNULL(equipment.name, "") FROM counter
	LEFT JOIN equipment ON counter.equipment_id = equipment.id WHERE counter.id=?`, id)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.EquipmentId,
		&c.Total,
		&c.UpdatedAt,
		&c.IsActive,
		&c.Equipment,
	)
	return c, err
}

func WCounterGetAll(withDeleted bool, deletedOnly bool) ([]WCounter, error) {
	query := `SELECT counter.*, IFNULL(equipment.name, "") FROM counter
	LEFT JOIN equipment ON counter.equipment_id = equipment.id`
	if deletedOnly {
		query += "  WHERE counter.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE counter.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCounter{}
	for rows.Next() {
		var c WCounter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
			&c.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func WCounterGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WCounter, error) {

	if !CounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT counter.*, IFNULL(equipment.name, "") FROM counter
	LEFT JOIN equipment ON counter.equipment_id = equipment.id WHERE counter.%s=?`, field)
	if deletedOnly {
		query += "  AND counter.is_active = 0"
	} else if !withDeleted {
		query += "  AND counter.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCounter{}
	for rows.Next() {
		var c WCounter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
			&c.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

func WCounterGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WCounter, error) {

	if !CounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT counter.*, IFNULL(equipment.name, "") FROM counter
	LEFT JOIN equipment ON counter.equipment_id = equipment.id WHERE counter.%s=?`, field)
	if deletedOnly {
		query += "  AND counter.is_active = 0"
	} else if !withDeleted {
		query += "  AND counter.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WCounter{}
	for rows.Next() {
		var c WCounter
		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.EquipmentId,
			&c.Total,
			&c.UpdatedAt,
			&c.IsActive,
			&c.Equipment,
		); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil

}

type WRecordToCounter struct {
	Id        int    `json:"id"`
	CounterId int    `json:"counter_id"`
	CreatedAt string `json:"created_at"`
	Number    int    `json:"number"`
	IsActive  bool   `json:"is_active"`
	Counter   string `json:"counter"`
}

func WRecordToCounterGet(id int) (WRecordToCounter, error) {
	var r WRecordToCounter
	row := db.QueryRow(`SELECT record_to_counter.*, IFNULL(counter.name, "") FROM record_to_counter
	LEFT JOIN counter ON record_to_counter.counter_id = counter.id WHERE record_to_counter.id=?`, id)
	err := row.Scan(
		&r.Id,
		&r.CounterId,
		&r.CreatedAt,
		&r.Number,
		&r.IsActive,
		&r.Counter,
	)
	return r, err
}

func WRecordToCounterGetAll(withDeleted bool, deletedOnly bool) ([]WRecordToCounter, error) {
	query := `SELECT record_to_counter.*, IFNULL(counter.name, "") FROM record_to_counter
	LEFT JOIN counter ON record_to_counter.counter_id = counter.id`
	if deletedOnly {
		query += "  WHERE record_to_counter.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE record_to_counter.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WRecordToCounter{}
	for rows.Next() {
		var r WRecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
			&r.Counter,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func WRecordToCounterGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WRecordToCounter, error) {

	if !RecordToCounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT record_to_counter.*, IFNULL(counter.name, "") FROM record_to_counter
	LEFT JOIN counter ON record_to_counter.counter_id = counter.id WHERE record_to_counter.%s=?`, field)
	if deletedOnly {
		query += "  AND record_to_counter.is_active = 0"
	} else if !withDeleted {
		query += "  AND record_to_counter.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WRecordToCounter{}
	for rows.Next() {
		var r WRecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
			&r.Counter,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil

}

func WRecordToCounterGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WRecordToCounter, error) {

	if !RecordToCounterTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT record_to_counter.*, IFNULL(counter.name, "") FROM record_to_counter
	LEFT JOIN counter ON record_to_counter.counter_id = counter.id WHERE record_to_counter.%s=?`, field)
	if deletedOnly {
		query += "  AND record_to_counter.is_active = 0"
	} else if !withDeleted {
		query += "  AND record_to_counter.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WRecordToCounter{}
	for rows.Next() {
		var r WRecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
			&r.Counter,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil

}

func WRecordToCounterGetBetweenCreatedAt(created_at1, created_at2 string, withDeleted bool, deletedOnly bool) ([]WRecordToCounter, error) {
	query := `SELECT record_to_counter.*, IFNULL(counter.name, "") FROM record_to_counter
	LEFT JOIN counter ON record_to_counter.counter_id = counter.id WHERE (record_to_counter.created_at BETWEEN ? AND ?)`
	if deletedOnly {
		query += "  AND record_to_counter.is_active = 0"
	} else if !withDeleted {
		query += "  AND record_to_counter.is_active = 1"
	}

	rows, err := db.Query(query, created_at1, created_at2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WRecordToCounter{}
	for rows.Next() {
		var r WRecordToCounter
		if err := rows.Scan(
			&r.Id,
			&r.CounterId,
			&r.CreatedAt,
			&r.Number,
			&r.IsActive,
			&r.Counter,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

type WWmcNumber struct {
	Id          int     `json:"id"`
	WhsId       int     `json:"whs_id"`
	MatherialId int     `json:"matherial_id"`
	ColorId     int     `json:"color_id"`
	Total       float64 `json:"total"`
	IsActive    bool    `json:"is_active"`
	Whs         string  `json:"whs"`
	Matherial   string  `json:"matherial"`
	Color       string  `json:"color"`
}

func WWmcNumberGet(id int) (WWmcNumber, error) {
	var w WWmcNumber
	row := db.QueryRow(`SELECT wmc_number.*, IFNULL(whs.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, "") FROM wmc_number
	LEFT JOIN whs ON wmc_number.whs_id = whs.id
	LEFT JOIN matherial ON wmc_number.matherial_id = matherial.id
	LEFT JOIN color ON wmc_number.color_id = color.id WHERE wmc_number.id=?`, id)
	err := row.Scan(
		&w.Id,
		&w.WhsId,
		&w.MatherialId,
		&w.ColorId,
		&w.Total,
		&w.IsActive,
		&w.Whs,
		&w.Matherial,
		&w.Color,
	)
	return w, err
}

func WWmcNumberGetAll(withDeleted bool, deletedOnly bool) ([]WWmcNumber, error) {
	query := `SELECT wmc_number.*, IFNULL(whs.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, "") FROM wmc_number
	LEFT JOIN whs ON wmc_number.whs_id = whs.id
	LEFT JOIN matherial ON wmc_number.matherial_id = matherial.id
	LEFT JOIN color ON wmc_number.color_id = color.id`
	if deletedOnly {
		query += "  WHERE wmc_number.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE wmc_number.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWmcNumber{}
	for rows.Next() {
		var w WWmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
			&w.Whs,
			&w.Matherial,
			&w.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil
}

func WWmcNumberGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WWmcNumber, error) {

	if !WmcNumberTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT wmc_number.*, IFNULL(whs.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, "") FROM wmc_number
	LEFT JOIN whs ON wmc_number.whs_id = whs.id
	LEFT JOIN matherial ON wmc_number.matherial_id = matherial.id
	LEFT JOIN color ON wmc_number.color_id = color.id WHERE wmc_number.%s=?`, field)
	if deletedOnly {
		query += "  AND wmc_number.is_active = 0"
	} else if !withDeleted {
		query += "  AND wmc_number.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWmcNumber{}
	for rows.Next() {
		var w WWmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
			&w.Whs,
			&w.Matherial,
			&w.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

func WWmcNumberGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WWmcNumber, error) {

	if !WmcNumberTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT wmc_number.*, IFNULL(whs.name, ""), IFNULL(matherial.name, ""), IFNULL(color.name, "") FROM wmc_number
	LEFT JOIN whs ON wmc_number.whs_id = whs.id
	LEFT JOIN matherial ON wmc_number.matherial_id = matherial.id
	LEFT JOIN color ON wmc_number.color_id = color.id WHERE wmc_number.%s=?`, field)
	if deletedOnly {
		query += "  AND wmc_number.is_active = 0"
	} else if !withDeleted {
		query += "  AND wmc_number.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WWmcNumber{}
	for rows.Next() {
		var w WWmcNumber
		if err := rows.Scan(
			&w.Id,
			&w.WhsId,
			&w.MatherialId,
			&w.ColorId,
			&w.Total,
			&w.IsActive,
			&w.Whs,
			&w.Matherial,
			&w.Color,
		); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, nil

}

type WNumbersToProduct struct {
	Id        int     `json:"id"`
	ProductId int     `json:"product_id"`
	Number    float64 `json:"number"`
	Pieces    int     `json:"pieces"`
	Size      float64 `json:"size"`
	Persent   float64 `json:"persent"`
	IsActive  bool    `json:"is_active"`
	Product   string  `json:"product"`
}

func WNumbersToProductGet(id int) (WNumbersToProduct, error) {
	var n WNumbersToProduct
	row := db.QueryRow(`SELECT numbers_to_product.*, IFNULL(product.name, "") FROM numbers_to_product
	LEFT JOIN product ON numbers_to_product.product_id = product.id WHERE numbers_to_product.id=?`, id)
	err := row.Scan(
		&n.Id,
		&n.ProductId,
		&n.Number,
		&n.Pieces,
		&n.Size,
		&n.Persent,
		&n.IsActive,
		&n.Product,
	)
	return n, err
}

func WNumbersToProductGetAll(withDeleted bool, deletedOnly bool) ([]WNumbersToProduct, error) {
	query := `SELECT numbers_to_product.*, IFNULL(product.name, "") FROM numbers_to_product
	LEFT JOIN product ON numbers_to_product.product_id = product.id`
	if deletedOnly {
		query += "  WHERE numbers_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  WHERE numbers_to_product.is_active = 1"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WNumbersToProduct{}
	for rows.Next() {
		var n WNumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
			&n.Product,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func WNumbersToProductGetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]WNumbersToProduct, error) {

	if !NumbersToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT numbers_to_product.*, IFNULL(product.name, "") FROM numbers_to_product
	LEFT JOIN product ON numbers_to_product.product_id = product.id WHERE numbers_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND numbers_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND numbers_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WNumbersToProduct{}
	for rows.Next() {
		var n WNumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
			&n.Product,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil

}

func WNumbersToProductGetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]WNumbersToProduct, error) {

	if !NumbersToProductTestForExistingField(field) {
		return nil, errors.New("field not exist")
	}
	query := fmt.Sprintf(`SELECT numbers_to_product.*, IFNULL(product.name, "") FROM numbers_to_product
	LEFT JOIN product ON numbers_to_product.product_id = product.id WHERE numbers_to_product.%s=?`, field)
	if deletedOnly {
		query += "  AND numbers_to_product.is_active = 0"
	} else if !withDeleted {
		query += "  AND numbers_to_product.is_active = 1"
	}
	rows, err := db.Query(query, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []WNumbersToProduct{}
	for rows.Next() {
		var n WNumbersToProduct
		if err := rows.Scan(
			&n.Id,
			&n.ProductId,
			&n.Number,
			&n.Pieces,
			&n.Size,
			&n.Persent,
			&n.IsActive,
			&n.Product,
		); err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil

}
