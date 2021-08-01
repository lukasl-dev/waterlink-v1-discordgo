package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// registerHandlers registers discordgo event handlers.
func registerHandlers() {
	session.AddHandler(handleReady)
	session.AddHandler(handleVoiceUpdate)
	session.AddHandler(handleMessageCreate)
}

// handleReady handles the ready event.
func handleReady(_ *discordgo.Session, ready *discordgo.Ready) {
	sessionID = ready.SessionID
}

// handleVoiceUpdate handles the voice server update event.
func handleVoiceUpdate(_ *discordgo.Session, update *discordgo.VoiceServerUpdate) {
	err := conn.UpdateVoice(update.GuildID, sessionID, update.Token, update.Endpoint)
	if err != nil {
		log.Printf("Updating voice server failed on guild %s: %s\n", update.GuildID, err)
	} else {
		log.Printf("Updated voice server of guild %s.\n", update.GuildID)
	}
}

// handleMessageCreate handles the message create event.
func handleMessageCreate(_ *discordgo.Session, create *discordgo.MessageCreate) {
	msg := strings.TrimSpace(create.Message.Content) // trim trailing and leading whitespaces
	if !isFromGuild(create.Message) || !isPlay(msg) {
		return
	}
	identifier := lookupIdentifier(msg)
	if len(identifier) == 0 {
		_, _ = session.ChannelMessageSend(create.ChannelID, "No track specified. Use: !play <URL>")
		return
	}
	if ok := joinMemberChannel(create.ChannelID, create.GuildID, create.Author.ID); !ok {
		return
	}
	play(create.ChannelID, create.GuildID, identifier)
}

// isFromGuild returns true whenever the message has been created on a guild.
func isFromGuild(msg *discordgo.Message) bool {
	return msg.GuildID != ""
}

// isPlay returns true whenever the message starts with '!play'.
func isPlay(msg string) bool {
	return strings.HasPrefix(msg, "!play")
}

// lookupIdentifier returns the identifier from a message.
func lookupIdentifier(msg string) string {
	return strings.TrimSpace(strings.TrimLeft(msg, "!play"))
}

// joinMemberChannel joins the user's channel and returns true when it's succeed.
func joinMemberChannel(channelID, guildID, userID string) bool {
	vcID := findMembersChannel(guildID, userID)
	if vcID == "" {
		_, _ = session.ChannelMessageSend(channelID, "You must be in a voice channel.")
		return false
	}
	if err := session.ChannelVoiceJoinManual(guildID, vcID, false, true); err != nil {
		_, _ = session.ChannelMessageSend(channelID, "Could not join your voice channel.")
		return false
	}
	return true
}

// findMembersChannel searches for the user's channel on the guild and returns
// the channel's id.
func findMembersChannel(guildID, userID string) string {
	guild, err := session.State.Guild(guildID)
	if err != nil {
		return ""
	}
	for _, state := range guild.VoiceStates {
		if strings.EqualFold(userID, state.UserID) {
			return state.ChannelID
		}
	}
	return ""
}

// play plays the passed identifier on the guild.
func play(channelID, guildID, identifier string) {
	resp, err := req.LoadTracks(identifier)
	if err != nil {
		_, _ = session.ChannelMessageSend(channelID, fmt.Sprint("Could not load track:", err))
		return
	}
	track := resp.Tracks[0]
	if err := conn.Play(guildID, track.ID); err != nil {
		_, _ = session.ChannelMessageSend(channelID, fmt.Sprint("Could not play track:", err))
		return
	}
	_, _ = session.ChannelMessageSend(channelID, fmt.Sprintf("Now playing %s.", track.Info.Title))
}
