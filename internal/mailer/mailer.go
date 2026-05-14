package mailer

import "embed"

const (
	FromName            = "GopherSocail"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile string, username, email string, data any, isSandbox bool) (int, error)
}
