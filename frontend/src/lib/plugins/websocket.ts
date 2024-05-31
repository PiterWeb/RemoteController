import { connect, Msg, NatsConnection, NatsError } from 'nats.ws';

export default class PluginWebsocket {
	static #port: number;
	static #connection?: NatsConnection;
	static #tries = 0;

	static async connect(port: number) {

        if (this.#connection && !this.#connection.isClosed()) {
            return;
        }

		this.#port = port;

		try {
			this.#connection = await connect({ servers: [`wss://127.0.0.1:${port}`] });
			this.#tries = 0;
		} catch (e) {
            this.#tries++;
			this.reconnect();
		}
	}

	static async disconnect() {
		if (!this.#connection) {
			return;
		}
		await this.#connection.close();
		this.#connection = undefined;
	}

	static async reconnect() {
		if (this.#connection) {
			await this.disconnect();
		}

		if (this.#tries <= 5) {
            await this.connect(this.#port);
		}
	}

	static send(channel: string, message: string) {
		if (!this.#connection || this.#connection.isClosed()) {
			return;
		}
		this.#connection.publish(channel, message);
	}

    static onMessage(channel: string, callback: (err: NatsError | null, msg: Msg) => void) {
        if (!this.#connection || this.#connection.isClosed()) {
            return;
        }
        const subscribe = this.#connection.subscribe(channel);

        subscribe.callback = callback;
    }

    static onMessageOnce(channel: string, callback: (err: NatsError | null, msg: Msg) => void) {
        if (!this.#connection || this.#connection.isClosed()) {
            return;
        }

        const subscribe = this.#connection.subscribe(channel);

        subscribe.callback = (err, msg) => {
            callback(err, msg);
            subscribe.unsubscribe(msg.sid);
        };
    }
}
