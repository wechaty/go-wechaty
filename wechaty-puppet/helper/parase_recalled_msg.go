package helper

import (
	"encoding/xml"
)

//<sysmsg type="revokemsg">
//<revokemsg>
//<session>wxid_qswi83jdiiw2</session>
//<msgid>12543453</msgid>
//<newmsgid>500043327888834838</newmsgid> // 元消息id
//<replacemsg><![CDATA["xxx" 撤回了一条消息]]></replacemsg>
//</revokemsg>
//</sysmsg>

// RecalledMsg 撤回消息的 xml 结构体
type RecalledMsg struct {
	XMLName   xml.Name `xml:"sysmsg"`
	Text      string   `xml:",chardata"`
	Type      string   `xml:"type,attr"`
	Revokemsg struct {
		Text       string `xml:",chardata"`
		Session    string `xml:"session"`
		Msgid      string `xml:"msgid"`
		Newmsgid   string `xml:"newmsgid"`
		Replacemsg string `xml:"replacemsg"`
	} `xml:"revokemsg"`
}

// ParseRecalledID 从 xml 中解析撤回的原始消息id
func ParseRecalledID(raw string) string {
	msg := &RecalledMsg{}
	err := xml.Unmarshal([]byte(raw), msg)
	if err != nil {
		return raw
	}
	if msg.Revokemsg.Newmsgid != "" {
		return msg.Revokemsg.Newmsgid
	}
	return raw
}
