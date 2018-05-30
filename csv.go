package main

// package main

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alfredxing/calc/constants"
	"github.com/alfredxing/calc/operators"
	"github.com/alfredxing/calc/operators/functions"

	//	. "github.com/alfredxing/calc/operators"
	"github.com/alfredxing/calc/compute"

	"github.com/lxn/win"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	_ "github.com/mattn/go-sqlite3"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

/*
	#include "./c/dirmon.c"
*/
import "C"

// 	_ "runtime/cgo"

const GraphMag = 0
const GraphPhase = 1
const GraphPhaseDiff = 2
const GraphRdb = 3

var imgW int64
var imgH int64
var MasterMeasure = -1

type measure struct {
	freq float64
	mag  float64
	deg  float64
}

type Foo struct {
	Index  int
	Bar    string
	Points int64

	Baz  float64
	Quux time.Time

	Name    string
	Date    time.Time
	Comment string

	checked bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

var model *FooModel

func hash_file_md5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}

func dbCountMeasures() int {
	var output string
	//	sql.Register("sqlite3", &SQLiteDriver{})
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	//	query, err := db.Prepare("SELECT count(measure_id) FROM measures")

	//checkErr(err)

	//defer query.Close()

	// Execute query using 'id' and place value into 'output'
	//err = query.QueryRow().Scan(&output)
	row := db.QueryRow("SELECT count(*) FROM measures")
	checkErr(err)
	row.Scan(&output)
	checkErr(err)

	// Catch errors
	switch {
	case err == sql.ErrNoRows:
		db.Close()
		return 0
	case err != nil:
		fmt.Printf("%s", err)
	default:
		{
			cnt, _ := strconv.Atoi(output)
			db.Close()
			return cnt
		}

	}
	return 0
}

func dbPointsCount(measureID int) int {
	var output string
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	query, err := db.Prepare("SELECT count(freq) FROM measure_data WHERE measure_id=? ORDER BY freq")

	checkErr(err)

	defer query.Close()

	// Execute query using 'id' and place value into 'output'
	err = query.QueryRow(measureID).Scan(&output)

	// Catch errors
	switch {
	case err == sql.ErrNoRows:
		db.Close()
		return 0
	case err != nil:
		fmt.Printf("%s", err)
	default:
		{
			cnt, _ := strconv.Atoi(output)
			db.Close()
			return cnt
		}

	}
	return 0
}

func dbPointsExpression(ds1 int, ds2 int, myexp string) (plotter.XYs, int) {
	//var myexp = "log(sqrt((MagA^2+MagB^2+sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))/(MagA^2+MagB^2-sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))))"
	//	var myexp = "log(sqrt((MagA ^ 2 + MagB ^ 2 + sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))) / (MagA ^ 2 + MagB ^ 2 - sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB))))))"
	var idx = 0
	// log(sqrt((MagA ^ 2 + MagB ^ 2 + sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))) / (MagA ^ 2 + MagB ^ 2 - sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB))))))
	//	myexp := "MagA ^ 2 + MagB ^ 2 - sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))"
	log.Print("dbPointsExpression DS1=", ds1, " DS2=", ds2, " EXP="+myexp)
	// xxxxx

	pts := make(plotter.XYs, dbPointsCount(ds1))

	var (
		cPh1 = &constants.Constant{
			Name:  "PhA",
			Value: 0,
		}
		cMag1 = &constants.Constant{
			Name:  "MagA",
			Value: 0,
		}
		cPh2 = &constants.Constant{
			Name:  "PhB",
			Value: 0,
		}
		cMag2 = &constants.Constant{
			Name:  "MagB",
			Value: 0,
		}
	)
	var freq int64
	var mag1 float64
	var ph1 float64
	var mag2 float64
	var ph2 float64

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)

	var q = "SELECT m1.freq, m1.magnitude, m1.degrees, m2.magnitude, m2.degrees FROM measure_data as m1 LEFT join measure_data as m2 WHERE m1.freq=m2.freq AND m1.measure_id=" + strconv.Itoa(ds1) + " AND m2.measure_id=" + strconv.Itoa(ds2) + " order by m1.freq"
	log.Print(q)

	rows, err := db.Query(q)
	defer rows.Close()
	defer db.Close()

	//	return pts, idx

	var nans = 0
	for rows.Next() {

		err = rows.Scan(&freq, &mag1, &ph1, &mag2, &ph2)
		cMag1.Value = mag1
		cMag2.Value = mag2
		cPh1.Value = ph1
		cPh2.Value = ph2
		constants.Register(cMag1)
		constants.Register(cMag2)
		constants.Register(cPh1)
		constants.Register(cPh2)

		res, err := compute.Evaluate(myexp)
		if err != nil {
			log.Print("Error: " + err.Error())
			return nil, 0
		} else {
			//			log.Print("Freq=", freq, " Res=", res)
			pts[idx].X = float64(freq)
			if math.IsNaN(res) {
				log.Print("NAN RESULT!!! Freq=", freq, " PhA=", ph1, " PhB=", ph2, "  MagA=", mag1, " MagB=", mag2, " res=", res)
				nans++
			} else {
				//				log.Print("OK  RESULT!!! Freq=", freq, " PhA=", ph1, " PhB=", ph2, "  MagA=", mag1, " MagB=", mag2, " res=", res)
			}

			pts[idx].Y = res
			idx++
		}

		//		log.Print("Freq=", freq, " Ph1=", ph1, " Ph2=", ph2, " res=", res)

	}
	log.Print("Nan Results: ", nans)
	return pts, idx
}

