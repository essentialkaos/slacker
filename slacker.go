// Package slacker provides methods for for bootstraping Slack bots
package slacker

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Status types
const (
	STATUS_NONE   uint = 0
	STATUS_TYPING      = 1
	STATUS_EMOJI       = 2
)

// VERSION is current package version
const VERSION = "9.0.0"

// ////////////////////////////////////////////////////////////////////////////////// //

// User is alias for slack.User struct
type User = slack.User

// CommandHandler is command handler
type CommandHandler func(user User, args []string) []string

// Bot is basic bot struct
type Bot struct {
	Token                string // Auth token
	BotName              string // Bot name
	AllowDM              bool   // Allow direct messages
	StatusType           uint8  // Processing status mark type
	Started              int64  // Bot start timestamp
	UserListUpdatePeriod int    // User list update period in seconds

	ErrorHandler          func(err error)
	ConnectHandler        func()
	HelloHandler          func() string
	UnknownCommandHandler func(user User, cmd string, args []string) string

	CommandHandlers map[string]CommandHandler

	rtm       *slack.RTM
	client    *slack.Client
	usersInfo *UsersInfo
	botID     string
	works     bool
}

// Basic users info
type UsersInfo struct {
	Users map[string]User

	lastUpdate int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewBot creates new bot
func NewBot(name, token string) *Bot {
	return &Bot{
		Token:                token,
		BotName:              name,
		StatusType:           STATUS_TYPING,
		AllowDM:              true,
		UserListUpdatePeriod: 86400,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetMention generates user mention link
func GetMention(u User) string {
	if u.ID == "" {
		return "Uknown User"
	}

	return "<@" + u.ID + ">"
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run starts bot
func (b *Bot) Run() error {
	if b.works {
		return nil
	}

	b.client = slack.New(b.Token)

	authResp, err := b.client.AuthTest()

	if err != nil {
		return err
	}

	// Return error if we got info not for our bot
	if authResp.User != b.BotName {
		return fmt.Errorf("Unknown user: %s ≠ %s", authResp.User, b.BotName)
	}

	b.botID = authResp.UserID
	b.Started = time.Now().Unix()
	b.usersInfo = &UsersInfo{}
	b.works = true

	b.FetchUsers()
	b.rtmLoop()

	return nil
}

// fetchUsers fetches users and creates map id->name
func (b *Bot) FetchUsers() error {
	users, err := b.client.GetUsers()

	if err != nil {
		return err
	}

	b.usersInfo.Users = make(map[string]User)

	for _, user := range users {
		b.usersInfo.Users[user.ID] = user
	}

	b.usersInfo.lastUpdate = time.Now().Unix()

	return nil
}

// GetUser tries to find user info by user name or ID
func (b *Bot) GetUser(nameOrID string) User {
	if strings.Contains(nameOrID, "@") {
		id := strings.Trim(nameOrID, "<@>")
		return b.usersInfo.Users[id]
	}

	for _, user := range b.usersInfo.Users {
		if user.Name == nameOrID {
			return user
		}
	}

	return User{}
}

// NormalizeInput normalizes links and usernames in a message
func (b *Bot) NormalizeInput(input string) string {
	if input == "" {
		return ""
	}

	var result []string

	inputSlice := strings.Split(input, " ")

	for _, t := range inputSlice {
		if strings.HasPrefix(t, "<http") && strings.Contains(t, "|") && strings.HasSuffix(t, ">") {
			result = append(result, t[strings.Index(t, "|")+1:len(t)-1])
			continue
		}

		if strings.HasPrefix(t, "<@U") && strings.HasSuffix(t, ">") {
			user := b.GetUser(t)

			if user.Name == "" {
				result = append(result, t)
			} else {
				result = append(result, "@"+user.Name)
			}

			continue
		}

		result = append(result, t)
	}

	return strings.Join(result, " ")
}

// SendMessage sends simple message to some user
func (b *Bot) SendMessage(to, message string) error {
	user := b.GetUser(to)

	if user.ID == "" {
		return fmt.Errorf("Can't find user %s", to)
	}

	return b.PostMessage(user.ID, message, slack.MsgOptionAsUser(true))
}

// PostMessage posts mesasge
func (b *Bot) PostMessage(channel, message string, options ...slack.MsgOption) error {
	_, _, err := b.client.PostMessage(channel, append(options, slack.MsgOptionText(message, false))...)

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

// rtmLoop is rtm processing loop
func (b *Bot) rtmLoop() {
	rtm := b.client.NewRTM()
	go rtm.ManageConnection()
	b.rtm = rtm

LOOP:
	for {
		if time.Now().Unix() >= b.usersInfo.lastUpdate+int64(b.UserListUpdatePeriod) {
			b.FetchUsers()
		}

		select {
		case event := <-rtm.IncomingEvents:
			switch event.Data.(type) {
			case *slack.ConnectedEvent:
				b.processConnectedEvent(event.Data.(*slack.ConnectedEvent))

			case *slack.ChannelJoinedEvent:
				b.processChannelJoinedEvent(event.Data.(*slack.ChannelJoinedEvent))

			case *slack.MessageEvent:
				b.processMessageEvent(event.Data.(*slack.MessageEvent))

			case *slack.RTMError:
				b.processRTMError(event.Data.(*slack.RTMError))

			case *slack.InvalidAuthEvent:
				b.processInvalidAuthEvent(event.Data.(*slack.InvalidAuthEvent))
				break LOOP
			}
		}
	}
}

// isBotCommand return true if it message for our bot
func (b *Bot) isBotCommand(message, channel string) bool {
	if len(channel) > 2 && channel[0:1] == "D" {
		return b.AllowDM
	}

	return strings.HasPrefix(message, "<@"+b.botID+">")
}

// processConnectEvent is Connected event handler
func (b *Bot) processConnectedEvent(ev *slack.ConnectedEvent) {
	if b.ConnectHandler != nil {
		b.ConnectHandler()
	}
}

// processChannelJoinedEvent is ChannelJoined event handler
func (b *Bot) processChannelJoinedEvent(ev *slack.ChannelJoinedEvent) {
	if b.HelloHandler != nil {
		response := b.HelloHandler()

		if response != "" {
			b.rtm.SendMessage(b.rtm.NewOutgoingMessage(response, ev.Channel.ID))
		}
	}
}

// processMessageEvent is Message event handler
func (b *Bot) processMessageEvent(ev *slack.MessageEvent) {
	if !b.isBotCommand(ev.Text, ev.Channel) || ev.User == b.botID {
		return
	}

	cmd, args := extractCommand(ev.Text)

	if b.CommandHandlers == nil || cmd == "" {
		return
	}

	b.setCommandStatus(ev)
	b.execHandler(ev, cmd, args)
}

// processRTMError is RTMError event handler
func (b *Bot) processRTMError(ev *slack.RTMError) {
	if b.ErrorHandler != nil {
		b.ErrorHandler(fmt.Errorf(ev.Error()))
	}
}

// processInvalidAuthEvent is InvalidAuth event handler
func (b *Bot) processInvalidAuthEvent(ev *slack.InvalidAuthEvent) {
	if b.ErrorHandler != nil {
		b.ErrorHandler(fmt.Errorf("Can't authorize with given token"))
	}

	b.Started = 0
	b.works = false
}

// execHandler execute command handler
func (b *Bot) execHandler(ev *slack.MessageEvent, cmd string, args []string) {
	user := b.usersInfo.Users[ev.User]
	handler := b.CommandHandlers[cmd]

	if handler == nil {
		if b.UnknownCommandHandler != nil {
			b.rtm.SendMessage(
				b.rtm.NewOutgoingMessage(
					b.UnknownCommandHandler(user, cmd, args),
					ev.Channel,
				),
			)
		}

		return
	}

	responses := handler(user, args)

	if len(responses) != 0 {
		for _, response := range responses {
			b.rtm.SendMessage(b.rtm.NewOutgoingMessage(response, ev.Channel))
		}
	}
}

// setCommandStatus set
func (b *Bot) setCommandStatus(ev *slack.MessageEvent) {
	switch b.StatusType {
	case STATUS_EMOJI:
		b.client.AddReaction(
			"white_check_mark",
			slack.NewRefToMessage(ev.Channel, ev.Timestamp),
		)

	case STATUS_TYPING:
		b.rtm.SendMessage(b.rtm.NewTypingMessage(ev.Channel))
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// extractCommand extracts command and arguments from user message
func extractCommand(message string) (string, []string) {
	if message == "" {
		return "", []string{}
	}

	// Remove bot id from message
	if strings.HasPrefix(message, "<@") {
		message = message[12:]
	}

	// Remove separators from message
	message = strings.TrimLeft(message, ": ")

	if message == "" {
		return "", []string{}
	}

	messageSlice := strings.Split(message, " ")

	return messageSlice[0], messageSlice[1:]
}
