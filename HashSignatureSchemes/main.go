package HashSignatureSchemes

import (
	"crypto/rand"
	"crypto/sha256"

	utils "github.com/Btcercises/HashSignatureSchemes/utils"
)

func GenerateKey() (utils.SecretKey, utils.PublicKey, error) {
	var sec utils.SecretKey
	var pub utils.PublicKey

	for i := 0; i < 256; i++ {
		var _, e = rand.Read(sec.ZeroPre[i][:])
		if e != nil {
			return sec, pub, e
		}
		_, e = rand.Read(sec.OnePre[i][:])
		if e != nil {
			return sec, pub, e
		}
	}
	for i := 0; i < 256; i++ {
		pub.ZeroHash[i] = sha256.Sum256(sec.ZeroPre[i][:])
		pub.OneHash[i] = sha256.Sum256(sec.OnePre[i][:])
	}
	return sec, pub, nil
}

func Sign(msg utils.Message, sec utils.SecretKey) utils.Signature {
	var sig utils.Signature

	var i uint32
	var slicedMsg byte
	var shift uint32
	var isEqZero bool

	for i = 0; i < 256; i++ {
		slicedMsg = msg[:][i/32]
		shift = i % 32
		isEqZero = (1<<i)&(slicedMsg&(1<<shift)) == 0

		if !isEqZero {
			sig.Preimage[i] = sec.OnePre[i]

		} else {
			sig.Preimage[i] = sec.ZeroPre[i]
		}
	}
	return sig
}

func Verify(msg utils.Message, pub utils.PublicKey, sig utils.Signature) bool {

	var i uint32
	var slicedMsg byte
	var shift uint32
	var isEqZero bool
	var SigSlicedPreImg []byte
	var hashesSigSlicedPreImg [32]byte

	for i = 0; i < 256; i++ {
		slicedMsg = msg[:][i/32]
		shift = i % 32
		isEqZero = (1<<i)&(slicedMsg&(1<<shift)) == 0

		SigSlicedPreImg = sig.Preimage[i][:]
		hashesSigSlicedPreImg = sha256.Sum256(SigSlicedPreImg)

		if !isEqZero {
			if pub.OneHash[i] != hashesSigSlicedPreImg {
				return false
			}

		} else {
			if pub.ZeroHash[i] != hashesSigSlicedPreImg {
				return false
			}
		}
	}

	return true
}
