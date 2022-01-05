package main


import (
	"fmt"
	"log"
	"os"
	"encoding/base64"
)

// keySize is a struct containing a keysize and the
// associated normalized hamming distance for the
// the average of all block pairs in the source of
// size equal to the keysize.
type keySize struct {
	size int
	dist float32
}

// hammingDist computes the number of differing bits in two
// byte slices as long as they are of equal length.
func hammingDist(strOne, strTwo []byte) (dist byte, err error) {
	if len(strOne) != len(strTwo) {
		return dist, fmt.Errorf("hamming distance can not be calculated for strings of differing length")
	}

	var scratch byte
	for i, _ := range strOne {
		scratch = strOne[i] ^ strTwo[i]
		for j := 0; j < 8; j++ {
			if scratch & (1 << j) > 0 {
				dist++
			}
		}
	}

	return dist, nil	
}

// mostProbSizes computes three key sizes between input minimum
// and maximum sizes such they are the key sizes which result
// in the smallest normalized hamming distance for each pair
// of blocks of length keysize.
func mostProbSizes(src []byte, minSize, maxSize int) (sizes [3]keySize, err error) {
	for size := minSize; size <= maxSize; size++ {
		// find number of pairs of blocks
		// we don't really care about precisely using every byte
		blockPairs := len(src) / (size*2)
		if blockPairs > 128 {
			blockPairs = 128
		}

		var sum float32
		for p := 0; p < 4; p++ {
			dist, err := hammingDist(src[p*(size*2):(p*(size*2))+size], src[(p*(size*2))+size:(p*(size*2))+(2*size)])
			if err != nil {
				return sizes, err
			}
			sum += float32(dist)
		}

		// normalize result by dividing by keysize
		sum /= float32(size)
		// replace in array of smallest three distances
		for i, s := range sizes {
			if sum < s.dist || s.dist == 0 {
				sizes[i] = keySize{size: size, dist: sum}
				break
			}
		}
	}

	return
}

// chuckCipherText simply splits a single slice of source
// bytes into a slice of slices of bytes such that each
// slice contains a number of bytes equal to the keysize.
// The last slice may be padded with 0's if the length
// of the source is not evenly divisible by the keysize.
func chunkCipherText(src []byte, keysize int) [][]byte {
	nBlk := len(src) / keysize
	if len(src) % keysize != 0 {
		nBlk++
	}

	blocks := make([][]byte, nBlk)
	for i := 0; i < nBlk-1; i++ {
		blocks[i] = make([]byte, keysize)
		copy(blocks[i], src[i*keysize:(i+1)*keysize])
	}
	// copy last block
	blocks[len(blocks)-1] = make([]byte, keysize)
	copy(blocks[len(blocks)-1], src[(nBlk-1)*keysize:])

	return blocks
}

// transpose performs a simple transpose of a byte matrix
func transpose(blocks [][]byte) [][]byte {
	// prepare transposed matrix of proper size
	tr := make([][]byte, len(blocks[0]))
	for r := range tr {
		tr[r] = make([]byte, len(blocks))
		// and populate with values
		for c := range tr[r] {
			copy(tr[r][c:c+1], blocks[c][r:r+1])
		}
	}

	return tr
}

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

// decryptSingleCharXor iterates through each possible ascii
// character and performs a single-char XOR operation on the
// input byte slice. For each key character used, it checks if
// the output has a higher score (based on English character
// frequency) than the previous best. Finally it returns the key
// and decoded line with the highest value.
func decryptSingleCharXor(src []byte, scoringMap map[byte]float32) (key byte) {
	var score float32
	scratch := make([]byte, len(src))
	// assume key is a single ascii character
	for k := 32; k < 127; k++ {
		var keyScore float32
		// inverse xor with this key and sum score
		copy(scratch, src)
		for i := range scratch {
			scratch[i] ^= byte(k)
			if ((scratch[i] > 64) && (scratch[i] < 91)) {
				// convert to lowercase for the purpose of scoring
				keyScore += scoringMap[scratch[i]+32]
			} else {
				keyScore += scoringMap[scratch[i]]
			}
		}

		// compare score with previous key
		// if new score is higher, replace previous
		if keyScore > score {
			score = keyScore
			key = byte(k)
		}
	}

	return
}


// decryptRepeatingCharXor iterates between the bytes of key as
// performs XOR on each byte of the source.
func decryptRepeatingCharXor(src, key []byte) []byte {
	dst := make([]byte, len(src))
	for i := range src {
		dst[i] = src[i] ^ key[i%len(key)]
	}

	return dst
}

func main() {
	// 1. read and decode base64 encoded file
	src, err := os.ReadFile("6.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = base64.StdEncoding.Decode(src, src)
	if err != nil {
		log.Fatal(err)
	}

	// 2. guess keysizes between 2 and 40
	// 3. for each keysize calculate normalized hamming distance
	// 4. proceed with 3 keys with smallest normalized hamming distance
	keys, err := mostProbSizes(src, 2, 40)
	if err != nil {
		log.Fatal(err)
	}


	// proceed with three best keysize candidates
	for _, k := range keys {
		//fmt.Printf("keysize of %v with normalized hamming distance of %v\n", k.size, k.dist)
		// 5. break ciphertext into blocks of length keysize
		blocks := chunkCipherText(src, k.size)
		// 6. transpose blocks
		tr := transpose(blocks)
		// 7. solve each block as if it were a single-char XOR
		key := make([]byte, len(tr))
		for i, blk := range tr {
			key[i] = decryptSingleCharXor(blk, englishCharFrequency)
		}
		// 8. combine single byte keys to create full key and test decryption
		// only print the first few characters for clarity
		fmt.Printf("key: %s\n", key)
		fmt.Printf("decrypted text:\n%s\n\n", decryptRepeatingCharXor(src, key)[:100])
	}


	// we can make an educated guess based on the result from the operations above
	// we see the second best key with a keysize of 29 chars resulted in text that
	// was nearly English
	// the key was: TER(IN$TOR X: BRIN" TH  NOI6E
	// testing with: Terminator x: Bring the noise

	keyGuess := "Terminator X: Bring the noise"
	fmt.Printf("guessed key: %v\n", keyGuess)
	fmt.Printf("decrypted text:\n%s\n\n", decryptRepeatingCharXor(src, []byte(keyGuess)))
}
