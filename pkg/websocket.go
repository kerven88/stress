package pkg

import (
	"fmt"
	"github.com/oldthreefeng/stress/utils"
	"sync"
	"time"
)

const (
	firstTime    = 1 * time.Second // 连接以后首次请求数据的时间
	intervalTime = 1 * time.Second // 发送数据的时间间隔
)

var (
	// 请求完成以后是否保持连接
	keepAlive bool
)

func init() {
	keepAlive = true
}

// WebSocket is  web socket go link
func WebSocket(chanId int, ch chan<- *RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *Request, ws *WebSocketC) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)

	defer func() {
		ws.Close()
	}()

	var (
		i uint64
	)

	// 暂停60秒
	t := time.NewTimer(firstTime)
	for {
		select {
		case <-t.C:
			t.Reset(intervalTime)

			// 请求
			webSocketRequest(chanId, ch, i, request, ws)

			// 结束条件
			i = i + 1
			if i >= totalNumber {
				goto end
			}
		}
	}

end:
	t.Stop()

	if keepAlive {
		// 保持连接
		chWaitFor := make(chan int, 0)
		<-chWaitFor
	}

	return
}

// webSocketRequest is 请求
func webSocketRequest(chanId int, ch chan<- *RequestResults, i uint64, request *Request, ws *WebSocketC) {

	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = HttpOk
	)

	// 需要发送的数据
	seq := fmt.Sprintf("%d_%d", chanId, i)
	err := ws.Write([]byte(`{"seq":"` + seq + `","cmd":"ping","data":{}}`))
	if err != nil {
		errCode = RequestErr // 请求错误
	} else {

		// time.Sleep(1 * time.Second)
		msg, err := ws.Read()
		if err != nil {
			errCode = ParseError
			fmt.Println("读取数据 失败~")
		} else {
			// fmt.Println(msg)
			errCode, isSucceed = request.GetVerifyWebSocket()(request, seq, msg)
		}
	}

	requestTime := uint64(utils.DiffNano(startTime))

	requestResults := &RequestResults{
		Time:      requestTime,
		IsSucceed: isSucceed,
		ErrCode:   errCode,
	}

	requestResults.SetId(chanId, i)

	ch <- requestResults

}
