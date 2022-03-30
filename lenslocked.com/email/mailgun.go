package email

import "C"
import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

const (
	welcomeSubject = "Welcome to LensLocked.com!"

	welcomeText = `Hi there!
		Welcome to LensLocked.com! We really hope you enjoy using our application.

		Best Regards
	`
	welcomeHTML = `Hi there!</br>
		Welcome to <a href="https://lenslocked.com"</a>! We really hope you enjoy using our application.</br>
		</br>
		Best Regards
	`
)

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

func WithMailgun(domain, privateAPIKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, privateAPIKey)
		c.mg = mg
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "support@lenslocked.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := c.mg.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := c.mg.Send(ctx, message)

	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
