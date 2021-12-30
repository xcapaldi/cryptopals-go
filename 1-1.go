package main


import (
        "fmt"
        "encoding/hex"
        "encoding/base64"
        "log"
)

func convertHexToBase64(src string) (string, error) {
        // decode from hex
        b, err := hex.DecodeString(src)
        if err != nil {
                return "", err
        }
        // encode to base64
        return base64.StdEncoding.EncodeToString(b), nil
}

func main() {
        src := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
        dst, err := convertHexToBase64(src)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("Base64 encoding:\n" , dst)
}
