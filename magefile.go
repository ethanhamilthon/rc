//go:build mage

package main

import (
	"fmt"
	"os/exec"
)

// Build compiles the project.
func Build() error {
	if err := Lint(); err != nil {
		return err
	}
	fmt.Println("Building the project...")
	return exec.Command("go", "build", "./...").Run()
}

// Test runs the tests.
func Test() error {
	fmt.Println("Running tests...")
	return exec.Command("go", "test", "./...").Run()
}

func Lint() error {
	fmt.Println("Linting the project...")
	cmd := exec.Command("golangci-lint", "run")

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	return err
}

func Move() error {
	if err := Build(); err != nil {
		return err
	}
	return exec.Command("sudo", "mv", "./hc", "/usr/local/bin/hc").Run()
}
