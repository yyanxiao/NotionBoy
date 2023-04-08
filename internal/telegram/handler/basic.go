package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"notionboy/db/ent/account"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"
	"notionboy/internal/pkg/utils/cache"
	"notionboy/internal/service/auth"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

var cacheClient = cache.DefaultClient()

func OnStart(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}
	if err := dao.SaveBasicAccount(context.Background(), account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10)); err != nil {
		logger.SugaredLogger.Errorw("SaveBasicAccount failed", "err", err)
	}
	return c.Send(config.MSG_START)
}

func OnBind(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}
	stage := fmt.Sprintf("%s:%d", account.UserTypeTelegram, sender.ID)
	oauthMgr := notion.GetOauthManager()
	url := oauthMgr.OAuthURL(stage)
	logger.SugaredLogger.Info("OAuthURL: ", url)
	text := config.MSG_BINDING + url
	return c.Reply(text)
}

func OnUnbind(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}

	if err := dao.ClearNotionAuthInfo(context.Background(), account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10)); err != nil {
		return c.Reply(config.MSG_UNBIND_FAILED + err.Error())
	}
	return c.Reply(config.MSG_UNBIND_SUCCESS)
}

func OnWebUI(c tele.Context) error {
	return c.Reply(config.GetConfig().Service.URL)
}

func OnMagicCode(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return c.Reply("User do not exist")
	}
	ctx := context.Background()
	acc, err := queryUserAccount(ctx, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return c.Reply("Query Account Error: " + err.Error())
	}

	code := uuid.New().String()
	cacheClient.Set(fmt.Sprintf("%s:%s", config.MAGIC_CODE_CACHE_KEY, code), acc, time.Duration(5)*time.Minute)

	return c.Reply(code)
}

func OnWhoAmI(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return fmt.Errorf("User do not exist")
	}
	ctx := context.Background()

	myInfo, err := auth.WhoAmI(ctx, account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10))
	if err != nil {
		return c.Reply("查询用户信息失败: " + err.Error())
	}
	return c.Reply(myInfo.String())
}

func OnSOS(c tele.Context) error {
	return c.Reply(fmt.Sprintf("欢迎添加作者微信，请搜索🔍:  %s", config.GetConfig().Wechat.AuthorID))
}

func OnApiKey(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return fmt.Errorf("User do not exist")
	}
	ctx := context.Background()
	key, err := auth.GenerateApiKey(ctx, account.UserTypeTelegram, strconv.FormatInt(sender.ID, 10))
	if err != nil {
		return c.Reply("生成 API Key 失败: " + err.Error())
	}

	return c.Reply(key)
}
