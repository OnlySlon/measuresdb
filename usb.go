package main

import (
	"log"
	"unsafe"

	winC "github.com/codyguo/win"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

const mainWindowWindowClass = `\o/ USB_Class \o/`

var notificationFilter = new(winC.DEV_BROADCAST_DEVICEINTERFACE)

func init() {
	walk.MustRegisterWindowClass(mainWindowWindowClass)
}

type USB struct {
	walk.MainWindow
}

func NewUSB() (*USB, error) {
	usb := new(USB)
	if err := walk.InitWindow(
		usb,
		nil,
		mainWindowWindowClass,
		win.WS_OVERLAPPEDWINDOW,
		win.WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	return usb, nil
}

func (usb *USB) RegisterDeviceNotification() bool {
	//var notificationFilter = new(winC.DEV_BROADCAST_DEVICEINTERFACE)
	notificationFilter.Dbcc_size = uint32(unsafe.Sizeof(*notificationFilter))
	notificationFilter.Dbcc_devicetype = winC.DBT_DEVTYP_DEVICEINTERFACE
	notificationFilter.Dbcc_classguid = winC.GUID_DEVINTERFACE_USB_DEVICE

	handle := winC.HANDLE(usb.Handle())
	ret := winC.RegisterDeviceNotificationW(handle, (uintptr)(unsafe.Pointer(notificationFilter)), winC.DEVICE_NOTIFY_WINDOW_HANDLE)

	return ret != 0
}

func (usb *USB) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_DEVICECHANGE:
		switch wParam {
		case winC.DBT_DEVICEARRIVAL:
			{
				log.Printf("USB Connected -> %v", msg)
				log.Print(notificationFilter.Dbcc_name)
			}
		case winC.DBT_DEVICEREMOVECOMPLETE:
			log.Printf("USB Disconnected -> %v", msg)
		}
	}
	return usb.FormBase.WndProc(hwnd, msg, wParam, lParam)
}
