package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/lordrusk/dctl/handler"
	"github.com/lordrusk/dctl/ui"
)

var OOOR = "Option out of range!"

var mainMenu = ui.Menu{
	"f": "Functions",
	"h": "Handlers",
	"t": "Tools",
	"x": "Exit",
}

func (c *client) mainLoop() error {
	for {
		k, _, err := c.Ask(mainMenu)
		if err != nil {
			errors.Wrap(err, "Failed to ask")
		}

		switch k {
		case "f":
			if err := c.functionsLoop(); err != nil {
				return err
			}
		case "h":
			if err := c.handlerLoop(); err != nil {
				return err
			}
		case "t":
			if err := c.toolLoop(); err != nil {
				return err
			}
		case "x":
			return nil
		default:
			fmt.Println(OOOR)
		}
	}
}

var functionsMenu = ui.Menu{
	"p": "Print Guilds",
	"P": "Print Guilds and Members",
	"m": "Print messages from a guild",
	"x": "Exit",
}

func (c *client) functionsLoop() error {
	for {
		k, _, err := c.Ask(functionsMenu)
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
			fmt.Println(OOOR)
		}
	}
}

var handlerMenu = ui.Menu{
	"inc": "Log all incoming messages",
	"x":   "Exit",
}

func (c *client) handlerLoop() error {
	for {
		k, _, err := c.Ask(handlerMenu)
		if err != nil {
			errors.Wrap(err, "Failed to ask")
		}

		switch k {
		case "inc":
			if handlers[k.(string)] == nil {
				handlers[k.(string)] = handler.New(c.incMsgsHF, c.errorEHF)
			}
			if handlers[k.(string)].Running() {
				handlers[k.(string)].Stop()
				fmt.Println("Incoming Messages Handler Stopped")
			} else {
				handlers[k.(string)].Start()
				fmt.Println("Incoming Messages Handler Started")
			}
		case "x":
			return nil
		default:
			fmt.Println(OOOR)
		}
	}
}

var toolMenu = ui.Menu{
	"cg": "Copy a guilds ID",
	"cc": "Copy a channels ID",
	"x":  "Exit",
}

func (c *client) toolLoop() error {
	for {
		k, _, err := c.Ask(toolMenu)
		if err != nil {
			errors.Wrap(err, "Failed to ask")
		}

		switch k {
		case "cg":
			if err := c.copyGuildID(); err != nil {
				fmt.Printf("Failed to copy guild id: %s\n", err)
			}
		case "cc":
			if err := c.copyChannelID(); err != nil {
				fmt.Printf("Failed to copy channel id: %s\n", err)
			}
		case "x":
			return nil
		default:
			fmt.Println(OOOR)
		}
	}
}
