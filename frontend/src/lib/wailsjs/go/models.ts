export namespace plugins {
	
	export class Plugin_arg {
	    name: string;
	    value: any;
	    value_list?: Plugin_arg[];
	
	    static createFrom(source: any = {}) {
	        return new Plugin_arg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.value_list = this.convertValues(source["value_list"], Plugin_arg);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Plugin {
	    name: string;
	    path: string;
	    init_client_args: Plugin_arg[];
	    init_host_args: Plugin_arg[];
	    background_args: Plugin_arg[];
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Plugin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.init_client_args = this.convertValues(source["init_client_args"], Plugin_arg);
	        this.init_host_args = this.convertValues(source["init_host_args"], Plugin_arg);
	        this.background_args = this.convertValues(source["background_args"], Plugin_arg);
	        this.enabled = source["enabled"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace webrtc {
	
	export class ICEServer {
	    urls: string[];
	    username?: string;
	    credential?: any;
	    credentialType?: number;
	
	    static createFrom(source: any = {}) {
	        return new ICEServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.urls = source["urls"];
	        this.username = source["username"];
	        this.credential = source["credential"];
	        this.credentialType = source["credentialType"];
	    }
	}

}

