package message

type EmailEvent struct {
    To      string `json:"to"      binding:"required,email"`
    Subject string `json:"subject" binding:"required"`
    Body    string `json:"body"    binding:"required"`
}