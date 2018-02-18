package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	var cmd *exec.Cmd

	if err := exec.Command("which", "rlwrap").Run(); err == nil {
		args := []string{"ijq-bare"}
		args = append(args, os.Args[1:]...)
		cmd = exec.Command("rlwrap", args...)
	} else {
		cmd = exec.Command("ijq-bare", os.Args[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
