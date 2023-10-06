package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/srebuff/paas_exporter/cmd"
	"os"
	"runtime"
)

var buildstamp = ""
var githash = ""
var goversion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "version" || args[1] == "--version" || args[1] == "-v") {
		fmt.Println(text.FgBlue.Sprintf("gops - Git Commit Hash: %s", githash))
		fmt.Println(text.FgBlue.Sprintf("Build Time : %s", buildstamp))
		fmt.Println(text.FgBlue.Sprintf("Golang Version : %s", goversion))
		return
	}
	cmd.Execute()
}
