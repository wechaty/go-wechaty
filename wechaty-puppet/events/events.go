// Package events provides simple EventEmitter support for Go Programming Language
package events

import (
  "log"
  "reflect"
  "sync"
)

const (
  // DefaultMaxListeners is the number of max listeners per event
  // default EventEmitters will print a warning if more than x listeners are
  // added to it. This is a useful default which helps finding memory leaks.
  // Defaults to 0, which means unlimited
  DefaultMaxListeners = 0

  // EnableWarning prints a warning when trying to add an event which it's len is equal to the maxListeners
  // Defaults to false, which means it does not prints a warning
  EnableWarning = false
)

type (
  // EventName is just a type of string, it's the event name
  EventName interface{}
  // Listener is the type of a Listener, it's a func which receives any,optional, arguments from the caller/emitter
  Listener func(...interface{})
  // Events the type for registered listeners, it's just a map[string][]func(...interface{})
  Events map[EventName][]Listener

  // EventEmitter is the message/or/event manager
  EventEmitter interface {
    // AddListener is an alias for .On(eventName, listener).
    AddListener(EventName, ...Listener)
    // Emit fires a particular event,
    // Synchronously calls each of the listeners registered for the event named
    // eventName, in the order they were registered,
    // passing the supplied arguments to each.
    Emit(EventName, ...interface{})
    // EventNames returns an array listing the events for which the emitter has registered listeners.
    // The values in the array will be strings.
    EventNames() []EventName
    // GetMaxListeners returns the max listeners for this emitter
    // see SetMaxListeners
    GetMaxListeners() int
    // ListenerCount returns the length of all registered listeners to a particular event
    ListenerCount(EventName) int
    // Listeners returns a copy of the array of listeners for the event named eventName.
    Listeners(EventName) []Listener
    // On registers a particular listener for an event, func receiver parameter(s) is/are optional
    On(EventName, ...Listener)
    // Once adds a one time listener function for the event named eventName.
    // The next time eventName is triggered, this listener is removed and then invoked.
    Once(EventName, ...Listener)
    // RemoveAllListeners removes all listeners, or those of the specified eventName.
    // Note that it will remove the event itself.
    // Returns an indicator if event and listeners were found before the remove.
    RemoveAllListeners(EventName) bool
    // RemoveListener removes given listener from the event named eventName.
    // Returns an indicator whether listener was removed
    RemoveListener(EventName, Listener) bool
    // Clear removes all events and all listeners, restores Events to an empty value
    Clear()
    // SetMaxListeners obviously this function allows the MaxListeners
    // to be decrease or increase. Set to zero for unlimited
    SetMaxListeners(int)
    // Len returns the length of all registered events
    Len() int
  }

  emitter struct {
    maxListeners int
    evtListeners Events
    mu           sync.Mutex
  }
)

// CopyTo copies the event listeners to an EventEmitter
func (e Events) CopyTo(emitter EventEmitter) {
  if len(e) > 0 {
    // register the events to/with their listeners
    for evt, listeners := range e {
      if len(listeners) > 0 {
        emitter.AddListener(evt, listeners...)
      }
    }
  }
}

// New returns a new, empty, EventEmitter
func New() EventEmitter {
  return &emitter{maxListeners: DefaultMaxListeners, evtListeners: Events{}}
}

var (
  _              EventEmitter = &emitter{}
  defaultEmitter              = New()
)

// AddListener is an alias for .On(eventName, listener).
func AddListener(evt EventName, listener ...Listener) {
  defaultEmitter.AddListener(evt, listener...)
}

func (e *emitter) AddListener(evt EventName, listener ...Listener) {
  if len(listener) == 0 {
    return
  }

  e.mu.Lock()
  defer e.mu.Unlock()

  if e.evtListeners == nil {
    e.evtListeners = Events{}
  }

  listeners := e.evtListeners[evt]

  if e.maxListeners > 0 && len(listeners) == e.maxListeners {
    if EnableWarning {
      log.Printf(`(events) warning: possible EventEmitter memory '
                    leak detected. %d listeners added. '
                    Use emitter.SetMaxListeners(n int) to increase limit.`, len(listeners))
    }
    return
  }

  if listeners == nil {
    listeners = make([]Listener, e.maxListeners)
  }

  e.evtListeners[evt] = append(listeners, listener...)
}

// Emit fires a particular event,
// Synchronously calls each of the listeners registered for the event named
// eventName, in the order they were registered,
// passing the supplied arguments to each.
func Emit(evt EventName, data ...interface{}) {
  defaultEmitter.Emit(evt, data...)
}

func (e *emitter) Emit(evt EventName, data ...interface{}) {
  if e.evtListeners == nil {
    return // has no listeners to emit/speak yet
  }
  if listeners := e.evtListeners[evt]; len(listeners) > 0 { // len() should be just fine, but for any case on future...
    for i := range listeners {
      l := listeners[i]
      if l != nil {
        l(data...)
      }
    }
  }
}

// EventNames returns an array listing the events for which the emitter has registered listeners.
// The values in the array will be strings.
func EventNames() []EventName {
  return defaultEmitter.EventNames()
}

func (e *emitter) EventNames() []EventName {
  if e.evtListeners == nil || e.Len() == 0 {
    return nil
  }

  names := make([]EventName, e.Len())
  i := 0
  for k := range e.evtListeners {
    names[i] = k
    i++
  }
  return names
}

// GetMaxListeners returns the max listeners for this emitter
// see SetMaxListeners
func GetMaxListeners() int {
  return defaultEmitter.GetMaxListeners()
}

