package zlib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/pkg/utils/cache"
	"sync"
	"time"
)

var (
	once       sync.Once
	client     *ZlibClient
	httpClient = &http.Client{}
)

const LIMIT = "50" // default limit

type Book struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Author        string   `json:"author,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	ISBN          string   `json:"isbn,omitempty"`
	FileSize      int64    `json:"filesize"`
	FileSizeHuman string   `json:"filesize_human"` // human readable file size
	Extension     string   `json:"extension"`
	Year          int      `json:"year,omitempty"`
	Pages         int      `json:"pages,omitempty"`
	IpfsCid       string   `json:"ipfs_cid"`
	IpfsLinks     []string `json:"ipfs_links"`
}

type ZlibClient struct {
	Host        string
	IfsGateways []string
}

type Searcher interface {
	Search(ctx context.Context, name string) ([]*Book, error)
}

func NewZlibClient(host string, ipfsGateways []string) Searcher {
	return &ZlibClient{Host: host, IfsGateways: ipfsGateways}
}

func DefaultZlibClient() Searcher {
	if client == nil {
		once.Do(func() {
			cfg := config.GetConfig().Zlib
			client = &ZlibClient{Host: cfg.Host, IfsGateways: cfg.IpfsGateways}
		})
	}
	return client
}

func (z *ZlibClient) Search(ctx context.Context, name string) ([]*Book, error) {
	var books []*Book
	var err error
	// check cache, if not exist, search from zlib
	cacheBooks, ok := cache.DefaultClient().Get(name)
	if !ok {
		books, err = z.getJSON(ctx, name)
		if err != nil {
			return nil, err
		}
		cache.DefaultClient().Set(name, books, 24*time.Hour)
	} else {
		books = cacheBooks.([]*Book)
	}
	for _, book := range books {
		book.IpfsLinks = buildIpfsLinks(book, z.IfsGateways)
		book.FileSizeHuman = convertHumanReadableFileSize(book.FileSize)
	}
	return books, nil
}

func buildIpfsLinks(book *Book, gateways []string) []string {
	var links []string
	for _, gateway := range gateways {
		u, _ := url.Parse(gateway)
		u.Path = "/ipfs/" + book.IpfsCid
		// use ipfs cid can find the file, no need for filename
		// q := u.Query()
		// q.Set("filename", fmt.Sprintf("%s_%s.%s", book.Title, book.Author, book.Extension))
		// u.RawQuery = q.Encode()
		links = append(links, u.String())
	}
	return links
}

func convertHumanReadableFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

type Result struct {
	Books []*Book `json:"books"`
}

func (z *ZlibClient) getJSON(ctx context.Context, name string) ([]*Book, error) {
	u, err := url.Parse(z.Host)
	if err != nil {
		logger.SugaredLogger.Errorw("parse url error", "url", z.Host, "err", err)
		return nil, err
	}
	u.Path = "/search"
	q := u.Query()
	q.Set("query", name)
	q.Set("limit", LIMIT)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		logger.SugaredLogger.Errorw("new http request error", "url", u.String(), "err", err)
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.SugaredLogger.Errorw("http get error", "url", u.String(), "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var books []*Book

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.SugaredLogger.Errorw("read body error", "url", u.String(), "err", err)
		return books, err
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		logger.SugaredLogger.Errorw("decode json error", "url", u.String(), "err", err)
		return books, err
	}

	return result.Books, nil
}
