package main

import (
	"os"
	"os/exec"
)

func executeCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		//log.Fatal("This error caused fatality ", err)
		return err
	}

	return nil
}

func main() {

}
