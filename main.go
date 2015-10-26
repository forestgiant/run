package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

type Cmd struct {
	Path  string   // Path to command
	Name  string   // Name of command
	Args  []string // Args to pass to command
	Sleep int      // in Miliseconds
}

func main() {
	jsonConfigPathDefault := "commands.json"

	// Check for command line configuration flags
	var (
		configPathUsage = fmt.Sprint("Path to config json.")
		configPathPtr   = flag.String("configpath", jsonConfigPathDefault, configPathUsage)
	)
	flag.Parse()

	// Setup commands fron commands.json
	cmds, cmdsErr := createCommandsFromJSON(*configPathPtr)
	if cmdsErr != nil {
		log.Fatal(cmdsErr)
	}

	// Run them in order
	for _, cmd := range cmds {
		sleep := time.Duration(cmd.Sleep) * time.Millisecond
		run(sleep, cmd.Path, cmd.Name, cmd.Args...)
	}

}

func run(sleep time.Duration, path string, command string, args ...string) {
	fmt.Println("[Installer] Running: ", command, args)
	runCommand := exec.Command(command, args...)

	// set path to command if passed in
	if path != "" {
		runCommand.Path = path
	}

	createPipeScanners(runCommand, command)

	if err := runCommand.Start(); err != nil {
		fmt.Println(err)
	}

	if sleep != 0 {
		fmt.Printf("[Installer] Sleeping: %s for %f seconds \n", command, sleep.Seconds())
		time.Sleep(sleep)
	}

	if err := runCommand.Wait(); err != nil {
		fmt.Println("Command finished with error: %v", err)
	}
}

// Created stdout, and stderr pipes for given *Cmd
// Only works with cmd.Start()
func createPipeScanners(cmd *exec.Cmd, prefix string) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Created scanners for in, out, and err pipes
	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	// Scan for text
	go func() {
		for errScanner.Scan() {
			fmt.Printf("[%s] %s\n", prefix, errScanner.Text())
		}
	}()

	go func() {
		for outScanner.Scan() {
			fmt.Printf("[%s] %s\n", prefix, outScanner.Text())
		}
	}()
}

func createCommandsFromJSON(jsonPath string) ([]*Cmd, error) {
	cmdFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal("Error:", err)
		return nil, err
	}

	// Unmarshal
	var cmds []*Cmd
	err = json.Unmarshal(cmdFile, &cmds)
	if err != nil {
		log.Fatalf("Can't unmarshal cmdFile.: %s", cmdFile)
		return nil, err
	}

	return cmds, nil
}
