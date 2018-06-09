package main

import (
	"database/sql"
	"strconv"
	//	"github.com/alfredxing/calc/operators"
	//	"github.com/alfredxing/calc/operators"
	//	"github.com/alfredxing/calc/operators/functions"
	//	"github.com/alfredxing/calc/operators/functions"
)

type Measure struct {
	Id              int
	hash            string
	name            string
	date            string
	fname           string
	points          int
	comment         string
	field1          string
	field2          string
	exression_id    int
	parent_dataset1 int
	parent_dataset2 int
}

func MeasureLoad(measureId int) (ms Measure) {
	var hash sql.NullString
	var name sql.NullString
	var date sql.NullString
	var fname sql.NullString
	var comment sql.NullString
	var field1 sql.NullString
	var field2 sql.NullString
	var expId int
	var pds1 int
	var pds2 int
	var points int

	var count = 0
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	rows, err := db.Query("SELECT hash,  name, date, fname, points, comment, field1, field2, exression_id, parent_dataset1, parent_dataset2 FROM measures WHERE measure_id=" + strconv.Itoa(measureId))
	checkErr(err)
	defer rows.Close(&hash, &name, &date, &fname, &points, &comment, &field1, &field2, &expId, &pds1, &pds2)
	rows.Scan(&count)

	var m Measure
	m.Id = measureId
	m.hash = hash.String()
	m.date = date.String()
	m.fname = fname.String()
	m.points = points
	m.comment = comment.String()
	m.field1 = field1.String()
	m.field2 = field2.String()
	m.exression_id = expId
	m.parent_dataset1 = pds1
	m.parent_dataset2 = pds2
	return m

}
