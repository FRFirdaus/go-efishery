package slacklog

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSlackHook(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	hook := &SlackHook{
		HookURL:        os.Getenv("TRAVIS_TEST_SLACK_HOOK_URL"),
		Username:       "Bot-Core",
		IconEmoji:      ":mega:",
		AcceptedLevels: LevelThreshold(logrus.WarnLevel),
		Channel:        "#bottest",
	}
	logrus.AddHook(hook)
	logrus.Debug("logging in Debug Mode")
	logrus.Info("Logging in Info Mode")
	if len(hook.Levels()) != 4 {
		t.Error("Error setting level, level length not less or more than [warning error fatal panic]")
	}
}

func TestSendSlackLog(t *testing.T) {
	hook := &SlackHook{
		HookURL:        os.Getenv("TRAVIS_TEST_SLACK_HOOK_URL"),
		Username:       "Bot-Core",
		IconEmoji:      ":mega:",
		AcceptedLevels: LevelThreshold(logrus.DebugLevel),
		Channel:        "#bottest",
	}
	err := hook.Fire(&logrus.Entry{
		Data: map[string]interface{}{
			"tag": "testing",
		},
		Logger:  &logrus.Logger{},
		Message: "Testing Slacklogrus Log",
	})

	if err != nil {
		t.Errorf("Could not fire slacklogrus: %v", err)
	}
}
