package prompt

type SelectPromptRunner interface {
	Run() (int, string, error)
	SetItems([]string)
}
