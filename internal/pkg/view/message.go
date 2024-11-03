package view

type Severity string

const (
	Warn  Severity = "Warn"
	Error Severity = "Error"
	Info  Severity = "Info"
)

type MessageViewModel struct {
	Msg      string
	Severity Severity
}

func NewMessage(msg string, severity Severity) MessageViewModel {
	return MessageViewModel{
		Msg:      msg,
		Severity: severity,
	}
}
