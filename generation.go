// This is where all code is for generating
// plaintext(UTF-8) from discord structures.
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/discord"
)

// gen stuff
var atRegex = regexp.MustCompile(`<.{1,2}[0-9]+>`)
var Tab = "	" // a single tab

// at parsing functions
func memberAt(id int, gid discord.GuildID, c *client) (string, error) { // member
	member, err := c.Member(gid, discord.UserID(id))
	if err != nil {
		return "", errors.Wrap(err, "Failed to get member")
	}
	return "@" + member.User.Username, nil
}

func channelAt(id int, gid discord.GuildID, c *client) (string, error) { // channel
	channel, err := c.Channel(discord.ChannelID(id))
	if err != nil {
		return "", errors.Wrap(err, "Failed to get channel")
	}
	return "#" + channel.Name, nil
}

var atMap = map[string]func(int, discord.GuildID, *client) (string, error){
	"@":  memberAt,
	"@!": memberAt,
	"@&": channelAt,
}

// str should look something like <@403697455331278848>
func atReplace(str string, guild *discord.Guild, c *client) (string, error) {
	var id int
	for at, f := range atMap {
		if strings.Contains(str, at) {
			var err error
			id, err = strconv.Atoi(strings.Trim(str, at+"<>"))
			if err != nil {
				return "", errors.New("Failed to convert id to int!")
			}
			return f(id, guild.ID, c)
		}
	}
	return "", errors.New("Unkown @ tag or not @!")
}

// trims the extention of a file
// index.md -> index
func trimExt(str string) string {
	sstr := strings.Split(str, ".")
	return strings.Join(sstr[:len(sstr)-1], ".")
}

// dlAtt is whether to download attachments with messages
func genMessage(guild *discord.Guild, channel *discord.Channel, message *discord.Message, c *client, dlAtt bool) (string, error) {
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
				sstr, err := atReplace(at, guild, c)
				if err != nil {

				}
				str = strings.Replace(str, at, sstr, -1)
				fin[at] = true
			}
		}
	}

	for _, attachment := range message.Attachments {
		if dlAtt { // download if needed
			if err := c.Download(attachment.URL, trimExt(*logFile)+"/"+strconv.Itoa(int(message.ID))); err != nil {
				c.Printf("Failed to download '%s' to '%s'", attachment.URL, trimExt(*logFile)+"/"+strconv.Itoa(int(message.ID)))
			}
		}
		str += fmt.Sprintf("\n%sName: %s, URL: %s", Tab, strings.TrimSpace(attachment.Filename), attachment.URL)
	}

	return str, nil
}