func dbPoints(measureID int, tpe int) plotter.XYs {
	//	pts := make(plotter.XYs, n)
	var count = 0
	count = dbPointsCount(measureID)
	//	fmt.Println("count=", count)

	//	log.Print("Requesting points type=" + strconv.Itoa(tpe) + " MeasureID:" + strconv.Itoa(measureID))
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	rows, err := db.Query("SELECT freq, magnitude,degrees FROM measure_data WHERE measure_id=" + strconv.Itoa(measureID) + " ORDER BY freq")
	checkErr(err)

	rows.Scan(&count)

	pts := make(plotter.XYs, count)

	var pidx = 0
	for rows.Next() {
		var freq int64
		var magnitude float64
		var degrees float64
		err = rows.Scan(&freq, &magnitude, &degrees)
		//		fmt.Println("freq=", freq, " ", magnitude, " ", pidx)
		//		pts[pidx].X = float64(freq)
		//		pts[pidx].Y = magnitude
		pts[pidx].X = float64(freq)
		if tpe == GraphMag {
			pts[pidx].Y = magnitude
		} else {
			pts[pidx].Y = degrees
		}

		pidx++
		checkErr(err)
	}
	rows.Close()
	db.Close()
	return pts
}

type tickerFunc func(min, max float64) []plot.Tick

func (tkfn tickerFunc) Ticks(min, max float64) []plot.Tick { return tkfn(min, max) }

func readableDuration(marker plot.Ticker) plot.Ticker {
	return tickerFunc(func(min, max float64) []plot.Tick {
		var out []plot.Tick
		for _, t := range marker.Ticks(min, max) {
			t.Label = strconv.FormatInt(int64(t.Value)/1000000, 10) + " MHz"
			//			log.Println("Min=", min, " ", max, " ", int64(t.Value)/1000000, " ", out)

			// time.Duration(t.Value).String()
			out = append(out, t)
		}
		return out
	})
}

/*
type LogTicks struct{}

var _ Ticker = LogTicks2{}

func (LogTicks2) Ticks(min, max float64) []Tick {
	if min <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}

	val := math.Pow10(int(math.Log10(min)))
	max = math.Pow10(int(math.Ceil(math.Log10(max))))
	var ticks []Tick
	for val < max {
		for i := 1; i < 10; i++ {
			if i == 1 {
				ticks = append(ticks, Tick{Value: val, Label: formatFloatTick(val, -1)})
			}
			ticks = append(ticks, Tick{Value: val * float64(i)})
		}
		val *= 10
	}
	ticks = append(ticks, Tick{Value: val, Label: formatFloatTick(val, -1)})

	return ticks
}
*/

