package config

type Config struct {
	Wechat      WechatConfig
	Service     ServiceConfig
	Database    DatabaseConfig
	NotionOauth NotionOauthConfig
	R2          R2Config
	Log         LogConfig
	DevToolsURL string
	ChatGPT     ChatGPTConfig
	Telegram    TelegramConfig
	Storage     S3Config
	Zlib        ZlibConfig
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

type NotionOauthConfig struct {
	ClientID     string
	ClientSecret string
	URLRedirect  string
	AuthURL      string
	AuthToken    string
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
	AuthorImageID  string
	AdminUserID    string
}

type R2Config struct {
	Token string
	Url   string
}

type ChatGPTConfig struct {
	Authorization string
	SessionToken  string
	ApiKey        string
	User          string
	Pass          string
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
