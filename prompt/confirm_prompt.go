package prompt

import "github.com/manifoldco/promptui"

type ConfirmPrompt struct {
	p *promptui.Prompt
}

func (cp *ConfirmPrompt) SetLabel(label string) {
	cp.p.Label = label
}

func (cp *ConfirmPrompt) Run() (string, error) {
	return cp.p.Run()
}

func NewConfirmPrompt() *ConfirmPrompt {
	return &ConfirmPrompt{
		p: &promptui.Prompt{
			IsConfirm: true,
		},
	}
}
