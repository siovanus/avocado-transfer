## 批量转账工具

主要实现一个地址给多个地址转账的功能，目前支持ont，ong和oep4.

### 配置文件

```
{
  "JsonRpcAddress":"http://polaris1.ont.io:20336",
  "Asset": "oep4",
  "ContractAddress": "be7ee30f1dfa27bfd85eb1319eeb038bc18bc4e7",
  "DataFile": "./input.txt",
  "WalletFile":"./wallet.dat",
  "Bonus": 100,
  "GasPrice":500,
  "GasLimit":20000000
}
```

JsonRpcAddress：配置ontology网络的rpc地址，主网为dappnode1.ont.io，测试网为polaris1.ont.io

Asset：资产类型，支持ont，ong，oep4

ContractAddress：当资产类型为oep4时有效，指具体的oep4资产合约地址

DataFile： 输入文件，内容格式后面详解

WalletFile：钱包文件路径

Bonus：转账总额，注意此处是总额，所有地址按照权重分配这些总额

### 输入文件

```
{"address":"AMX5e18a2PYJnrc9ytDiu8jzRhtJV9zfHT","value":300}
{"address":"AUvucjvShztgdRqbDuN1QVXxn9ywq1xgvJ","value":300}
{"address":"AWXYJwYsTQ29n7YjXF8X5LYu4YXNwLkg9H","value":100}
{"address":"ARgDnVhtxwTsQHqzh6vFrfCHXHiRnQXdWa","value":100}
{"address":"AbKv5mugr46upB6duPQLvh85gjuw1wB4bG","value":100}
{"address":"AVGLir2MLeWt31uV7rBKF21GYTUFHo2wEK","value":100}
{"address":"AbJ3V4gRaGgspfPEBNom2kBwyHtQ8QoW6W","value":100}
{"address":"AGR8VURsRHaZ4EWmGKEqbp3qpSqgAixVuk","value":100}
{"address":"AcvDjEcGNSD6Jyvmyzw7UJYbpin6mAQzCq","value":100}
{"address":"AdHVxwV9mYjPjA4zCvw8Vi1TWC2BF6iKKC","value":100}
{"address":"ASjfbM8EENX46RjwEf2Ku2fNGvmmz1DUdc","value":100}
{"address":"AZA9cVnvXpG7gZrWMQCCNATARvuworhhyq","value":100}
{"address":"ANiKkw49MENtJPiMXwC6t1UwiZCMv8GG3s","value":100}
```

格式如上，分别为地址和权重，按照权重值分配Bonus转账总额。分批转账，500个地址一批。

### 运行

编译 go build main.go

运行./main，按照提示输入密码