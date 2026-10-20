export enum Events {
  connection_state = "connection_state",
  error = "ERROR", 
  streaming_signal_client = "streaming-signal-client",
  streaming_signal_server = "streaming-signal-server"
}

function createWS(path: Events) {
  const baseUrl = location.port === "5173" ? `ws://${location.hostname}:8090` : location.origin.replace("http", "ws")
  return new WebSocket(`${baseUrl}/ws/${path}`)
}

type Clbk = (...data: any) => void

function createListener(callback: Clbk, onceCancel = () => {}) {
  return ({ data }: { data: any }) => {
    callback(data)
    onceCancel()    
  }
}

const eventsListenersMap = new Map<Events, WebSocket[]>()

function addWSEventToMap(e: Events, ws: WebSocket) {
  eventsListenersMap.set(e, [...eventsListenersMap.get(e) ?? [], ws])
}

function removeWSEventToMap(e: Events, ws: WebSocket) {
  eventsListenersMap.set(e, eventsListenersMap.get(e)?.filter(w => w !== ws) ?? [])
}

export function EventsOn(e: Events, callback: Clbk): () => () => void {  
  try {

    let closedByUser = false
    
    const ws = createWS(e)
    addWSEventToMap(e, ws)
      
    const listener = createListener(callback)

    let unlistener = () => {
      closedByUser = true
      removeWSEventToMap(e, ws)
      ws.removeEventListener("message", listener)
    }
    
    ws.addEventListener("message", listener)
    ws.addEventListener("close", (ev) => {
      if (closedByUser || ev.reason === "EventsOff") return

      unlistener = EventsOn(e, callback)
    })
    
    return () => unlistener
    
  } finally {
    return () => () => {}
  }
}

export function EventsOnce(e: Events, callback: Clbk): () => () => void {
  try {

    let closedByUser = false
    
    const ws = createWS(e)
    addWSEventToMap(e, ws)
    
    let unlistener = () => {
      closedByUser = true
      removeWSEventToMap(e, ws)
      ws.removeEventListener("message", listener)
    }
    
    const listener = createListener(callback, unlistener)
    
    ws.addEventListener("message", listener)
    ws.addEventListener("close", (ev) => {
      if (closedByUser || ev.reason === "EventsOff") return

      unlistener = EventsOn(e, callback)
    })
    
    return () => unlistener
    
  } finally {
    return () => () => {}
  }
}

export function EventsOff(e: Events, ...additionalEventNames: Events[]): void {
  eventsListenersMap.get(e)?.forEach((w) => {
    w.close(200, "EventsOff")
  })

  eventsListenersMap.set(e, [])

  additionalEventNames.forEach((aEvent) => {
    eventsListenersMap.get(aEvent)?.forEach((w) => {
      w.close(200, "EventsOff")
    })

    eventsListenersMap.set(aEvent, [])
  })
}

export function EventsEmit(e: Events, ...data: any): void {
  const ws = createWS(e)

  ws.addEventListener("open", () => {
    ws.send(JSON.stringify(data))
    ws.close()
  })
}