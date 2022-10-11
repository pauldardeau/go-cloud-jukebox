package jukebox

import (
    "bytes"
    "compress/zlib"
    "fmt"
    "io"
)

func CompressBuffer(uncompressed []byte) ([]byte, error) {
    level := zlib.BestCompression
    var buffer bytes.Buffer
    w, e := zlib.NewWriterLevel(&buffer, level)
    if e != nil {
        fmt.Printf("error: unable to create new zlib writer for level=%d\n", level)
        fmt.Printf("error: %v\n", e)
	return nil, e
    } else {
        w.Write(uncompressed)
        w.Close()
	return buffer.Bytes(), nil
    }
}

func UncompressBuffer(compressed []byte) ([]byte, error) {
    b := bytes.NewReader(compressed[:])
    r, err := zlib.NewReader(b)
    if err != nil {
        fmt.Printf("error: unable to create new zlib reader\n")
	fmt.Printf("error: %v\n", err)
	return nil, err
    } else {
        var out bytes.Buffer
        io.Copy(&out, r)
	return out.Bytes(), nil
    }
}

