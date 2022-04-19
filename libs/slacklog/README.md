## Example

```go
package main

import (
	"github.com/sirupsen/logrus"
	"bitbucket.org/efishery/go-efishery/libs/slacklog"
)

const (
	// slack webhook url
	hookURL = "https://hooks.slack.com/TXXXXX/BXXXXX/XXXXXXXXXX"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.AddHook(&slacklog.SlackHook{
		HookURL:        hookURL,
		AcceptedLevels: slacklog.LevelThreshold(logrus.WarnLevel),
		Channel:        "#general",
		IconEmoji:      ":ghost:",
		Username:       "slacklog",
		Disabled:		true
		Timeout:        5 * time.Second, // request timeout for calling slack api
	})

	logrus.WithFields(logrus.Fields{"foo": "bar", "foo2": "bar2"}).Warn("this is a warn level message")
	logrus.Debug("this is a debug level message")
	logrus.Info("this is an info level message")
	logrus.Error("this is an error level message")
	if(err != nil){
		logrus.WithError(err).Error("Message and error")
	}
}
```

### Extra fields
You can also add some extra fields to be sent with every slack message
```go
extra := map[string]interface{}{
			"service": "service-1",
			"maintener": "<@asdasd>",
		}
	
logrus.AddHook(&slacklog.SlackHook{
		//HookURL:        "https://hooks.slack.com/services/abc123/defghijklmnopqrstuvwxyz",
		Options: 			extra,
})
```

## Parameters

#### Required
  * HookURL

#### Optional
  * IconEmoji
  * IconURL
  * Username
  * Channel
  * Async
  * Options

## Credits 

This project based on [logrus_slack](https://github.com/bluele/logrus_slack)

## Author

**Erwan Akse**
* <erwan.akse@efishery.com>