package pubsub

import (
	"fmt"
	"testing"
)

func TestUnixMillis(t *testing.T) {
	dataTest := int64(1647477693)
	tsParsed := parseUnixMilliSecond(dataTest)
	ts := toUnixMilliSecond(tsParsed)

	if dataTest != ts {
		fmt.Println(dataTest, ts)
		t.Error("Invalid Unix timestamp")
		return
	}
}