func drawModel(mw *MyMainWindow, tpe int) {

	p, err := plot.New()

	if err != nil {
		panic(err)
	}
	if mw.tabWidget.Pages().Len() > 0 {
		//		mw.tabWidget.Pages().RemoveAt(0)
	}

	var fname string
	if tpe == GraphMag {
		fname = "./mag.png"
		p.Title.Text = "Magnitude graph"
		p.X.Label.Text = "Freq"
		p.Y.Label.Text = "Magnitude (db)"

	} else {
		fname = "./phase.png"
		p.Title.Text = "Phase Δ graph"
		p.X.Label.Text = "Freq"
		p.Y.Label.Text = "Degrees"

	}

	//
	//	p.X.Min = 5500000000
	//	p.X.Max = 6500000000
	//p.Y.Min = -29
	//p.Y.Min = -32

	// Verify master record
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

	//	var plottingPointArgs []interface{}
	vs := []interface{}{}
	// ----------------------- MAGNITUDE -------------------------
	if tpe == GraphMag {
		log.Print("-------------------- Magnitude GRAPH---------------------------")
		for i := range model.items {
			if model.items[i].checked {
				log.Print("Add dataset:" + model.items[i].Name)
				vs = append(vs, model.items[i].Name)
				vs = append(vs, dbPoints(model.items[i].Index, tpe))
			}
		}
	}

	// ------------------- PHASE DIFF ----------------------------
	// func dbPointsExpression(ds1 int, ds2 int, exp string) (plotter.XYs, int)
	if tpe == GraphPhase {
		log.Print("-------------------- PHASE GRAPH---------------------------")
		for i := range model.items {
			if model.items[i].checked && model.items[i].Index != MasterMeasure {
				log.Print("------Add Phase dataset:" + model.items[i].Name)
				vs = append(vs, model.items[i].Name)
				pts, cnt := dbPointsExpression(MasterMeasure, model.items[i].Index, "phdelta(PhA PhB)")
				log.Print("Records: ", cnt)
				vs = append(vs, pts)
			}
		}
		//		return
	}

	if tpe == GraphRdb {
		fname = "./r.png"
		p.Title.Text = "r(db) graph"
		p.X.Label.Text = "Freq"
		p.Y.Label.Text = "r(db)"

		log.Print("-------------------- R GRAPH---------------------------")
		for i := range model.items {
			if model.items[i].checked && model.items[i].Index != MasterMeasure {
				log.Print("------Add R dataset:" + model.items[i].Name)
				vs = append(vs, model.items[i].Name)
				pts, cnt := dbPointsExpression(MasterMeasure, model.items[i].Index, "20 * log(sqrt((MagA ^ 2 + MagB ^ 2 + sqrt(MagA^4+MagB^4+2*(MagA^2)*(MagB^2)*cos(2*(PhA-PhB)))) / (MagA ^ 2 + MagB ^ 2 - sqrt(MagA^4+MagB^4+2*(MagA^2)*(MagB^2)*cos(2*(PhA-PhB))))))")
				log.Print("Records: ", cnt)
				vs = append(vs, pts)
			}
		}
		//		return
	}

	log.Print("Draw graph...")
	err = plotutil.AddLines(p,
		vs...,
	)
	if err != nil {
		log.Print(err)
		return
		//panic(err)
	}

	p.X.Tick.Width = 2
	p.X.Tick.Marker = readableDuration(p.X.Tick.Marker)
	//	p.X.Tick.Marker = commaTicks
	//p.Y.Tick.Marker = plot.LogTicks{}
	// plot.ConstantTicks([]plot.Tick{{-31, "0"}, {-30, ""}, {-29, "zzz"}, {75, "-sss"}, {100, "ddd"}})

	p.Add(plotter.NewGrid())

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

	//	p.Save(15*vg.Inch, 8*vg.Inch, "points.pdf")
	// Save the plot to a PNG file.

	log.Print("Draw image:", vg.Points(float64(imgW)), vg.Points(float64(imgH)))
	if err := p.Save(vg.Points(float64(imgW))*0.75, vg.Points(float64(imgH))*0.74, fname); err != nil {
		panic(err)
	}

	openImage(mw, fname, tpe)

}

func draw() {
	// 1783
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//	p.X.Min = 5500000000
	//	p.X.Max = 6500000000
	p.Y.Min = -29
	p.Y.Min = -32
	p.Title.Text = "--------TEST--------"
	p.X.Label.Text = "Freq"
	p.Y.Label.Text = "Magnitude (db)"

	fmt.Println(dbPoints(1677, 1))

	err = plotutil.AddLines(p,
		"Measure1", dbPoints(759, 1),
		"Measure2", dbPoints(760, 1),
	)
	if err != nil {
		panic(err)
	}
	p.X.Tick.Width = 2
	p.X.Tick.Marker = readableDuration(p.X.Tick.Marker)
	//	p.X.Tick.Marker = commaTicks
	//p.Y.Tick.Marker = plot.LogTicks{}
	// plot.ConstantTicks([]plot.Tick{{-31, "0"}, {-30, ""}, {-29, "zzz"}, {75, "-sss"}, {100, "ddd"}})

	p.Add(plotter.NewGrid())

	//	p.Save(15*vg.Inch, 8*vg.Inch, "points.pdf")
	// Save the plot to a PNG file.

	log.Print("Draw image:", vg.Points(float64(imgW)), vg.Points(float64(imgH)))
	//if err := p.Save(vg.Points(float64(imgW))/1, vg.Points(float64(imgH))/1, "points.png"); err != nil {
	// 100x100 132
	//
	if err := p.Save(vg.Points(1452)*0.75, vg.Points(965)*0.75, "points.png"); err != nil {

		panic(err)
	}

}
func checkErr(err error) {
	if err != nil {
		log.Print(err)
		panic(err)
	}
}

