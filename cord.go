// simple discord functions
package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/lordrusk/dctl/dmenu"
)

var barrier = "----------"

// prints all guilds
func (c *client) printGuilds() error {
	guilds, err := c.Guilds()
	if err != nil {
		return errors.Wrap(err, "Failed to get guilds")
	}
	strs := make([]string, len(guilds))
	for pos, guild := range guilds {
		strs[pos] = guild.Name
	}
	c.Println(strings.Join(strs, "\n"))
	return nil
}

// prints all guilds and all members according
func (c *client) printGuildPeople() error {
	var b strings.Builder
	guilds, err := c.Guilds()
	if err != nil {
		return errors.Wrap(err, "Failed to get guilds")
	}
	for _, guild := range guilds {
		b.WriteString(guild.Name + ": ")
		members, err := c.Members(guild.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to get members")
		}
		for _, member := range members {
			b.WriteString(member.User.Username + ", ")
		}
		b.WriteString("\n")
	}
	c.Print(strings.ReplaceAll(b.String(), ", \n", ".\n"))
	return nil
}

// TODO fix this
func (c *client) printMessages() error {
	guilds, err := c.Guilds()
	if err != nil {
		return errors.Wrap(err, "could not fetch guilds")
	}
	guildNames := make(map[string]*discord.Guild)
	opts := make([]string, len(guilds)+1)
	fmt.Printf("Which guild do you want?\n%s\n", barrier)
	for pos, guild := range guilds {
		guildNames[guild.Name] = &guild
		opts[pos] = guild.Name
		fmt.Printf("%s\n", guild.Name)
	}
	opts[len(guilds)] = "x: Back"
	fmt.Println("x: Back")
	var guild *discord.Guild
	for {
		if *useDmenu {
			str, err := dmenu.PromptSlice(opts, defaultDmenuOpts)
			if err != nil {
				return errors.Wrap(err, "unable to use dmenu")
			}
			if guildNames[str] != nil {
				guild = guildNames[str]
				break
			}
		} else {
			str, err := c.Input(false)
			if err != nil {
				return errors.Wrap(err, "could not get user input")
			}
			if str == "x" {
				return nil
			}
			if guildNames[str] != nil {
				guild = guildNames[str]
				break
			}
		}
	}

	channels, err := c.Channels(guild.ID)
	if err != nil {
		return errors.Wrap(err, "could not fetch channels")
	}
	channelNames := make(map[string]*discord.Channel)
	opts = make([]string, len(channels)+1)
	fmt.Printf("Which channel do you want?\n%s\n", barrier)
	for pos, channel := range channels {
		channelNames[channel.Name] = &channel
		opts[pos] = channel.Name
		fmt.Printf("%s\n", channel.Name)
	}
	fmt.Println("x: Back")
	var channel *discord.Channel
	for {
		if *useDmenu {
			str, err := dmenu.PromptSlice(opts, defaultDmenuOpts)
			if err != nil {
				return errors.Wrap(err, "unable to use dmenu")
			}
			if channelNames[str] != nil {
				channel = channelNames[str]
				break
			}
		} else {
			str, err := c.Input(false)
			if err != nil {
				return errors.Wrap(err, "could not get user input")
			}
			if str == "x" {
				return nil
			}
			if channelNames[str] != nil {
				channel = channelNames[str]
				break
			}
		}
	}

	messages, err := c.Messages(channel.ID)
	if err != nil {
		return errors.Wrap(err, "could not fetch messages")
	}
	msgs := make([]string, *limit)
	for pos, msg := range messages {
		str, err := genMessage(guild, channel, &msg)
		if err != nil {
			return errors.Wrap(err, "failed to generate message")
		}
		msgs[pos] = str
	}
	for i := 0; i < 100; i++ { // this is done so all message are logged in order
		c.Println(msgs[i])
	}
	return nil
}
