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
  "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Contact struct {
  _interface.Accessory

  Id      string
  payload schemas.ContactPayload
}

func NewContact(accessory _interface.Accessory, id string) *Contact {
  return &Contact{
    Accessory: accessory,
    Id:        id,
  }
}

func (r *Contact) Ready(forceSync bool) {
  return
}

func (r *Contact) isReady() bool {
  return true
}

func (r *Contact) Sync() {
  r.Ready(true)
}

func (r *Contact) String() string {
  return fmt.Sprintf("Contact<%s>", r.identity())
}

func (r *Contact) identity() string {
  identity := "loading..."
  if r.payload.Alias != "" {
    identity = r.payload.Alias
  } else if r.payload.Name != "" {
    identity = r.payload.Name
  } else if r.Id != "" {
    identity = r.Id
  }
  return identity
}
