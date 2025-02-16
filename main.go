package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		fmt.Println("🔹 Launching GomniRun CLI Mode...")
		runExecutable("cmd/cli/main.go")
	} else {
		fmt.Println("🔹 Launching GomniRun GUI Mode...")
		runExecutable("cmd/fyne-ui/main.go")
	}
}

// runExecutable compiles and runs a Go file dynamically
func runExecutable(file string) {
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start %s: %v", file, err)
	}
}
