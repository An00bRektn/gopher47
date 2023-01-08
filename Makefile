all:
	printf "\033[0;36m[*]\033[0m Building Gopher47 agent (windows and linux)...\n"
	GOOS=windows GOARCH=amd64 go build -o bin/gopher47_win.exe
	GOOS=linux GOARCH=amd64 go build -o bin/gopher47_linux

linux:
	printf "\033[0;36m[*]\033[0m Building Gopher47 agent (linux)...\n"
	GOOS=linux GOARCH=amd64 go build -o bin/gopher47

windows:
	printf "\033[0;36m[*]\033[0m Building Gopher47 agent (windows)...\n"
	GOOS=windows GOARCH=amd64 go build -o bin/gopher47