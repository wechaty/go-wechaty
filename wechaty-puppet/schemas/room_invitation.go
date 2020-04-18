package schemas

type RoomInvitationPayload struct {
  Id           string   `json:"id"`
  InviterId    string   `json:"inviterId"`
  Topic        string   `json:"topic"`
  Avatar       string   `json:"avatar"`
  Invitation   string   `json:"invitation"`
  MemberCount  int      `json:"memberCount"`
  MemberIdList []string `json:"memberIdList"`
  Timestamp    int64    `json:"timestamp"`
  ReceiverId   string   `json:"receiverId"`
}