func file2db(fpath string) {

	//	extension := filepath.Ext(fpath)
	filename := filepath.Base(fpath)

	log.Print("File2db: " + fpath)
	hash, err := hash_file_md5(fpath)
	if err != nil {
		log.Print("Unable to calculate md5 hash")
		return
	}

	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]

	//	fmt.Println("File: " + filename + "  Hash: " + hash)
	// First check MD5 already exist in db

	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	rows, err := db.Query("SELECT * FROM measures WHERE hash='" + hash + "'")
	checkErr(err)

	defer rows.Close()
	defer db.Close()
	for rows.Next() {
		//log.Print("File with same hash already exist!\n")
		return
	}

	//	stmt, err = db.Prepare("")

	/* ----------------------------------------------- START LOAD DATA */
	f, _ := os.Open(fpath)

	var begins = 0
	var records_cnt_mag = 0
	var records_cnt_phase = 0

	records := make([]measure, 20000)
	var ln = ""
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var point = 0
	var magmode = 0
	ts := time.Now()
	for scanner.Scan() {
		ln = scanner.Text()

		if ln == "BEGIN" {
			point = 0
			begins++
		}

		if ln == "! DATA UNIT dB" {
			magmode = 0
		}

		if ln == "! DATA UNIT Lin Mag" {
			magmode = 1
		}

		if strings.Index(ln, "! TIMESTAMP ") > -1 {
			log.Print("GOT TIMESTAMP!!!!!!" + ln)
			var layout = "! TIMESTAMP Monday, 2 January 2006 15:04:05"

			ts, err = time.Parse(layout, ln)

			if err != nil {
				fmt.Println("ERROR!!!", err)
			}

		}
		split := strings.Split(ln, ",")
		//		fmt.Printf("%q LENSPLIT=%u\n", strings.Split(ln, ","), len(split))
		if len(split) == 2 {
			f := split[0]
			val := split[1]
			//			fmt.Println(f)
			//			fmt.Println(val)
			freq, err := strconv.ParseFloat(f, 64)
			if err != nil {
				continue
			}
			mag, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			//			log.Println(point)
			records[point].freq = freq
			if records[point].mag != 0 {
				records[point].deg = mag
				records_cnt_mag++
			} else {
				if magmode == 1 {
					// LINEAR
					records[point].mag = mag // 10 * math.Log10(mag)
					//					log.Print("Magmode 1 in=", mag, " out=", records[point].mag)
				} else {
					// DB
					records[point].mag = math.Pow(10, mag) / 10
				}

				records_cnt_phase++
			}
			point++
			//records = append(records, s)
			//			fmt.Printf("%i ", point)
			///			records[point] = s
			///data[point] = append(data[point], s)
		}

	}

	if begins != 2 || records_cnt_mag == 0 || records_cnt_mag != records_cnt_phase {
		log.Print("Wrong file format: " + filename + " cnt_mag=" + strconv.Itoa(records_cnt_mag) + " cnt_phase=" + strconv.Itoa(records_cnt_phase))
		return
	}

	//stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,DateTime('now'),?,?)")
	stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,?,?,?)")
	checkErr(err)

	log.Print("Insert new MEASURE: hash=" + hash + " name=" + name + " filename=" + filename)
	res, err := stmt.Exec(hash, name, ts, filename, records_cnt_mag)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	stmtData, err := db.Prepare("INSERT INTO measure_data(measure_id, freq, magnitude, degrees) values(?,?,?,?)")

	db.Exec("BEGIN TRANSACTION")
	for i := 0; i < records_cnt_mag; i++ {

		//		log.Print("Insert data. Record #" + strconv.Itoa(i) + " ID=" + strconv.FormatInt(id, 10))
		_, err := stmtData.Exec(id, records[i].freq, records[i].mag, records[i].deg)
		checkErr(err)
	}
	db.Exec("END TRANSACTION")

}

func process_dir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		//		fmt.Println(filepath.Ext(f.Name()))
		if filepath.Ext(f.Name()) == ".csv" {
			//			var filename = f.Name()
			//			var extension = filepath.Ext(filename)
			//			var name = filename[0 : len(filename)-len(extension)]
			//			fmt.Println(name)
			//			fmt.Println(hash_file_md5(dir + f.Name()))
			file2db(dir + "\\" + f.Name())

		}
	}
}

