package main

import (
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

	txWriter := NewTxWriter(conn)

	if err := txWriter.WriteHeader(opStartPush, make([]byte, 0)); err != nil {
		fmt.Fprintf(os.Stderr, "Writing header: %s\n", err.Error())
		return
	}
	if err := txWriter.WriteData(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "Writing data: %s\n", err.Error())
		return
	}
}
