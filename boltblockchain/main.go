// boltblockchain project main.go
package main

func main() {
	bc := NewBlockChain() //创建区块链

	defer bc.db.Close() //延迟关闭数据库

	cli := CLI{bc} //cli结构体创建
	cli.Run()      //使用命令进行区块链交易
}
