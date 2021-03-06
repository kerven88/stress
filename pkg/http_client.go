package pkg

import (
	"crypto/tls"
	"fmt"
	"github.com/oldthreefeng/stress/utils"
	"io"
	"net/http"
	"time"
)

// HttpRequest is HTTP 请求
// method 方法 GET POST
// url 请求的url
// body 请求的body
// headers 请求头信息
// timeout 请求超时时间
func HttpRequest(method, url string, body io.Reader, headers utils.ConcurrentMap, timeout time.Duration) (resp *http.Response, requestTime uint64, err error) {

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	req, err := http.NewRequest(method, url, body)

	// 使用 --compressed。 则使用gzip压缩算法去请求。
	if Compressed {
		req.Header.Add("Accept-Encoding", "gzip")
	}
	if err != nil {

		return
	}

	// 设置默认为utf-8编码

	if _, ok := headers.Get("Content-Type"); !ok {
		if headers == nil {
			headers = utils.New(Concurrency)
		}
		headers.Set("Content-Type","application/x-www-form-urlencoded; charset=utf-8")
	}

	for key, value := range headers.Items() {
		req.Header.Set(key, value.(string))
	}

	startTime := time.Now()
	resp, err = client.Do(req)
	requestTime = uint64(utils.DiffNano(startTime))
	if err != nil {
		fmt.Println("请求失败:", err)

		return
	}

	// bytes, err := json.Marshal(req)
	// fmt.Printf("%#v \n", req)

	return
}
