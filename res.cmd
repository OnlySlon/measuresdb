windres -o rsrc.syso version.rc
#rsrc -arch amd64 -manifest measuredb.manifest -o rsrc.syso
#rsrc -arch amd64 -b ./img/create.png -o rsrc.syso
go build -ldflags="-H windowsgui"
./upx.exe measuresdb.exe
