package dto

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type EmailAttachment struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

type SendEmailRequest struct {
	To          []EmailAddress    `json:"to"`
	From        EmailAddress      `json:"from"`
	Subject     string            `json:"subject"`
	Text        string            `json:"text,omitempty"`
	HTML        string            `json:"html,omitempty"`
	CC          []EmailAddress    `json:"cc,omitempty"`
	BCC         []EmailAddress    `json:"bcc,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	CustomArgs  map[string]string `json:"custom_args,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

type SendEmailResponse struct {
	MessageIDs []string `json:"message_ids"`
}
