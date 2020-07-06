package schemas

import "regexp"

type RoomMemberQueryFilter struct {
	Name         string
	RoomAlias    string
	ContactAlias string
}

type RoomQueryFilter struct {
	Id          string
	Topic       string
	TopicRegexp *regexp.Regexp
}

func (r *RoomQueryFilter) Empty() bool {
	return r.Id == "" && r.Topic == "" && r.TopicRegexp == nil
}

func (r *RoomQueryFilter) All() bool {
	return r.Id != "" && r.Topic != "" && r.TopicRegexp != nil
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

type RoomPayloadFilterFunction func(payload *RoomPayload) bool
