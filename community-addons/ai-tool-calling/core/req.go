package core

type PromptReq struct {
	Prompt string `json:"prompt"`
}

func (p PromptReq) Validate() error {
	return validatePrompt(p.Prompt)
}
