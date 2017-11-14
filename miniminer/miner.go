package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

// i forgot how to write code in golang :)
func go_tutorial() {
	example := []byte(`{"data":[],"nonce":45}`)
	sha := sha256.Sum256(example)
	hexstr := hex.EncodeToString(sha[:])

	// it works!
	fmt.Println(hexstr)
}

var TEST_JSON = `{ "block": { "nonce": null, "data":[] }, "difficulty":8}`
var USE_TEST_JSON = false

func main() {
	// read stdin
	in := bufio.NewReader(os.Stdin)
	challenge, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if USE_TEST_JSON {
		challenge = TEST_JSON
	}

	// test read string is json
	var js map[string]interface{}
	err = json.Unmarshal([]byte(challenge), &js)
	if err != nil {
		panic(err)
	}

	fmt.Println("challenge:", challenge)
	//fmt.Println("json:", js)
	//test, _ := json.Marshal(js)
	//fmt.Println("json string:", string(test))

	block, _ := json.Marshal(js["block"])
	difficulty := int(js["difficulty"].(float64))

	startMining(string(block), difficulty)
}

func startMining(block string, difficulty int) {
	fmt.Println("start mining")

	fmt.Println("block", block)
	fmt.Println("diffy", difficulty)

	block1 := block[:len(block)-5]
	block2 := block[len(block)-1:]

	var nonce int
found:
	for i := 0; ; i++ {
		t := block1 + strconv.Itoa(i) + block2
		//fmt.Println(t)
		t2 := sha256.Sum256([]byte(t))

		var zeros = 0
		for _, v := range t2 {
			z := bits.LeadingZeros8(v)
			zeros += z
			//fmt.Println("z", z, "zero", zeros)

			if zeros == difficulty {
				nonce = i
				break found
			}
			if zeros > difficulty {
				break
			}
			if z != 8 {
				break
			}
		}
	}

	fmt.Println("nonce found:", nonce)
}
