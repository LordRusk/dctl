// simple discord functions
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/atotto/clipboard"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/lordrusk/dctl/ui"
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

func (c *client) copyGuildID() error {
	guilds, err := c.Guilds()
	if err != nil {
		return errors.Wrap(err, "Failed to get guilds")
	}
	m := make(ui.Menu)
	for _, guild := range guilds {
		m[guild.Name] = guild.ID
	}
	m["x"] = "Exit"

	for {
		k, _, err := c.Ask(m)
		if err != nil {
			return errors.Wrap(err, "Failed to ask")
		}
		if k.(string) == "x" {
			return nil
		}

		if m[k.(string)] == nil {
			fmt.Println("Unrecognized guild")
			continue
		}
		if err := clipboard.WriteAll(strconv.Itoa(int(m[k.(string)].(discord.GuildID)))); err != nil { // divert your eyes from this shit
			return errors.Wrap(err, "Failed to copy ID to clipboard")
		}
		fmt.Printf("%s copied to clipboard!\n", m[k.(string)])
		return nil
	}
}

func (c *client) copyChannelID() error {
	fmt.Print("Enter guild's id: ")
	id, err := c.NextInt()
	if err != nil {
		return errors.Wrap(err, "Failed to get input")
	}
	channels, err := c.Channels(discord.GuildID(id))
	if err != nil {
		return errors.Wrap(err, "Failed to get channels")
	}

	m := make(ui.Menu)
	for _, channel := range channels {
		if channel.Type == 0 || channel.Type == 5 {
			m[channel.Name] = channel.ID
		}
	}
	m["x"] = "Exit"

	for {
		k, _, err := c.Ask(m)
		if err != nil {
			return errors.Wrap(err, "Failed to ask")
		}
		if k.(string) == "x" {
			return nil
		}

		if m[k.(string)] == nil {
			fmt.Println("Unrecognized channel")
			continue
		}
		if err := clipboard.WriteAll(strconv.Itoa(int(m[k.(string)].(discord.ChannelID)))); err != nil { // divert your eyes from this shit
			return errors.Wrap(err, "Failed to copy ID to clipboard")
		}
		fmt.Printf("%s copied to clipboard!\n", m[k.(string)])
		return nil
	}
}

func (c *client) printMessages() error {
	fmt.Print("Enter guild's id: ")
	gid, err := c.NextInt()
	if err != nil {
		return errors.Wrap(err, "Unable to get input")
	}
	fmt.Print("Enter channel's id: ")
	cidstr, err := c.NextString()
	if err != nil {
		return errors.Wrap(err, "Unable to get input")
	}

	guilds, err := c.Guilds()
	if err != nil {
		return errors.Wrap(err, "Unable to get guilds")
	}
	var g *discord.Guild
	for _, guild := range guilds {
		if guild.ID == discord.GuildID(gid) {
			g = &guild
			break
		}
	}
	if g == nil {
		return errors.New(fmt.Sprintf("No guild with id '%d'", gid))
	}

	channels, err := c.Channels(g.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to get channels")
	}

	msgs := make([]string, *limit)
	cid, err := strconv.Atoi(cidstr)
	if err != nil {
		if strings.ToLower(cidstr) != "all" {
			return errors.New(fmt.Sprintf("No channel with id '%s'", cidstr))
		}
		for _, channel := range channels {
			if channel.Type != 0 && channel.Type != 5 {
				continue
			}
			messages, err := c.Messages(channel.ID)
			if err != nil {
				return errors.Wrap(err, "Failed to get messages")
			}
			msgss := make([]string, len(messages)) // I'm not good at naming I know
			for pos, message := range messages {
				msg, err := genMessage(g, &channel, &message)
				if err != nil {
					return errors.Wrap(err, "Failed to generate message")
				}
				msgss[pos] = msg
			}
			msgs = append(msgs, msgss[:]...)
		}
	} else {
		var ch *discord.Channel
		for _, channel := range channels {
			if channel.ID == discord.ChannelID(cid) {
				ch = &channel
				break
			}
		}
		if ch == nil {
			return errors.New(fmt.Sprintf("No channel with id '%d'", cid))
		}

		messages, err := c.Messages(ch.ID)
		if err != nil {
			return errors.Wrap(err, "Failed to get messages")
		}
		for pos, message := range messages { // this is done to keep messages ordered
			msg, err := genMessage(g, ch, &message)
			if err != nil {
				return errors.Wrap(err, "Failed to generate message")
			}
			msgs[pos] = msg
		}
	}
	for i := 0; i < len(msgs); i++ {
		c.Println(msgs[i])
	}

	return nil
}
