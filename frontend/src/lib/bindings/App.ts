import {bindings} from './models';
import {webrtc} from './models';

// @ts-check
type Response<T> = {
  result: T
}

type Request<T> = {
  args: T
}

async function callBinding<U, T>(path: string, method: string = "GET", args?: Request<U>): Promise<T> {
  
  const body = args ? JSON.stringify(args) : undefined

  const baseUrl = location.port === "5173" ? `http://${location.hostname}:8090` : location.origin
  
  const res = await fetch(`${baseUrl}/bindings${path}`, {
    mode: "cors",
    method,
    body: body
  })

  if (!res.ok) throw new Error()

  const contentLength = Number(res.headers.get("Content-Length"))

  
  if (contentLength === 0) {
    const t: T = undefined as T
    return t
  }

  const resJSON = await res.json() as Response<T>
  
  return resJSON.result    
  
}

export function GetCurrentOS(): Promise<string> {
  return callBinding("/os")
}

export function GetUsedPorts(): Promise<bindings.UsedPorts> {
  return callBinding("/ports");
}

export function IsGamepadEnabled(): Promise<boolean> {
  return callBinding("/gamepad");
}

export function IsKeyboardEnabled(): Promise<boolean> {
  return callBinding("/keyboard");
}

export function IsMouseEnabled(): Promise<boolean> {
  return callBinding("/mouse");
}

export function IsWhipEnabled(): Promise<boolean> {
  return callBinding("/whip");
}

export function LogPrintln(info: string): Promise<void> {
  return callBinding("/log", "POST", {args: info});
}

export function NotifyCloseClient(): Promise<void> {
  return callBinding("/notify/peer/close", "PUT")
}

export function NotifyCreateClient(): Promise<void> {
  return callBinding("/notify/peer/open", "PUT");
}

export function OpenViGEmWizard(): Promise<void> {
  return callBinding("/vigem");
}

export function ToogleGamepad(): Promise<void> {
  return callBinding("/gamepad", "PUT")
}

export function ToogleKeyboard(): Promise<void> {
  return callBinding("/keyboard", "PUT");
}

export function ToogleMouse(): Promise<void> {
  return callBinding("/mouse", "PUT");
}

export function ToogleWhip(): Promise<void> {
  return callBinding("/whip", "PUT");
}

export function TryClosePeerConnection(): Promise<void> {
  return callBinding("/conn/host", "DELETE")
}

export function TryCreateHost(ice_servers:Array<webrtc.ICEServer>,offer_encoded:string): Promise<string> {
  return callBinding("/conn/host", "POST", {
    args: {
      ice_servers,
      offer_encoded
    }
  })
}
