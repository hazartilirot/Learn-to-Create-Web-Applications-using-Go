package views

const (
	AlertError   = "danger"
	AlertWarning = "warning"
	AlertInfo    = "info"
	AlertSuccess = "success"
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
