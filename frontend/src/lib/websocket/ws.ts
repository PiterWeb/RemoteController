class WS extends WebSocket {

    private static url = "ws://localhost:8081/ws"

    static #instance: WS | null;
    private constructor() {
        super(WS.url)
    }


    public static get instance(): WS {
        if (!WS.#instance) {
            WS.#instance = new WS();
        }
        
        return WS.#instance;
    }

    public close(code?: number, reason?: string) {
        this.close(code, reason)
        WS.#instance = null;
    }
}

let ws = WS.instance

export default ws