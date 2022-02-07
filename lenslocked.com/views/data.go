package views

const (
	AlertError   = "danger"
	AlertWarning = "warning"
	AlertInfo    = "info"
	AlertSuccess = "success"

	AlertMsgGeneric = "Something went wrong, please try again or contact us if the problem persists"
)

/*Alert is used to pass the seriousness level of a message and its context to a user */
type Alert struct {
	Level   string
	Message string
}

/*Data is the top-level structure views expect data to come in*/
type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.AlertMsg(pErr.Public())
	} else {
		d.AlertMsg(AlertMsgGeneric)
	}
}

func (d *Data) AlertError(msg string) {
	d.AlertMsg(msg)
}

func (d *Data) AlertMsg(msg string) {
	d.Alert = &Alert{
		Level:   AlertError,
		Message: msg,
	}
}

type PublicError interface {
	Error() string
	Public() string
}
