package main

// package main

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	_ "github.com/mattn/go-sqlite3"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// 	_ "runtime/cgo"

type measure struct {
	freq float64
	mag  float64
	deg  float64
}

type Foo struct {
	Index int
	Bar   string
	Baz   float64
	Quux  time.Time

	Name    string
	Date    time.Time
	Points  int64
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

func DataLoad(fname string) []measure {
	f, _ := os.Open(fname)

	var begins = 0
	data := make(map[int]measure)

	records := make([]measure, 15000)
	fmt.Printf("%v\n", data)
	var ln = ""
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var point = 0
	for scanner.Scan() {
		ln = scanner.Text()
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

			records[point].freq = freq
			if records[point].mag != 0 {
				records[point].deg = mag
			} else {
				records[point].mag = mag
			}
			point++
			//records = append(records, s)
			//			fmt.Printf("%i ", point)
			///			records[point] = s
			///data[point] = append(data[point], s)
		}
		if ln == "BEGIN" {
			fmt.Println("\n-------------------------!\n %i", point)
			point = 0
			begins++

		}
		//		fmt.Println(ln)

	}
	fmt.Printf("%v\n", records)
	return records
}

func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
		//		fmt.Println("Random Points i=", i)
	}
	return pts
}

func dbCountMeasures() int {
	return 0
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
			//			db.Close()
			return cnt
		}

	}
	return 0
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

	defer rows.Close()
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
		pts[pidx].Y = magnitude

		pidx++
		checkErr(err)
	}

	//	db.Close()
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

	p.Save(15*vg.Inch, 8*vg.Inch, "points.pdf")
	// Save the plot to a PNG file.
	if err := p.Save(15*vg.Inch, 8*vg.Inch, "points.png"); err != nil {

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
	for scanner.Scan() {
		ln = scanner.Text()
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
				records[point].mag = mag
				records_cnt_phase++
			}
			point++
			//records = append(records, s)
			//			fmt.Printf("%i ", point)
			///			records[point] = s
			///data[point] = append(data[point], s)
		}
		if ln == "BEGIN" {
			//			fmt.Println("\n-------------------------!\n %i", point)
			point = 0
			begins++

		}
		//		fmt.Println(ln)

	}

	if begins != 2 || records_cnt_mag == 0 || records_cnt_mag != records_cnt_phase {
		log.Print("Wrong file format: " + filename + " cnt_mag=" + strconv.Itoa(records_cnt_mag) + " cnt_phase=" + strconv.Itoa(records_cnt_phase))
		return
	}

	stmt, err := db.Prepare("INSERT INTO measures(hash, name, date, fname, points) values(?,?,DateTime('now'),?,?)")
	checkErr(err)

	log.Print("Insert new MEASURE: hash=" + hash + " name=" + name + " filename=" + filename)
	res, err := stmt.Exec(hash, name, filename, records_cnt_mag)
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
	db.Close()
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
		return item.Bar

	case 2:
		return item.Baz

	case 3:
		return item.Quux
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
			return c(a.Baz < b.Baz)

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
	rows, err := db.Query("SELECT measure_id, name, date, fname, comment FROM measures")
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

		err = rows.Scan(&measure_id, &name, &date, &fname, &comment)
		checkErr(err)

		//dte, _ := time.Parse("2014-09-12 11:45:26", date.String)
		log.Print("----", pidx, measure_id)
		m.items[pidx] = &Foo{
			Index:   measure_id,
			Name:    name.String,
			Date:    date,
			Points:  points.Int64,
			Comment: comment.String,
			Bar:     name.String, // strings.Repeat("*", rand.Intn(5)+1),
			Baz:     rand.Float64() * 1000,
			Quux:    time.Unix(rand.Int63n(now.Unix()), 0),
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

	mw := new(MyMainWindow)
	var openAction *walk.Action
	////////////////////////////////

	//	main_z()
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

	model := NewFooModel()

	var tv *walk.TableView

	var splitter *walk.Splitter

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
					Separator{},
					Action{
						Text:        "Exit",
						OnTriggered: func() { mw.Close() },
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
		},
		MinSize: Size{320, 240},
		Size:    Size{800, 600},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				AssignTo: &splitter,
				Children: []Widget{

					TabWidget{
						AssignTo: &mw.tabWidget,
					},
					TableView{
						AssignTo:              &tv,
						AlternatingRowBGColor: walk.RGB(239, 239, 239),
						CheckBoxes:            true,
						ColumnsOrderable:      true,
						MultiSelection:        true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "Name"},
							{Title: "Points", Alignment: AlignFar},
							{Title: "Date", Format: "2006-01-02 15:04:05", Width: 150},
						},
						StyleCell: func(style *walk.CellStyle) {
							item := model.items[style.Row()]

							if item.checked {
								if style.Row()%2 == 0 {
									style.BackgroundColor = walk.RGB(159, 215, 255)
								} else {
									style.BackgroundColor = walk.RGB(143, 199, 239)
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
						OnSelectedIndexesChanged: func() {
							fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw.MainWindow)
	if err != nil {
		log.Fatal(err)
	}

	//	lv.PostAppendText("Hello!\n")
	//	lv.PostAppendText("This is a log\n")
	//	lv.SetClientSize(walk.Size{500, 500})
	log.SetOutput(lv)

	log.Print("|      Showtime         |")
	log.Print("Database.......... OK", lv)
	log.Print("Measures.......... ", dbCountMeasures())
	main_z()

	openImage(mw, "./points.png")
	//	return
	//lv.SetClientSize(walk.Size{500, 30})
	//lv.SetHeight(50)
	//lv.SetX(50)

	go func() {
		for i := 0; i < 10000; i++ {
			time.Sleep(10000 * time.Millisecond)
			log.Println("Tic" + "\r\n")
		}
	}()

	mw.MainWindow.Run()
}

type MyMainWindow struct {
	*walk.MainWindow
	tabWidget    *walk.TabWidget
	prevFilePath string
}

func (mw *MyMainWindow) openAction_Triggered() {
	if err := mw.openImage(); err != nil {
		log.Print(err)
	}
}

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

func openImage(mw *MyMainWindow, fpath string) error {

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

	if page.SetTitle(path.Base(strings.Replace(fpath, "\\", "/", -1))); err != nil {
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

func (mw *MyMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw, "About", "RelizIT RF measures archive", walk.MsgBoxIconInformation)
}
