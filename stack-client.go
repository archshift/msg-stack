package main

import (
    "encoding/binary"
    "bufio"
    "fmt"
    "net"
    "os"
)

func ProgramPop() {

}

func ProgramPush() {
	conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        return
	}

    reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

    err1 := binary.Write(writer, binary.BigEndian, NewNetHeader(opPush))
    if err1 != nil {
        fmt.Fprintf(os.Stderr, "Writing header: %s\n", err1.Error())
        return
    }

	n, err2 := reader.WriteTo(writer)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Writing %d bytes: %s\n", n, err2.Error())
        return
	}

    writer.Flush()
}