package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"sync"
	"time"
	"vbc/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
}

func NewWebsocketUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *WebsocketUsecase {
	uc := &WebsocketUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket 连接映射
var clients = make(map[string]*websocket.Conn) // userID -> WebSocket
var mu sync.Mutex

// Redis 客户端
//var redisClient = redis.NewClient(&redis.Options{
//	Addr: "localhost:6379",
//})

// 消息结构
type WebsocketMessage struct {
	MessageID string `json:"messageID"`
	UserID    string `json:"userID"`
	Content   string `json:"content"`
}

var messagePool = sync.Pool{
	New: func() interface{} {
		return new(WebsocketMessage)
	},
}

// SendMessage 发送消息到客户端 currentTimes: 当前次数
func (c *WebsocketUsecase) SendMessage(userID, content string, currentTimes int) {
	if currentTimes >= 3 { // 消息总共会尝试3次， 3次后没有确认就直接丢了
		return
	}
	messageID := fmt.Sprintf("%d", time.Now().UnixNano()) // 生成唯一消息 ID
	msg := WebsocketMessage{MessageID: messageID, UserID: userID, Content: content}
	msgBytes, _ := json.Marshal(msg)

	c.log.Info("aaa:", userID, " : ", InterfaceToString(clients), len(clients))

	mu.Lock()
	conn, exists := clients[userID]
	mu.Unlock()

	c.log.Info("SendMessage: ", InterfaceToString(exists), " ", InterfaceToString(conn), " userId: ", InterfaceToString(userID))

	if exists {
		conn.WriteMessage(websocket.TextMessage, msgBytes)
		// 存入 Redis，等待 ACK（超时 15s）
		c.CommonUsecase.RedisClient().Set(context.Background(), "msg:"+messageID, msgBytes, 120*time.Second)
		// 10s 后检查 ACK
		go func() {
			time.Sleep(10 * time.Second)
			ack, _ := c.CommonUsecase.RedisClient().Get(context.Background(), "ack:"+messageID).Result()
			if ack == "" {
				c.log.Info("消息丢失，重新发送:", messageID, currentTimes)
				currentTimes += 1
				c.SendMessage(userID, content, currentTimes) // 重新发送
			}
		}()
	}
}

// HttpHandleWS 处理
func (c *WebsocketUsecase) HttpHandleWS(ctx *gin.Context) {
	userID := ctx.Query("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		c.log.Info("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[userID] = conn
	mu.Unlock()

	c.log.Info("用户连接成功:", userID, InterfaceToString(clients))

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			//c.log.Error(err)
			break
		}

		// 处理心跳 ping-pong
		if string(msgBytes) == "ping" {
			c.log.Info("💓 收到 ping，返回 pong")
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}

		// 解析 JSON
		//var msg WebsocketMessage
		msg := messagePool.Get().(*WebsocketMessage)
		json.Unmarshal(msgBytes, &msg)
		c.log.Info("ReadMessage:", InterfaceToString(msg))

		// 处理 ACK
		if msg.Content == "ACK" {
			c.log.Info("收到 ACK:", msg.MessageID)
			c.CommonUsecase.RedisClient().Set(context.Background(), "ack:"+msg.MessageID, "ok", 30*time.Second) // 标记消息已确认
			c.CommonUsecase.RedisClient().Del(context.Background(), "msg:"+msg.MessageID)                       // 删除缓存
		}
		messagePool.Put(msg) // 释放对象到池中
	}
	//c.log.Info("deleted:", userID)
	mu.Lock()
	delete(clients, userID)
	mu.Unlock()
}

func (c *WebsocketUsecase) HttpSendMessage(ctx *gin.Context) {
	userID := ctx.Query("userID")
	content := ctx.Query("message")
	c.log.Info("userID: ", userID, " content: ", content)
	c.SendMessage(userID, content, 0)
	ctx.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

//// Gin 服务器
//func main() {
//	r := gin.Default()
//
//	// WebSocket 连接
//	r.GET("/ws", handleWS)
//
//	// 发送消息
//	r.GET("/send", func(c *gin.Context) {
//		userID := c.Query("userID")
//		content := c.Query("message")
//		sendMessage(userID, content)
//		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
//	})
//
//	//log.Println("WebSocket Server running on :8080")
//	//r.Run(":8080")
//}
