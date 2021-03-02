import {
	EndpointMap,
	TypedFetch,
	TypeFetchReqyestInit,
	TypedFetchResponse,
} from './typings';

export function createTypedFetch<T extends EndpointMap>(originalFetch: typeof fetch): TypedFetch<T> {
	return <U extends keyof T>(url: U, init: TypeFetchReqyestInit<T[U]['request']>) => {
		const query = new URLSearchParams(Object(init.body));
		const fetchInit = {...Object(init)} as RequestInit;
		let fetchUrl = url as string;

		if (init.method === 'GET') {
			fetchUrl += `?${query}`;
		} else {
			fetchInit.body = query;
		}

		return originalFetch(fetchUrl, fetchInit) as TypedFetchResponse<U, T[U]['response']>;
	};
}
