package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var types map[string]bool
var typeList string

func init() {
	types = make(map[string]bool)
	types["build"] = true
	types["ci"] = true
	types["docs"] = true
	types["feat"] = true
	types["fix"] = true
	types["perf"] = true
	types["refactor"] = true
	types["style"] = true
	types["test"] = true

	sl := make([]string, 0)
	for k := range types {
		sl = append(sl, k)
	}

	typeList = strings.Join(sl, ",")
}

func main() {
	var scope = flag.String("s", "", "scope, like modulea, moduleb")
	var changeType = flag.String("t", "fix", fmt.Sprintf("change of type: %s", typeList))
	flag.Parse()

	if !types[*changeType] {
		fmt.Printf("changeType must be in: %s\n", typeList)
		return
	}
	if flag.NArg() < 1 {
		fmt.Printf("Must specifiy the head")
		return
	}

	var args []string
	args = append(args, "commit", "-m")
	if *scope != "" {
		args = append(args, fmt.Sprintf("%s(%s): %s", *changeType, *scope, flag.Arg(0)))
	} else {
		args = append(args, fmt.Sprintf("%s: %s", *changeType, flag.Args()[0]))
	}

	if flag.NArg() > 1 {
		args = append(args, flag.Args()[1:]...)
	}
	if ret, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		fmt.Printf("%s\n", ret)
		log.Fatal(err)
	}
}
