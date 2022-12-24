package notion

import (
	"strings"
	"sync"

	"github.com/jomei/notionapi"
)

var validCodeLanguages = []string{
	"abap", "arduino", "bash", "basic", "c", "clojure", "coffeescript",
	"c++", "c#", "css", "dart", "diff", "docker", "elixir", "elm",
	"erlang", "flow", "fortran", "f#", "gherkin", "glsl", "go", "graphql",
	"groovy", "haskell", "html", "java", "javascript", "json", "julia",
	"kotlin", "latex", "less", "lisp", "livescript", "lua", "makefile",
	"markdown", "markup", "matlab", "mermaid", "nix", "objective-c",
	"ocaml", "pascal", "perl", "php", "plain text", "powershell", "prolog",
	"protobuf", "python", "r", "reason", "ruby", "rust", "sass", "scala",
	"scheme", "scss", "shell", "sql", "swift", "typescript", "vb.net",
	"verilog", "vhdl", "visual basic", "webassembly", "xml", "yaml", "java/c/c++/c#",
}

var (
	once                     sync.Once
	validCodeLanguageMapping map[string]struct{}
)

type ChatContent struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	UserID   string `json:"user_id"`
}

func (c *ChatContent) BuildBlocks() []notionapi.Block {
	blocks := make([]notionapi.Block, 0)
	blocks = append(blocks, notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: c.Question,
					},
					Annotations: &notionapi.Annotations{
						Bold: true,
					},
				},
			},
		},
	})
	if c.UserID != "" {
		blocks = append(blocks, notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Mention: &notionapi.Mention{
							Type: "user",
							User: &notionapi.User{
								ID: notionapi.UserID(c.UserID),
							},
						},
					},
				},
			},
		})
	}
	blocks = formatAnswer(blocks, c.Answer)
	return blocks
}

func formatAnswer(blocks []notionapi.Block, answer string) []notionapi.Block {
	content := make([]string, 0)
	codeStart := true
	language := "plain text"
	for _, txt := range strings.Split(answer, "\n") {
		if strings.HasPrefix(txt, "```") {
			if codeStart {
				blocks = append(blocks, buildTextBlock(content))
				content = []string{}
				codeStart = false
				if txt[3:] != "" && isValidCodeLanguage(txt[3:]) {
					language = txt[3:]
				}
			} else {
				blocks = append(blocks, buildCodeBlock(content, language))
				content = []string{}
				codeStart = true
			}
		} else {
			content = append(content, txt)
		}
	}
	if len(content) > 0 {
		if codeStart {
			blocks = append(blocks, buildTextBlock(content))
		} else {
			blocks = append(blocks, buildCodeBlock(content, language))
		}
	}
	return blocks
}

func buildTextBlock(content []string) notionapi.Block {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: strings.Join(content, "\n"),
					},
				},
			},
		},
	}
}

func buildCodeBlock(content []string, language string) notionapi.Block {
	return notionapi.CodeBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeCode,
		},
		Code: notionapi.Code{
			Language: language,
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: strings.Join(content, "\n"),
					},
				},
			},
		},
	}
}

func isValidCodeLanguage(language string) bool {
	if validCodeLanguageMapping == nil {
		once.Do(func() {
			validCodeLanguageMapping = make(map[string]struct{})
			for _, lan := range validCodeLanguages {
				validCodeLanguageMapping[lan] = struct{}{}
			}
		})
	}

	_, ok := validCodeLanguageMapping[language]
	return ok
}
