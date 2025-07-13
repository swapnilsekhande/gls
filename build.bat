$env:GOOS = "windows"
$env:GOARCH = "386"
go build -o build/plc_reader_win32.exe main.go