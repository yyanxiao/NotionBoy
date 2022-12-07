package server

import (
	"net/http"

	notion "notionboy/internal/pkg/notion"
)

var oauthMgr notion.OauthInterface

func initNotion() {
	oauthMgr = notion.GetOauthManager()
	http.HandleFunc("/notion/oauth", notionOauth)
	http.HandleFunc("/notion/oauth/callback", notionOauthCallback)
}

func notionOauth(w http.ResponseWriter, r *http.Request) {
	state := getQueryParam(r, "state")
	url := oauthMgr.OAuthProcess(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func notionOauthCallback(w http.ResponseWriter, r *http.Request) {
	state := getQueryParam(r, "state")
	code := getQueryParam(r, "code")
	if code == "" || state == "" {
		renderError(w, http.StatusBadRequest, "code or state is empty", nil)
		return
	}
	ctx := r.Context()
	msg, err := oauthMgr.OAuthCallback(ctx, code, state)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderHtml(w, msg, http.StatusOK)
}
