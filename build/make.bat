set GOARCH=arm
set GOOS=linux
go build -o ..\bin\s0counter ..\cmd\s0counter.go

rem set GOARCH=386
rem set GOOS=windows
rem go build -o ..\bin\s0counter.exe ..\cmd\s0counter.go