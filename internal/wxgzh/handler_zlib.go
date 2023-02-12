/*
zlib handler for wechat

Three commands: /zlib, zlibm, zlibs
- /zlib search book from zlib
  - when seach done, will cache the result for 30 minutes
  - will return first 5 books

- zlibm read from cache and return next page
  - when user send zlibm, will return 5 books from next page

- zlibs save all result seached from zlib to notion
  - when user send zlibs, will save all result to notion
*/
package wxgzh

import (
	"context"
	"fmt"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/db/dao"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"notionboy/internal/zlib"
	"strings"
	"time"

	notion "notionboy/internal/pkg/notion"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

const (
	// PAGE_LIMIT is the limit of books per page
	PAGE_LIMIT = 5
	// CACHE_TTL is the ttl of cache for zlib search result
	CACHE_TTL = 30 * time.Minute
)

type ZlibResultCache struct {
	name      string
	page      int
	totalPage int
	books     []*zlib.Book
}

func searchZlib(ctx context.Context, msg *message.MixMessage, mr chan *message.Reply) {
	name := strings.TrimSpace(msg.Content[5:])
	logger.SugaredLogger.Debugw("zlib search name", "name", name)
	books, err := zlib.DefaultZlibClient().Search(ctx, name)
	if err != nil {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(fmt.Sprintf("Search from zlib error: %s", err))}
		return
	}
	if len(books) == 0 {
		mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ZLIB_NO_RESULT)}
		return
	}

	cacheRes := &ZlibResultCache{
		name:      name,
		page:      1,
		totalPage: len(books) / 5,
		books:     books,
	}
	cache.DefaultClient().Set(getCacheKey(msg), cacheRes, CACHE_TTL)
	res := buildZlibSearchResult(cacheRes)
	mr <- &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(res)}
}

func searchZlibNextPage(ctx context.Context, msg *message.MixMessage) *message.Reply {
	cacheRes, ok := cache.DefaultClient().Get(getCacheKey(msg))
	if !ok {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("没有找到相关记录，请重新搜索， 例如 「/zlib 鲁迅」")}
	}
	res := cacheRes.(*ZlibResultCache)

	if res.page >= res.totalPage {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("没有更多内容了！")}
	}
	res.page++
	cache.DefaultClient().Set(getCacheKey(msg), res, CACHE_TTL)
	resStr := buildZlibSearchResult(res)

	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(resStr)}
}

func getCacheKey(msg *message.MixMessage) string {
	return fmt.Sprintf("zlib_%s", msg.FromUserName)
}

func buildZlibSearchResult(rc *ZlibResultCache) string {
	sb := strings.Builder{}

	buildDownloadLink := func(cid string) string {
		return fmt.Sprintf("%s/files/ipfs/%s", config.GetConfig().Service.URL, cid)
	}

	for i := PAGE_LIMIT * (rc.page - 1); i < len(rc.books) && i < PAGE_LIMIT*rc.page; i++ {
		book := rc.books[i]
		sb.WriteString(fmt.Sprintf("📚 %s (📅  %d)\n", book.Title, book.Year))
		if book.Author != "" {
			sb.WriteString(fmt.Sprintf("✍ %s\n", book.Author))
		}
		if book.Publisher != "" {
			sb.WriteString(fmt.Sprintf("出版社: %s\n", book.Publisher))
		}
		sb.WriteString(fmt.Sprintf("⬇️  %s (%s, %s)\n", buildDownloadLink(book.IpfsCid), book.Extension, book.FileSizeHuman))
		sb.WriteString("--------\n\n")
	}
	if rc.page+1 <= rc.totalPage {
		sb.WriteString("回复 zlibm 查看更多, 回复 zlibs 保存所有的结果到 Notion")
	} else {
		sb.WriteString("到底了， 没有更多的内容！回复 zlibs 保存所有的结果到 Notion")
	}
	sb.WriteString(config.MSG_ZLIB_TIPS_CN)
	return sb.String()
}

func searchZlibSaveToNotion(ctx context.Context, msg *message.MixMessage) *message.Reply {
	cacheRes, ok := cache.DefaultClient().Get(getCacheKey(msg))
	if !ok {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("没有找到相关记录，请重新搜索， 例如 「/zlib 鲁迅」")}
	}
	res := cacheRes.(*ZlibResultCache)

	acc, err := dao.QueryAccountByWxUser(ctx, msg.GetOpenID())
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ERROR_ACCOUNT_NOT_FOUND)}
	}
	if acc.ID == 0 {
		return bindNotion(ctx, msg)
	}
	n := &notion.Notion{BearerToken: acc.AccessToken, DatabaseID: acc.DatabaseID}

	if msg.MsgType != message.MsgTypeText {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(config.MSG_ZLIB_UNSUPPOERT)}
	}
	nContent := &notion.Content{
		Zlib:    &notion.ZlibContent{Books: res.books},
		Text:    res.name,
		Tags:    []string{"zlib", "wechat", res.name},
		Account: acc,
	}
	// 创建初始 Record
	var zlibPageId string
	var resp string
	resp, zlibPageId, err = n.CreateRecord(ctx, &notion.Content{
		Text:    "Zlib 专属页面",
		Account: acc,
	})
	if err != nil {
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(fmt.Sprintf("Create zlib page error: %s", err))}
	}
	n.PageID = zlibPageId

	nContent.Process(ctx)
	go n.UpdateRecord(context.TODO(), nContent)

	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(resp)}
}
