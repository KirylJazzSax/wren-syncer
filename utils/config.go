package utils

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Username       string `mapstructure:"USERNAME"`
	JiraToken      string `mapstructure:"JIRA_TOKEN"`
	LinkAuthHeader string `mapstructure:"LINK_AUTH_HEADER"`
	JiraHost       string `mapstructure:"JIRA_HOST"`
	LinkHost       string `mapstructure:"LINK_HOST"`
	RequestTimeout string `mapstructure:"REQUEST_TIMEOUT"`
	JiraUsername   string `mapstructure:"JIRA_USERNAME"`
}
