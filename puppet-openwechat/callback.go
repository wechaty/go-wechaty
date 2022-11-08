package puppet_openwechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

func (p PuppetOpenWechat) onScan() {
	p.bot.UUIDCallback = func(uuid string) {
		//qrLink := "http://login.weixin.qq.com/l/"+uuid
		log.Trace(openwechat.GetQrcodeUrl(uuid))
		qrLink := "https://login.weixin.qq.com/l/" + uuid
		p.Emit(schemas.PuppetEventNameScan, &schemas.EventScanPayload{
			BaseEventPayload: schemas.BaseEventPayload{},
			Status:           schemas.ScanStatusWaiting,
			QrCode:           qrLink,
		})
	}
	p.bot.ScanCallBack = func(body []byte) {
		p.Emit(schemas.PuppetEventNameScan, &schemas.EventScanPayload{
			BaseEventPayload: schemas.BaseEventPayload{
				Data: string(body),
			},
			Status: schemas.ScanStatusScanned,
			QrCode: "",
		})
	}
}

func (p PuppetOpenWechat) onMsg() {
	p.bot.MessageOnSuccess(func(msg *openwechat.Message) {
		payload, err := rawMsgToPayload(msg)
		if err != nil {
			log.Errorf("PuppetOpenWechat onMsg err: %s", err.Error())
			return
		}
		p.Puppet.AddMessagePayloadToCache(msg.MsgId, payload)
		p.Emit(schemas.PuppetEventNameMessage, &schemas.EventMessagePayload{MessageId: msg.MsgId})
	})
}
