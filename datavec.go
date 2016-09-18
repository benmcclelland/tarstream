package tarstream

import (
	"io"
	"os"
)

// MemVec is a buffer vec type
type memVec struct {
	Data []byte
	Good bool
}

// PathVec is a filename vec type
type pathVec struct {
	Path string
	info os.FileInfo
	file *os.File
}

// PadVec is a padding (0s) vec type
type padVec struct {
	Size int64
}

// Datavec is an interface for all vector types
type Datavec interface {
	GetSize() int64
	Open() error
	Close()
	ReadAt(b []byte, off int64) (int, error)
}

// GetSize gets the size of the memory vec
func (m memVec) GetSize() int64 {
	return int64(len(m.Data))
}

// Open opens a memory vec
func (m memVec) Open() error {
	return nil
}

// Close closes the memory vec
func (m memVec) Close() {
}

// ReadAt reads at an offset of a memory vec
func (m memVec) ReadAt(b []byte, off int64) (int, error) {
	var end int64
	if int64(len(m.Data))-off > int64(len(b)) {
		end = off + int64(len(b))
	} else {
		end = off + int64(len(m.Data))
	}
	if end > int64(len(m.Data)) {
		end = int64(len(m.Data))
	}

	n := copy(b, m.Data[off:end])
	if n == 0 {
		return n, io.EOF
	}
	return n, nil
}

//GetSize gets the file size of the path vec
func (p pathVec) GetSize() int64 {
	return p.info.Size()
}

// Open opens a file represented by a path vec
func (p *pathVec) Open() error {
	var err error
	p.file, err = os.Open(p.Path)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the file represented by the path vec
func (p *pathVec) Close() {
	p.file.Close()
}

// ReadAt reads the file represented by path vec at the given offset
func (p *pathVec) ReadAt(b []byte, off int64) (int, error) {
	n, err := p.file.ReadAt(b, off)
	if err == io.EOF {
		return n, nil
	}
	if n == 0 {
		return n, io.EOF
	}
	return n, err
}

// GetSize gets the size of the padding vec
func (p padVec) GetSize() int64 {
	return p.Size
}

// Open opens the padding vec
func (p padVec) Open() error {
	return nil
}

// Close closes the padding vec
func (p padVec) Close() {
}

// ReadAt read the padding vec at a given offset (which is always 0s)
func (p padVec) ReadAt(b []byte, off int64) (int, error) {
	n := int(p.Size - off)

	if n > len(b) {
		n = len(b)
	}

	if n == 0 {
		return 0, io.EOF
	}
	b = make([]byte, n)
	return n, nil
}
