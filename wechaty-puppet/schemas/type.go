package schemas

// EmitStruct receive puppet emit event
type EmitStruct struct {
  EventName PuppetEventName
  Payload   interface{}
}

// EventParams wechaty emit params
type EventParams struct {
  EventName PuppetEventName
  Params    []interface{}
}
