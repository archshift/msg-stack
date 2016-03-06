package main

import (
    "encoding/binary"
    "strings"
)

type NetOp uint16

const (
	opPop NetOp = iota
	opPush
)

type NetHeader struct {
    Magic   Magic32
	Version uint32
	Op      NetOp
}

func NewNetHeader(op NetOp) NetHeader {
    magic, err := NewMagic32("MSGS")
    if err != nil {
        panic("Creating NetHeader magic failed!")
    }
    return NetHeader {
        Magic: magic,
        Version: programVersion,
        Op: op,
    }
}

type Magic32 uint32

func NewMagic32(magic string) (Magic32, error) {
    var val Magic32
    err := binary.Read(strings.NewReader(magic), binary.LittleEndian, &val)
    return val, err
}