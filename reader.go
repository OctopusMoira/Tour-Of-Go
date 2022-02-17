package main

import "golang.org/x/tour/reader"

type MyReader struct{}
type MyReaderError bool

// TODO: Add a Read([]byte) (int, error) method to MyReader.

func (MyReaderError) Error() string{
	return "b is nil"
}

func (MyReader) Read(b []byte) (int, error) {
	if cap(b) < 1 {
		return 0, MyReaderError(true)
	}
	for i := range b {
		b[i] = 65
	}
	return len(b), nil
}

func main() {
	reader.Validate(MyReader{})
}

