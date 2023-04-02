package config

// const for message
const (
	MSG_BIND_SUCCESS = `
<html>
	<head>
		<meta charset="utf-8" />
		<title>NotionBoy</title>
		<style>
			h1 {
				font-size: 48px;
				font-weight: 500;
				line-height: 1.2;
				margin: 0;
				padding: 0;
				text-align: center;
			}
			p {
				text-align: center;
				text-overflow: clip;
				color: black;
				font-size: 30px;
				margin: 0;
				padding: 0;
			}
		</style>
	</head>
	<body>
		<h1>恭喜 🎉 成功绑定 Notion</h1>
		<p>
			请关闭网页，返回微信即可开始使用 NotionBoy 回到 Notion
			可查看新建的欢迎信息
		</p>
	</body>
</html>
`

	MSG_BINDING = `欢迎使用 NotionBoy, 了解 NotionBoy 的使用指南，请参考: https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw


您还未绑定，请访问下面的链接，然后选择允许访问的页面，并点击 "允许"，即可绑定 Notion。
(升级后的 NotionBoy 支持下面的链接进行绑定，不需要再手动配置 Bot, 也不需要手动输入绑定信息)

注意，只选择一个页面，如果选择了多个页面，也只会使用选择的第一个页面。

如果微信打开有问题，请复制 URL 到浏览器打开。

`

	MSG_UNBIND_SUCCESS  = "成功解除 Notion 绑定！"
	MSG_UNBIND_FAILED   = "解除 Notion 绑定失败！失败原因: "
	MSG_UNSUPPOERT      = "不支持的消息类型!"
	MSG_CHAT_UNSUPPOERT = "Chat 只支持文本消息类型!"
	MSG_ZLIB_UNSUPPOERT = "Zlib 只支持文本消息类型!"
	MSG_PROCESSING      = "正在处理，请稍后去 Notion 查看"
	MSG_HELP            = `欢迎使用 NotionBoy, 了解 NotionBoy 的使用指南, 请参考: https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw

命令大全:
- /bind 命令可以用于绑定 Notion 账户, 使 NotionBoy 能够访问 Notion 中的内容。
- /unbind 命令可以用于解绑 Notion 账户, 使 NotionBoy 不再能够访问 Notion 中的内容。
- /chat 命令可以与 ChatGPT 畅聊, ChatGPT 是一种自然语言生成模型, 能够通过对话方式回答用户的问题, 例如: 「/chat 给我讲个笑话吧」
- /zlib 命令可以搜索 Z-Library 中的电子书, 加上 #ext(e.g: #pdf) 可以指定搜索的文件类型。例如「/zlib 如何阅读一本书」
- /magiccode 命令可以获取一个 Magic Code, Magic Code 可以用于网页登录。
- /whoami 命令查看个人信息。
- /sos 命令可以获取作者的微信, 我会尽量解答你的问题

基本操作
- 发送任意文字、图片或者视频到 NotionBoy 时, NotionBoy 会将内容保存到 Notion 中
- 如果发送到内容中包含 # 开头的内容, 会被自动识别成标签, 并在 Notion 中添加这个标签
- 如果发送的内容中包含 #全文和一个 URL, 则会自动保存此 URL 的全文内容到 Notion 中, 例如: 「#全文 https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw #NotionBoy」
`
	MSG_ERROR_ACCOUNT_NOT_FOUND = `查询账户失败:
- 如未绑定请回复「绑定」进行绑定
- 如已绑定请先回复「解绑」解除与 Notion 的绑定，再回复「绑定」进行绑定
- 如需帮助，请回复「帮助」

如果是想使用 Zlib 搜索, 可以不用绑定 Notion, 请按照下面的格式回复（不包括「」）即可
「/zlib 书名或者作者」
`
	MSG_ERROR_QUOTA_NOT_FOUND = `查询账户失败: 没有找到 Quota 信息, 请联系作者!`
	MSG_ERROR_QUOTA_LIMIT     = `额度已经用完, 请点击公众号菜单栏服务中的 VIP 进行充值`

	MSG_WELCOME = `#NotionBoy 欢迎🎉使用 Notion Boy!`

	MSG_RESET_CHATGPT_HISTORY = `已重置 ChatGPT 历史, 请输入「/chat 内容」重新开始`
	MSG_EMPTY_MESSAGE         = `消息为空, 请重新输入!`
	MSG_ZLIB_NO_RESULT        = `没有找到相关的图书, 请重新搜索! 可以尝试减少搜索关键词, 例如:「/zlib 如何阅读一本书」变成「/zlib 阅读」`
	MSG_ZLIB_TIPS             = `

Tips: When searching with a keyword containing #ext, you can specify the file type. For example, "/zlib How to read #pdf" will only search for books in pdf format.`
	MSG_ZLIB_TIPS_CN = `

Tips: 搜索关键字中包含 #ext 时可以指定文件类型，例如 「/zlib 鲁迅 #pdf」只会搜索 pdf 格式的书籍`

	MSG_USING_NOTION_TEST_ACCOUNT = "正在使用测试的 Notion 账户，数据只会保存 7 天，过期后会自动删除。如果需要长期保存，请回复「绑定」来绑定您的 Notion 账号。\n\n"
)

// const for command
const (
	CMD_BIND                = "/BIND"
	CMD_UNBIND              = "/UNBIND"
	CMD_FULLTEXT            = "全文"
	CMD_HELP_ZH             = "帮助"
	CMD_HELP                = "/HELP"
	CMD_SOS                 = "/SOS"
	CMD_CHAT                = "/CHAT"
	CMD_CHAT_RESET          = "RESET"
	CMD_ZLIB                = "/ZLIB"
	CMD_ZLIB_NEXT           = "ZLIBM"
	CMD_ZLIB_SAVE_TO_NOTION = "ZLIBS"
	CMD_UI                  = "/WEBUI"
	CMD_MAGIC_CODE          = "/MAGICCODE"
	CMD_WHOAMI              = "/WHOAMI"
	CMD_API_KEY             = "/APIKEY"
	MAGIC_CODE_CACHE_KEY    = "MAGIC_CODE_CACHE_KEY"
	QRCODE_CACHE_KEY        = "QRCODE_CACHE_KEY"
)

const (
	DB_DRIVER_SQLITE   = "sqlite"
	DB_DRIVER_MYSQL    = "mysql"
	DB_DRIVER_POSTGRES = "postgres"
)

const (
	CTX_KEY_QUOTA = "ctx_key_quota"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var DATABASE_ID contextKey = "database"

const (
	AUTH_HEADER_X_API_KEY  = "x-api-key"
	AUTH_HEADER_TOKEN_TYPE = "Bearer"
	AUTH_HEADER_COOKIE     = "cookie"
	AUTH_HEADER_TOKEN      = "token"
	AUTH_HEADER_PATH       = "path"
	AUTH_USER_ID           = "user_id"
	AUTH_USER_ACC          = "acc"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	ContextKeyUserId      ContextKey = AUTH_USER_ID
	ContextKeyUserAccount ContextKey = AUTH_USER_ACC
	ContentKeyTransaction ContextKey = "transaction"
	ContextKeyUserAgent   ContextKey = "grpcgateway-user-agent"
)
