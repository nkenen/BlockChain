// Block.go
package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//区块结构体
type Block struct {
	Timestamp     int64  //时间戳
	Data          []byte //交易信息
	PrevBlockHash []byte //前一个区块哈希值
	Hash          []byte //自身哈希值
	Nonce         int    //工作量
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0} //初始化区块
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() //工作量计算哈希值

	block.Hash = hash[:] //哈希值
	block.Nonce = nonce  //工作量证明

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{}) //创建创世区块
}

//将区块序列化，为放入数据库做准备
func (b *Block) Serialize() []byte {
	var result bytes.Buffer           //byte类型缓冲区
	encode := gob.NewEncoder(&result) //新建序列byte类型缓冲区

	err := encode.Encode(b) //序列化区块
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//将传入的序列化数据返回序列化
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block) //反序列化
	if err != nil {
		log.Panic(err)
	}

	return &block
}
