import type { Config } from './config.ts';
import { createHash } from 'node:crypto';

// todo - figure out how to remove this any
//eslint-disable-next-line @typescript-eslint/no-explicit-any
type QueryValue = any | any[];

// todo - schema validation here or above?
export const request = async <Response, Data>(
  config: Config,
  method: string,
  endpoint: string,
  queryParams: Record<string, QueryValue> = {},
  data?: Data,
): Promise<Response> => {
  const url = new URL(endpoint, config.API_HOST);

  const params = {
    ...Object.fromEntries(
      Object.entries(queryParams).map(([k, v]) => {
        const val = Array.isArray(v) ? v.join(',') : v;
        return [k, String(val)];
      }),
    ),
    ...authenticate(config),
  };
  url.search = new URLSearchParams(params).toString();

  const headers = new Headers({ Accept: 'application/json' });

  let body = {};
  if (method === 'PUT' || (method === 'POST' && data)) {
    headers.set('Content-Type', 'application/json');
    body = { body: JSON.stringify(data) };
  }

  const response = await fetch(url, {
    method,
    headers,
    ...body,
  });

  if (!response.ok) {
    console.log(response.headers); // todo remove
    throw new Error(
      `Request not ok with: ${response.status} ${response.statusText}`,
    );
  }

  return (await response.json()) as Response;
};

const authenticate = (config: Config): Record<string, string> => {
  const ts = `${Date.now()}`;

  const hash = createHash('md5')
    .update(`${ts}${config.API_PRIVATE_KEY}${config.API_PUBLIC_KEY}`)
    .digest('hex');

  return {
    ts,
    hash,
    apikey: config.API_PUBLIC_KEY,
  };
};
