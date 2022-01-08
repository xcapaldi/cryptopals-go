package main


import (
        "fmt"
        "log"
        "os"
        "encoding/hex"
        "bufio"
        "bytes"
)

func scoreLine(line []byte, blockSize int) (score int) {
        nBlks := len(line)/blockSize
        // normally we would care about the last block but
        // in this case, the last block will already have been
        // compared with all prior blocks so no need to pad
        for i := 0; i < nBlks - 1; i++ {
                // only need to compare with blocks ahead
                for  j := i + 1; j < nBlks - 1; j++ {
                        if bytes.Compare(
                                line[i*blockSize:(i+1)*blockSize],
                                line[j*blockSize:(j+1)*blockSize]) == 0 {
                                score++
                        }
                }

                // deal with last block separately
                if bytes.Compare(
                        line[i*blockSize:(i+1)*blockSize],
                        line[(nBlks-2)*blockSize:]) == 0 {
                        score++
                }
        }

        return
}

func main() {
        file, err := os.Open("8.txt")
        if err != nil {
                log.Fatal(err)
        }
        var score int
        scanner := bufio.NewScanner(file)
        scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
                decodedLine, err := hex.DecodeString(scanner.Text())
                if err != nil {
                        log.Fatal(err)
                }
                if lScore := scoreLine(decodedLine, 16); lScore > score {
                        score = lScore
                        fmt.Printf("%v duplicate blocks in the following line: %v", score, decodedLine)
                }
        }
}
