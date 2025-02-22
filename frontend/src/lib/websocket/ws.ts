class WS extends WebSocket {

    private static url = "ws://localhost:8080/ws"

    static #instance: WS | null;
    private constructor() {
        try {super(WS.url)} catch {}
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

const ws = WS.instance

export default ws