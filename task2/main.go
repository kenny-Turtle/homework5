package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// 获取输入的合约地址
	contractAddr := flag.String("contract", "", "合约地址")
	flag.Parse()

	// 获取环境变量
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("请设置环境变量 ETH_RPC_URL")
	}

	// 连接以太坊客户端
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// DialContext需要用https://开头的URL，否则会报错
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("连接以太坊客户端失败: %v", err)
	}

	// 连接合约
	addr := common.HexToAddress(*contractAddr)
	counter, err := NewCounter(addr, client)
	if err != nil {
		log.Fatalf("连接合约失败: %v", err)
	}

	// ===== 合约操作 =======
	// 读取当前计数
	log.Println("...读取当前计数...")
	current, err := counter.X(nil)
	if err != nil {
		log.Fatalf("读取当前计数失败: %v", err)
	}
	log.Printf("当前计数: %d", current)

	// 创建签名交易
	privKeyHex := os.Getenv("PRIVATE_ACCOUNT1_KEY")
	if privKeyHex == "" {
		log.Fatal("请设置环境变量 PRIVATE_ACCOUNT1_KEY")
	}
	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		log.Fatalf("私钥错误: %v", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("获取链ID失败: %v", err)
	}

	// 创建签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("创建签名器失败: %v", err)
	}

	// 调用incBy(2)增加计数
	log.Println("...调用incBy(2)增加计数...")
	tx, err := counter.IncBy(auth, big.NewInt(2))
	if err != nil {
		log.Fatalf("调用incBy失败: %v", err)
	}
	log.Printf("交易已发送: %s", tx.Hash().Hex())
	log.Println("...等待交易确认...")

	// 等待上链
	log.Println("等待交易打包上链...")
	waitTxConfirm(client, tx.Hash())

	// 再次读取当前计数
	log.Println("...再次读取当前计数...")
	newCount, err := counter.X(nil)
	if err != nil {
		log.Fatalf("读取当前计数失败: %v", err)
	}
	log.Printf("新的计数: %d", newCount)
}

func waitTxConfirm(client *ethclient.Client, txHash common.Hash) {
	for {
		// 查询交易是否被打包
		_, isPending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			log.Println("查询中...", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 如果已经打包（不是pending）
		if !isPending {
			log.Println("✅ 交易已上链！")
			break
		}

		// 还在等待
		log.Println("⌛ 交易打包中...请稍等")
		time.Sleep(2 * time.Second)
	}

	// 额外等 2s 确保状态完全同步
	time.Sleep(2 * time.Second)
}
