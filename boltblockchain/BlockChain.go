// BlockChain.go
package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const blocksBucket = "blocks" //区块桶
const dbFile = "blockchain"   //数据库名称

type BlockChain struct {
	tip []byte   //序列化的区块
	db  *bolt.DB //数据库
}

//在数据库增加新的区块
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	//读取数据库
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l")) //标志"l"的值为最后一个哈希值，以此放入下一个区块，形成区块链

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	//向数据库上传新的区块数据
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize()) //将哈希值和区块序列化数据放入区块桶中
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash) //把新区块的哈希值放入“l”标志地址
		bc.tip = newBlock.Hash

		return nil
	})
}

//区块链的迭代器
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db} //区块链迭代器

	return bci
}

//创建新的区块链，最开始进行
func NewBlockChain() *BlockChain {
	var tip []byte                          //序列化的缓冲区
	db, err := bolt.Open(dbFile, 0600, nil) //数据库的新建、打开

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket)) //读取区块桶数据

		//如若为空，则创建创世区块
		if b == nil {
			genesis := NewGenesisBlock()                    //创建创世区块
			b, err := tx.CreateBucket([]byte(blocksBucket)) //创建区块桶
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize()) //放入创世区块
			err = b.Put([]byte("l"), genesis.Hash)         //为下一个区块准备前一块链
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := BlockChain{tip, db}

	return &bc
}
