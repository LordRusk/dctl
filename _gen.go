// handles generating readable plantext(utf8) out of
// discord structs and such.
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// gen stuff
var atRegex = regexp.MustCompile(`<.{1,2}[0-9]+>`)
var Tab = "	" // a single tab

// at parsing functions
func memberAt(id int, gid discord.GuildID) (string, error) { // member
	member, err := c.Member(gid, discord.UserID(id))
	if err != nil {
		return "", fmt.Errorf("Failed to get member: %s\n", err)
	}
	return "@" + member.User.Username, nil
}

func channelAt(id int, gid discord.GuildID) (string, error) { // channel
	channel, err := c.Channel(discord.ChannelID(id))
	if err != nil {
		return "", fmt.Errorf("Failed to get channel: %s\n", err)
	}
	return "#" + channel.Name, nil
}

var atMap = map[string]func(int, discord.GuildID) (string, error){
	"@":  memberAt,
	"@!": memberAt,
	"@&": channelAt,
}

// str should look something like <@403697455331278848>
func atReplace(str string, guild *discord.Guild) (string, error) {
	var id int
	for at, f := range atMap {
		if strings.Contains(str, at) {
			var err error
			id, err = strconv.Atoi(strings.Trim(str, at+"<>"))
			if err != nil {
				return "", fmt.Errorf("Failed to convert id to int!")
			}
			return f(id, guild.ID)
		}
	}
	return "", errors.New("Unkown @ tag or not @!")
}

func genMessage(guild *discord.Guild, channel *discord.Channel, message *discord.Message) (string, error) {
	var str string
	if message.Timestamp.IsValid() {
		h, m, s := message.Timestamp.Time().Clock()
		str = fmt.Sprintf("%s: %s %d:%d:%d %s -> %s | %d embeds | ID: %s", guild.Name, channel.Name, h, m, s, message.Author.Username, message.Content, len(message.Embeds), message.ID)
	} else {
		str = fmt.Sprintf("%s: %s: %s -> %s | %d embeds | ID: %s", guild.Name, channel.Name, message.Author.Username, message.Content, len(message.Embeds), message.ID)
	}

	if atRegex.MatchString(str) {
		strs := atRegex.FindAllString(str, -1)
		fin := make(map[string]bool)
		for _, at := range strs {
			if !fin[at] {
				sstr, err := atReplace(at, guild)
				if err != nil {

				}
				str = strings.Replace(str, at, sstr, -1)
				fin[at] = true
			}
		}
	}

	for _, attachment := range message.Attachments {
		str += fmt.Sprintf("\n%sName: %s, URL: %s", Tab, strings.TrimSpace(attachment.Filename), attachment.URL)
	}

	return str, nil
}
