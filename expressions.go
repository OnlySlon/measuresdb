package main

import (
	"log"

	"database/sql"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Expression struct {
	Id         int
	Name       string // just name
	ExpVal     string // evaluation value
	Remarks    string // evaluatin remarks
	Axisname   string // Y axis name
	Apply      bool   // Create graph from this eval
	ExportName string // export value
	DbRes      bool   // result in dB
	Graph      string // put this to graph
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

func ExpressionNew(exp Expression) {

	log.Print("ExpressionNew!!!")
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,?,?,?)")
	stmt, err := db.Prepare("INSERT INTO expressions(name, exp, axisname, apply, comment, exportname, dbres, graph) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ")

	checkErr(err)

	res, err := stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph)

	checkErr(err)
	iid, _ := res.LastInsertId()
	log.Print("New Expression -------- OK!!!! ID:", iid)

}

func ExpressionUpdate(exp Expression) {

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	stmt, err := db.Prepare("UPDATE expressions SET name=?, exp=?, axisname=?, apply=?, comment=?, exportname=?, dbres=?, graph=? WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph, exp.Id)
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
	var dbres sql.NullInt64
	var exportname sql.NullString
	var graph sql.NullString

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

	query, err := db.Prepare("SELECT id, name, exp, axisname, apply, comment, dbres, exportname, graph FROM expressions where name=?")
	checkErr(err)
	// AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	//	log.Print(q)

	//rows, err := db.Query(q)

	err = query.QueryRow(name).Scan(&eid, &ename, &eexp, &eaxisname, &eapply, &ecomment, &dbres, &exportname, &graph)
	if err == nil {
		exp.Id = int(eid.Int64)
		exp.Name = ename.String
		exp.ExpVal = eexp.String
		exp.Remarks = ecomment.String
		exp.Axisname = eaxisname.String
		exp.Apply = !(eapply.Int64 == 0)
		exp.DbRes = !(dbres.Int64 == 0)
		exp.ExportName = exportname.String
		exp.Graph = graph.String
	} else {
		log.Print("Can't load expression " + name)
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
						MinSize: Size{100, 50},
						Text:    Bind("ExpVal"),
					},

					Label{
						Text: "Result is dB:",
					},
					CheckBox{
						Checked: Bind("DbRes"),
					},

					Label{
						Text: "Export name:",
					},
					LineEdit{
						Text: Bind("ExportName"),
					},

					Label{
						Text: "Y axis name:",
					},
					LineEdit{
						Text: Bind("Axisname"),
					},

					Label{
						Text: "Apply to graph:",
					},
					CheckBox{
						Checked: Bind("Apply"),
					},
					Label{
						Text: "Graph Name:",
					},
					LineEdit{
						Text: Bind("Graph"),
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
