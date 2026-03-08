> 一个使用 Go 语言编写的通过企业微信推送消息的工具，可作为 Web 接口接收请求并推送通知。



## 🚀 快速开始



### 1. 安装依赖



确保你已安装 Go 并配置好环境。然后运行以下命令安装依赖库：



```bash

go get github.com/gin-gonic/gin

go get github.com/go-playground/validator/v10

go get github.com/spf13/viper

```



### 2. 配置环境变量



你需要设置以下环境变量来配置企业微信相关参数：



```bash

export WX_CORPID="你的企业ID"

export WX_CORPSECRET="你的应用凭证密钥"

export WX_TOUSER="接收消息的成员ID"

export WX_AGENTID="应用ID"

export WX_SECRETKEY="用于验证请求的密钥（可自定义）"

export WX_ADDR="监听地址（默认：127.0.0.1:8234）"

```



### 3. 启动服务



```bash

go run main.go

```



服务将会在 `http://127.0.0.1:8234/2wx` 上监听 POST 请求。



### 4. 发送消息示例



你可以使用 curl 或 Postman 向接口发送 POST 请求：

```bash
curl -X POST "http://127.0.0.1:8234/2wx" -H "Content-Type: application/json" -d '{
    "title": "测试标题",
    "msg": "这是一条测试消息",
    "key": "你的SECRETKEY"
}'
```



## 📦 功能特点



- ✅ 支持通过 HTTP 接口接收消息并自动推送至企业微信

- ✅ 验证请求来源的合法性（使用 secret key）

- ✅ 自动刷新企业微信 access_token

- ✅ 使用 Gin 框架构建高性能的 Web 接口



## 📌 请求参数说明



| 参数名 | 类型 | 是否必填 | 说明 |

|--------|------|----------|------|

| title | string | 否 | 消息标题（可选） |

| msg | string | 是 | 消息正文（必须） |

| key | string | 是 | 请求验证密钥（必须匹配配置的 WX_SECRETKEY） |



## 📝 常见问题



### 1. 企业微信相关参数怎么获取？



请参考 [企业微信官方文档](https://work.weixin.qq.com/api/doc) 获取 `corpid`、`corpsecret`、`agentid` 等信息。



### 2. 接收消息的用户如何指定？



用户 ID 为企业微信用户在通讯录中 `UserID`，可在管理后台中获取。



### 3. 如何使用作为 CLI 工具发送消息？



你可以直接运行以下命令（前提是你配置了环境变量）：

```bash
go run main.go "这是一条测试消息"
```



## 📁 目录结构


```bash

├── main.go

└── README.md

```


