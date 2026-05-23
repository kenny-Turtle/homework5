package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 2.查询区块
// 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
// 输出查询结果到控制台。
// go run main.go -mode block -number 123456

// 3.发送交易
// 准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
// 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
// 构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
// 对交易进行签名，并将签名后的交易发送到网络。
// 输出交易的哈希值。
// go run main.go -mode transfer -to 0xRecipientAddress -amount 0.1

func main() {
	// 获取输入信息
	mode := flag.String("mode", "block", "模式: block 或 transfer")
	blockNumber := flag.Uint64("number", 0, "区块号")
	toHex := flag.String("to", "", "接收方地址 (仅 transfer 模式)")
	amount := flag.String("amount", "", "转账金额 (仅 transfer 模式)")
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

	// ===============================================================
	// 查询指定区块信息
	if *mode == "block" && *blockNumber > 0 {
		num := big.NewInt(int64(*blockNumber))
		block, err := fetchBlockByNumber(ctx, client, num)
		if err != nil {
			log.Fatalf("查询区块失败: %v", err)
		}
		printBlockInfo(fmt.Sprintf("区块 #%d", *blockNumber), block)
	}
	// ===============================================================
	// 发送交易
	if *mode == "transfer" {
		// 实现发送交易的逻辑
		handleTransfer(ctx, client, *toHex, *amount)
	}

}

// 查询指定区块信息,重试3次
func fetchBlockByNumber(ctx context.Context, client *ethclient.Client, number *big.Int) (*types.Block, error) {
	var lastErr error
	for i := 0; i < 3; i++ {
		// 每次重试使用新的上下文，避免上下文被取消
		reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		block, err := client.BlockByNumber(reqCtx, number)
		cancel()

		if err == nil {
			return block, nil
		}
		lastErr = err
		log.Printf("查询区块失败 (尝试 %d/3): %v", i+1, err)
		time.Sleep(2 * time.Second) // 等待一段时间后重试
	}
	return nil, fmt.Errorf("查询区块失败: %v", lastErr)
}

// 输出区块信息
// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
func printBlockInfo(title string, block *types.Block) {
	fmt.Println("======================================")
	fmt.Println(title)
	fmt.Println("======================================")
	fmt.Printf("Block: %+v\n", block)

	// 基本信息
	fmt.Printf("Number       : %d\n", block.Number().Uint64())
	fmt.Printf("Hash         : %s\n", block.Hash().Hex())
	fmt.Printf("Parent Hash  : %s\n", block.ParentHash().Hex())

	// 时间信息
	blockTime := time.Unix(int64(block.Time()), 0)
	fmt.Printf("Time         : %s\n", blockTime.Format(time.RFC3339))
	fmt.Printf("Time (Local) : %s\n", blockTime.Local().Format("2006-01-02 15:04:05 MST"))

	// Gas 信息
	gasUsed := block.GasUsed()
	gasLimit := block.GasLimit()
	gasUsagePercent := float64(gasUsed) / float64(gasLimit) * 100
	fmt.Printf("Gas Used     : %d (%.2f%%)\n", gasUsed, gasUsagePercent)
	fmt.Printf("Gas Limit    : %d\n", gasLimit)

	// 交易信息
	txCount := len(block.Transactions())
	fmt.Printf("Tx Count     : %d\n", txCount)

	// 区块根信息（Merkle 树根）
	fmt.Printf("State Root   : %s\n", block.Root().Hex())
	fmt.Printf("Tx Root      : %s\n", block.TxHash().Hex())
	fmt.Printf("Receipt Root : %s\n", block.ReceiptHash().Hex())

	// 区块大小估算（简化版，实际大小还包括其他字段）
	if txCount > 0 {
		fmt.Printf("\nFirst Tx Hash: %s\n", block.Transactions()[0].Hash().Hex())
		if txCount > 1 {
			fmt.Printf("Last Tx Hash : %s\n", block.Transactions()[txCount-1].Hash().Hex())
		}
	}

	// 难度信息（PoW 相关，PoS 后基本固定）
	fmt.Printf("Difficulty   : %s\n", block.Difficulty().String())

	// 区块奖励相关信息
	coinbase := block.Coinbase()
	if coinbase != (common.Address{}) {
		fmt.Printf("Coinbase     : %s\n", coinbase.Hex())
	}

	fmt.Println("======================================")
	fmt.Println()
}

// 处理转账逻辑
func handleTransfer(ctx context.Context, client *ethclient.Client, toHex, amountStr string) {
	// 校验
	if toHex == "" || amountStr == "" {
		log.Fatal("请提供接收方地址和转账金额")
	}

	// 检查环境变量
	privKeyHex := os.Getenv("PRIVATE_ACCOUNT1_KEY")
	if privKeyHex == "" {
		log.Fatal("请设置环境变量 PRIVATE_ACCOUNT1_KEY")
	}
	// 解析私钥
	privKey, err := crypto.HexToECDSA(trim0x(privKeyHex))
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	// 获取发送方地址
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法断言公钥类型")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddr := common.HexToAddress(toHex)

	// 获取链ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("获取链ID失败: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}

	// 估算Gas Limit（转账ETH固定21000）
	gasLimit := uint64(21000)
	// 获取建议的Gas价格（EIP-1559）
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatalf("获取建议的 Gas Tip Cap 失败: %v", err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("获取最新区块头失败: %v", err)
	}
	baseFee := header.BaseFee
	if baseFee == nil {
		log.Println("当前网络不支持 EIP-1559，使用传统的 Gas Price")
		baseFee, err = client.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("获取建议的 Gas Price 失败: %v", err)
		}
	}

	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	// 检查余额是否足够：ETH + gas
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatalf("获取余额失败: %v", err)
	}
	// 将amountStr转成float64
	ethValue, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Fatalf("解析转账金额失败: %s", amountStr)
	}
	// 直接乘以1e18转换成wei
	weiValue := new(big.Float).Mul(big.NewFloat(ethValue), big.NewFloat(1e18))
	// 转成big.Int
	amount, _ := weiValue.Int(nil)
	totalCost := new(big.Int).Add(amount, new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit))))

	if balance.Cmp(totalCost) < 0 {
		log.Fatalf("余额不足: 需要至少 %s wei (转账金额 + 预估的 Gas 费用)，当前余额: %s wei", totalCost.String(), balance.String())
	}

	// 构造交易
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr, // 接收人地址
		Value:     amount,  // 要转的 ETH
		Data:      nil,     // 空
	}
	tx := types.NewTx(txData)
	// 签名交易
	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		log.Fatalf("签名交易失败: %v", err)
	}

	// 发送交易
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("发送交易失败: %v", err)
	}

	// 输出交易的哈希值
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("交易已发送，Tx Hash: %s\n", signedTx.Hash().Hex())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("From          : %s\n", fromAddress.Hex())
	fmt.Printf("To            : %s\n", toAddr.Hex())
	fmt.Printf("Amount        : %s wei\n", amount.String())
	fmt.Printf("Gas Limit     : %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap   : %s wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap   : %s wei\n", gasFeeCap.String())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// trim0x 移除十六进制字符串前缀 "0x"
func trim0x(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}
