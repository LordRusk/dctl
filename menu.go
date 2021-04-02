package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/lordrusk/dctl/ui"
)

var mainMenu = ui.Menu{
	"p": "Print Guilds",
	"P": "Print Guilds and Members",
	"m": "Print messages from a guild",
	"x": "Exit",
}

func (c *client) mainLoop() error {
	for {
		k, _, err := c.Ask(mainMenu)
		if err != nil {
			errors.Wrap(err, "Failed to ask")
		}
		switch k {
		case "p":
			if err := c.printGuilds(); err != nil {
				// return errors.Wrap(err, "failed to print guilds")
				fmt.Printf("failed to print guilds: %s\n", err)
			}
		case "P":
			if err := c.printGuildPeople(); err != nil {
				// return errors.Wrap(err, "failed to print guilds and members")
				fmt.Printf("failed to print guilds and members: %s\n", err)
			}
		case "m":
			if err := c.printMessages(); err != nil {
				fmt.Printf("failed to print messages: %s\n", err)
			}
		case "x":
			return nil
		default:
			fmt.Println("Option out of range!")
		}
	}
}
