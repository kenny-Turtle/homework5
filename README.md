

## 任务 1：区块链读写 任务目标
使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。 具体任务

要求
```
环境搭建
安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
查询区块
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
输出查询结果到控制台。
发送交易
准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
对交易进行签名，并将签名后的交易发送到网络。
输出交易的哈希值。
```


设置环境变量
export ETH_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/xxxxxxx
export PRIVATE_ACCOUNT1_KEY=xxxxxxxx


```
kenny@anonymous task1 % go run main.go -mode block -number 123456
======================================
区块 #123456
======================================
Block: &{header:0x2627301bd188 uncles:[] transactions:[] withdrawals:[] accessList:<nil> hash:{_:[] _:{} v:<nil>} size:{_:{} _:{} v:0} ReceivedAt:0001-01-01 00:00:00 +0000 UTC ReceivedFrom:<nil>}
Number       : 123456
Hash         : 0x2056507046b07a5d7ed4f124a7febce2aec7295b464746523787b8c2acc627dc
Parent Hash  : 0x93bff867b68a2822ee7b6e0a4166cfdf5fc4782d60458fae1185de9b2ecdba16
Time         : 2021-11-12T19:16:28+08:00
Time (Local) : 2021-11-12 19:16:28 CST
Gas Used     : 0 (0.00%)
Gas Limit    : 8000000
Tx Count     : 0
State Root   : 0x8126ca1ac23042b8c719fe8c708be4f3a106972cb3500da8c7632c9036e8eb35
Tx Root      : 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
Receipt Root : 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
Difficulty   : 658615317
Coinbase     : 0x2f14582947E292a2eCd20C430B46f2d27CFE213c
======================================





kenny@anonymous task1 % go run main.go -mode transfer -to 0xF2659305D8573B8B2cB0B3bCbd2b11fA66edCCbf -amount 0.1  
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
交易已发送，Tx Hash: 0x3c3893b97289c265721aa7479aa3139e8f956cd110a7bb48243db7fb8dbaeb38
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
From          : 0xc7F340b38178cfF13a32d33570Da4dA771F4ED13
To            : 0xF2659305D8573B8B2cB0B3bCbd2b11fA66edCCbf
Amount        : 100000000000000000 wei
Gas Limit     : 21000
Gas Tip Cap   : 1440000 wei
Gas Fee Cap   : 2121139250 wei
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```