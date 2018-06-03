package main

import (
	"image/color"
	"log"
	"strconv"

	"database/sql"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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
	DbGraph    bool
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
	stmt, err := db.Prepare("INSERT INTO expressions(name, exp, axisname, apply, comment, exportname, dbres, graph, dbgraph) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ")

	checkErr(err)

	res, err := stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph, exp.DbGraph)

	checkErr(err)
	iid, _ := res.LastInsertId()
	log.Print("New Expression -------- OK!!!! ID:", iid)

}

func ExpressionUpdate(exp Expression) {

	log.Print("EEEEXP UPDATE:", exp)
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)
<<<<<<< HEAD
	stmt, err := db.Prepare("UPDATE expressions SET name=?, exp=?, axisname=?, apply=?, comment=?, exportname=?, dbres=?, graph=?, dbgraph=? WHERE id=?")
=======
	stmt, err := db.Prepare("UPDATE expressions SET name=?, exp=?, axisname=?, apply=?, comment=?, exportname=?, dbres=?, graph=? WHERE id=?")
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
	checkErr(err)

	_, err = stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph, exp.DbGraph, exp.Id)
	checkErr(err)
	log.Print("Update exp rows!!!")

	db.Close()
}

//mw *MyMainWindow
<<<<<<< HEAD
func ExpressionDraw(mw *MyMainWindow) {
	var egraph sql.NullString
	var eid sql.NullInt64
	var exp sql.NullString
	var exportname sql.NullString

	// Master match handler
	masterMatch := false
	fistCheck := -1
	for i := range model.items {
		if model.items[i].checked {
			if fistCheck == -1 {
				fistCheck = model.items[i].Index
			}
			if model.items[i].Index == MasterMeasure {
				masterMatch = true
			}
		}
	}
	if fistCheck == -1 {
		log.Print("No datasets found")
		return
	}
	if masterMatch == false {
		MasterMeasure = fistCheck
	}
	// ----------------------------------------------------

=======
func ExpressionDraw() {
	var egraph sql.NullString
	var eid sql.NullInt64
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)

<<<<<<< HEAD
	// ------------ Load all graph for drawing ---------------
	rows, err := db.Query("SELECT id, graph, exp, exportname FROM expressions WHERE apply=1")

	//	return pts, idx g
=======
	//	query, err := db.Prepare("SELECT Graph from expressions where apply=1 GROUP BY Graph;")
	//	checkErr(err)

	rows, err := db.Query("SELECT id, graph from expressions where apply=1")

	//	return pts, idx
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
	graphs := make(map[string][]int64)

	//	var nans = 0
	for rows.Next() {
<<<<<<< HEAD
		err = rows.Scan(&eid, &egraph, &exp, &exportname)
=======
		err = rows.Scan(&eid, &egraph)
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
		if len(egraph.String) > 0 {
			log.Print("Begin draw graph '" + egraph.String + "'")
			graphs[egraph.String] = append(graphs[egraph.String], eid.Int64)
		}
	}
	rows.Close()
	log.Print(graphs)

<<<<<<< HEAD
	/*
			MyEvalAdd("B", "5+5")
		MyEvalAdd("A", "sin(30)+10 + B")
	*/
	// --------- Load all exported values to compute module
	rows, err = db.Query("SELECT exportname, exp FROM expressions GROUP BY exportname")
	MyEvalClear()
	for rows.Next() {
		err = rows.Scan(&exportname, &exp)

		// func dbPointsExpression(ds1 int, ds2 int, myexp string) (plotter.XYs, int) {
		if len(exportname.String) > 0 {
			log.Print("EXP " + exportname.String + "=" + exp.String)
			MyEvalAdd(exportname.String, exp.String)
		}
	}

	//var fname string
	var cnt = 0
	for graph := range graphs {
		vs := []interface{}{}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		var fname = "./" + graph + ".png"
		p.Title.Text = ""
		p.X.Label.Text = "Frequency"

		p.Y.Label.Text = "Magnitude (db)"
		p.Add(plotter.NewGrid())
		GraphSetTheme(p)
		cnt++
		log.Print("-------")
		for expid := range graphs[graph] {
			log.Print("Dooo ", graphs[graph][expid])
			gid := graphs[graph][expid]
			exp := ExressionLoad("id=" + strconv.Itoa(int(gid)))
			p.Title.Text += exp.Name + " "
			p.Y.Label.Text = exp.Axisname
			// Apply expression to all datasets
			for i := range model.items {
				if model.items[i].checked && model.items[i].Index != MasterMeasure {
					vs = append(vs, model.items[i].Name)
					pts, cnt := dbPointsExpression(MasterMeasure, model.items[i].Index, exp.ExpVal, exp.DbGraph) // load all points from db and apply expression
					vs = append(vs, pts)
					log.Print("Records: ", cnt)
				}
			}
		}
		log.Print("Draw graph...")
		err = plotutil.AddLines(p,
			vs...,
		)
		if err := p.Save(vg.Points(float64(imgW))*0.75, vg.Points(float64(imgH))*0.74, fname); err != nil {
			panic(err)
		}
		openImage(mw, fname, graph)
=======
	//var fname string
	for graph := range graphs {
		/*
			p, err := plot.New()
			if err != nil {
				panic(err)
			}
		*/
		log.Print("-------")
		for expid := range graphs[graph] {

			log.Print("Dooo ", graphs[graph][expid])
			gid := graphs[graph][expid]
			exp := ExressionLoad(strconv.Itoa(int(gid)), "id")

		}
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78

	}

	db.Close()

}

<<<<<<< HEAD
func ExressionLoad(where string) Expression {
=======
func ExressionLoad(name string, by string) Expression {
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
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
	var dbgraph sql.NullInt64

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
<<<<<<< HEAD
	log.Print("ExressionLoad load: " + where)
	query, err := db.Prepare("SELECT id, name, exp, axisname, apply, comment, dbres, exportname, graph, dbgraph FROM expressions WHERE " + where)
=======

	query, err := db.Prepare("SELECT id, name, exp, axisname, apply, comment, dbres, exportname, graph FROM expressions where ?=?")
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
	checkErr(err)
	// AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	//	log.Print(q)

	//rows, err := db.Query(q)

<<<<<<< HEAD
	err = query.QueryRow().Scan(&eid, &ename, &eexp, &eaxisname, &eapply, &ecomment, &dbres, &exportname, &graph, &dbgraph)
=======
	err = query.QueryRow(by, name).Scan(&eid, &ename, &eexp, &eaxisname, &eapply, &ecomment, &dbres, &exportname, &graph)
>>>>>>> 048cbc8221098281f5de2cd0c702811afedb6e78
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
		exp.DbGraph = !(dbgraph.Int64 == 0)
	} else {
		log.Print("Can't load expression "+where, err)
		exp.Id = -1
		return exp
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
						Text: "Graph in bB:",
					},
					CheckBox{
						Checked: Bind("DbGraph"),
					},
					Label{
						Text: "Gtaph Name:",
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

func GraphSetTheme(p *plot.Plot) {
	p.X.Tick.Width = 2
	p.X.Tick.Marker = readableDuration(p.X.Tick.Marker)
	p.X.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Label.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	p.Y.Tick.Label.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	p.BackgroundColor = color.RGBA{R: 55, G: 58, B: 60, A: 255}
	p.Title.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Title.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Legend.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
}
