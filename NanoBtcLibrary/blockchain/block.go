package blockchain

import (
	"bytes"
	"encoding/binary"
	"time"

	transactions "github.com/Btcercises/NanoBtcLibrary/blockchain/transactions"
	utils "github.com/Btcercises/NanoBtcLibrary/blockchain/utils"
)

type Hash256 []byte
type MagicId uint32

type BlockHeader struct {
	Hash             []byte
	Version          int
	HashPrev         Hash256
	HashMerkle       Hash256
	Timestamp        time.Time
	TargetDifficulty uint32
	Nonce            []byte
	MerkleRoot       [32]byte
}

type Block struct {
	BlockHeader
	MagicId          MagicId
	Length           uint32
	TransactionCount uint64
	Transactions     []transactions.Transaction
	StartPos         uint64
}

func NewBlock(version int,
	prevBlock []byte,
	merkleRoot []byte,
	timestamp time.Time,
	bits []byte,
	nonce []byte,
	total uint32,
	hash Hash256,
	flags []byte) *BlockHeader {
	result := &BlockHeader{
		Version:   version,
		Timestamp: timestamp,
		Hash:      hash,
	}
	copy(result.HashPrev[:32], prevBlock)
	copy(result.MerkleRoot[:32], merkleRoot)
	copy(result.Nonce[:4], nonce)
	return result
}

func (blockHeader *BlockHeader) HashBlock() Hash256 {
	if blockHeader.Hash != nil {
		return blockHeader.Hash
	}

	bin := make([]byte, 0)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(blockHeader.Version))
	bin = append(bin, version...)

	bin = append(bin, blockHeader.HashPrev...)
	bin = append(bin, blockHeader.HashMerkle...)

	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, uint32(blockHeader.Timestamp.Unix()))
	bin = append(bin, timestamp...)

	targetDifficulty := make([]byte, 4)
	binary.LittleEndian.PutUint32(targetDifficulty, blockHeader.TargetDifficulty)
	bin = append(bin, targetDifficulty...)

	nonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonce, uint32(blockHeader.Nonce))
	bin = append(bin, nonce...)

	blockHeader.Hash = utils.DoubleSha256(bin)
	return blockHeader.Hash
}

func (block *Block) IsValid() bool {
	flagBits := util.BytesToBitField(block.Flags)
	hashes := make([][]byte, len(block.Hashes))
	for i, hash := range block.Hashes {
		hashes[i] = make([]byte, len(hash))
		copy(hashes[i], hash)
		util.ReverseByteArray(hashes[i])
	}
	tree := NewTree(int(block.Total))
	tree.PopulateTree(flagBits, hashes)
	return bytes.Equal(util.ReverseByteArray(tree.Root()), block.MerkleRoot[:])
}
