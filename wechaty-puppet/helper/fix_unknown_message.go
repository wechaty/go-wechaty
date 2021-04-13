package helper

import (
	"encoding/xml"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"log"
)

// FixUnknownMessage 修复wechaty服务端不能解析的消息，尝试服务端去解析
func FixUnknownMessage(payload *schemas.MessagePayload) {
	if payload.Type != schemas.MessageTypeUnknown {
		return
	}
	msg := &Msg{}
	err := xml.Unmarshal([]byte(payload.Text), msg)
	if err != nil {
		log.Printf("FixUnknownMessage raw:%s || err: %s", payload.Text, err)
		return
	}
	if msg.Appmsg.Type.Text == "36" {
		payload.Type = schemas.MessageTypeMiniProgram
		payload.FixMiniApp = true
	}
}

// ParseMiniApp 解析小程序xml
func ParseMiniApp(payload *schemas.MessagePayload) (*schemas.MiniProgramPayload, error) {
	msg := &Msg{}
	err := xml.Unmarshal([]byte(payload.Text), msg)
	if err != nil {
		return nil, fmt.Errorf("ParseMiniApp raw:%s || err: %s", payload.Text, err)
	}
	if msg.Appmsg.Type.Text != "36" {
		return nil, fmt.Errorf("ParseMiniApp 不是小程序消息 xml type: %s", msg.Appmsg.Type.Text)
	}
	return &schemas.MiniProgramPayload{
		Appid:       msg.Appmsg.Weappinfo.Appid.Text,
		Description: "",
		PagePath:    msg.Appmsg.Weappinfo.Pagepath.Text,
		ThumbUrl:    msg.Appmsg.Appattach.Cdnthumburl.Text,
		Title:       msg.Appmsg.Title.Text,
		Username:    msg.Appmsg.Weappinfo.Username.Text,
		ThumbKey:    msg.Appmsg.Appattach.Cdnthumbaeskey.Text,
	}, nil
}

