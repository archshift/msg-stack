package main

import (
	"bufio"
	// "flag"
	"fmt"
	// "hash/crc32"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const programVersion = 0

type NetOp uint16

const (
	opPop NetOp = iota
	opPush
)

type NetHeader struct {
	version uint32
	op      NetOp
}

func ProgramPb() {

}

func ProgramPop() {

}

func ProgramPush() {
	scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial("tcp", "localhost:8080")
	writer := bufio.NewWriter(conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	for scanner.Scan() {
		nn, err := writer.Write([]byte(scanner.Text()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Writing %d bytes: %s\n", nn, err.Error())
		}
		writer.WriteRune('\n')
		writer.Flush()
	}
}

func MakeFilename(dirname string, timestamp time.Time) string {
	timestampStr := strconv.FormatInt(timestamp.Unix(), 10)

	dir, err := os.Open(dirname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Opening dir %s: %s\n", dirname, err.Error())
		panic("")
	}
	defer dir.Close()

	dirChildren, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reading dir %s: %s\n", dirname, err.Error())
		panic("")
	}

	var highestNum int
	for index := range dirChildren {
		if strings.HasPrefix(dirChildren[index], timestampStr) {
			fileNumStr := strings.SplitN(dirChildren[index], "_", 2)[1]
			n, err := strconv.ParseInt(fileNumStr, 10, 32)
			if err == nil && int(n) > highestNum {
				highestNum = int(n)
			}
		}
	}

	return fmt.Sprintf("%s/%s_%d", dirname, timestampStr, highestNum+1)
}

func HandleConnection(conn net.Conn, dirMutex *sync.Mutex) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Could not handle connection!\n")
		}
	}()

	var file *os.File
	var err error
	func() {
		dirMutex.Lock()
		defer dirMutex.Unlock()

		file, err = os.Create(MakeFilename(".", time.Now()))
		if err != nil {
			panic("")
		}
	}()

	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(conn)
	_, err = reader.WriteTo(writer)

	if err != nil {
		panic("")
	}

	writer.Flush()
}

func ProgramServe() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	dirMutex := &sync.Mutex{}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		go HandleConnection(conn, dirMutex)
	}
}

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
