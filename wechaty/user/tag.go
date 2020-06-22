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

import _interface "github.com/wechaty/go-wechaty/wechaty/interface"

type Tag struct {
	_interface.IAccessory
	id string
}

// NewTag ...
func NewTag(id string, accessory _interface.IAccessory) *Tag {
	return &Tag{accessory, id}
}

func (t *Tag) ID() string {
	return t.id
}

// Add tag for contact
func (t *Tag) Add(to _interface.IContact) error {
	return t.GetPuppet().TagContactAdd(t.id, to.ID())
}

// Remove this tag from Contact
func (t *Tag) Remove(from _interface.IContact) error {
	return t.GetPuppet().TagContactRemove(t.id, from.ID())
}