func main_z() {

	var data_lookup = "C:\\MY\\zz\\"

	process_dir(data_lookup)
	draw()

}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Name

	case 2:
		return item.Points

	case 3:
		return item.Date
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Bar < b.Bar)

		case 2:
			return c(a.Points < b.Points)

		case 3:
			return c(a.Quux.Before(b.Quux))
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *FooModel) ResetRows() {
	// Create some random data. zzzzzzzzzz

	var count = 0
	//	fmt.Println("count=", count)

	//	log.Print("Requesting points type=" + strconv.Itoa(tpe) + " MeasureID:" + strconv.Itoa(measureID))
	db, err := sql.Open("sqlite3", "./db.sqlite3")

	checkErr(err)
	rows, err := db.Query("SELECT measure_id, name, date, fname, comment, points FROM measures")
	checkErr(err)

	defer rows.Close()
	rows.Scan(&count)

	var pidx = 0

	m.items = make([]*Foo, 109)
	now := time.Now()
	for rows.Next() {
		var measure_id int
		var name sql.NullString
		var date time.Time
		var fname sql.NullString
		var points sql.NullInt64
		var comment sql.NullString

		err = rows.Scan(&measure_id, &name, &date, &fname, &comment, &points)
		checkErr(err)

		//dte, _ := time.Parse("2014-09-12 11:45:26", date.String)
		//		log.Print("----", pidx, measure_id)
		m.items[pidx] = &Foo{
			checked: false,
			Index:   measure_id,
			Name:    name.String,
			Date:    date,
			Points:  points.Int64,
			Comment: comment.String,
			//			Name:     name.String, // strings.Repeat("*", rand.Intn(5)+1),
			Baz:  rand.Float64() * 1000,
			Quux: time.Unix(rand.Int63n(now.Unix()), 0),
		}
		//			Quux:    time.Unix(rand.Int63n(now.Unix()), 0),

		pidx++
	}

	db.Close()
	//return pts
	/*
		now := time.Now()

		for i := range m.items {
			m.items[i] = &Foo{
				Index: i,
				Bar:   strings.Repeat("*", rand.Intn(5)+1),
				Baz:   rand.Float64() * 1000,
				Quux:  time.Unix(rand.Int63n(now.Unix()), 0),
			}
		}
	*/
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func main() {
	var combo *walk.ComboBox
	var ExpEditButton *walk.PushButton
	//	dbPointsExpression(759, 760, "test")
	//	return

	app := walk.App()
	app.SetOrganizationName("RelizIT")
	app.SetProductName("MeasureDB")

	settings := walk.NewIniFileSettings("settings.ini")
	if err := settings.Load(); err != nil {
		log.Fatal(err)
	} else {
		log.Print(settings.Get("testzz"))
	}

	settings.Put("test2", "fuck off")
	// ------------------------register math functions
	var (
		phdelta = &operators.Operator{
			Name:          "phdelta",
			Precedence:    0,
			Associativity: operators.L,
			Args:          2,
			Operation: func(args []float64) float64 {
				if (math.Abs(args[0] - args[1])) < 180 {
					//					log.Print("match #1")
					return math.Abs(args[0] - args[1])
				} else {
					//					log.Print("match #2")
					return 360 - math.Abs(args[0]-args[1])
				}

			},
		}
	)
	functions.Register(phdelta)

	res, err := compute.Evaluate("phdelta(24 179)")

	if err != nil {
		return
	}
	fmt.Printf("%s\n", strconv.FormatFloat(res, 'G', -1, 64))

	log.Print(res)
	//	return

	mw := new(MyMainWindow)
	var openAction *walk.Action
	var MagGraphAction *walk.Action
	var PhaseGraphAction *walk.Action
	////////////////////////////////

	main_z()
	//	MasterMeasure = 1472
	//	dbPointsExpression(MasterMeasure, 1473, "log(sqrt((MagA ^ 2 + MagB ^ 2 + sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB)))) / (MagA ^ 2 + MagB ^ 2 - sqrt(MagA^4+MagB^4+2*cos(2*(PhA-PhB))))))")
	//	return
	//////////////
	//	return
	boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	goodIcon, _ := walk.Resources.Icon("../img/check.ico")
	badIcon, _ := walk.Resources.Icon("../img/stop.ico")

	barBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas, err := walk.NewCanvasFromImage(barBitmap)
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas.GradientFillRectangle(walk.RGB(255, 0, 0), walk.RGB(0, 255, 0), walk.Horizontal, walk.Rectangle{0, 0, 100, 1})

	canvas.Dispose()

	model = NewFooModel()

	var tv *walk.TableView
	var composite *walk.Composite
	var splitter *walk.Splitter
	var showContextMenu *walk.Action

	expobj := new(Expression)
	///////////////////////////////////////

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "RelizIT measures DB",

		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Image:       "../img/open.png",
						OnTriggered: mw.openAction_Triggered,
					},
					Action{
						AssignTo:    &MagGraphAction,
						Text:        "&Create Graph",
						Image:       "../img/plus.png",
						OnTriggered: mw.MagGraphAction_Triggered,
					},
					Action{
						AssignTo:    &PhaseGraphAction,
						Text:        "&Create ",
						Image:       "../img/document-properties.png",
						OnTriggered: mw.PhaseGraphAction_Triggered,
					},
					Separator{},
					Action{
						AssignTo:    &PhaseGraphAction,
						Text:        "&Options ",
						Image:       "../img/document-properties.png",
						OnTriggered: mw.PhaseGraphAction_Triggered,
					},
					Separator{},
					Action{
						Text:        "Exit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},

			Menu{
				Text: "E&xpressions",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "Expressions Editor",
						Image:       "../img/open.png",
						OnTriggered: mw.openAction_Triggered,
					},
					Action{
						AssignTo:    &MagGraphAction,
						Text:        "Constants Editor",
						Image:       "../img/plus.png",
						OnTriggered: mw.MagGraphAction_Triggered,
					},
				},
			},

			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: mw.aboutAction_Triggered,
					},
				},
			},
		},

		ToolBarItems: []MenuItem{
			ActionRef{&openAction},
			ActionRef{&MagGraphAction},
			ActionRef{&PhaseGraphAction},
		},
		MinSize: Size{320, 240},

		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				AssignTo: &splitter,
				Children: []Widget{
					TabWidget{
						AssignTo: &mw.tabWidget,

						OnSizeChanged: func() {
							log.Printf("Main window size changed", mw.tabWidget.Width())
							imgW = int64(mw.tabWidget.Width())
							imgH = int64(mw.tabWidget.Height())

						},
					},
					Composite{
						AssignTo: &composite,
						MaxSize:  Size{450, 300},

						Layout: Grid{Columns: 2},
						Children: []Widget{
							ComboBox{
								AssignTo: &combo,
								Value:    Bind("SpeciesId", SelRequired{}),
								//BindingMember: "Id",
								DisplayMember: "Name",
								Model:         KnownExpressions(),
								OnCurrentIndexChanged: func() {
									log.Print(combo.CurrentIndex())
									log.Print(combo.Text())
								},
							},

							PushButton{
								AssignTo:   &ExpEditButton,
								ColumnSpan: 2,
								Text:       "Edit Expression",

								OnClicked: func() {
									return
									//exp := new(Expression)
									var exp Expression
									log.Print("Epression load")

									//									exp := ExressionLoad(combo.Text())
									log.Print("Call dialog")
									return
									if cmd, err := RunExpressionDialog(mw, &exp); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										log.Printf("%+v", expobj)
									}
								},
							},
							TableView{
								AssignTo:              &tv,
								AlternatingRowBGColor: walk.RGB(239, 239, 239),
								CheckBoxes:            true,
								ColumnsOrderable:      true,
								MultiSelection:        true,
								ContextMenuItems: []MenuItem{
									Action{
										AssignTo: &showContextMenu,
										Text:     "Draw dataset",
										OnTriggered: func() {
											log.Printf("Жмак!")
											mw.MagGraphAction_Triggered()
										},
									},

									Action{
										AssignTo: &showContextMenu,
										Text:     "Set Master Dataset",
										OnTriggered: func() {
											idx := tv.CurrentIndex()
											if model.Checked(idx) && MasterMeasure != model.items[idx].Index {
												log.Print("Master measure:"+model.items[idx].Name, "  ", model.items[idx].Index)
												MasterMeasure = model.items[idx].Index
											}
										},
									},
									Action{
										AssignTo: &showContextMenu,
										Text:     "MenutItem",
										OnTriggered: func() {
											log.Printf("Жмак!")
										},
									},
								},
								Columns: []TableViewColumn{
									{Title: "#"},
									{Title: "Name"},
									{Title: "Points", Width: 70}, //Alignment: AlignFar
									{Title: "Date", Format: "2006-01-02 15:04:05", Width: 130},
								},

								StyleCell: func(style *walk.CellStyle) {
									item := model.items[style.Row()]

									if item.checked {
										if style.Row()%2 == 0 {
											style.BackgroundColor = walk.RGB(159, 215, 255)
											//									style.BackgroundColor = walk.RGB(159, 0, 0)
										} else {
											style.BackgroundColor = walk.RGB(143, 199, 239)
											//									style.BackgroundColor = walk.RGB(159, 0, 255)
										}
									}

									switch style.Col() {
									case 1:
										if canvas := style.Canvas(); canvas != nil {
											bounds := style.Bounds()
											bounds.X += 2
											bounds.Y += 2
											bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Bar)))
											bounds.Height -= 4
											///									canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Bar), 1}, 127)

											bounds.X += 4
											bounds.Y += 2
											//									canvas.DrawText(item.Bar, tv.Font(), 0, bounds, walk.TextLeft)
										}

									case 2:

										if item.Baz >= 900.0 {
											style.TextColor = walk.RGB(0, 191, 0)
											style.Image = goodIcon
										} else if item.Baz < 100.0 {
											style.TextColor = walk.RGB(255, 0, 0)
											style.Image = badIcon
										}

									case 3:
										if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
											style.Font = boldFont
										}
									}
								},
								Model: model,
								OnCurrentIndexChanged: func() {
									//							model.ccccccccccccc
									//									log.Printf("OnCurrentIndexChanged: %v\n", tv.CurrentIndex())

								},
								OnSelectedIndexesChanged: func() {

									//							log.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
									/*
										idx := tv.CurrentIndex()
										if model.Checked(idx) && MasterMeasure != model.items[idx].Index {
											log.Print("Master measure:"+model.items[idx].Name, "  ", model.items[idx].Index)
											MasterMeasure = model.items[idx].Index
										}
									*/

									//							checked := model.Checked(idx)
									//							model.SetChecked(idx, true)
									//							model.PublishRowChanged(idx)
								},
								OnItemActivated: func() {
									//idx := pg.roleBaseDataTV.CurrentIndex()

									//log.Printf("SelectedIndexes: %v\n", tv.Swi())
								},
							},
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	//	style := uint32(win.GetWindowLong(mw.MainWindow.Handle(), win.GWL_STYLE))
	//	style &= win.WS_MAXIMIZE
	//	style |= win.WS_MAXIMIZE

	//	win.SetWindowLong(mw.MainWindow.Handle(), win.GWL_STYLE, int32(style))
	///	mw.MainWindow.setAnd

	win.ShowWindow(mw.MainWindow.Handle(), win.SW_SHOWMAXIMIZED)

	lv, err := NewLogView(mw.MainWindow)
	if err != nil {
		log.Fatal(err)
	}
	//tv.SetWidth(100)
	//	mw.tabWidget.SetWidth((mw.Width() / 4) * 3)
	//	mw.tabWidget.SetWidth(1200)

	//	mw.SetFullscreen(true)
	imgW = int64(mw.tabWidget.Width())
	imgH = int64(mw.tabWidget.Height())

	//	lv.PostAppendText("Hello!\n")
	//	lv.PostAppendText("This is a log\n")
	//	lv.SetClientSize(walk.Size{500, 500})

	//log.SetOutput(lv)

	log.Printf("-----------------------------imgW=", imgW, " imgH=", imgH, lv)

	log.Print("|      Showtime         |")
	log.Print("Database.......... OK")
	log.Print("Measures.......... ", dbCountMeasures())

	//	openImage(mw, "./points.png")
	//	return
	//lv.SetClientSize(walk.Size{500, 30})
	//lv.SetHeight(50)
	//lv.SetX(50)

	/////////////////
	/*
		res, err := compute.Evaluate("sin(1)*15.1/cos(15)")
		if err != nil {
			return
		}
		fmt.Printf("%s\n", strconv.FormatFloat(res, 'G', -1, 64))
	*/
	//////////////////

	/*
		// Groutine
			go func() {
				// wwwwwwwwwwwwwwwwwwww
				log.Print("Directory monitor started.")

				usb, err := NewUSB()
				if err != nil {
					log.Fatal(err)
				}

				usb.RegisterDeviceNotification()
				usb.Run()
				//C.WatchDirectory(C.CString("C:\\MY\\Z"))
				//C.WatchDirectory((*C.char)((syscall.StringToUTF16Ptr("C:\\MY\\Z"))))
				//		syscall.StringToUTF16Ptr
				for i := 0; i < 10000; i++ {
					time.Sleep(10000 * time.Millisecond)
					//			log.Println("Tic" + "\r\n")
				}
			}()
	*/
	defer settings.Save()

	mw.MainWindow.Run()

}

