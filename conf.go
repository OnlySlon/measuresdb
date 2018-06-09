package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Configuration struct {
	UsbDriveLetter string
	MonitoringDir  string
}

var conf Configuration

func ConfigTest() {
	ConfigLoad()
	conf.MonitoringDir = "./csv/"
	conf.UsbDriveLetter = "F:"
	ConfigSave()
}

func ConfigLoad() {

	log.Print("----LOAD-CONFIG----")
	settings := walk.NewIniFileSettings("settings.ini")
	settings.SetPortable(true)

	if err := settings.Load(); err != nil {
		log.Print(err)
	} else {
		log.Print(settings.Get("testzz"))
		conf.MonitoringDir = ""
		MonitoringDir, ok := settings.Get("MonitoringDir")
		if ok {
			conf.MonitoringDir = MonitoringDir
		}
		UsbDriveLetter, ok := settings.Get("UsbDriveLetter")
		if ok {
			conf.UsbDriveLetter = UsbDriveLetter
		}
	}

}

func ConfigSave() {
	log.Print("----SAVE-CONFIG----" + conf.MonitoringDir)
	settings := walk.NewIniFileSettings("settings.ini")
	settings.SetPortable(true)
	if err := settings.Load(); err != nil {
		log.Print(err)
	} else {
		settings.Put("MonitoringDir", conf.MonitoringDir)
		settings.Put("UsbDriveLetter", conf.UsbDriveLetter)
		settings.Save()
	}
}

func ConfigClick(owner walk.Form) {
	if cmd, err := ConfigRunDialog(owner); err != nil {
		log.Print(err)
	} else if cmd == walk.DlgCmdOK {
		ConfigSave()
	}

}

func ConfigRunDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	log.Print("RunExpressionDialog!!!")

	return Dialog{
		AssignTo:      &dlg,
		Title:         "Configuration",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			DataSource:     &conf,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{640, 200},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Usb Drive letter:",
					},
					LineEdit{
						Text: Bind("UsbDriveLetter"),
					},

					Label{
						//						ColumnSpan: 2,
						Text: "Monitoring Dir:",
					},
					LineEdit{
						Text: Bind("MonitoringDir"),
					},
					Label{
						ColumnSpan: 2,
						Text:       "* For apply changes need restart program",
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
							} else {
								dlg.Accept()
								ConfigSave()
								ConfigLoad()
								dlg.Accept()

							}

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
