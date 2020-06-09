/**
 * Go Wechaty - https://github.com/wechaty/go-wechaty
 *
 * Authors: Huan LI (李卓桓) <https://github.com/huan>
 *          Bojie LI (李博杰) <https://github.com/SilkageNet>
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

import (
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type UrlLink struct {
	payload *schemas.UrlLinkPayload
}

func NewUrlLink(payload *schemas.UrlLinkPayload) *UrlLink {
	return &UrlLink{payload: payload}
}

func (ul *UrlLink) String() string {
	return fmt.Sprintf("UrlLink<%s>", ul.Url())
}

func (ul *UrlLink) Url() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Url
}

func (ul *UrlLink) Title() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Title
}

func (ul *UrlLink) ThumbnailUrl() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.ThumbnailUrl
}

func (ul *UrlLink) Description() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Description
}

// Payload UrlLink payload
func (ul *UrlLink) Payload() schemas.UrlLinkPayload {
	return *ul.payload
}
