package config

import "time"

type Config struct {
	Wechat         WechatConfig
	Service        ServiceConfig
	Database       DatabaseConfig
	NotionOauth    OAuthConfig
	R2             R2Config
	Log            LogConfig
	DevToolsURL    string
	ChatGPT        ChatGPTConfig
	Telegram       TelegramConfig
	Storage        S3Config
	Zlib           ZlibConfig
	Readability    ReadabilityConfig
	NotionTestPage NotionTestPage
	JWT            JWTConfig
	OAuth          OAuthConfigMap
}

type LogConfig struct {
	Level string
}

type ServiceConfig struct {
	Name string
	Host string
	Port string
	URL  string
}

type Notion struct {
	BearerToken string
	DatabaseID  string
}

type DatabaseConfig struct {
	Driver string
	Sqlite SqliteConfig
	MySQL  MySQLConfig
}

type MySQLConfig struct {
	Host string
	Port int
	User string
	Pass string
	DB   string
}

type SqliteConfig struct {
	File string
}

type WechatConfig struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	AuthorID       string
	AdminUserID    string
}

type R2Config struct {
	Token string
	Url   string
}

type ChatGPTConfig struct {
	ApiKey        string
	Model         string
	DefaultPromot string
}

type TelegramConfig struct {
	Token string
}

type S3Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
	Domain     string // domain is cloudflare proxy domain
}

type ZlibConfig struct {
	Host         string
	IpfsGateways []string
}

type ReadabilityConfig struct {
	Host string
}

type NotionTestPage struct {
	DatabaseID string
	Token      string
}

type JWTConfig struct {
	SigningKey string
	Expiration time.Duration
}

type OAuthConfigMap struct {
	Notion OAuthConfig
	Github OAuthConfig
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	URLRedirect  string
	AuthURL      string
	AuthToken    string
	State        string
}
