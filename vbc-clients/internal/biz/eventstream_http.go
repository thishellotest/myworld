package biz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

type EventstreamHttpUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	JWTUsecase    *JWTUsecase
	CommonUsecase *CommonUsecase
}

func NewEventstreamHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	CommonUsecase *CommonUsecase) *EventstreamHttpUsecase {
	return &EventstreamHttpUsecase{
		log:           log.NewHelper(logger),
		conf:          conf,
		JWTUsecase:    JWTUsecase,
		CommonUsecase: CommonUsecase,
	}
}

func GenSubscribeKey(gid string) string {
	return fmt.Sprintf("eventstream:user:%s", gid)
}

func (c *EventstreamHttpUsecase) WriteMessage(w gin.ResponseWriter, message string) {
	m := "data: " + message
	w.Write([]byte(m))
	w.Write([]byte("\n\n"))
	w.Flush()
}

func GenAuthorization(jwt string) string {
	return "Bearer " + jwt
}

// It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

func (c *EventstreamHttpUsecase) Handle(ctx *gin.Context) {

	jwt := ctx.Param("jwt")
	c.JWTUsecase.HandleJWTAuth(GenAuthorization(jwt), ctx)
	if ctx.IsAborted() {
		return
	}
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	lib.DPrintln(userFacade)

	w := ctx.Writer

	//ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()
	//// 订阅用户专属频道
	//pubsub := c.CommonUsecase.RedisClient().Subscribe(ctxWithTimeout, GenSubscribeKey(userFacade.Gid()))
	//defer pubsub.Close()
	//
	//ch := pubsub.Channel()

	// 定期发送心跳消息，避免连接超时
	heartbeatTicker := time.NewTicker(5 * time.Second)
	defer heartbeatTicker.Stop()
	notify := ctx.Writer.CloseNotify()
	for {
		select {
		//case msg := <-ch:
		//	c.WriteMessage(w, msg.Payload)
		case <-heartbeatTicker.C:
			c.WriteMessage(w, "heartbeat")
			c.log.Debug("heartbeat:", userFacade.ToFabUser())
			// 发送心跳信号
		case <-ctx.Done(): // 连接断开
			c.log.Debug("用户断开:", userFacade.ToFabUser(), ctx.Err())
			return
		//case <-ctx.Request.Context().Done():
		//	c.log.Debug("Request context done. Reason: ", ctx.Err())
		//	c.log.Debug("用户断开22:", userFacade.ToFabUser())
		//	return
		case <-notify:
			c.log.Debug("用户断开1:", userFacade.ToFabUser())
			return
		}
	}
}

func (c *EventstreamHttpUsecase) BizHandle(str string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)
	data.Set("data.val", "aaa")
	return data, nil
}
