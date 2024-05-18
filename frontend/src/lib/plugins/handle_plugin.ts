import { connect } from "nats.ws";
import { GetPlugins } from "$lib/wailsjs/go/desktop/App";

export async function getPlugins() {
    return await GetPlugins();
}

// Handle the plugin initialization
function handleInitPlugin(pluginName: string, channel: RTCDataChannel) {
  
}

// Return a function that will handle messages from the plugin
function handleMessagePlugin() {

}