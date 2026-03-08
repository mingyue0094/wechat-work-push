package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// 使用 viper 从环境变量或者配置文件读取参数
var (
	_corpid      string
	_corpsecret  string
	_touser      string
	_agentid     int
	_secretkey   string
	_token       string
	_json_slice  struct {
		Access_token string
		Errcode      int
		Errmsg       string
		Expires_in   int
	}
	_token_expires_in int
	_Addr string
)

// 定义结构体
type RequestData struct {
	Title string `json:"title" validate:"omitempty"`
	Msg   string `json:"msg" validate:"required"`
	Key   string `json:"key" validate:"required"`
}

// text_msg 结构体
type Text_Msg struct {
	Content string `json:"content"`
}

// post body 结构体
type Post_Body struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    Text_Msg `json:"text"`
	Safe    int `json:"safe"`
}

// 定义函数
func _get_token() (string, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + _corpid + "&corpsecret=" + _corpsecret
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("get token failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &_json_slice); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %w", err)
	}

	if _json_slice.Errcode != 0 {
		return "", fmt.Errorf("get token API failed with ErrCode: %d, ErrMsg: %s", _json_slice.Errcode, _json_slice.Errmsg)
	}

	_token = _json_slice.Access_token
	_token_expires_in = _json_slice.Expires_in

	return _token, nil
}



func _send_msg(msg string, _token string) (string, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + _token

	fmt.Printf("Sending message to WeChat Work: %s\n", msg)

	post_body := Post_Body{
		Touser:  _touser,
		Msgtype: "text",
		Agentid: _agentid,
		Text: Text_Msg{
			Content: msg,
		},
		Safe: 0,
	}

	json_data, err := json.Marshal(post_body)
	if err != nil {
		return "marshal body failed:", fmt.Errorf(" %w", err)
	}

	post_body_str := strings.NewReader(string(json_data))

	req, _ := http.NewRequest("POST", url, post_body_str)

	req.Header.Add("cache-control", "no-cache")

	resp, _ := http.DefaultClient.Do(req)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp)
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return "unexpected status",fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	return string(body), nil
}

// 企业微信推送
func Wx_push_qy(msg string) (string, error) {
	// 获取当前时间戳 整数型
	current_timestamp := time.Now().Unix()
	
	// 如果 token 已过期，重新获取 token
	if _token == "" || _token_expires_in <= 0 || current_timestamp >= int64(_token_expires_in) {
		var err error
		_token, err = _get_token()
		if err != nil {
			return "failed to get token", fmt.Errorf("%w", err)
		}
	}

	return _send_msg(msg, _token)
}

// 通过环境变量/配置加载参数
func initViper() {

	// 设置 viper 可从环境变量中读取
	viper.SetEnvPrefix("WX")
	viper.AutomaticEnv()

	// 从环境变量加载
	_corpid = viper.GetString("corpid")
	_corpsecret = viper.GetString("corpsecret")
	_touser = viper.GetString("touser")
	_agentid = viper.GetInt("agentid")
	_secretkey = viper.GetString("secretkey")
	_Addr := viper.GetString("addr")

}



func main() {
	// 初始化 viper 并加载配置
	initViper()
	if _corpid == "" || _corpsecret == "" || _touser == "" || _agentid == 0 || _secretkey == "" {
		fmt.Println("请配置环境变量或者配置文件，示例: WX_CORPID=xxx WX_CORPSECRET=xxx WX_TOUSER=xxx WX_AGENTID=1234 WX_SECRETKEY=123456")
		os.Exit(1)
	}

	// 启动参数解析
	addr := flag.String("addr", "127.0.0.1:8234", "Address and port to listen on")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()


	if _Addr != "" {
		addr = &_Addr
	}

	if *help {
		fmt.Fprintf(os.Stderr, "用法: %s [选项] [消息]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		fmt.Fprintf(os.Stderr, "  -help      显示本帮助信息\n")
		fmt.Fprintf(os.Stderr, "  -addr ip:port\n")
		fmt.Fprintf(os.Stderr, "        设置监听的地址和端口 (默认: 127.0.0.1:8234)\n")
		fmt.Fprintf(os.Stderr, "说明:\n")
		fmt.Fprintf(os.Stderr, "  如果未提供消息，则启动一个HTTP服务，监听指定地址。\n")
		fmt.Fprintf(os.Stderr, "  如果提供了消息，则直接发送消息，不启动服务。\n")
		fmt.Fprintf(os.Stderr, "  使用环境变量配置：\n")
		fmt.Fprintf(os.Stderr, "    WX_CORPID=...\n")
		fmt.Fprintf(os.Stderr, "    WX_CORPSECRET=...\n")
		fmt.Fprintf(os.Stderr, "    WX_TOUSER=...\n")
		fmt.Fprintf(os.Stderr, "    WX_AGENTID=...\n")
		fmt.Fprintf(os.Stderr, "    WX_SECRETKEY=...\n")
		os.Exit(1)
	}

	if len(flag.Args()) > 0 {
		msg := flag.Args()[0]
		_,err := Wx_push_qy(msg)
		if err != nil {
			fmt.Printf("发送消息失败: %v\n", err)
		}
		return
	}

	// 启动 Gin Web Server
	r := gin.Default()


	// 设置路由
	r.POST("/2wx", func(c *gin.Context) {
		var req RequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// 使用 validator 验证参数是否必填
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// 验证 secret key 是否正确
		if req.Key != _secretkey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}


		fmt.Println("%s\n%s", req.Title, req.Msg)

		response,err := Wx_push_qy(fmt.Sprintf("%s\n%s", req.Title, req.Msg))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message", "details": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "200", "msg": response })
		}
	})




	// 启动服务
	fmt.Printf("Server is running at http://%s/2wx\n\n\n", *addr)
	fmt.Println("Examples:")
	fmt.Println("POST: http://"+*addr+"/2wx")
	fmt.Println("with JSON: {\"title\": \"Hi\", \"msg\": \"Hello\", \"key\": \""+_secretkey+"\"}")


	err := r.Run(*addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
