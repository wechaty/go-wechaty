/**
 * Go Wechaty - https://github.com/wechaty/go-wechaty
 *
 * Authors: Huan LI (李卓桓) <https://github.com/huan>
 *          Xiaoyu DING （丁小雨） <https://github.com/dingdayu>
 *          Bojie LI (李博杰) <https://github.com/SilkageNet>
 *
 * 2020-now @ Copyright Wechaty
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
package wechaty

import (
	"github.com/wechaty/go-wechaty/src/wechaty/user"
)

// Wechaty
type Wechaty struct {
	Message user.Message
	Image   user.Images
	Room    user.Room
	Contact user.Contact
}

// NewWechaty
// instance by golang.
func NewWechaty() *Wechaty {
	return &Wechaty{}
}

func (w *Wechaty) OnScan(f func(qrCode, status string)) *Wechaty {
	return w
}

// OnLogin
// todo:: fake code. user should be struct
func (w *Wechaty) OnLogin(func(user string)) *Wechaty {
	return w
}

// OnMessage
// todo:: fake code. message should be struct
func (w *Wechaty) OnMessage(func(message string)) *Wechaty {
	return w
}

// Start
func (w *Wechaty) Start() *Wechaty {
	return w
}