type MyMainWindow struct {
	*walk.MainWindow
	tabWidget    *walk.TabWidget
	prevFilePath string
}

func (mw *MyMainWindow) openAction_Triggered() {
	/*
		if err := mw.openImage(); err != nil {
			log.Print(err)
		}*/

}

func (mw *MyMainWindow) MagGraphAction_Triggered() {

	log.Print("graph action triggered")
	log.Print(model.items)

	/*
		if mw.tabWidget.Pages().Len() == 2 {
			mw.tabWidget.Pages().RemoveAt(0)
			mw.tabWidget.Pages().RemoveAt(0)

			}*/

	// HERE THE FUCK!

	if mw.tabWidget.Pages().Len() == 3 {
		//		mw.tabWidget.Pages().RemoveAt(0)
		//		mw.tabWidget.Pages().RemoveAt(0)
		//		mw.tabWidget.Pages().RemoveAt(2)
	}
	//	mw.tabWidget.Pages().Clear()

	drawModel(mw, GraphMag)
	drawModel(mw, GraphPhase)
	drawModel(mw, GraphRdb)

	//mw.tabWidget.Pages().At(0).SetCur
	mw.tabWidget.SetCurrentIndex(0)
	//win.SendMessage(mw.tabWidget.Handle(), win.TCM_SETCURFOCUS, 0)
	for i := range model.items {
		if model.items[i].checked {
			log.Print(model.items[i].Name, model.items[i].Index)
			//			zzzzzzzzzz

		}
	}
}

