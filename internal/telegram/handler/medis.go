package handler

import (
	"context"
	"fmt"
	"time"

	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/notion"

	"github.com/jomei/notionapi"
	tele "gopkg.in/telebot.v3"
)

func OnMedia(c tele.Context) error {
	ctx := context.Background()
	acc, err := queryUserAccount(ctx, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Query Account Error: %v", err)
		return c.Reply(config.MSG_ERROR_ACCOUNT_NOT_FOUND)
	}

	m := c.Message().Media()
	if m.MediaFile().FileSize > 20*1024*1024 {
		return c.Reply("Sorry, Telegram does not support get files larger than 20 MB, https://core.telegram.org/bots/api#file")
	}

	text := ""
	var media *notion.MediaContent

	switch c.Message().Media().MediaType() {
	case "photo":
		text, media = processPhoto(c.Message().Photo)
	case "audio":
		text, media = processAudio(c.Message().Audio)
	case "document":
		text, media = processDocument(c.Message().Document)
	case "video":
		text, media = processVideo(c.Message().Video)
	case "animation":
		text, media = processAnimation(c.Message().Animation)
	case "voice":
		text, media = processVoice(c.Message().Voice)
	case "videoNote":
		text, media = processVideoNote(c.Message().VideoNote)
	case "sticker":
		text, media = processSticker(c.Message().Sticker)
	default:
		return c.Reply("Unsupported message type")
	}

	medias := make([]*notion.MediaContent, 0)
	medias = append(medias, media)
	nContent := &notion.Content{
		Text:    text,
		Medias:  medias,
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

func buildMeidaURL(fileID string) string {
	return fmt.Sprintf("%s/files/tg/%s", config.GetConfig().Service.URL, fileID)
}

func buildMediaText(mediaType string) string {
	return fmt.Sprintf("#%s %s", mediaType, time.Now().Format(time.RFC3339))
}

func processPhoto(m *tele.Photo) (string, *notion.MediaContent) {
	text := buildMediaText("Image")
	media := &notion.MediaContent{
		Type: notionapi.BlockTypeImage.String(),
		URL:  buildMeidaURL(m.FileID) + ".jpg",
	}
	return text, media
}

func processVoice(m *tele.Voice) (string, *notion.MediaContent) {
	return processFile(m, "Voice")
}

func processAudio(m *tele.Audio) (string, *notion.MediaContent) {
	return processFile(m, "Audio")
}

func processAnimation(m *tele.Animation) (string, *notion.MediaContent) {
	return processFileWithName(m, "animation", m.FileName)
}

func processDocument(m *tele.Document) (string, *notion.MediaContent) {
	if m.MIME == "application/pdf" {
		text := fmt.Sprintf("#%s %s %s", "pdf", time.Now().Format(time.RFC3339), m.FileName)
		media := &notion.MediaContent{
			Type: notionapi.BlockTypePdf.String(),
			URL:  buildMeidaURL(m.MediaFile().FileID) + ".pdf",
		}
		return text, media
	}

	return processFileWithName(m, "Document", m.FileName)
}

func processVideo(m *tele.Video) (string, *notion.MediaContent) {
	return processFileWithName(m, "Video", m.FileName)
}

func processVideoNote(m *tele.VideoNote) (string, *notion.MediaContent) {
	return processFile(m, "VideoNote")
}

func processSticker(m *tele.Sticker) (string, *notion.MediaContent) {
	return processFile(m, "Sticker")
}

func processFile(m tele.Media, t string) (string, *notion.MediaContent) {
	text := buildMediaText(t)
	media := &notion.MediaContent{
		Type: notionapi.BlockTypeFile.String(),
		URL:  buildMeidaURL(m.MediaFile().FileID),
	}
	return text, media
}

func processFileWithName(m tele.Media, t, name string) (string, *notion.MediaContent) {
	text := fmt.Sprintf("#%s %s %s", t, time.Now().Format(time.RFC3339), name)
	media := &notion.MediaContent{
		Type: notionapi.BlockTypeFile.String(),
		URL:  buildMeidaURL(m.MediaFile().FileID),
	}
	return text, media
}
