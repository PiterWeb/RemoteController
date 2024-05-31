import { GetPlugins, GetNATSPort } from "$lib/wailsjs/go/desktop/App";
import PluginWebsocket from "$lib/plugins/websocket";


export async function getPlugins() {
    return await GetPlugins();
}

PluginWebsocket.connect(await GetNATSPort())

// Handle the plugin initialization
function handleInitPlugin(pluginName: string, channel: RTCDataChannel) {
  
    // channel.onmessage = (_ch: RTCDataChannel, ev: MessageEvent) => {
    //     PluginWebsocket.send(pluginName, ev.data);
    // };
    // PluginWebsocket.onMessage(pluginName, (err, msg) => {
    //     if (err) {
    //         console.error(err);
    //         return;
    //     }
    //     channel.send(msg.data);
    // });

}

// Return a function that will handle messages from the plugin
function handleMessagePlugin(this: RTCDataChannel,ev: MessageEvent) {
    
}