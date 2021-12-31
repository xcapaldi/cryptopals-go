package main


import (
        "fmt"
        "log"
        "encoding/hex"
        "os"
        "bufio"
)

// englishCharFrequency holds the ascii (or utf-8)
// byte value for lowercase letters in the English
// language along with their corresponding percentage
// frequency. To access the value of uppercase
// characters, one can just modify the byte value.
// Uppercase characters have values from 65 to 90.
// Lowercase characters have value from 97 to 122.
// All other values hold symbols that we wouldn't
// want to modify. So to get the frequency of an
// uppercase character, simply add 32 to the byte
// value before looking up in the map.
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

// readLines takes in a filepath and reads lines from it
// one-by-one (appending to an output list) until it
// reaches the end of the file. Before appending, each
// line is decoded from hex into a byte slice for easy
// operation later.
func readLines(path string) (lines [][]byte, err error) {
        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
                decodedLine, err := hex.DecodeString(scanner.Text())
                if err != nil {
                        return nil, err
                }
                lines = append(lines, decodedLine)
        }
        if err = scanner.Err(); err != nil {
                return nil, err
        }

        return lines, nil
}

// scoreLine iterates through every possible ascii character
// and performs a single-char XOR operation on the input line.
// For each key character used, it checks if the output has a
// higher score (based on English character frequency) than the
// previous best. Finally it returns the key and decoded line
// with the highest value. We are more thorough than 1-3 and
// actually check full range of possible keys.
func scoreLine(line []byte, scoreMap map[byte]float32) (key byte, decoded []byte, score float32) {
        decoded = make([]byte, len(line))
        scratch := make([]byte, len(line))
        for k := 32; k < 127; k++ {
                copy(scratch, line)
                var kScore float32

                // single-char XOR and sum score
                for i, _ := range scratch {
                        scratch[i] ^= byte(k)
                        if ((scratch[i] > 64) && (scratch[i] < 91)) {
                                // convert uppercase to lowercase
                                kScore += scoreMap[scratch[i]+32]
                        } else {
                                kScore += scoreMap[scratch[i]]
                        }
                }

                // compare score with previous key and supplant if higher
                if kScore > score {
                        score = kScore
                        key = byte(k)
                        copy(decoded, scratch)
                }
        }

        return key, decoded, score
        }


// scoreLines scores each line by finding the single-char key
// for XOR cipher that results in the highest scoring (based
// on character frequency in English language).
func scoreLines(lines [][]byte, scoreMap map[byte]float32) (bestKey byte, decoded []byte, linum int) {
        decoded = make([]byte, len(lines[0]))
        var bestScore float32

        for i, line := range lines {
                key, scratch, score := scoreLine(line, scoreMap)
                if score > bestScore {
                        bestScore = score
                        bestKey = key
                        copy(decoded, scratch)
                        linum = i
                }

        }

        return

}

func main () {
        path := "4.txt"
        // read in the lines and decode from hex
        lines, err := readLines(path)
        if err != nil {
                log.Fatal(err)
        }
        // score the lines finding the most likely
        // encoded
        bestKey, decoded, linum := scoreLines(lines, englishCharFrequency)
        fmt.Printf("line %v was encoding via single-character XOR cipher with %s as the key. The decoded line is: %s\n", linum, string([]byte{bestKey}), string(decoded))
}
