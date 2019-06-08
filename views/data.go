package views

import "log"

// Data is the top-level structure that views expect data
// to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// PublicError establishes public-facing error strings, effectively
// white-listing error messages on types that support the
// PublicError interface
type PublicError interface {
	error
	Public() string
}

// Constants for use with alerts
const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error
	// is encountered by our backend
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// SetAlert sanitizes error messages with white-listed strings, and log the message
func (d *Data) SetAlert(err error) {
	var msg string

	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		msg = AlertMsgGeneric
	}
	log.Println(err)

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
