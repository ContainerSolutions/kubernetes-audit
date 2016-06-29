package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	f, err := os.Create("audit.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating audit log file")
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println("WELCOME TO THE KUBENETES AUDITOR! IT DOESN'T WORK YET I'M AFRAID")

	for {
		fmt.Println("tick")
		time.Sleep(1)

	}
}
