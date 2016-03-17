package slacker

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	STATUS_NONE   uint = 0
	STATUS_TYPING      = 1
	STATUS_EMOJI       = 2
)

// Bot is basic bot struct
type Bot struct {
	Token      string // Auth token
	BotName    string // Bot name
	AllowDM    bool   // Allow direct messages
	StatusType uint8  // Processing status mark type
	Started    int64  // Bot start timestamp

	// Async handlers
	ErrorHandler   func(err error)
	ConnectHandler func()
	HelloHandler   func() string
	CommandHandler func(command string, args []string) []string

	client *slack.Client
	botID  string
	works  bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewBot return new bot struct
func NewBot(name, token string) *Bot {
	return &Bot{
		Token:      token,
		BotName:    name,
		StatusType: STATUS_TYPING,
		AllowDM:    true,
	}
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
	b.works = true

	b.rtmLoop()

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// rtmLoop is rtm processing loop
func (b *Bot) rtmLoop() {
	rtm := b.client.NewRTM()
	go rtm.ManageConnection()

LOOP:
	for {
		select {
		case event := <-rtm.IncomingEvents:
			switch event.Data.(type) {
			case *slack.ConnectedEvent:
				if b.ConnectHandler != nil {
					b.ConnectHandler()
				}

			case *slack.ChannelJoinedEvent:
				if b.HelloHandler != nil {
					joinedEvent := event.Data.(*slack.ChannelJoinedEvent)
					response := b.HelloHandler()

					if response != "" {
						rtm.SendMessage(rtm.NewOutgoingMessage(response, joinedEvent.Channel.ID))
					}
				}

			case *slack.MessageEvent:
				msgEvent := event.Data.(*slack.MessageEvent)

				if !b.isBotCommand(msgEvent.Text, msgEvent.Channel) ||
					b.CommandHandler == nil || msgEvent.User == b.botID {
					continue
				}

				switch b.StatusType {
				case STATUS_EMOJI:
					b.client.AddReaction(
						"white_check_mark",
						slack.NewRefToMessage(msgEvent.Channel, msgEvent.Timestamp),
					)

				case STATUS_TYPING:
					rtm.SendMessage(rtm.NewTypingMessage(msgEvent.Channel))
				}

				cmd, args := extractCommand(msgEvent.Text)
				responses := b.CommandHandler(cmd, args)

				if len(responses) != 0 {
					for _, response := range responses {
						rtm.SendMessage(rtm.NewOutgoingMessage(response, msgEvent.Channel))
					}
				}

			case *slack.RTMError:
				errEvent := event.Data.(*slack.RTMError)
				b.ErrorHandler(fmt.Errorf(errEvent.Error()))

			case *slack.InvalidAuthEvent:
				b.ErrorHandler(fmt.Errorf("Can't authorize with given token"))
				b.Started = 0
				b.works = false
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
