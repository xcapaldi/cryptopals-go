package main


import (
        "fmt"
        "encoding/hex"
        "log"
)

var englishCharFrequency = map[byte]float32 {
        97:8.34,    // a
        98:1.54,    // b
        99:2.73,    // c
        100:4.14,   // d
        101:12.60,  // e
        102:2.03,   // f
        103:1.92,   // g
        104:6.11,   // h
        105:6.71,   // i
        106:0.23,   // j
        107:0.87,   // k
        108:4.24,   // l
        109:2.53,   // m
        110:6.80,   // n
        111:7.70,   // o
        112:1.66,   // p
        113:0.09,   // q
        114:5.68,   // r
        115:6.11,   // s
        116:9.37,   // t
        117:2.85,   // u
        118:1.06,   // v
        119:2.34,   // w
        120:0.20,   // x
        121:2.04,   // y
        122:0.06,   // z
}

func decode(src string, scoringMap map[byte]float32) (key byte, decoded []byte, err error) {
        // decode from hex
        d, err := hex.DecodeString(src)
        if err != nil {
                return byte(0), nil, err
        }
        decoded = make([]byte, len(d))
        var score float32
        // assume key is a single ascii character
        for k := 0; k < 126; k++ {
                var keyScore float32
                // inverse the xor with this key and sum the score
                for i, _ := range d {
                        d[i] ^= byte(k)
                        if d[i] < 91 {
                                // convert to lowercase for the purposes of scoring
                                keyScore += scoringMap[d[i]+32]
                        } else {
                                keyScore += scoringMap[d[i]]
                        }
                }

                // compare score with previous key
                // if new score is higher replace the "best" key and decoding
                if keyScore > score {
                        score = keyScore
                        key = byte(k)
                        copy(decoded, d)
                }
        }

        return key, decoded, nil
}

func main() {
        src := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

        key, decoded, err := decode(src, englishCharFrequency)
        if err != nil {
                log.Fatal(err)
        }

        fmt.Printf("key: %s\ndecoded message: %s\n", string([]byte{key}), string(decoded))
}
