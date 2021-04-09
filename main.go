package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/lordrusk/dctl/logger"
	// "github.com/diamondburned/arikawa/v2/voice"
)

// this fixes there being an extra
// space whenever the day is
// less than two digits long
var snow = strings.Split(strings.ReplaceAll(time.Now().Format(time.UnixDate), "  ", " "), " ")
var logFile = flag.String("f",
	fmt.Sprintf("~/.local/share/dctl/%s-%s-%s-%s.log",
		snow[5], snow[1], snow[2], snow[3]), "")

var token = flag.String("t", "", "Set the token (skips login)")
var isBot = flag.Bool("b", false, "Whether or not the account is a bot")

// unused as of now
var useDmenu = flag.Bool("d", false, "Use dmenu for list inputs (dmenu must be installed)")
var defaultDmenuOpts = []string{"-l", "10"} // nice looking defaults

var limit = flag.Uint("l", 100, "Set the limit for discord requests") // default api value

func main() {
	flag.Parse()
	l, err := logger.New(os.Stdin, "", 0, *logFile)
	if err != nil {
		fmt.Printf("Failed to get client: %s\n", err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			fmt.Printf("Failed to close logger: %s\n", err)
		}
	}()
	c := newClient(l, setupDefaultHandlerFunc)

	if *token == "" {
		if err := c.uiLogin(); err != nil {
			fmt.Printf("Failed logging in: %s\n", err)
			return
		}
	} else {
		if *isBot {
			*token = "Bot " + *token
		}
		s, err := state.New(*token)
		if err != nil {
			fmt.Println("Could not get state: %s\n", err)
			os.Exit(1)
		}
		c.State = s
	}
	if c.State == nil { // c.State should never be nil at this point
		fmt.Println("somehow failed to get the state?")
		os.Exit(1)
	}

	addIntents(c)

	if err := c.mainLoop(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

// i love the discord api 💜💜😘
func addIntents(c *client) {
	c.Gateway.AddIntents(gateway.IntentGuilds)
	// c.Gateway.AddIntents(gateway.IntentGuildMembers)
	// c.Gateway.AddIntents(gateway.IntentGuildBans)
	// c.Gateway.AddIntents(gateway.IntentGuildEmojis)
	// c.Gateway.AddIntents(gateway.IntentGuildIntegrations)
	// c.Gateway.AddIntents(gateway.IntentGuildWebhooks)
	// c.Gateway.AddIntents(gateway.IntentGuildInvites)
	c.Gateway.AddIntents(gateway.IntentGuildVoiceStates)
	// c.Gateway.AddIntents(gateway.IntentGuildPresences)
	c.Gateway.AddIntents(gateway.IntentGuildMessages)
	// c.Gateway.AddIntents(gateway.IntentGuildMessageReactions)
	// c.Gateway.AddIntents(gateway.IntentGuildMessageTyping)
	c.Gateway.AddIntents(gateway.IntentDirectMessages)
	// c.Gateway.AddIntents(gateway.IntentDirectMessageReactions)
	// c.Gateway.AddIntents(gateway.IntentDirectMessageTyping)

	// voice.AddIntents(c.Gateway) // for voice
}
