package main

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

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

func HandleConnectionPop(cReader *bufio.Reader, cWriter *bufio.Writer, dirMutex *sync.Mutex) error {
    return nil;
}

func HandleConnectionPush(cReader *bufio.Reader, cWriter *bufio.Writer, dirMutex *sync.Mutex) error {
    var err error

    var file *os.File
	func() {
		dirMutex.Lock()
		defer dirMutex.Unlock()

		file, err = os.Create(MakeFilename(".", time.Now()))
	}()
    if err != nil {
        return err
    }
    defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = cReader.WriteTo(writer)

	if err != nil {
		return err
	}

	writer.Flush()
    return nil
}

func HandleConnection(conn net.Conn, dirMutex *sync.Mutex) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Could not handle connection! %s\n", r)
		}
	}()

    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)
    var err error

    var netHeader NetHeader
    err = binary.Read(reader, binary.BigEndian, &netHeader)
    if err != nil {
        panic(err.Error())
    }

    switch netHeader.Op {
        case opPop:
            err = HandleConnectionPop(reader, writer, dirMutex)
        case opPush:
            err = HandleConnectionPush(reader, writer, dirMutex)
    }
    if err != nil {
        panic(err.Error())
    }
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