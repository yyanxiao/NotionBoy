package server

import (
	"net/http"

	"notionboy/internal/wxgzh"

	wechat "github.com/silenceper/wechat/v2"
)

var (
	wc        *wechat.Wechat
	wxAccount *wxgzh.OfficialAccount
)

func initWx(mux *http.ServeMux) {
	wc = wechat.NewWechat()
	wxAccount = wxgzh.NewOfficialAccount(wc)
	// 获取ak
	mux.HandleFunc("/api/v1/oa/basic/get_access_token", wxGetAccessToken)
	// 获取微信callback IP
	mux.HandleFunc("/api/v1/oa/basic/get_callback_ip", wxGetCallbackIP)
	// 获取微信API接口 IP
	mux.HandleFunc("/api/v1/oa/basic/get_api_domain_ip", wxGetAPIDomainIP)
	// 清理接口调用次数
	mux.HandleFunc("/api/v1/oa/basic/clear_quota", wxClearQuota)
	// 处理消息

	// mux.HandleFunc("/", wxProcessMsg)
}

func wxGetAccessToken(w http.ResponseWriter, r *http.Request) {
	ak, err := wxAccount.GetAccessToken()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "", err)
		return
	}
	renderHtml(w, ak, http.StatusOK)
}

func wxGetCallbackIP(w http.ResponseWriter, r *http.Request) {
	ipList, err := wxAccount.GetCallbackIP()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "", err)
		return
	}
	renderHtml(w, ipList, http.StatusOK)
}

func wxGetAPIDomainIP(w http.ResponseWriter, r *http.Request) {
	ipList, err := wxAccount.GetAPIDomainIP()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "", err)
		return
	}
	renderHtml(w, ipList, http.StatusOK)
}

func wxClearQuota(w http.ResponseWriter, r *http.Request) {
	msg, err := wxAccount.GetAccessToken()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "", err)
		return
	}
	renderHtml(w, msg, http.StatusOK)
}

func wxProcessMsg(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	wxAccount.Serve(w, r)
}
