package gracefully_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	simpleServerGo         = "./testdata/simple_server/main.go"
	simpleServerBin        = "./testdata/simple_server/main"
	afterShutdownServerGo  = "./testdata/after_shutdown/main.go"
	afterShutdownServerBin = "./testdata/after_shutdown/main"
	afterShutdownFile      = "./testdata/after_shutdown/after_shutdown.txt"
)

func setup() {
	cleanup()

	must(goCompile(simpleServerGo))
	must(goCompile(afterShutdownServerGo))
}

func cleanup() {
	os.Remove(simpleServerBin)
	os.Remove(afterShutdownServerBin)
	os.Remove(afterShutdownFile)
}

func goCompile(f string) error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(c, "go", "build", "-o", strings.TrimSuffix(f, ".go"), f)
	fmt.Printf("GOTTEN CMD :: %s", cmd.String())
	return cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
