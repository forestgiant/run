# run
Go (golang) tool to read system commands from a JSON file and execute them.

# Install
`go get -u github.com/forestgiant/run`

# Usage
Create a commands.json
```
[{
	"path": "",
	"name": "ls",
	"args": ["-all", "-G"],
	"sleep": 500
}, {
	"path": "",
	"name": "ls",
	"args": ["-aG"],
	"sleep": 500
}, {
	"name": "pwd"
}]
```

Now run all the commands in order:
`run ./commands.json`
