IN FREQ1 VAL1 PHASE1 (Z AA) A B
   FREQ2 VAL2 PHASE1 (AI AJ)

S32(mod) = 0
S41(mod) = 0


Ex(dB) = 10 * LOG10(VAL1)
Ey(dB) = 10 * LOG10(VAL2)

Ey-Ex = (Ey - Ex)


export PATH=$PATH:/c/Users/Slon/go/bin
export PATH=$PATH:/c/msys64/mingw64/bin/

/c/Users/Slon/bin:/mingw64/bin:/usr/local/bin:/usr/bin:/bin:/mingw64/bin:/usr/bin:/c/Users/Slon/bin:/c/Program Files (x86)/Intel/Intel(R) Management Engine Components/iCLS:/c/Program Files/Intel/Intel(R) Management Engine Components/iCLS:/c/WINDOWS/system32:/c/WINDOWS:/c/WINDOWS/System32/Wbem:/c/WINDOWS/System32/WindowsPowerShell/v1.0:/c/Program Files (x86)/Intel/Intel(R) Management Engine Components/DAL:/c/Program Files/Intel/Intel(R) Management Engine Components/DAL:/c/Program Files (x86)/Intel/Intel(R) Management Engine Components/IPT:/c/Program Files/Intel/Intel(R) Management Engine Components/IPT:/c/Program Files (x86)/NVIDIA Corporation/PhysX/Common:/c/WINDOWS/system32:/c/WINDOWS:/c/WINDOWS/System32/Wbem:/c/WINDOWS/System32/WindowsPowerShell/v1.0:/c/WINDOWS/System32/OpenSSH:/c/Go/bin:/c/Program Files/SlickEditV16.0.0x64/win:/c/Users/Slon/AppData/Local/Microsoft/WindowsApps:%GOPATH%/bin:/usr/bin/vendor_perl:/usr/bin/core_perl:/c/Users/Slon/go/bin



   CREATE TABLE `measure_data` (
        `measure_id` INTEGER SECONDARY KEY,
        `freq`  INTEGER NULL,
        `magnitude` REAL NULL,
        `degrees` REAL  NULL
    );


   CREATE TABLE `measures` (
        `measure_id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `hash`  TEXT  NULL,
        `name` REAL NULL,
        `date` DATETIME  NULL,
	`fname` TEXT  NULL
    );


{
	"resource": "/c:/Users/Slon/WORK/test/csv/csv.go",
	"owner": "go-error",
	"code": "undefined",
	"severity": 8,
	"message":
 "flag provided but not defined: -linkmode internal\n
usage: build [-o output] [-i] [build flags] [packages]\nRun 'go help build' for details.\n",

	"startLineNumber": 1,
	"startColumn": 1,
	"endLineNumber": 1,
	"endColumn": 13
}


\The "goplay" command is not available.  Use "go get -v github.com/haya14busa/goplay/cmd/goplay" to install.


CREATE TABLE `measure_data` (
	`measure_id`	INTEGER SECONDARY KEY,
	`freq`	INTEGER,
	`magnitude`	REAL,
	`degrees`	REAL
);

CREATE TABLE `measures` (
	`measure_id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`hash`	TEXT,
	`name`	REAL,
	`date`	DATETIME,
	`fname`	TEXT,
	`points`	INTEGER
);

go build -ldflags='-linkmode=internal' && ./csv.exe
$ go build -ldflags="-linkmode=internal -extld=none" -x && ./csv.exe > log 2>&1

rsrc -manifest=csv.exe.manifest -arch amd64 -o rsrc.syso



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

----------------------
ZZ     - freq
AA 227 - phase1 
AA 20 -  mag1      B

AJ - mag2

C = 10*LOG10(Mag1)
G = 10*LOG10(Mag2)

H(Ey-Ex)  =  (G-C)

E3n=SQRT(mag1^2 +D3^2 + 2*mag1*D3*COS(PI()/180*(I3-J3)))


