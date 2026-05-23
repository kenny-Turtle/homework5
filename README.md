

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




## 任务 2：合约代码生成 任务目标

描述
```
使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。 具体任务

编写智能合约
使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
编译智能合约，生成 ABI 和字节码文件。
使用 abigen 生成 Go 绑定代码
安装 abigen 工具。
使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
使用生成的 Go 绑定代码与合约交互
编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
调用合约的方法，例如增加计数器的值。
输出调用结果。
```

### 开发步骤
#### 1.创建hardhat3项目
```
kenny@anonymous task2 % mkdir hardhat3
kenny@anonymous task2 % cd hardhat3 
kenny@anonymous hardhat3 % npm init -y
Wrote to /Users/kenny/ideaProjects/web3/MetaNode/golang-learning-kenny/homework5/task2/hardhat3/package.json:

{
  "name": "hardhat3",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC"
}



kenny@anonymous hardhat3 % npm install --save-dev hardhat@3.1.0

added 59 packages in 5s

14 packages are looking for funding
  run `npm fund` for details
kenny@anonymous hardhat3 % npx hardhat --init
```
#### 2. 部署合约
hh3框架自带一个Counter合约，
新建.env 配置下url和私钥
然后进行合约部署
```
kenny@anonymous hardhat3 % npx hardhat ignition deploy ./ignition/modules/Counter.ts --network sepolia                                
✔ Confirm deploy to network sepolia (11155111)? … yes
Hardhat Ignition 🚀

Resuming existing deployment from ./ignition/deployments/chain-11155111

Deploying [ CounterModule ]

Batch #1
  Executed CounterModule#Counter

Batch #2
  Executed CounterModule#Counter.incBy

[ CounterModule ] successfully deployed 🚀

Deployed Addresses

CounterModule#Counter - 0x2e843Ae3a37B00498E8f58b5E1ecE3822954Da02
```

#### 3. abigen操作
安装abigen
```
kenny@anonymous task2 % go install github.com/ethereum/go-ethereum/cmd/abigen@latest
kenny@anonymous task2 % abigen --version
abigen version 1.17.3-stable
```
生成go文件
```
kenny@anonymous task2 % abigen \
--abi <(jq -r '.abi' ./hardhat3/artifacts/contracts/Counter.sol/Counter.json) \
--pkg=main \
--type=Counter \
--out=counter.go
```

#### 4. 执行
```
kenny@anonymous task2 % export ETH_RPC_URL=https://sepolia.drpc.org
go run . -contract 0x2e843Ae3a37B00498E8f58b5E1ecE3822954Da02
2026/05/23 17:03:07 ...读取当前计数...
2026/05/23 17:03:08 当前计数: 7
2026/05/23 17:03:08 ...调用incBy(2)增加计数...
2026/05/23 17:03:11 交易已发送: 0x1c5f4eb0557d214d1f62d37ff0ce4f2005f0e18ca3556abd8f59e83729cb30e6
2026/05/23 17:03:11 ...等待交易确认...
2026/05/23 17:03:11 等待交易打包上链...
2026/05/23 17:03:11 ⌛ 交易打包中...请稍等
2026/05/23 17:03:14 ⌛ 交易打包中...请稍等
2026/05/23 17:03:16 ⌛ 交易打包中...请稍等
2026/05/23 17:03:18 ⌛ 交易打包中...请稍等
2026/05/23 17:03:21 ⌛ 交易打包中...请稍等
2026/05/23 17:03:23 ⌛ 交易打包中...请稍等
2026/05/23 17:03:26 ⌛ 交易打包中...请稍等
2026/05/23 17:03:28 ⌛ 交易打包中...请稍等
2026/05/23 17:03:30 ⌛ 交易打包中...请稍等
2026/05/23 17:03:33 ⌛ 交易打包中...请稍等
2026/05/23 17:03:35 ⌛ 交易打包中...请稍等
2026/05/23 17:03:38 ✅ 交易已上链！
2026/05/23 17:03:40 ...再次读取当前计数...
2026/05/23 17:03:40 新的计数: 9
```