// *most* everything pertaining to handlers
package main

import (
	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/dctl/handler"
)

var handlers = make(map[string]handler.Handler)

func (c *client) errorEHF(err error) { c.Printf("Handler failed: %s\n", err) }

func (c *client) incMsgsHF(killCh chan struct{}, errCh chan error) {
	cancel := c.AddHandler(func(m *gateway.MessageCreateEvent) {
		guild, err := c.Guild(m.GuildID)
		if err != nil {
			errCh <- errors.Wrap(err, "Failed to get guild")
		}
		channel, err := c.Channel(m.ChannelID)
		if err != nil {
			errCh <- errors.Wrap(err, "Failed to get channel")
		}
		msg, err := genMessage(guild, channel, &m.Message, c, true)
		if err != nil {
			errCh <- errors.Wrap(err, "Failed to generate message")
		}
		c.Println(msg)
	})

	go func() {
		select {
		case <-killCh:
			cancel()
		}
	}()
}
