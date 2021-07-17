package wxgzh

import (
	"fmt"

	"notionboy/config"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
)

func Run() {
	r := gin.Default()

	wc := wechat.NewWechat()
	account := NewOfficialAccount(wc)
	r.Any("/", account.Serve)
	svc := config.GetConfig().Service
	r.Run(fmt.Sprintf("%s:%s", svc.Host, svc.Port))
}
