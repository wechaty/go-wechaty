package payload

type RoomMemberResponse struct {
	WxID   string `json:"wx_id"`
	Status int    `json:"status"`
}

// CreateRoomResponse create room response
type CreateRoomResponse struct {
	Status        int                  `json:"status"`
	RoomId        string               `json:"roomId"`
	Message       string               `json:"message"`
	CreateMessage []RoomMemberResponse `json:"createMessage"`
}
