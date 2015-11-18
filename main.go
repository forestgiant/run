package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jesselucas/executil"
)

const VERSION = "0.0.3"

var (
	wg sync.WaitGroup
)

type cmd struct {
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
		versionUsage    = "Prints current version" + " (v. " + VERSION + ")"
		versionPtr      = flag.Bool("version", false, versionUsage)
	)
	// Set up short hand flags
	flag.BoolVar(versionPtr, "v", false, versionUsage+" (shorthand)")
	flag.Parse()

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}

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
	err := executil.CmdStart(command, args...)
	if err != nil {
		fmt.Println(err)
	}

	if sleep != 0 {
		fmt.Printf("[Installer] Sleeping: %s for %f seconds \n", command, sleep.Seconds())
		time.Sleep(sleep)
	}
}

func createCommandsFromJSON(jsonPath string) ([]*cmd, error) {
	cmdFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal("Error:", err)
		return nil, err
	}

	// Unmarshal
	var cmds []*cmd
	err = json.Unmarshal(cmdFile, &cmds)
	if err != nil {
		log.Fatalf("Can't unmarshal cmdFile.: %s", cmdFile)
		return nil, err
	}

	return cmds, nil
}
