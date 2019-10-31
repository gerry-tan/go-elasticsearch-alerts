package wechat

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWechat(t *testing.T) {
	var m AlertMethodConfig = AlertMethodConfig{
		Corpid:  "wwf66d9b7ea0f9bdc6",
		Secret:  "1STRPTXAp5VMkKETTECokf_4E1ZqbiFhWotoIJfjZUU",
		ToUser:  "@all",
		ToParty: "",
		MsgType: "text",
		Agentid: 1000002,
		Text:    map[string]string{"content": "测试"},
	}

	token, err := Get_token("wwf66d9b7ea0f9bdc6", "1STRPTXAp5VMkKETTECokf_4E1ZqbiFhWotoIJfjZUU")
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println("获取token成功：", token)
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	fmt.Println(string(buf))
	err = Send_msg(token.Access_token, buf)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Println("发送消息成功", string(buf))
	}
}
