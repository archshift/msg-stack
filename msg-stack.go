package main

import (
	// "flag"
	"fmt"
	"os"
)

const programVersion = 0

func SelectProgram(programName string) {
	switch programName {
	case "pb":
		ProgramPb()
	case "pop":
		ProgramPop()
	case "push":
		ProgramPush()
	case "serve":
		ProgramServe()
	}
}

func main() {
	var programName string

	if len(os.Args) < 2 {
		fmt.Println("Insert help message here!")
		return
	} else {
		programName = os.Args[1]
	}

	SelectProgram(programName)
}
