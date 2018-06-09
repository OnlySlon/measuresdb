package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"database/sql"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
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

func DxpressionDelete(name string) {

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	defer db.Close()
	// stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,?,?,?)")
	stmt, err := db.Prepare("DELETE FROM expressions WHERE name=?")
	checkErr(err)
	stmt.Exec(name)

}

func ExpressionNew(exp Expression) {

	//	log.Print("ExpressionNew!!!")
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,?,?,?)")
	stmt, err := db.Prepare("INSERT INTO expressions(name, exp, axisname, apply, comment, exportname, dbres, graph, dbgraph) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ")

	checkErr(err)

	_, err = stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph, exp.DbGraph)

	checkErr(err)
	//	iid, _ := res.LastInsertId()
	//	log.Print("New Expression -------- OK!!!! ID:", iid)

}

func ExpressionUpdate(exp Expression) {

	//	log.Print("EEEEXP UPDATE:", exp)
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)
	stmt, err := db.Prepare("UPDATE expressions SET name=?, exp=?, axisname=?, apply=?, comment=?, exportname=?, dbres=?, graph=?, dbgraph=? WHERE id=?")
	checkErr(err)

	_, err = stmt.Exec(exp.Name, exp.ExpVal, exp.Axisname, exp.Apply, exp.Remarks, exp.ExportName, exp.DbRes, exp.Graph, exp.DbGraph, exp.Id)
	checkErr(err)
	//	log.Print("Update exp rows!!!")

	db.Close()
}

func MasterMeasureFix() {
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

}

//mw *MyMainWindow
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

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)

	// ------------ Load all graph for drawing ---------------
	rows, err := db.Query("SELECT id, graph, exp, exportname FROM expressions WHERE apply=1")

	//	return pts, idx g
	graphs := make(map[string][]int64)

	//	var nans = 0
	for rows.Next() {
		err = rows.Scan(&eid, &egraph, &exp, &exportname)
		if len(egraph.String) > 0 {
			//			log.Print("Begin draw graph '" + egraph.String + "'")
			graphs[egraph.String] = append(graphs[egraph.String], eid.Int64)
		}
	}
	rows.Close()
	//	log.Print(graphs)

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
			//			log.Print("EXP " + exportname.String + "=" + exp.String)
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

		cnt++
		//		log.Print("-------")
		for expid := range graphs[graph] {
			//			log.Print("Dooo ", graphs[graph][expid])
			gid := graphs[graph][expid]
			exp := ExressionLoad("id=" + strconv.Itoa(int(gid)))
			p.Title.Text += exp.Name + " "
			p.Y.Label.Text = exp.Axisname
			// Apply expression to all datasets
			for i := range model.items {
				if model.items[i].checked && model.items[i].Index != MasterMeasure {
					vs = append(vs, model.items[i].Name)
					pts, _ := dbPointsExpression(MasterMeasure, model.items[i].Index, exp.ExpVal, exp.DbGraph) // load all points from db and apply expression
					vs = append(vs, pts)
					//					log.Print("Records: ", cnt)
				}
			}
			GraphSetTheme(p, exp.DbGraph)
		}

		//		log.Print("Draw graph...")
		err = plotutil.AddLines(p,
			vs...,
		)
		if err := p.Save(vg.Points(float64(imgW))*0.75, vg.Points(float64(imgH))*0.72, fname); err != nil {
			panic(err)
		}
		openImage(mw, fname, graph)

	}

	db.Close()

}

func ExressionLoad(where string) Expression {
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
	//	log.Print("ExressionLoad load: " + where)
	query, err := db.Prepare("SELECT id, name, exp, axisname, apply, comment, dbres, exportname, graph, dbgraph FROM expressions WHERE " + where)
	checkErr(err)
	// AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	//	log.Print(q)

	//rows, err := db.Query(q)

	err = query.QueryRow().Scan(&eid, &ename, &eexp, &eaxisname, &eapply, &ecomment, &dbres, &exportname, &graph, &dbgraph)
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
	//	log.Print(q)

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

func DatasetFromExression(applyexp Expression) {
	//	var exp Expression
	var exportname sql.NullString
	var exp sql.NullString
	MasterMeasureFix()

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT exportname, exp FROM expressions GROUP BY exportname")
	MyEvalClear()
	for rows.Next() {
		err = rows.Scan(&exportname, &exp)

		// func dbPointsExpression(ds1 int, ds2 int, myexp string) (plotter.XYs, int) {
		if len(exportname.String) > 0 {
			//			log.Print("EXP " + exportname.String + "=" + exp.String)
			MyEvalAdd(exportname.String, exp.String)
		}
	}

	for i := range model.items {
		if model.items[i].checked {
			log.Print(model.items[i].Name, model.items[i].Index)
		}
	}

	for i := range model.items {
		if model.items[i].checked && model.items[i].Index != MasterMeasure {
			pts, _ := dbPointsExpression(MasterMeasure, model.items[i].Index, applyexp.ExpVal, applyexp.DbGraph) // load all points from db and apply expression
			//					log.Print("Records: ", cnt)
			log.Print(len(pts))

			if len(pts) > 0 {
				stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points, comment, field1, field2, exression_id, parent_dataset1, parent_dataset2 ) values(?,?,?,?,?,?,?,?,?,?)")
				checkErr(err)

				log.Print("Insert new MEASURE: hash=" + hash + " name=" + name + " filename=" + filename)
				res, err := stmt.Exec("", name, ts, filename, records_cnt_mag)

			}

		}
	}

}

func GraphSetTheme(p *plot.Plot, db bool) {
	p.Add(MyNewGrid())
	p.X.LineStyle = draw.LineStyle{
		Color: color.White,
		Width: vg.Points(0.5),
	}
	p.X.Tick.Width = 1
	p.X.Tick.Marker = readableMhz(p.X.Tick.Marker)
	if db {
		p.Y.Tick.Marker = readabledB(p.Y.Tick.Marker)
	} else {
		p.Y.Tick.Marker = markerY(p.Y.Tick.Marker)
	}

	p.X.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Label.Font.Size = 15
	p.Y.Label.Font.Size = 15
	p.Y.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Label.Font.Size = 14
	p.X.Tick.Label.Font.Size = 10
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

func expressionTest() {
	/*
		settings := walk.NewIniFileSettings("settings.ini")
		if err := settings.Load(); err != nil {
			log.Fatal(err)
		} else {
			log.Print(settings.Get("testzz"))
		}
	*/
	// ------------------------register math functions
	//Z=Cos(PhA)

	SetMagPhase(A, 11, 22)
	SetMagPhase(B, 33, 46)

	//res, err := MyEvaluate("phdelta(PhD(A) PhD(B))")
	res, err := MyEvaluate("sum(1 2)")
	log.Print("--------------------------------------->>>>>>>>>>>>>>>>>>>>>>>>")

	if err != nil {
		log.Print("Error: " + err.Error())
		return
	}
	fmt.Printf("%s\n", strconv.FormatFloat(res, 'G', -1, 64))

	log.Print(res)

	res, err = MyEvaluate("phdelta(Phd(A),Phd(B))")

	if err != nil {
		log.Print("Error: " + err.Error())
		return
	}
	fmt.Printf("%s\n", strconv.FormatFloat(res, 'G', -1, 64))

	log.Print(res)
}
