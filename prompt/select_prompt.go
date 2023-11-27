package prompt

import "github.com/manifoldco/promptui"

type SelectPrompt struct {
	s *promptui.Select
}

func (sp *SelectPrompt) Run() (int, string, error) {
	return sp.s.Run()
}

func (sp *SelectPrompt) SetItems(items []string) {
	sp.s.Items = items
}

func NewSelectPrompt() *SelectPrompt {
	return &SelectPrompt{
		s: &promptui.Select{},
	}
}
