package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Hash [32]byte

var (
	serverHostname = "127.0.0.1:6262"
)

type Block struct {
	PrevHash Hash
	Name     string
	Nonce    string
}

func (self Block) ToString() string {
	return fmt.Sprintf("%x %s %s", self.PrevHash, self.Name, self.Nonce)
}

func (self Hash) ToString() string {
	return fmt.Sprintf("%x", self)
}

func (self Block) Hash() Hash {
	return sha256.Sum256([]byte(self.ToString()))
}

func BlockFromString(s string) (Block, error) {
	var bl Block

	if len(s) < 66 || len(s) > 100 {
		return bl, fmt.Errorf("Invalid string length %d, expect 66 to 100", len(s))
	}
	subStrings := strings.Split(s, " ")

	if len(subStrings) != 3 {
		return bl, fmt.Errorf("got %d elements, expect 3", len(subStrings))
	}

	hashbytes, err := hex.DecodeString(subStrings[0])
	if err != nil {
		return bl, err
	}
	if len(hashbytes) != 32 {
		return bl, fmt.Errorf("got %d byte hash, expect 32", len(hashbytes))
	}

	copy(bl.PrevHash[:], hashbytes)

	bl.Name = subStrings[1]

	bl.Nonce = strings.TrimSpace(subStrings[2])

	return bl, nil
}

func GetTipFromServer() (Block, error) {
	var bl Block

	connection, err := net.Dial("tcp", serverHostname)
	if err != nil {
		return bl, err
	}
	fmt.Printf("connected to server %s\n", connection.RemoteAddr().String())

	sendbytes := []byte("TRQ\n")

	_, err = connection.Write(sendbytes)
	if err != nil {
		return bl, err
	}

	bufReader := bufio.NewReader(connection)

	blockLine, err := bufReader.ReadBytes('\n')
	if err != nil {
		return bl, err
	}

	fmt.Printf("read from server:\n%s\n", string(blockLine))

	bl, err = BlockFromString(string(blockLine))
	if err != nil {
		return bl, err
	}

	return bl, nil
}

func SendBlockToServer(bl Block) error {
	connection, err := net.Dial("tcp", serverHostname)
	if err != nil {
		return err
	}
	fmt.Printf("connected to server %s\n", connection.RemoteAddr().String())

	sendbytes := []byte(fmt.Sprintf("%s\n", bl.ToString()))

	_, err = connection.Write(sendbytes)
	if err != nil {
		return err
	}

	bufReader := bufio.NewReader(connection)
	ResponseLine, err := bufReader.ReadBytes('\n')
	if err != nil {
		return err
	}

	fmt.Printf("Server resposnse: %s\n", string(ResponseLine))

	return connection.Close()
}

func (self Block) Mine(targetBits uint8) {
	nonce := 1
	self.Nonce = string(nonce)
	verified := CheckWork(self, targetBits)

	for !verified {
		nonce += 1
		self.Nonce = strconv.Itoa(nonce)
		verified = CheckWork(self, targetBits)
		fmt.Print("Not this time!")
	}
	if verified {
		SendBlockToServer(self)
		fmt.Printf("block nonce \n%c", nonce)
	}
}

func CheckWork(bl Block, targetBits uint8) bool {
	hash := bl.Hash()
	var v bool
	for i := 0; i < int(targetBits); i++ {
		if hash[i/8]>>uint(7-(i%8))&0x01 == 0 {
			v = true
		} else {
			v = false
		}
		if v == false {
			break
		}
	}
	return v
}

func PingServer(ch chan Block) {

	var oldtip Block
	var newtip Block

	oldtip, _ = GetTipFromServer()
	ch <- oldtip
	go func() {
		for {
			newtip, _ = GetTipFromServer()
			if newtip != oldtip {
				ch <- newtip
				oldtip = newtip
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func main() {

	var block Block
	var prevBlock Block

	ch := make(chan Block)
	go PingServer(ch)
	prevBlock = <-ch

	block.PrevHash = prevBlock.Hash()
	block.Name = "Shawki"

	go block.Mine(1)

	for {
		prevBlock = <-ch
		block.PrevHash = prevBlock.Hash()
		go block.Mine(1)
	}
}
