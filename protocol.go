package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

type NetOp uint16

const (
	opChunk NetOp = iota
	opEndChunks
	opStartPop
	opStartPush
)

const NetHeaderMagic string = "MSGS"

type NetHeader struct {
	Magic    Magic32
	Version  uint32
	Op       NetOp
	DataSize uint32
}

func NewNetHeader(op NetOp, dataSize uint32) NetHeader {
	magic, err := NewMagic32(NetHeaderMagic)
	if err != nil {
		panic("Creating NetHeader magic failed!")
	}
	return NetHeader{
		Magic:    magic,
		Version:  programVersion,
		Op:       op,
		DataSize: dataSize,
	}
}

type Magic32 [4]byte

func NewMagic32(magic string) (Magic32, error) {
	var val Magic32
	if len(magic) != 4 {
		return val, errors.New("Magic string has wrong size!")
	}
	copy(val[:], []byte(magic))
	return val, nil
}

func (m Magic32) Verify(magic string) bool {
	return string(m[:]) == magic
}

type TxWriter struct {
	writer io.Writer
}

func NewTxWriter(writer io.Writer) TxWriter {
	return TxWriter{
		writer: writer,
	}
}

func (w *TxWriter) WriteHeader(op NetOp, extraData []byte) error {
	err := binary.Write(w.writer, binary.BigEndian, NewNetHeader(op, uint32(len(extraData))))
	if err != nil {
		return err
	}

	n, err := w.writer.Write(extraData)
	if err != nil {
		return err
	}

	_ = n
	return nil
}

func (w *TxWriter) WriteData(reader io.Reader) error {
	bufReader := bufio.NewReaderSize(reader, 1024*1024)
	for {
		if _, err := bufReader.Peek(1); err == io.EOF {
			break
		}

		size := bufReader.Buffered()
		err := binary.Write(w.writer, binary.BigEndian, NewNetHeader(opChunk, uint32(size)))
		if err != nil {
			return err
		}

		n, err := io.CopyN(w.writer, bufReader, int64(size))
		if err != nil {
			return err
		}

		_ = n
	}

	err := w.SkipData()
	return err
}

func (w *TxWriter) SkipData() error {
	err := binary.Write(w.writer, binary.BigEndian, NewNetHeader(opEndChunks, 0))
	return err
}

type RxReader struct {
	reader *bufio.Reader
}

func NewRxReader(reader io.Reader) RxReader {
	return RxReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *RxReader) ReadHeader() (NetHeader, []byte, error) {
	var netHeader NetHeader
	if err := binary.Read(r.reader, binary.BigEndian, &netHeader); err != nil {
		return netHeader, nil, err
	}

	extraData := make([]byte, netHeader.DataSize)
	if _, err := io.ReadFull(r.reader, extraData); err != nil {
		return netHeader, nil, err
	}

	return netHeader, extraData, nil
}

func (r *RxReader) NextData() (io.Reader, error) {
	var netHeader NetHeader
	if err := binary.Read(r.reader, binary.BigEndian, &netHeader); err != nil {
		return nil, err
	}

	if !netHeader.Magic.Verify(NetHeaderMagic) {
		return nil, errors.New("Invalid magic value!")
	}

	switch netHeader.Op {
	case opChunk:
		return io.LimitReader(r.reader, int64(netHeader.DataSize)), nil
	case opEndChunks:
		r.reader.Discard(int(netHeader.DataSize))
		return nil, io.EOF
	default:
		return nil, errors.New("Invalid operation for chunk!")
	}
}

func (r *RxReader) SkipData() error {
loop:
	for {
		var netHeader NetHeader
		if err := binary.Read(r.reader, binary.BigEndian, &netHeader); err != nil {
			return err
		}

		if !netHeader.Magic.Verify(NetHeaderMagic) {
			return errors.New("Stream corrupted: found invalid magic value!")
		}

		r.reader.Discard(int(netHeader.DataSize))
		switch netHeader.Op {
		case opChunk:
			continue
		case opEndChunks:
			break loop
		default:
			return errors.New("Stream corrupted: found invalid operation for chunk!")
		}
	}
	return nil
}
