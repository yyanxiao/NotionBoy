package prompt

import (
	"notionboy/api/pb/model"
	"notionboy/db/ent"
)

type PromptDTO struct {
	P *ent.Prompt
}

func NewPromptDTO(p *ent.Prompt) *PromptDTO {
	return &PromptDTO{p}
}

func (p *PromptDTO) ToProto() *model.Prompt {
	return &model.Prompt{
		Id:       p.P.UUID.String(),
		Act:      p.P.Act,
		Prompt:   p.P.Prompt,
		IsCustom: p.P.IsCustom,
	}
}
