package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/lordrusk/dctl/ui"
)

// this was ripped from arikawa/session/session.go
var emailLogin = func(email, pass, mfa string) (*state.State, error) {
	client := api.NewClient("")
	l, err := client.Login(email, pass)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	if l.Token != "" && !l.MFA {
		s, err := state.New(l.Token)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get state")
		}
		return s, nil
	}

	// discord requests MFA, so we need the MFA token.
	if mfa == "" {
		return nil, session.ErrMFA
	}

	// retry logging in with a 2FA token
	l, err = client.TOTP(mfa, l.Ticket)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login with 2FA")
	}

	s, err := state.New(l.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get state")
	}
	return s, nil
}

var uiLoginMenu = ui.Menu{
	"t": "Token Login",
	"e": "Email and Passwod Login (wip)",
	"x": "Exit",
}

// if returned error is nil, c.State will not be nil
func (c *client) uiLogin() error {
	for {
		k, _, err := c.Ask(uiLoginMenu)
		if err != nil {
			return errors.Wrap(err, "Failed to ask")
		}
		switch k {
		case "t":
			var err error
			*token, err = c.NextString("Enter token (use -b for bots): ")
			if err != nil {
				return errors.Wrap(err, "Failed getting user input")
			}
			if *isBot {
				*token = "Bot " + *token
			}

			s, err := state.New(*token)
			if err != nil {
				return errors.Wrap(err, "Failed to get state")
			}

			c.State = s
			return nil
		case "e":
			email, err := c.NextString("Enter email: ")
			if err != nil {
				return errors.Wrap(err, "Failed getting user input for email")
			}
			pass, err := c.NextString("Enter password: ")
			if err != nil {
				return errors.Wrap(err, "Failed getting user input for pass")
			}
			mfa, err := c.NextString("Enter MFA (leave blank if unsure): ")
			if err != nil {
				return errors.Wrap(err, "Failed getting user input for mfa")
			}
			s, err := emailLogin(email, pass, mfa)
			if err != nil {
				return errors.Wrap(err, "Failed to login with email")
			}
			c.State = s
			return nil
		case "x":
			os.Exit(0)
		default:
			fmt.Println("Option out of range!")
		}
	}

	return errors.New("You shouldn't be here...")
}
