// This is where all code is for generating
// plaintext(UTF-8) from discord structures.
package main

import (
	"fmt"
	"strings"

	// "github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/discord"
)

func genMessage(g *discord.Guild, ch *discord.Channel, m *discord.Message) (string, error) {
	return fmt.Sprintf("%s: %s: %s\n", g.Name, ch.Name, strings.TrimSpace(m.Content)), nil
}
