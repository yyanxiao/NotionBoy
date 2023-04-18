package server

import (
	"net/http"

	"notionboy/internal/pkg/logger"
	"notionboy/internal/server/handler"
)

type wechatPayCallbackErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func wechatPayCallback(w http.ResponseWriter, r *http.Request) {
	err := handler.WechatPayCallback(r.Context(), r)
	if err != nil {
		logger.SugaredLogger.Errorw("WechatPayCallback failed", "err", err)
		renderJson(w, wechatPayCallbackErrorResp{
			Code:    "FAIL",
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
