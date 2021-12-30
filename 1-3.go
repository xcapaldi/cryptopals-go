package main


import (
        "fmt"
        "encoding/hex"
        "log"
)

var englishCharFrequency = map[byte]float32 {
        97:8.34,
        98:1.54,
        99:2.73,
        100:4.14,
        101:12.60,
        102:2.03,
        103:1.92,
        104:6.11,
        105:6.71,
        106:0.23,
        107:0.87,
        108:4.24,
        109:2.53,
        110:6.80,
        111:7.70,
        112:1.66,
        113:0.09,
        114:5.68,
        115:6.11,
        116:9.37,
        117:2.85,
        118:1.06,
        119:2.34,
        120:0.20,
        121:2.04,
        122:0.06,
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
