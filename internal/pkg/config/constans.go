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

NotionBoy 提供以下命令可供使用:
- 帮助: 直接输入「帮助」或者「help」 可以获取最新的帮助教程
- 绑定: 直接输入「绑定」进行绑定 Notion 的操作
- 解绑: 直接输入「解绑」进行解除 Notion 绑定的操作
- 全文: 在存储链接的时候，加上「#全文」这个标签, 可以剪辑文章全文到 Notion 中, 例如: 「#全文 https://mp.weixin.qq.com/s/ib7HrRMIXwZjJyYFOwBQrw #NotionBoy」
- ChatGPT: 回复「/chat 聊天内容」，可以和机器人聊天, 例如: 「/chat 给我讲个笑话吧」
- Zlib: 回复 「/zlib 书名或者作者」可以获取图书的下载信息，例如「/zlib 如何阅读一本书」
- SOS: 回复「SOS」获取作者的微信, 我会尽量解答你的问题
`
	MSG_ERROR_ACCOUNT_NOT_FOUND = `查询账户失败:
- 如未绑定请回复「绑定」进行绑定
- 如已绑定请先回复「解绑」解除与 Notion 的绑定，再回复「绑定」进行绑定
- 如需帮助，请回复「帮助」

如果是想使用 Zlib 搜索, 可以不用绑定 Notion, 请按照下面的格式回复（不包括「」）即可
「/zlib 书名或者作者」
`
	MSG_WELCOME = `#NotionBoy 欢迎🎉使用 Notion Boy!`

	MSG_RESET_CHATGPT_HISTORY = `已重置 ChatGPT 历史, 请输入「/chat 内容」重新开始`
	MSG_EMPTY_MESSAGE         = `消息为空, 请重新输入!`
)

// const for command
const (
	CMD_BIND                = "绑定"
	CMD_UNBIND              = "解绑"
	CMD_FULLTEXT            = "全文"
	CMD_HELP_ZH             = "帮助"
	CMD_HELP                = "HELP"
	CMD_SOS                 = "SOS"
	CMD_CHAT                = "#CHAT"
	CMD_CHAT_SLASH          = "/CHAT"
	CMD_CHAT_RESET          = "RESET"
	CMD_ZLIB                = "/ZLIB"
	CMD_ZLIB_NEXT           = "ZLIBM"
	CMD_ZLIB_SAVE_TO_NOTION = "ZLIBS"
)

const (
	DB_DRIVER_SQLITE   = "sqlite"
	DB_DRIVER_MYSQL    = "mysql"
	DB_DRIVER_POSTGRES = "postgres"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var DATABASE_ID contextKey = "database"
