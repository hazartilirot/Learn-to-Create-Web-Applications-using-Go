package views

import (
	"github.com/username/project-name/models"
	"net/http"
	"time"
)

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
	User  *models.User
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

func persistAlert(w http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}
func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}

	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}

	return &alert
}

/*RedirectAlert accepts all the normal params for an http.Redirect and performs a redirect, but only
after persisting the provided alert in a cookie so that it can be displayed when the new page is loaded*/
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)
}
