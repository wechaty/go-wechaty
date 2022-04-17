package schemas

import "time"

type RoomInvitationPayload struct {
	Id           string    `json:"id"`
	InviterId    string    `json:"inviterId"`
	Topic        string    `json:"topic"`
	Avatar       string    `json:"avatar"`
	Invitation   string    `json:"invitation"`
	MemberCount  int       `json:"memberCount"`
	MemberIdList []string  `json:"memberIdList"`
	Timestamp    time.Time `json:"timestamp"`
	ReceiverId   string    `json:"receiverId"`
}
