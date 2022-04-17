/**
 * Go Wechaty - https://github.com/wechaty/go-wechaty
 *
 * Authors: Huan LI (李卓桓) <https://github.com/huan>
 *          Chao Fei () <https://github.com/dchaofei>
 *
 * 2020-now @ Copyright Wechaty <https://github.com/wechaty>
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package user

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type MiniProgram struct {
	payload *schemas.MiniProgramPayload
}

func NewMiniProgram(payload *schemas.MiniProgramPayload) *MiniProgram {
	return &MiniProgram{payload: payload}
}

func (mp *MiniProgram) AppID() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Appid
}

func (mp *MiniProgram) Description() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Description
}

func (mp *MiniProgram) PagePath() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.PagePath
}

func (mp *MiniProgram) ThumbUrl() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.ThumbUrl
}

func (mp *MiniProgram) Title() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Title
}

func (mp *MiniProgram) Username() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Username
}

func (mp *MiniProgram) ThumbKey() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.ThumbKey
}

func (mp *MiniProgram) ShareId() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.ShareId
}

func (mp *MiniProgram) IconUrl() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.IconUrl
}

func (mp *MiniProgram) Payload() schemas.MiniProgramPayload {
	return *mp.payload
}

func (mp *MiniProgram) payloadIsNil() bool {
	return mp.payload == nil
}
