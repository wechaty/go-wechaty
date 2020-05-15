package user

import (
	"fmt"
	"testing"
)

func TestMessage_MentionList(t *testing.T) {
	u := NewMessage("", nil).(*Message)
	at := u.multipleAt("hello@a@b@c")
	fmt.Println(at)
}