SUPPORTED
trig: sin, cos, tan, cot, sec, csc (1 / math.Sin(args[0]), asin, acos, atan, acot, asec, acsc
abs
log (10), ln, lg(2)
sqrt
phdiff
+ - neg * / % ^

phdiff () =IF(ABS(I6-L6)<180;ABS(I6-L6);360-ABS(I6-L6))

MagA PhA MagB PhB

MagA^2 + MagB^2 + sqrt(MagA^4 + MagB^4 + 2*cos(2*(PhA-PhB))
MagA^2 + MagB^2 - sqrt(MagA^4 + MagB^4 + 2*cos(2*(PhA-PhB))

AAAAAAAA
sqrt((MagA^2 + MagB^2 + sqrt(MagA^4 + MagB^4 + 2*cos((phdelta(PhA PhB)))/(MagA^2 + MagB^2 - sqrt(MagA^4 + MagB^4 + 2*MagA^2*MagB^2*cos(phdelta(PhA PhB))))) / (MagA^2 + MagB^2 - sqrt(MagA^4 + MagB^4 + 2*MagA^2*MagB^2*cos(phdelta(PhA PhB))))))

sqrt((MagA^2 + MagB^2 + sqrt(MagA^4 + MagB^4 + 2*MagA^2*MagB^2*cos(2*(PhA - PhB))

SELECT m1.freq, m1.magnitude, m1.degrees,m2.magnitude, m2.degrees FROM measure_data as m1 LEFT join measure_data as m2 WHERE m1.freq=m2.freq AND m1.measure_id=748 AND m2.measure_id=747

----------------

sqrt((MagA^2+MagB^2+sqrt(MagA^4+MagB^4+2*MagA^2*MagB^2*cos(torad(2*phdelta(PhA PhB))))) / (MagA^2+MagB^2-sqrt(MagA^4+MagB^4+2*MagA^2*MagB^2*cos(torad(2*phdelta(PhA PhB))))))

---------------
function y = db2mag(ydb)
  y = 10.^(ydb/20);

function ydb = mag2db(y)
  ydb = 20*log10(y);
endfunction


2018/05/28 19:34:04 Freq=6360875000 Ph1=108.297742005317 Ph2=-154.820263512559 res=96.88199448212401
2018/05/28 19:34:04 Freq=6364250000 Ph1=78.00738906025 Ph2=-175.370865563921 res=106.621745375829
2018/05/28 19:34:04 Freq=6367625000 Ph1=57.4853829172765 Ph2=164.139585295373 res=-106.6542023780965
2018/05/28 19:34:04 Freq=6371000000 Ph1=36.564623506595 Ph2=143.714210441235 res=-107.14958693464
2018/05/28 19:34:04 Freq=6374375000 Ph1=12.0971577923648 Ph2=123.47379430992 res=-111.3766365175552
2018/05/28 19:34:04 Freq=6377750000 Ph1=-12.2231178153525 Ph2=103.367474345291 res=-115.5905921606435
2018/05/28 19:34:04 Freq=6381125000 Ph1=-34.7584584132803 Ph2=83.3536793212013 res=-118.11213773448159
2018/05/28 19:34:04 Freq=6384500000 Ph1=-54.1337896233897 Ph2=63.3596921125955 res=-117.4934817359852
2018/05/28 19:34:04 Freq=6387875000 Ph1=-75.6357401192399 Ph2=43.1980166543113 res=-118.8337567735512
2018/05/28 19:34:04 Freq=6391250000 Ph1=-95.0688507198399 Ph2=23.0782861142001 res=-118.14713683404
2018/05/28 19:34:04 Freq=6394625000 Ph1=-116.426036287604 Ph2=2.96917479939855 res=-119.39521108700255
2018/05/28 19:34:04 Freq=6398000000 Ph1=-138.382820787411 Ph2=-17.2803010781616 res=-121.10251970924939
2018/05/28 19:34:04 Freq=6401375000 Ph1=-159.637664337244 Ph2=-37.5657271020197 res=-122.07193723522428
2018/05/28 19:34:04 Freq=6404750000 Ph1=138.935244139813 Ph2=-57.7021995735328 res=163.3625562866542
2018/05/28 19:34:04 Freq=6408125000 Ph1=139.084408686856 Ph2=-78.0803009886217 res=142.8352903245223
2018/05/28 19:34:04 Freq=6411500000 Ph1=120.279700972961 Ph2=-98.4672643454214 res=141.2530346816176
2018/05/28 19:34:04 Freq=6414875000 Ph1=96.5837294172663 Ph2=-119.006796341495 res=144.40947424123868
2018/05/28 19:34:04 Freq=6418250000 Ph1=81.477340053818 Ph2=-139.430337781335 res=139.09232216484702
2018/05/28 19:34:04 Freq=6421625000 Ph1=48.5856466239776 Ph2=-159.917777992017 res=151.4965753840054
2018/05/28 19:34:04 
Freq=6425000000 Ph1=24.1067392478334 Ph2=179.599888107782 res=-155.4931488599486
2018/05/28 19:34:04 

=IF(ABS(I6-L6)<180;ABS(I6-L6);360-ABS(I6-L6))
Records: 201
2018/05/28 19:34:04 Draw graph...



		function Amplitude2dB(amplitude) {
			return 20 * Math.log(amplitude) / Math.LN10;
		}
		function dB2Amplitude(db) {
			return Math.pow(10, db / 20);
		}

TODO -
- Редактор формул - done
- Apply формул при выборе датасетов - done
- Мониторинг USB драйва
- тип датасетов - входящие/Посчитанные
- экспорт данных/картинок

[2018.05.30 16:30:21] DEBUG SYS 26577 F:btcp_deasm_error: TSK-#3. REMOTE: Wrong filesize (InFooter:512323261 InPush: 512306602)



[2018.05.30 16:30:19] MARK  SYS 7224 F:btcp_file_lock_check: File:'
/data/incoming/b804_2d_VF/JurassicWorld2_FTR-6-2D_S_RU-XX_RU-12-TD_51-Dbox_2K_UP_20180526_MPS_IOP_VF/JurassicWorld2_FTR-6-2D_S_RU-XX_51-Dbox_0526_03.mxf
' has no lock
[2018.05.30 16:30:19] MARK  SYS 7224 F:btcp_file_lock_set: File lock (fd #18) set
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_ffoter_read: footer struct loaded. Chunks=489 ChunkSize=1048576 file_id=108242 fSize=512323261 PhysImageSize=6
2
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_bitmap_info: TSK-#3 TPE:RXIMG  ChunkSize:1048576 ImgSizeBits:489 ZeroBits: 4. Complete: 485 Mb of 4 Mb
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_ffoter_read: Image loaded. Size: 62 bytes
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_bitmap_info: TSK-#3 TPE:RXIMG  ChunkSize:1048576 ImgSizeBits:489 ZeroBits: 4. Complete: 485 Mb of 4 Mb
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_deasm_push: Image Loaded
[2018.05.30 16:30:19] DEBUG SYS 7224 F:btcp_sqlite_commit: ---- TRANSACTION COMMIT ----
[2018.05.30 16:30:19] ERROR SYS 7224 F:btcp_deasm_push: Wrong filesize (InFooter:512323261 InPush: 512306602)




2018/05/31 16:02:31 dbPointsExpression DS1=1965 DS2=1966 
EXP=20*Log10(Sqr((MagA^2+MagB^2+Sqr(MagA^4+MagB^4+2*MagA^2*MagB^2*CosD(2*phdelta(PhA PhB)))) / (MagA^2+MagB^2-Sqr(MagA^4+MagB^4+2*MagA^2*MagB^2*CosD(2*phdelta(PhA PhB))))))




Freq= 5700000000  PhA= 96.087583188858  PhB= -0.279421434867479   MagA= 0.03733024048686795  MagB= 0.039644654862311446  res= 1.1349877737411853
Freq= 5701000000  PhA= 95.9810698764451  PhB= -0.270097102555912   MagA= 0.037425721931247574  MagB= 0.0396050603605059  res= 1.131059658587619

