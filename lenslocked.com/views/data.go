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
