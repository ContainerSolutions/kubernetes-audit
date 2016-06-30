package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {

	f, err := os.Create("audit.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating audit log file")
		os.Exit(1)
	}
	defer f.Close()

	resp, err := http.Get("http://localhost:8080/api/v1/events")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get events from API server: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading event: %v\n", err)
		}
		fmt.Fprintln(os.Stdout, "%v", line)
	}

}
