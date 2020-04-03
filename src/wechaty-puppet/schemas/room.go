package schemas

type RoomMemberQueryFilter struct {
  Name         string
  RoomAlias    string
  ContactAlias string
}

type RoomQueryFilter struct {
  Id    string
  Topic string
}

type RoomPayload struct {
  Id           string
  Topic        string
  Avatar       string
  MemberIdList []string
  OwnerId      string
  AdminIdList  []string
}

type RoomMemberPayload struct {
  Id        string
  RoomAlias string
  InviterId string
  Avatar    string
  Name      string
}
