package slacklog

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/bluele/slack"
	"github.com/sirupsen/logrus"
)

// SlackHook is a logrus Hook for dispatching messages to the specified
// channel on Slack.
type SlackHook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	HookURL        string // Webhook URL

	// slack post parameters
	Username  string // display name
	Channel   string // `#channel-name`
	IconEmoji string // emoji string ex) ":ghost:":
	IconURL   string // icon url

	FieldHeader string        // a header above field data
	Timeout     time.Duration // request timeout
	Async       bool          // if async is true, send a message asynchronously.
	Options     map[string]interface{}
	Disabled    bool

	hook *slack.WebHook
}

// Fire -  Sent event to slack
func (sh *SlackHook) Fire(e *logrus.Entry) error {
	if sh.Disabled {
		return nil
	}
	if sh.hook == nil {
		sh.hook = slack.NewWebHook(sh.HookURL)
	}

	payload := &slack.WebHookPostPayload{
		Username:  sh.Username,
		Channel:   sh.Channel,
		IconEmoji: sh.IconEmoji,
		IconUrl:   sh.IconURL,
	}
	color, _ := LevelColorMap[e.Level]

	attachment := slack.Attachment{}
	payload.Attachments = []*slack.Attachment{&attachment}
	// fetch all entries and add as attachment
	allEntries := sh.newEntry(e)
	// If there are fields we need to render them at attachments
	if len(allEntries.Data) > 0 {
		// Add a header above field data
		attachment.Text = sh.FieldHeader

		for k, v := range allEntries.Data {
			field := &slack.AttachmentField{}
			field.Title = k
			if str, ok := v.(string); ok {
				field.Value = str
			} else {
				field.Value = fmt.Sprint(v)
			}
			// If the field is <= 20 then we'll set it to short
			if len(field.Value) <= 20 {
				field.Short = true
			}
			attachment.Fields = append(attachment.Fields, field)
		}
	} else {
		attachment.Text = e.Message
	}
	attachment.Fallback = e.Message
	attachment.Color = color

	sort.SliceStable(attachment.Fields, func(i, j int) bool {
		iTitle, jTitle := attachment.Fields[i].Title, attachment.Fields[j].Title
		return iTitle > jTitle
	})

	if sh.Async {
		go sh.postMessage(payload)
		return nil
	}

	return sh.postMessage(payload)
}

func (sh *SlackHook) postMessage(payload *slack.WebHookPostPayload) error {
	if sh.Timeout <= 0 {
		return sh.hook.PostMessage(payload)
	}

	ech := make(chan error, 1)
	go func(ch chan error) {
		ch <- nil
		ch <- sh.hook.PostMessage(payload)
	}(ech)
	<-ech

	select {
	case err := <-ech:
		return err
	case <-time.After(sh.Timeout):
		return TimeoutError
	}
}

// newEntry adds a new entry to the Logger
func (sh *SlackHook) newEntry(entry *logrus.Entry) *logrus.Entry {
	data := map[string]interface{}{}

	for k, v := range sh.Options {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	data["Where"] = entry.Caller
	data["Time"] = entry.Time
	data["Message"] = entry.Message
	return &logrus.Entry{
		Logger:  entry.Logger,
		Time:    entry.Time,
		Level:   entry.Level,
		Data:    data,
		Message: entry.Message,
	}
}

// Levels sets which levels to sent to slack
func (sh *SlackHook) Levels() []logrus.Level {
	if sh.AcceptedLevels == nil {
		return AllLevels
	}
	return sh.AcceptedLevels
}

var LevelColorMap = map[logrus.Level]string{
	logrus.DebugLevel: "#9B30FF",
	logrus.InfoLevel:  "good",
	logrus.WarnLevel:  "warning",
	logrus.ErrorLevel: "danger",
	logrus.FatalLevel: "danger",
	logrus.PanicLevel: "danger",
}

// Supported log levels
var AllLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

var TimeoutError = errors.New("Request timed out")

// LevelThreshold - Returns every logging level above and including the given parameter.
func LevelThreshold(l logrus.Level) []logrus.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []logrus.Level{}
}
