package main


import (
        "fmt"
        "log"
        "os"
        "encoding/base64"
        "crypto/aes"
)

func main() {
        // read and decode base64 encoded file
        src, err := os.ReadFile("7.txt")
        if err != nil {
                log.Fatal(err)
        }
        _, err = base64.StdEncoding.Decode(src, src)
        if err != nil {
                log.Fatal(err)
        }

        key := "YELLOW SUBMARINE"

        // pad the source to ensure it is a
        // multiple of the block size
        pad := len(src) % len(key)
        if pad != 0 {
                for i := 0; i < len(key)-pad; i++ {
                        src = append(src, 0)
                }
        }

        // create AES cipher
        c, err := aes.NewCipher([]byte(key))
        if err != nil {
                log.Fatal(err)
        }

        // find how many blocks
        blks := len(src) / len(key)

        for j := 0; j < blks; j++ {
                c.Decrypt(src[j*len(key):(j+1)*len(key)],
                        src[j*len(key):(j+1)*len(key)])
        }

        fmt.Printf("%s", src)
}
