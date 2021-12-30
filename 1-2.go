package main


import (
        "fmt"
        "encoding/hex"
        "log"
)

func xorHexStrings(a, b string) ([]byte, error) {
        // decode from hex
        aDecoded, err := hex.DecodeString(a)
        if err != nil {
                return nil, err
        }
        bDecoded, err := hex.DecodeString(b)
        if err != nil {
                return nil, err
        }

        // check that strings are same length
        if len(aDecoded) != len(bDecoded) {
                return nil, fmt.Errorf("hex strings are not equal length")
        }

        // XOR
        d := make([]byte, len(aDecoded))
        for i, _ := range d {
                d[i] = aDecoded[i] ^ bDecoded[i]
        }

        return d, nil
}

func main() {
        a := "1c0111001f010100061a024b53535009181c"
        b := "686974207468652062756c6c277320657965"

        dst, err := xorHexStrings(a, b)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("%x", dst)
}
