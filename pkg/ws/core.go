package ws

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	zLog "zgo/pkg/logs"
)

var (
	WebsocketHub *Hub
	upGrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
)

type Client struct {
	Uuid               string //ws连接ID
	Key                string //浏览器随机key
	Hub                *Hub
	Socket             *websocket.Conn // ws连接
	Message            chan []byte     //ws连接的消息管道
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailTimes int
	sync.RWMutex
	State uint8 //状态

}

func (c *Client) OpenWs(writer http.ResponseWriter, request *http.Request) (*Client, bool) {
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				zLog.Log.Err(val).Msg("websocket open err")
			}

		}
	}()
	//连接升级
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(writer, request, nil)
	if err != nil {
		zLog.Log.Err(err).Msg("webSocket协议升级失败")
		return nil, false
	}
	nClient := newClient(ws)
	nClient.Key = request.Header.Get("Sec-WebSocket-Key")
	return nClient, true
}

func (c *Client) ReadPump(callbackOnMessage func(messageType int, receiveData []byte), callbackOnError func(err error), callbackOnclose func()) {
	defer func() {
		err := recover()
		if err != nil {
			callbackOnclose()
		}
	}()

	for {
		if 1 == c.State {
			mt, ReceivedData, err := c.Socket.ReadMessage()
			if err != nil {
				callbackOnError(err)
				break
			}
			callbackOnMessage(mt, ReceivedData)
		} else {
			callbackOnError(errors.New("client state err"))
			break
		}
	}
}

func (c *Client) SendMessage(messageType int, message string) error {
	c.Lock()
	defer c.Unlock()

	//消息单次写入超时时间，单位：秒
	if err := c.Socket.SetWriteDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
		return err
	}

	if err := c.Socket.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	}

	return nil

}

func newClient(ws *websocket.Conn) *Client {
	//todo message长度可配置
	returnC := &Client{
		Uuid:          uuid.New().String(),
		Socket:        ws,
		Message:       make(chan []byte, 1000),
		PingPeriod:    time.Second * 20,
		ReadDeadline:  time.Second * 100,
		WriteDeadline: time.Second * 35,
		State:         1,
	}
	returnC.SendMessage(websocket.TextMessage, `{"code":2001,"msg":"ws连接成功","data":""}`)
	//设置最长读取长度
	returnC.Socket.SetReadLimit(16777217)

	returnC.Hub = WebsocketHub
	returnC.Hub.Register <- returnC

	return returnC
}

// Heartbeat 按照websocket标准协议实现隐式心跳,Server端向Client远端发送ping格式数据包,浏览器收到ping标准格式，自动将消息原路返回给服务器
func (c *Client) Heartbeat() {
	//  1. 设置一个时钟，周期性的向client远端发送心跳数据包
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				zLog.Log.Err(val).Msg("ErrorsWebsocketBeatHeartFail")
			}
		}
		ticker.Stop() // 停止该client的心跳检测
	}()
	//2.浏览器收到服务器的ping格式消息，会自动响应pong消息，将服务器消息原路返回过来
	if c.ReadDeadline == 0 {
		_ = c.Socket.SetReadDeadline(time.Time{})
	} else {
		_ = c.Socket.SetReadDeadline(time.Now().Add(c.ReadDeadline))
	}
	c.Socket.SetPongHandler(func(receivedPong string) error {
		if c.ReadDeadline > time.Nanosecond {
			_ = c.Socket.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		} else {
			_ = c.Socket.SetReadDeadline(time.Time{})
		}
		//fmt.Println("浏览器收到ping标准格式，自动将消息原路返回给服务器：", received_pong)  // 接受到的消息叫做pong，实际上就是服务器发送出去的ping数据包
		return nil
	})
	//3.自动心跳数据
	for {
		select {
		case <-ticker.C:
			if c.State == 1 {
				if err := c.SendMessage(websocket.PingMessage, "Server->Ping->Client"); err != nil {
					c.HeartbeatFailTimes++
					if c.HeartbeatFailTimes > 4 {
						c.State = 0
						zLog.Log.Err(err).Msg("ErrorsWebsocketBeatHeartsMoreThanMaxTimes")

						return
					}
				} else {
					if c.HeartbeatFailTimes > 0 {
						c.HeartbeatFailTimes--
					}
				}
			} else {
				return
			}

		}
	}
}
