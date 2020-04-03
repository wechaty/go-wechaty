package schemas

type RoomInvitationPayload struct {
  Id           string
  InviterId    string
  Topic        string
  Avatar       string
  Invitation   string
  MemberCount  int
  MemberIdList []string
  Timestamp    int64
  ReceiverId   string
}
