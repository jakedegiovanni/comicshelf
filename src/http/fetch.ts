export default async function fetchJson<T>(
  url: string,
  opts?: globalThis.RequestInit,
): Promise<T> {
  const response: globalThis.Response = await fetch(url, opts);

  if (!response.ok) throw new Error(response.statusText);

  // todo - we should actually validate here rather than just case
  return (await response.json()) as T;
}
