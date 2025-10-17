package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run %s: %v", name, err)
	}
}

func InitCmd() {
	fmt.Println(">>> Running swag init")
	runCmd("swag", "init")

	fmt.Println(">>> Running goose up")
	runCmd("goose", "up")
}
