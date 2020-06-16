package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

func Hash256(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), sha256.New())
}

func LittleEndianToInt32(b []byte) uint32 {
	if len(b) > 4 {
		panic("Value is too large!")
	}
	if len(b) < 4 {
		b = append(b, make([]byte, 4-len(b))...)
	}
	var result uint32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}
func Int32ToLittleEndian(num uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int64ToLittleEndian(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Int16ToLittleEndian(num uint16) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func EncodeVarInt(i int) []byte {
	if i < 0xfd {
		return []byte{byte(i)}
	}
	if i < 0x10000 {
		result := make([]byte, 3)
		copy(result[1:], Int16ToLittleEndian(uint16(i)))
		result[0] = 0xfd
		return result
	}
	if i < 0x100000000 {
		result := make([]byte, 5)
		copy(result[1:], Int32ToLittleEndian(uint32(i)))
		result[0] = 0xfe
		return result
	}

	result := make([]byte, 9)
	copy(result[1:], Int64ToLittleEndian(uint64(i)))
	result[0] = 0xff
	return result
}

func LittleEndianToInt64(b []byte) uint64 {
	if len(b) > 8 {
		panic("Value is too large!")
	}
	if len(b) < 8 {
		b = append(b, make([]byte, 8-len(b))...)
	}
	var result uint64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func LittleEndianToInt16(b []byte) uint16 {
	if len(b) > 2 {
		panic("Value is too large!")
	}
	if len(b) < 2 {
		b = append(b, make([]byte, 2-len(b))...)
	}
	var result uint16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func ReverseByteArray(arr []byte) []byte {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		arr[i], arr[length-i-1] = arr[length-i-1], arr[i]
	}
	return arr
}

func ReadVarInt(r *bytes.Reader) int {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	var bufsize int
	switch b {
	case 0xfd:
		bufsize = 2
	case 0xfe:
		bufsize = 4
	case 0xff:
		bufsize = 8
	default:
		return int(b)
	}
	buffer := make([]byte, bufsize)
	r.Read(buffer)
	return int(LittleEndianToInt64(buffer))
}
