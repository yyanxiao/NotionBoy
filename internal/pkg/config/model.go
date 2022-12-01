package config

type Config struct {
	Wechat      WechatConfig
	Service     ServiceConfig
	Database    DatabaseConfig
	NotionOauth NotionOauthConfig
	R2          R2Config
	Log         LogConfig
	DevToolsURL string
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
}

type R2Config struct {
	Token string
	Url   string
}