func (e *emitter) GetMaxListeners() int {
  return e.maxListeners
}

// ListenerCount returns the length of all registered listeners to a particular event
func ListenerCount(evt EventName) int {
  return defaultEmitter.ListenerCount(evt)
}

func (e *emitter) ListenerCount(evt EventName) int {
  if e.evtListeners == nil {
    return 0
  }
  length := 0

  if evtListeners := e.evtListeners[evt]; evtListeners != nil { // length() should be just fine, but for any case on future...
    for _, l := range evtListeners {
      if l == nil {
        continue
      }
      length++
    }
  }

  return length
}

// Listeners returns a copy of the array of listeners for the event named eventName.
func Listeners(evt EventName) []Listener {
  return defaultEmitter.Listeners(evt)
}

func (e *emitter) Listeners(evt EventName) []Listener {
  if e.evtListeners == nil {
    return nil
  }
  var listeners []Listener
  if evtListeners := e.evtListeners[evt]; evtListeners != nil {
    // do not pass any inactive/removed listeners(nil)
    for _, l := range evtListeners {
      if l == nil {
        continue
      }

      listeners = append(listeners, l)
    }

    if len(listeners) > 0 {
      return listeners
    }
  }

  return nil
}

// On registers a particular listener for an event, func receiver parameter(s) is/are optional
func On(evt EventName, listener ...Listener) {
  defaultEmitter.On(evt, listener...)
}

func (e *emitter) On(evt EventName, listener ...Listener) {
  e.AddListener(evt, listener...)
}

// Once adds a one time listener function for the event named eventName.
// The next time eventName is triggered, this listener is removed and then invoked.
func Once(evt EventName, listener ...Listener) {
  defaultEmitter.Once(evt, listener...)
}

func (e *emitter) Once(evt EventName, listener ...Listener) {
  if len(listener) == 0 {
    return
  }

  var modifiedListeners []Listener

  if e.evtListeners == nil {
    e.evtListeners = Events{}
  }

  for i, l := range listener {

    idx := len(e.evtListeners) + i // get the next index (where this event should be added) and adds the i for the 'capacity'

    func(listener Listener, index int) {
      fired := false
      // remove the specific listener from the listeners before fire the real listener
      modifiedListeners = append(modifiedListeners, func(data ...interface{}) {
        if e.evtListeners == nil {
          return
        }
        if !fired {
          // make sure that we don't get a panic(index out of array or nil map here
          if e.evtListeners[evt] != nil && (len(e.evtListeners[evt]) > index || index == 0) {

            e.mu.Lock()
            //e.evtListeners[evt] = append(e.evtListeners[evt][:index], e.evtListeners[evt][index+1:]...)
            // we do not must touch the order because of the pre-defined indexes, we need just to make this listener nil in order to be not executed,
            // and make the len of listeners increase when listener is not nil, not just the len of listeners.
            // so set this listener to nil
            e.evtListeners[evt][index] = nil
            e.mu.Unlock()
          }
          fired = true
          listener(data...)
        }

      })
    }(l, idx)

  }
  e.AddListener(evt, modifiedListeners...)
}

// RemoveAllListeners removes all listeners, or those of the specified eventName.
// Note that it will remove the event itself.
// Returns an indicator if event and listeners were found before the remove.
func RemoveAllListeners(evt EventName) bool {
  return defaultEmitter.RemoveAllListeners(evt)
}

func (e *emitter) RemoveAllListeners(evt EventName) bool {
  if e.evtListeners == nil {
    return false // has nothing to remove
  }
  e.mu.Lock()
  defer e.mu.Unlock()
  if listeners := e.evtListeners[evt]; listeners != nil {
    l := e.ListenerCount(evt) // in order to not get the len of any inactive/removed listeners
    delete(e.evtListeners, evt)
    if l > 0 {
      return true
    }
  }

  return false
}

// RemoveListener removes the specified listener from the listener array for the event named eventName.
func (e *emitter) RemoveListener(evt EventName, listener Listener) bool {
  if e.evtListeners == nil {
    return false
  }

  if listener == nil {
    return false
  }

  e.mu.Lock()
  defer e.mu.Unlock()

  listeners := e.evtListeners[evt]

  if listeners == nil {
    return false
  }

  idx := -1
  listenerPointer := reflect.ValueOf(listener).Pointer()

  for index, item := range listeners {
    itemPointer := reflect.ValueOf(item).Pointer()
    if itemPointer == listenerPointer {
      idx = index
      break
    }
  }

  if idx < 0 {
    return false
  }

  var modifiedListeners []Listener = nil

  if len(listeners) > 1 {
    modifiedListeners = append(listeners[:idx], listeners[idx+1:]...)
  }

  e.evtListeners[evt] = modifiedListeners

  return true
}

// Clear removes all events and all listeners, restores Events to an empty value
func Clear() {
  defaultEmitter.Clear()
}

func (e *emitter) Clear() {
  e.evtListeners = Events{}
}

// SetMaxListeners obviously this function allows the MaxListeners
// to be decrease or increase. Set to zero for unlimited
func SetMaxListeners(n int) {
  defaultEmitter.SetMaxListeners(n)
}

func (e *emitter) SetMaxListeners(n int) {
  if n < 0 {
    if EnableWarning {
      log.Printf("(events) warning: MaxListeners must be positive number, tried to set: %d", n)
      return
    }
  }
  e.maxListeners = n
}

// Len returns the length of all registered events
func Len() int {
  return defaultEmitter.Len()
}

func (e *emitter) Len() int {
  if e.evtListeners == nil {
    return 0
  }
  return len(e.evtListeners)
}
