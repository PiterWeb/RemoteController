import { toogleLoading } from '$lib/loading/loading_hook';
import log from '$lib/logger/logger';
import { showToast, ToastType } from '$lib/toast/toast_hook';
import { ConnectToHostWeb, CreateClientWeb } from '$lib/webrtc/client_webrtc_hook';
import { SendClientCode, server, ReceiveClientCode, SendHostCode } from '@libreremoteplay/signals';
import { _ } from 'svelte-i18n';
import { get } from 'svelte/store';

export const DEFAULT_EASY_CONNECT_SERVER_IP_DOMAIN = 'libreremoteplay-easyconnect-server.local:8091';
export const DEFAULT_EASY_CONNECT_ID = 1;

export const easyConnectServerIpDomain = $state({ value: DEFAULT_EASY_CONNECT_SERVER_IP_DOMAIN });
export const easyConnectID = $state({ value: DEFAULT_EASY_CONNECT_ID });

export async function handleEasyConnectClient() {
	const easyConnectServerUrl = (() => {
		if (easyConnectServerIpDomain.value.length < 1) {
			return '';
		}
		return `ws://${easyConnectServerIpDomain.value}/ws`;
	})();

	log(`Easy Connect Server URL: ${easyConnectServerUrl}`);

	if (!URL.canParse(easyConnectServerUrl)) {
		throw new Error('Easy Connect Server URL is not set');
	}

	const serverInstance = server(easyConnectServerUrl);

	const clientCode = await CreateClientWeb({ easyConnect: true });

	log(`Client Code: ${clientCode}`);

	if (!clientCode) {
		throw new Error('Failed to create client code');
	}

	toogleLoading();

	try {
		const hostCode = await SendClientCode(
			serverInstance,
			{ data: clientCode },
			easyConnectID.value
		);
		
		log(`Host Code: ${hostCode.data}`);

		const { data: hostCodeStr } = hostCode;

		if (hostCodeStr.length < 1) {
			showToast(get(_)('code-is-empty'), ToastType.ERROR);
			throw new Error('Failed to send client code to host');
		}

		await ConnectToHostWeb(hostCodeStr);
	} finally {
		toogleLoading();
	}
}

export async function handleEasyConnectHost() {
	const easyConnectServerUrl = (() => {
		if (easyConnectServerIpDomain.value.length < 1) {
			return '';
		}
		return `ws://${easyConnectServerIpDomain.value}/ws`;
	})();

	log(`Easy Connect Server URL: ${easyConnectServerUrl}`);

	if (!URL.canParse(easyConnectServerUrl)) {
		throw new Error('Easy Connect Server URL is not set');
	}

	const serverInstance = server(easyConnectServerUrl);

	toogleLoading();

	try {
		const clientCode = await ReceiveClientCode(serverInstance, easyConnectID.value);

		const { data: clientCodeData } = clientCode;

		if (clientCodeData.length < 1) {
			showToast(get(_)('code-is-empty'), ToastType.ERROR);
			throw new Error('Failed to receive client code from host');
		}

		const { CreateHost } = await import('$lib/webrtc/host_webrtc_hook');

		// Pending to return hostCode from CreateHost
		const hostCode = await CreateHost({ clientCode: clientCodeData, easyConnect: true });

		if (!hostCode) throw new Error('Host code undefined')
		
		await SendHostCode(serverInstance, { data: hostCode }, easyConnectID.value);
	} finally {
		toogleLoading();
	}
}