// Msg 小程序xml消息体
type Msg struct {
	XMLName xml.Name `xml:"msg"`
	Text    string   `xml:",chardata"`
	Appmsg  struct {
		Text   string `xml:",chardata"`
		Appid  string `xml:"appid,attr"`
		Sdkver string `xml:"sdkver,attr"`
		Title  struct {
			Text string `xml:",chardata"`
		} `xml:"title"`
		Des struct {
			Text string `xml:",chardata"`
		} `xml:"des"`
		Username struct {
			Text string `xml:",chardata"`
		} `xml:"username"`
		Action struct {
			Text string `xml:",chardata"`
		} `xml:"action"`
		Type struct {
			Text string `xml:",chardata"`
		} `xml:"type"`
		Showtype struct {
			Text string `xml:",chardata"`
		} `xml:"showtype"`
		Content struct {
			Text string `xml:",chardata"`
		} `xml:"content"`
		URL struct {
			Text string `xml:",chardata"`
		} `xml:"url"`
		Lowurl struct {
			Text string `xml:",chardata"`
		} `xml:"lowurl"`
		Dataurl struct {
			Text string `xml:",chardata"`
		} `xml:"dataurl"`
		Lowdataurl struct {
			Text string `xml:",chardata"`
		} `xml:"lowdataurl"`
		Contentattr struct {
			Text string `xml:",chardata"`
		} `xml:"contentattr"`
		Streamvideo struct {
			Text           string `xml:",chardata"`
			Streamvideourl struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideourl"`
			Streamvideototaltime struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideototaltime"`
			Streamvideotitle struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideotitle"`
			Streamvideowording struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideowording"`
			Streamvideoweburl struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideoweburl"`
			Streamvideothumburl struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideothumburl"`
			Streamvideoaduxinfo struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideoaduxinfo"`
			Streamvideopublishid struct {
				Text string `xml:",chardata"`
			} `xml:"streamvideopublishid"`
		} `xml:"streamvideo"`
		CanvasPageItem struct {
			Text          string `xml:",chardata"`
			CanvasPageXml struct {
				Text string `xml:",chardata"`
			} `xml:"canvasPageXml"`
		} `xml:"canvasPageItem"`
		Appattach struct {
			Text     string `xml:",chardata"`
			Attachid struct {
				Text string `xml:",chardata"`
			} `xml:"attachid"`
			Cdnthumburl struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumburl"`
			Cdnthumbmd5 struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumbmd5"`
			Cdnthumblength struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumblength"`
			Cdnthumbheight struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumbheight"`
			Cdnthumbwidth struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumbwidth"`
			Cdnthumbaeskey struct {
				Text string `xml:",chardata"`
			} `xml:"cdnthumbaeskey"`
			Aeskey struct {
				Text string `xml:",chardata"`
			} `xml:"aeskey"`
			Encryver struct {
				Text string `xml:",chardata"`
			} `xml:"encryver"`
			Fileext struct {
				Text string `xml:",chardata"`
			} `xml:"fileext"`
			Islargefilemsg struct {
				Text string `xml:",chardata"`
			} `xml:"islargefilemsg"`
		} `xml:"appattach"`
		Extinfo struct {
			Text string `xml:",chardata"`
		} `xml:"extinfo"`
		Androidsource struct {
			Text string `xml:",chardata"`
		} `xml:"androidsource"`
		Thumburl struct {
			Text string `xml:",chardata"`
		} `xml:"thumburl"`
		Mediatagname struct {
			Text string `xml:",chardata"`
		} `xml:"mediatagname"`
		Messageaction struct {
			Text string `xml:",chardata"`
		} `xml:"messageaction"`
		Messageext struct {
			Text string `xml:",chardata"`
		} `xml:"messageext"`
		Emoticongift struct {
			Text        string `xml:",chardata"`
			Packageflag struct {
				Text string `xml:",chardata"`
			} `xml:"packageflag"`
			Packageid struct {
				Text string `xml:",chardata"`
			} `xml:"packageid"`
		} `xml:"emoticongift"`
		Emoticonshared struct {
			Text        string `xml:",chardata"`
			Packageflag struct {
				Text string `xml:",chardata"`
			} `xml:"packageflag"`
			Packageid struct {
				Text string `xml:",chardata"`
			} `xml:"packageid"`
		} `xml:"emoticonshared"`
		Designershared struct {
			Text        string `xml:",chardata"`
			Designeruin struct {
				Text string `xml:",chardata"`
			} `xml:"designeruin"`
			Designername struct {
				Text string `xml:",chardata"`
			} `xml:"designername"`
			Designerrediretcturl struct {
				Text string `xml:",chardata"`
			} `xml:"designerrediretcturl"`
		} `xml:"designershared"`
		Emotionpageshared struct {
			Text string `xml:",chardata"`
			Tid  struct {
				Text string `xml:",chardata"`
			} `xml:"tid"`
			Title struct {
				Text string `xml:",chardata"`
			} `xml:"title"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			IconUrl struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl"`
			SecondUrl struct {
				Text string `xml:",chardata"`
			} `xml:"secondUrl"`
			PageType struct {
				Text string `xml:",chardata"`
			} `xml:"pageType"`
		} `xml:"emotionpageshared"`
		Webviewshared struct {
			Text             string `xml:",chardata"`
			ShareUrlOriginal struct {
				Text string `xml:",chardata"`
			} `xml:"shareUrlOriginal"`
			ShareUrlOpen struct {
				Text string `xml:",chardata"`
			} `xml:"shareUrlOpen"`
			JsAppId struct {
				Text string `xml:",chardata"`
			} `xml:"jsAppId"`
			PublisherId struct {
				Text string `xml:",chardata"`
			} `xml:"publisherId"`
		} `xml:"webviewshared"`
		TemplateID struct {
			Text string `xml:",chardata"`
		} `xml:"template_id"`
		Md5 struct {
			Text string `xml:",chardata"`
		} `xml:"md5"`
		Weappinfo struct {
			Text     string `xml:",chardata"`
			Pagepath struct {
				Text string `xml:",chardata"`
			} `xml:"pagepath"`
			Username struct {
				Text string `xml:",chardata"`
			} `xml:"username"`
			Appid struct {
				Text string `xml:",chardata"`
			} `xml:"appid"`
			Version struct {
				Text string `xml:",chardata"`
			} `xml:"version"`
			Weappiconurl struct {
				Text string `xml:",chardata"`
			} `xml:"weappiconurl"`
			Appservicetype struct {
				Text string `xml:",chardata"`
			} `xml:"appservicetype"`
			Videopageinfo struct {
				Text       string `xml:",chardata"`
				Thumbwidth struct {
					Text string `xml:",chardata"`
				} `xml:"thumbwidth"`
				Thumbheight struct {
					Text string `xml:",chardata"`
				} `xml:"thumbheight"`
				Fromopensdk struct {
					Text string `xml:",chardata"`
				} `xml:"fromopensdk"`
			} `xml:"videopageinfo"`
		} `xml:"weappinfo"`
		Statextstr struct {
			Text string `xml:",chardata"`
		} `xml:"statextstr"`
		FinderFeed struct {
			Text     string `xml:",chardata"`
			ObjectId struct {
				Text string `xml:",chardata"`
			} `xml:"objectId"`
			ObjectNonceId struct {
				Text string `xml:",chardata"`
			} `xml:"objectNonceId"`
			FeedType struct {
				Text string `xml:",chardata"`
			} `xml:"feedType"`
			Nickname struct {
				Text string `xml:",chardata"`
			} `xml:"nickname"`
			Username struct {
				Text string `xml:",chardata"`
			} `xml:"username"`
			Avatar struct {
				Text string `xml:",chardata"`
			} `xml:"avatar"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			MediaCount struct {
				Text string `xml:",chardata"`
			} `xml:"mediaCount"`
			LocalId struct {
				Text string `xml:",chardata"`
			} `xml:"localId"`
			MediaList struct {
				Text string `xml:",chardata"`
			} `xml:"mediaList"`
			MegaVideo struct {
				Text     string `xml:",chardata"`
				ObjectId struct {
					Text string `xml:",chardata"`
				} `xml:"objectId"`
				ObjectNonceId struct {
					Text string `xml:",chardata"`
				} `xml:"objectNonceId"`
			} `xml:"megaVideo"`
		} `xml:"finderFeed"`
		FinderLive struct {
			Text         string `xml:",chardata"`
			FinderLiveID struct {
				Text string `xml:",chardata"`
			} `xml:"finderLiveID"`
			FinderUsername struct {
				Text string `xml:",chardata"`
			} `xml:"finderUsername"`
			FinderObjectID struct {
				Text string `xml:",chardata"`
			} `xml:"finderObjectID"`
			Nickname struct {
				Text string `xml:",chardata"`
			} `xml:"nickname"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			FinderNonceID struct {
				Text string `xml:",chardata"`
			} `xml:"finderNonceID"`
			HeadUrl struct {
				Text string `xml:",chardata"`
			} `xml:"headUrl"`
			Media struct {
				Text     string `xml:",chardata"`
				CoverUrl struct {
					Text string `xml:",chardata"`
				} `xml:"coverUrl"`
				Height struct {
					Text string `xml:",chardata"`
				} `xml:"height"`
				Width struct {
					Text string `xml:",chardata"`
				} `xml:"width"`
			} `xml:"media"`
			LiveStatus struct {
				Text string `xml:",chardata"`
			} `xml:"liveStatus"`
		} `xml:"finderLive"`
		FinderMegaVideo struct {
			Text     string `xml:",chardata"`
			ObjectId struct {
				Text string `xml:",chardata"`
			} `xml:"objectId"`
			ObjectNonceId struct {
				Text string `xml:",chardata"`
			} `xml:"objectNonceId"`
			Nickname struct {
				Text string `xml:",chardata"`
			} `xml:"nickname"`
			Avatar struct {
				Text string `xml:",chardata"`
			} `xml:"avatar"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			Username struct {
				Text string `xml:",chardata"`
			} `xml:"username"`
			MediaCount struct {
				Text string `xml:",chardata"`
			} `xml:"mediaCount"`
			MediaList struct {
				Text string `xml:",chardata"`
			} `xml:"mediaList"`
			FinderFeed struct {
				Text     string `xml:",chardata"`
				ObjectId struct {
					Text string `xml:",chardata"`
				} `xml:"objectId"`
				ObjectNonceId struct {
					Text string `xml:",chardata"`
				} `xml:"objectNonceId"`
			} `xml:"finderFeed"`
		} `xml:"finderMegaVideo"`
		Findernamecard struct {
			Text     string `xml:",chardata"`
			Username struct {
				Text string `xml:",chardata"`
			} `xml:"username"`
			Avatar struct {
				Text string `xml:",chardata"`
			} `xml:"avatar"`
			Nickname struct {
				Text string `xml:",chardata"`
			} `xml:"nickname"`
			AuthJob struct {
				Text string `xml:",chardata"`
			} `xml:"auth_job"`
			AuthIcon struct {
				Text string `xml:",chardata"`
			} `xml:"auth_icon"`
			AuthIconURL struct {
				Text string `xml:",chardata"`
			} `xml:"auth_icon_url"`
		} `xml:"findernamecard"`
		FinderTopic struct {
			Text  string `xml:",chardata"`
			Topic struct {
				Text string `xml:",chardata"`
			} `xml:"topic"`
			TopicType struct {
				Text string `xml:",chardata"`
			} `xml:"topicType"`
			IconUrl struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			PatMusicId struct {
				Text string `xml:",chardata"`
			} `xml:"patMusicId"`
			Location struct {
				Text          string `xml:",chardata"`
				PoiClassifyId struct {
					Text string `xml:",chardata"`
				} `xml:"poiClassifyId"`
				Longitude struct {
					Text string `xml:",chardata"`
				} `xml:"longitude"`
				Latitude struct {
					Text string `xml:",chardata"`
				} `xml:"latitude"`
			} `xml:"location"`
		} `xml:"finderTopic"`
		FinderColumn struct {
			Text   string `xml:",chardata"`
			CardId struct {
				Text string `xml:",chardata"`
			} `xml:"cardId"`
			Title struct {
				Text string `xml:",chardata"`
			} `xml:"title"`
			SecondTitle struct {
				Text string `xml:",chardata"`
			} `xml:"secondTitle"`
			IconUrl struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl"`
			IconUrl1 struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl_1"`
			IconUrl2 struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl_2"`
			IconUrl3 struct {
				Text string `xml:",chardata"`
			} `xml:"iconUrl_3"`
			Cardbuffer struct {
				Text string `xml:",chardata"`
			} `xml:"cardbuffer"`
		} `xml:"finderColumn"`
		FinderEndorsement struct {
			Text  string `xml:",chardata"`
			Scene struct {
				Text string `xml:",chardata"`
			} `xml:"scene"`
		} `xml:"finderEndorsement"`
		Directshare struct {
			Text string `xml:",chardata"`
		} `xml:"directshare"`
		Gamecenter struct {
			Text     string `xml:",chardata"`
			Namecard struct {
				Text    string `xml:",chardata"`
				IconUrl struct {
					Text string `xml:",chardata"`
				} `xml:"iconUrl"`
				Name struct {
					Text string `xml:",chardata"`
				} `xml:"name"`
				Desc struct {
					Text string `xml:",chardata"`
				} `xml:"desc"`
				Tail struct {
					Text string `xml:",chardata"`
				} `xml:"tail"`
				JumpUrl struct {
					Text string `xml:",chardata"`
				} `xml:"jumpUrl"`
			} `xml:"namecard"`
		} `xml:"gamecenter"`
		PatMsg struct {
			Text     string `xml:",chardata"`
			ChatUser struct {
				Text string `xml:",chardata"`
			} `xml:"chatUser"`
			Records struct {
				Text      string `xml:",chardata"`
				RecordNum struct {
					Text string `xml:",chardata"`
				} `xml:"recordNum"`
			} `xml:"records"`
		} `xml:"patMsg"`
		FinderLiveInvite struct {
			Text         string `xml:",chardata"`
			FinderLiveID struct {
				Text string `xml:",chardata"`
			} `xml:"finderLiveID"`
			FinderUsername struct {
				Text string `xml:",chardata"`
			} `xml:"finderUsername"`
			FinderObjectID struct {
				Text string `xml:",chardata"`
			} `xml:"finderObjectID"`
			Nickname struct {
				Text string `xml:",chardata"`
			} `xml:"nickname"`
			Desc struct {
				Text string `xml:",chardata"`
			} `xml:"desc"`
			FinderNonceID struct {
				Text string `xml:",chardata"`
			} `xml:"finderNonceID"`
			HeadUrl struct {
				Text string `xml:",chardata"`
			} `xml:"headUrl"`
			CoverUrl struct {
				Text string `xml:",chardata"`
			} `xml:"coverUrl"`
			LiveMicId struct {
				Text string `xml:",chardata"`
			} `xml:"liveMicId"`
			LiveMicSdkUserId struct {
				Text string `xml:",chardata"`
			} `xml:"liveMicSdkUserId"`
		} `xml:"finderLiveInvite"`
		Websearch struct {
			Text string `xml:",chardata"`
		} `xml:"websearch"`
	} `xml:"appmsg"`
	Fromusername struct {
		Text string `xml:",chardata"`
	} `xml:"fromusername"`
	Scene struct {
		Text string `xml:",chardata"`
	} `xml:"scene"`
	Appinfo struct {
		Text    string `xml:",chardata"`
		Version struct {
			Text string `xml:",chardata"`
		} `xml:"version"`
		Appname struct {
			Text string `xml:",chardata"`
		} `xml:"appname"`
	} `xml:"appinfo"`
	Commenturl struct {
		Text string `xml:",chardata"`
	} `xml:"commenturl"`
}
