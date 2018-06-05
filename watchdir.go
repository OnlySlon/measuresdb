package main

/*
#include <stdio.h>
#include <windows.h>
#include <stdlib.h>
extern void goCallbackFileChange(void);
// https://msdn.microsoft.com/de-de/library/windows/desktop/aa365261(v=vs.85).aspx
static inline void WatchDirectory(const char* dir) {
  DWORD waitStatus;
  HANDLE handle;
	// FILE_NOTIFY_CHANGE_FILE_NAME � File creating, deleting and file name changing
	// FILE_NOTIFY_CHANGE_DIR_NAME � Directories creating, deleting and file name changing
	// FILE_NOTIFY_CHANGE_ATTRIBUTES � File or Directory attributes changing
	// FILE_NOTIFY_CHANGE_SIZE � File size changing
	// FILE_NOTIFY_CHANGE_LAST_WRITE � Changing time of write of the files
	// FILE_NOTIFY_CHANGE_SECURITY � Changing in security descriptors
  handle = FindFirstChangeNotification(
  	dir,
		TRUE,
		FILE_NOTIFY_CHANGE_FILE_NAME | FILE_NOTIFY_CHANGE_SIZE | FILE_NOTIFY_CHANGE_DIR_NAME
	);
  if (handle == INVALID_HANDLE_VALUE){
    printf("\n ERROR: FindFirstChangeNotification function failed.\n");
    ExitProcess(GetLastError());
  }
  if ( handle == NULL ) {
    printf("\n ERROR: Unexpected NULL from FindFirstChangeNotification.\n");
    ExitProcess(GetLastError());
  }
  while (TRUE) {
    // printf("\nWaiting for notification...\n");
		waitStatus = WaitForSingleObject(handle, INFINITE);
		switch (waitStatus) {
      case WAIT_OBJECT_0:
				// printf("A file was created, renamed, or deleted in the directory\n");
				goCallbackFileChange();
				// continue monitoring
				FindNextChangeNotification(handle);
        break;
      case WAIT_TIMEOUT:
        // printf("\nNo changes in the timeout period.\n");
        break;
      default:
        printf("\n ERROR: Unhandled status.\n");
        ExitProcess(GetLastError());
        break;
    }
  }
}
*/
import "C"
import (
	"os"
	"path/filepath"
	"unsafe"
)

// Snapshot struct holds information about files in the watched directory
type Snapshot struct {
	CallbackChan chan os.FileInfo
	Root         string
	DirSnapshot  map[string]os.FileInfo
}

// TODO: for watching more as one directory this needs to be a map
var snap Snapshot

// DirectoryChangeNotification expected path to the directory to watch as string
// and a FileInfo channel for the callback notofications
// Notofication is fired each time file in the directory is changed or some new
// file (or sub-directory) is created
func DirectoryChangeNotification(path string, callbackChan chan os.FileInfo) {

	snap = Snapshot{
		CallbackChan: callbackChan,
		Root:         path,
		DirSnapshot:  createSnapshot(path),
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.WatchDirectory(cpath)
}

//export goCallbackFileChange
func goCallbackFileChange() {
	fileInfo := findChange()
	if fileInfo != nil {
		snap.CallbackChan <- fileInfo
	}
}

func findChange() os.FileInfo {
	var fileInfo os.FileInfo
	filepath.Walk(snap.Root, func(file string, f os.FileInfo, err error) error {
		if _, ok := snap.DirSnapshot[file]; !ok {
			snap.DirSnapshot[file] = f
			fileInfo = f
			return nil
		}
		if snap.DirSnapshot[file].Size() != f.Size() {
			snap.DirSnapshot[file] = f
			fileInfo = f
			return nil
		}
		return nil
	})
	return fileInfo
}

func createSnapshot(path string) map[string]os.FileInfo {
	fileList := make(map[string]os.FileInfo)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList[path] = f
		return nil
	})
	return fileList
}
