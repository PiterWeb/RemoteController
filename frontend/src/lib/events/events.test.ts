import { Events, EventsEmit, EventsOn, EventsOnce } from "$lib/events/events"
import { test } from "vitest"

test('event listen, emit and then cancel', (t) => {

  let i = 0
  
  let unlisten = EventsOn(Events.error, () => {
    i++
  })

  // Handled events
  EventsEmit(Events.error)
  EventsEmit(Events.error)

  setTimeout(() => {
    unlisten()()

    // Unhandled events
    EventsEmit(Events.error)
    EventsEmit(Events.error)
    EventsEmit(Events.error)
    EventsEmit(Events.error)
    
    setTimeout(() => {
      t.expect(i).toBe(2)
    }, 1000)
    
  }, 1000)
  
})

test('once listener cancel before emit event', (t) => {
  let i = 0

  let unlistener = EventsOnce(Events.error, () => i++)

  unlistener()()
  EventsEmit(Events.error)
  setTimeout(() => {
    t.expect(i).toBe(0)
  }, 1000)
})

test('multiple listener, one emit', (t) => {

  let i = 0;

  EventsOnce(Events.error, () => i++)
  EventsOnce(Events.error, () => i++)

  EventsEmit(Events.error)

  setTimeout(() => {
    t.expect(i).toBe(2)
  }, 1000)
  
})