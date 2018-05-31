package main

import (
	"log"

	"database/sql"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Expression struct {
	Id       int
	Name     string
	ExpVal   string
	Remarks  string
	Axisname string
	Apply    bool
}

type Species struct {
	Id   int
	Name string
	hz   int
}

func KnownExpressions2() []*Species {
	return []*Species{
		{1, "Expression1", 3},
		{2, "Expression2", 2},
		{3, "Expression3", 2},
		{4, "Expression4", 1},
		{5, "Expression5", 0},
	}
}

/*

func ExpressionNew(exp Expression) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	// stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,?,?,?)")
	stmt, err := db.Prepare("INSERT INTO expressions SET name=?, exp=?, axisname=?, apply=?, comment=? WHERE id=?")

}
*/

func ExpressionUpdate(exp Expression) {

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	stmt, err := db.Prepare("UPDATE expressions SET name=?, exp=?, axisname=?, apply=?, comment=? WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.Id)
	checkErr(err)
	log.Print("Update exp rows!!!")

	db.Close()
}

func ExressionLoad(name string) Expression {
	var exp Expression
	var eid sql.NullInt64
	var ename sql.NullString
	var eexp sql.NullString
	var eaxisname sql.NullString
	var ecomment sql.NullString
	var eapply sql.NullInt64

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	/*
		checkErr(err)
		query, err := db.Prepare("SELECT count(freq) FROM measure_data WHERE measure_id=? ORDER BY freq")

		checkErr(err)

		defer query.Close()

		// Execute query using 'id' and place value into 'output'
		err = query.QueryRow(measureID).Scan(&output)

	*/

	query, err := db.Prepare("SELECT id, name, exp, axisname, apply, comment FROM expressions where name=?")
	checkErr(err)
	// AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	//	log.Print(q)

	//rows, err := db.Query(q)

	err = query.QueryRow(name).Scan(&eid, &ename, &eexp, &eaxisname, &eapply, &ecomment)
	if err == nil {
		exp.Id = int(eid.Int64)
		exp.Name = ename.String
		exp.ExpVal = eexp.String
		exp.Remarks = ecomment.String
		exp.Axisname = eaxisname.String
		exp.Apply = !(eapply.Int64 == 0)
	}

	db.Close()
	return exp
}

func KnownExpressions(dbox *DropDownBox) []Species {

	var exps []Species
	var eid sql.NullInt64
	var ename sql.NullString
	var exp Species

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)

	var q = "SELECT id, name FROM expressions"
	// AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	log.Print(q)

	rows, err := db.Query(q)
	defer rows.Close()

	//	return pts, idx

	//	var nans = 0
	dbox.model.Clean()
	for rows.Next() {

		err = rows.Scan(&eid, &ename)
		exp = Species{
			Name: ename.String,
			Id:   int(eid.Int64),
		}
		exps = append(exps, exp)
		dbox.model.Add(ename.String)
	}

	db.Close()
	return exps
}

func RunExpressionDialog(owner walk.Form, animal *Expression) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	log.Print("RunExpressionDialog!!!")

	return Dialog{
		AssignTo:      &dlg,
		Title:         "Exression Editor",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			DataSource:     animal,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{640, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Name:",
					},
					LineEdit{
						Text: Bind("Name"),
					},

					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},

					Label{
						//						ColumnSpan: 2,
						Text: "Expression:",
					},
					TextEdit{
						//						ColumnSpan: 2,
						MinSize: Size{100, 50},
						Text:    Bind("ExpVal"),
					},

					Label{
						Text: "Y axis name:",
					},
					LineEdit{
						Text: Bind("Axisname"),
					},

					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},

					Label{
						ColumnSpan: 2,
						Text:       "Remarks:",
					},
					TextEdit{
						ColumnSpan: 2,
						MinSize:    Size{100, 50},
						Text:       Bind("Remarks"),
					},

					Label{
						Text: "Apply to graph:",
					},
					CheckBox{
						Checked: Bind("Apply"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
