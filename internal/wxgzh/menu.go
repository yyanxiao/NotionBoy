package wxgzh

import (
	"fmt"

	"notionboy/internal/pkg/config"

	"github.com/silenceper/wechat/v2/officialaccount/menu"
)

type Btn int

func (b Btn) String() string {
	return fmt.Sprintf("btn_%d", b)
}

const (
	BtnBind Btn = iota
	BtnUnbind
	BtnMagicCode
	BtnHelpNote
	BtnHelpFulltext
	BtnHelpZlib
	BtnhelpSOS
	BtnWhoAMI
	BtnApiKey
)

func buildMenuButton() []*menu.Button {
	svcButton := buildServiceMenuButton()
	helpButton := buildHelpMenuButton()
	cmdButton := buildCmdMenuButton()
	buttons := []*menu.Button{svcButton, cmdButton, helpButton}
	return buttons
}

func buildCmdMenuButton() *menu.Button {
	bindButton := menu.NewClickButton("绑定 Notion", BtnBind.String())
	unbindButton := menu.NewClickButton("解绑 Notion", BtnUnbind.String())
	magicCodeButton := menu.NewClickButton("MagicCode", BtnMagicCode.String())
	whoAMIButton := menu.NewClickButton("个人信息", BtnWhoAMI.String())
	return menu.NewSubButton("常用命令", []*menu.Button{bindButton, unbindButton, magicCodeButton, whoAMIButton})
}

func buildServiceMenuButton() *menu.Button {
	chatGPTButton := menu.NewViewButton("ChatGPT",
		fmt.Sprintf("%s%s", config.GetConfig().Service.URL, "/web/chat.html"))
	vipButton := menu.NewViewButton("升级 VIP",
		fmt.Sprintf("%s%s", config.GetConfig().Service.URL, "/web/price.html"))
	svcButton := menu.NewSubButton("服务", []*menu.Button{vipButton, chatGPTButton})
	return svcButton
}

func buildHelpMenuButton() *menu.Button {
	noteButton := menu.NewClickButton("做笔记", BtnHelpNote.String())
	fulltextButton := menu.NewClickButton("全文剪藏", BtnHelpFulltext.String())
	zlibButton := menu.NewClickButton("Zlib", BtnHelpZlib.String())
	sosButton := menu.NewClickButton("联系作者", BtnhelpSOS.String())
	return menu.NewSubButton("帮助", []*menu.Button{noteButton, fulltextButton, zlibButton, sosButton})
}
