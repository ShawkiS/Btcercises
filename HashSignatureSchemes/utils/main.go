package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block [32]byte

type SecretKey struct {
	ZeroPre [256]Block
	OnePre  [256]Block
}

type PublicKey struct {
	ZeroHash [256]Block
	OneHash  [256]Block
}

func (self PublicKey) ToHex() string {
	var s string
	for _, zero := range self.ZeroHash {
		s += zero.ToHex()
	}
	for _, one := range self.OneHash {
		s += one.ToHex()
	}
	return s
}

func HexToPubkey(s string) (PublicKey, error) {
	var p PublicKey

	expectedLength := 256 * 2 * 64

	if len(s) != expectedLength {
		return p, fmt.Errorf(
			"Pubkey string %d characters, expect %d", expectedLength)
	}

	bts, err := hex.DecodeString(s)
	if err != nil {
		return p, err
	}
	buf := bytes.NewBuffer(bts)

	for i, _ := range p.ZeroHash {
		p.ZeroHash[i] = BlockFromByteSlice(buf.Next(32))
	}
	for i, _ := range p.OneHash {
		p.OneHash[i] = BlockFromByteSlice(buf.Next(32))
	}

	return p, nil
}

type Message Block

func (self Block) ToHex() string {
	return fmt.Sprintf("%064x", self[:])
}

func (self Block) Hash() Block {
	return sha256.Sum256(self[:])
}

func (self Block) IsPreimage(arg Block) bool {
	return self.Hash() == arg
}

func BlockFromByteSlice(by []byte) Block {
	var bl Block
	copy(bl[:], by)
	return bl
}

type Signature struct {
	Preimage [256]Block
}

func (self Signature) ToHex() string {
	var s string
	for _, b := range self.Preimage {
		s += b.ToHex()
	}

	return s
}

func HexToSignature(s string) (Signature, error) {
	var sig Signature

	expectedLength := 256 * 64

	if len(s) != expectedLength {
		return sig, fmt.Errorf(
			"Pubkey string %d characters, expect %d", expectedLength)
	}

	bts, err := hex.DecodeString(s)
	if err != nil {
		return sig, err
	}
	buf := bytes.NewBuffer(bts)

	for i, _ := range sig.Preimage {
		sig.Preimage[i] = BlockFromByteSlice(buf.Next(32))
	}
	return sig, nil
}

func GetMessageFromString(s string) Message {
	return sha256.Sum256([]byte(s))
}
