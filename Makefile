build: win64 win32 linux

win64:
	env GOARCH=amd64 GOOS=windows go build -o ./bin/go-rt.win64.exe

win32:
	env GOARCH=386 GOOS=windows go build -o ./bin/go-rt.win32.exe

linux:
	env GOARCH=amd64 GOOS=linux go build -o ./bin/go-rt
