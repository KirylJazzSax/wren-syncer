package prompt

type ConfirmPromptRunner interface {
	Run() (string, error)
	SetLabel(string)
}
