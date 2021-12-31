package main


import (
        "fmt"
)

func encrypt(src []byte, key []byte) {
        for i, _ := range src {
                src[i] ^= key[i%len(key)]
        }
}

func main() {
        src := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
        bSrc := []byte(src)
        key := "ICE"

        encrypt(bSrc, []byte(key))
        fmt.Printf("%x", bSrc)
}