func (mw *MyMainWindow) PhaseGraphAction_Triggered() {

	log.Print("graph action triggered")
	/*
		log.Print(model.items)
		drawModel(mw, GraphPhase)
		for i := range model.items {
			if model.items[i].checked {
				log.Print(model.items[i].Name, model.items[i].Index)
				//			zzzzzzzzzz

			}
		}
	*/
}

/*
func (mw *MyMainWindow) openImage() error {
	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.Title = "Select an Image"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return err
	} else if !ok {
		return nil
	}

	mw.prevFilePath = dlg.FilePath

	img, err := walk.NewImageFromFile(dlg.FilePath)
	if err != nil {
		return err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			img.Dispose()
		}
	}()

	page, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	if page.SetTitle(path.Base(strings.Replace(dlg.FilePath, "\\", "/", -1))); err != nil {
		return err
	}
	page.SetLayout(walk.NewHBoxLayout())

	defer func() {
		if !succeeded {
			page.Dispose()
		}
	}()

	imageView, err := walk.NewImageView(page)
	if err != nil {
		return err
	}

	defer func() {
		if !succeeded {
			imageView.Dispose()
		}
	}()

	imageView.SetMode(walk.ImageViewModeShrink)

	if err := imageView.SetImage(img); err != nil {
		return err
	}

	if err := mw.tabWidget.Pages().Add(page); err != nil {
		return err
	}

	if err := mw.tabWidget.SetCurrentIndex(mw.tabWidget.Pages().Len() - 1); err != nil {
		return err
	}

	succeeded = true

	return nil
}
*/

