package wechat

import (
	"bytes"
	"context"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
)

var wechatSession *bytes.Buffer

func (b *WechatBot) login(uuid string) error {
	if err := b.loginWithSession(); err == nil {
		return nil
	}

	if uuid == "" {
		var err error
		uuid, err = b.Caller.GetLoginUUID()
		if err != nil {
			logger.SugaredLogger.Errorw("Login wechat generate uuid error", "err", err)
			return err
		}
	}
	return b.loginWithID(uuid)
}

func (b *WechatBot) loginWithSession() error {
	readLoginSession(b.ctx)
	if err := b.HotLogin(wechatSession, false); err != nil {
		logger.SugaredLogger.Errorw("Login wehcat with session failed", "err", err)
		return err
	}
	logger.SugaredLogger.Info("Login wechat with session success")
	return nil
}

func (b *WechatBot) loginWithID(uuid string) error {
	if err := b.LoginWithUUID(uuid); err != nil {
		logger.SugaredLogger.Errorw("Login wechat with uuid failed", "err", err)
		return err
	}
	return nil
}

func readLoginSession(ctx context.Context) {
	session, err := dao.GetSession(ctx)
	if err != nil {
		logger.SugaredLogger.Debug("Get wechat session failed, init without session")
		wechatSession = &bytes.Buffer{}
	} else {
		logger.SugaredLogger.Debug("Get wechat session success, init with session")
		wechatSession = bytes.NewBuffer(session.Session)
	}
}

func saveLoginSession(ctx context.Context, body []byte) {
	if err := dao.SaveSession(context.TODO(), wechatSession.Bytes()); err != nil {
		logger.SugaredLogger.Errorw("Save wechat session failed", "err", err)
	}
	logger.SugaredLogger.Info("Save wechat session success")
}

func loginCallBack(body []byte) {
	saveLoginSession(context.TODO(), body)
	logger.SugaredLogger.Info("Login wechat successful finished")
}
