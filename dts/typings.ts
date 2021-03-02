export type EndpointAccess = {
	headers?: object;
	params: object;
};

export type Endpoint = {
	access: EndpointAccess;

    request: {
        method: 'GET' | 'POST';
		
		headers?: {
			[key: string]: string | undefined;
		};

        params: {
			[key: string]: string | number | boolean | undefined
		};
    };

    response: {
        headers?: {
			[key: string]: string | undefined;
		};

        body: {
            status: number;
            body: any;
        };
    };
};

export type EndpointMap = {
    [url: string]: Endpoint;
};

export type TypedFetch<T extends EndpointMap> = <
    U extends keyof T,
>(
    url: U,
    init: TypeFetchReqyestInit<T[U]['request']>,
) => TypedFetchResponse<U, T[U]['response']>;

export type TypeFetchReqyestInit<
    R extends Endpoint['request'],
> = (
    & Pick<RequestInit, 
        | 'cache'
        | 'credentials'
        | 'integrity'
        | 'keepalive'
        | 'mode'
        | 'redirect'
        | 'referrer'
        | 'referrerPolicy'
        | 'signal'
        | 'window'
    >
    & {
        method: R['method'];
        body: R['params'];
    }
);

export type TypedFetchResponse<U, R extends Endpoint['response']> = Promise<(
    & Pick<Response,
        | 'arrayBuffer'
        | 'blob'
        | 'body'
        | 'bodyUsed'
        | 'clone'
        | 'formData'
        | 'headers'
        | 'ok'
        | 'redirected'
        | 'status'
        | 'statusText'
        | 'text'
        | 'trailer'
        | 'type'
    >
    & {
		url: U;
        json(): Promise<R['body']>;
    }
)>;
