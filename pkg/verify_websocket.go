package pkg

import (
	"encoding/json"
	"fmt"
)

/***************************  返回值为json  ********************************/

// WebSocketResponseJson is 返回数据结构体
type WebSocketResponseJson struct {
	Seq      string `json:"seq"`
	Cmd      string `json:"cmd"`
	Response struct {
		Code    int         `json:"code"`
		CodeMsg string      `json:"codeMsg"`
		Data    interface{} `json:"data"`
	} `json:"response"`
}

// WebSocketJson is 通过返回的Body 判断
// 返回示例: {"seq":"1566276523281-585638","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}
// code 取body中的返回code
func WebSocketJson(request *Request, seq string, msg []byte) (code int, isSucceed bool) {

	responseJson := &WebSocketResponseJson{}
	err := json.Unmarshal(msg, responseJson)
	if err != nil {
		code = ParseError
		fmt.Printf("请求结果 json.Unmarshal msg:%s err:%v", string(msg), err)
	} else {

		if seq != responseJson.Seq {
			code = ParseError
			fmt.Println("请求和返回seq不一致 ~请求:", seq, responseJson.Seq, string(msg))
		} else {
			code = responseJson.Response.Code
			// body 中code返回200为返回数据成功
			if code == 200 {
				isSucceed = true
			}
		}
	}

	// 开启调试模式
	if request.GetDebug() {
		fmt.Printf("请求结果 seq:%s body:%s \n", seq, string(msg))
	}

	return
}
