
export interface ServersConfig {
	[group: string]: ICEServer;
}

export interface ICEServer {
	urls: string[];
	username?: string;
	credential?: string;
}