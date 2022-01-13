package main


import (
        "fmt"
)

func pkcs7(src []byte, blkSize int) []byte {
        // find by how much we need to pad
        pad := blkSize - (len(src) % blkSize)
        if pad == blkSize {
                // no padding required
                return src
        }

        // pad
        padSlice := make([]byte, pad)
        for i := range padSlice {
                padSlice[i] = byte(pad)
        }

        dst := make([]byte, len(src) + pad)
        copy(dst, src)
        copy(dst[len(src):], padSlice)
        return dst
}

func main() {
        src := "YELLOW SUBMARINE"
        padded := pkcs7([]byte(src), 20)
        fmt.Printf("%s\n", padded)
        fmt.Printf("%v\n", padded)
}
