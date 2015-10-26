package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

var (
	errors chan error
	wg     sync.WaitGroup
)

func main() {
	// Establish our exit channels
	errors = make(chan error)

	run(500, "ls", "-all", "-G")

}

func run(sleep time.Duration, command string, args ...string) {
	sleep = sleep * time.Millisecond
	wg.Add(1)
	fmt.Println("[Installer] Running: ", command, args)
	runCommand := exec.Command(command, args...)
	createPipeScanners(runCommand, command)

	if err := runCommand.Start(); err != nil {
		errors <- err
	}

	if sleep != 0 {
		fmt.Printf("[Installer] Sleeping: %s for %f seconds \n", command, sleep.Seconds())
		time.Sleep(sleep)
	}
	wg.Done()

	if err := runCommand.Wait(); err != nil {
		errors <- err
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
