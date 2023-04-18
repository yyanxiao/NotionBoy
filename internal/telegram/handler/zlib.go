package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"
	"notionboy/internal/zlib"

	tele "gopkg.in/telebot.v3"
)

const (
	LIMIT                             = 5 // limit for every message show how many books
	SIZE                              = 3 // size of pagination
	INLINE_UNIQUE_ZLIB_SEARCHER       = "zlib_searcher"
	INLINE_UNIQUE_ZLIB_SAVE_TO_NOTION = "zlib_save_to_notion"
)

func OnZlibSaveToNotion(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*60*time.Second)
	defer cancel()
	// ctx := context.Background()
	name := c.Callback().Data
	books, err := zlib.DefaultZlibClient().Search(ctx, name)
	if err != nil {
		return c.Reply(fmt.Sprintf("Search from zlib error: %s", err))
	}

	acc, err := queryUserAccount(ctx, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return c.Reply(config.MSG_ERROR_ACCOUNT_NOT_FOUND)
	}

	nContent := &notion.Content{
		Zlib:    &notion.ZlibContent{Books: books},
		Text:    name,
		Tags:    []string{"zlib", "telegram", name},
		Account: acc,
	}

	nt := &notion.Notion{BearerToken: acc.AccessToken, DatabaseID: acc.DatabaseID}
	res, pageID, err := nt.CreateRecord(ctx, &notion.Content{
		Text:    "内容正在更新，请稍等",
		Account: acc,
	})
	nContent.Process(ctx)
	if err == nil {
		nt.PageID = pageID
		nt.UpdateRecord(ctx, nContent)
	}

	return c.Reply(res)
}

func OnZlib(c tele.Context) error {
	isCallback := c.Callback() != nil
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var msg string

	// check if is callback from inline button
	// query from inline button will have pagination info
	if isCallback {
		msg = c.Callback().Data
	} else {
		msg = c.Message().Payload
	}
	msgs := strings.Split(msg, " ")
	logger.SugaredLogger.Debugw("zlib", "msgs", msgs)

	// build query with ext
	query := msgs[0]
	for i := 1; i < len(msgs); i++ {
		if msgs[i][0] == '#' {
			query += " " + msgs[i]
		}
	}

	books, err := zlib.DefaultZlibClient().Search(ctx, query)
	if err != nil {
		return c.Reply(fmt.Sprintf("Search from zlib error: %s", err))
	}

	if len(books) == 0 {
		return c.Reply("No book found")
	}

	page := uint(1)
	// if is callback, get page from callback data
	if c.Callback() != nil {
		u, err := strconv.ParseInt(msgs[len(msgs)-1], 10, 0)
		if err != nil {
			logger.SugaredLogger.Errorw("parse page error for callback", "err", err)
			page = uint(1)
		} else {
			page = uint(u)
		}
	}

	reply := buildReplyBooks(books[LIMIT*(page-1) : LIMIT*page])
	logger.SugaredLogger.Infof("reply: %s", reply)

	markUp := buildPaginationInlineMarkup(books, msgs[0], page)

	err = c.EditOrReply(reply, tele.ModeHTML, tele.NoPreview, markUp)
	if err != nil {
		logger.SugaredLogger.Errorw("reply error", "err", err)
	}

	return err
}

func buildReplyBooks(books []*zlib.Book) string {
	sb := strings.Builder{}
	for _, book := range books {
		downloadLink := strings.Builder{}
		for idx, link := range book.IpfsLinks {
			downloadLink.WriteString(fmt.Sprintf(`<a href="%s">%d </a>`, link, idx+1))
		}

		strings.Join(book.IpfsLinks, " ")
		sb.WriteString(fmt.Sprintf(`
📚 <b>%s</b> (📅 %d)
✍ %s
⬇️ %s (%s, %s)
		`,
			book.Title,
			book.Year,
			book.Author,
			downloadLink.String(),
			book.Extension,
			book.FileSizeHuman,
		))
	}
	sb.WriteString(config.MSG_ZLIB_TIPS)
	return sb.String()
}

func buildPaginationInlineMarkup(books []*zlib.Book, name string, page uint) *tele.ReplyMarkup {
	length := uint(len(books))

	maxPage := length / 5

	hasPre := func(n uint) bool {
		return n > 1
	}
	hasNext := func(n uint) bool {
		return n < maxPage
	}

	buildInlineButton := func(n uint, text string) tele.InlineButton {
		return tele.InlineButton{
			Unique: INLINE_UNIQUE_ZLIB_SEARCHER,
			Text:   text,
			Data:   fmt.Sprintf("%s %d", name, n),
		}
	}

	var buttons []tele.InlineButton
	i := page
	if hasPre(page) {
		buttons = append(buttons, buildInlineButton(page-1, fmt.Sprintf("<< (%d)", page-1)))
	}
	for ; i < page+SIZE; i++ {
		if i > maxPage {
			break
		}
		buttons = append(buttons, buildInlineButton(i, strconv.Itoa(int(i))))
	}

	if hasNext(i) {
		buttons = append(buttons, buildInlineButton(i, fmt.Sprintf("(%d) >> ", i)))
	}

	return &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			buttons,
			{
				{
					Unique: INLINE_UNIQUE_ZLIB_SAVE_TO_NOTION,
					Text:   "Save to Notion",
					Data:   name,
				},
			},
		},
	}
}
