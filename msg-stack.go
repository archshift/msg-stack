package main

import "bufio"
// import "flag"
import "fmt"
// import "hash/crc32"
import "net"
import "os"

func program_pb() {
	
}

func program_pop() {

}

func program_push() {
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
		writer.Flush()
	}
}

func handle_connection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(os.Stdout)
	n, err := reader.WriteTo(writer)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reading %d bytes: %s\n", n, err.Error())
	}
	
	writer.Flush()
}

func program_serve() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		go handle_connection(conn)
	}
}

func select_program(program_name string) {
	switch program_name {
	case "pb":
		program_pb()
	case "pop":
		program_pop()
	case "push":
		program_push()
	case "serve":
		program_serve()
	}
}

func main() {
	var program_name string

	if len(os.Args) < 2 {
		fmt.Println("Insert help message here!")
		return
	} else {
		program_name = os.Args[1]
	}

	select_program(program_name)
}