func openImage(mw *MyMainWindow, fpath string, tpe int) error {

	img, err := walk.NewImageFromFile(fpath)
	if err != nil {
		return err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			img.Dispose()
		}
	}()

	page, err := walk.NewTabPage()
	if err != nil {
		return err
	}

	if tpe == GraphMag {
		page.SetTitle("Magnitude Graph")
	} else {
		page.SetTitle("Phase Δ Graph")
	}

	if tpe == GraphRdb {
		page.SetTitle("r(db) Graph")
	}
	/*
		if page.SetTitle(path.Base(strings.Replace(fpath, "\\", "/", -1))); err != nil {
			return err
		}
	*/
	page.SetLayout(walk.NewHBoxLayout())

	defer func() {
		if !succeeded {
			page.Dispose()
		}
	}()

	imageView, err := walk.NewImageView(page)
	if err != nil {
		return err
	}

	defer func() {
		if !succeeded {
			imageView.Dispose()
		}
	}()

	imageView.SetMode(walk.ImageViewModeShrink)

	if err := imageView.SetImage(img); err != nil {
		return err
	}

	// HERE
	if err := mw.tabWidget.Pages().Add(page); err != nil {
		return err
	}

	if err := mw.tabWidget.SetCurrentIndex(mw.tabWidget.Pages().Len() - 1); err != nil {
		return err
	}
	succeeded = true

	return nil
}

func (mw *MyMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw, "About", "RelizIT RF measures archive\n\nInput arguments:\nPhA, PhB, MagA, MagB\n\nFunctions list:\nsin, cos, tan, cot, sec, csc,sin, acos, atan, acot, asec, acsc\nabs, log, ln, lg, sqrt, phdelta\n\nOperators:\n+ - neg * / % ^", walk.MsgBoxIconInformation)
}
