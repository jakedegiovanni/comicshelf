import crypto from 'node:crypto';

interface DataWrapper<T> {
  code: number;
  status: string;
  copyright: string;
  attribution_text: string;
  attribution_html: string;
  etag: string;
  data: DataContainer<T>;
}

interface DataContainer<T> {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: T[];
}

interface Item {
  name: string;
  resourceURI: string;
}

interface Uri {
  type: string;
  url: string;
}

interface MarvelDate {
  type: string;
  date: string;
}

interface Thumbnail {
  path: string;
  extension: string;
}

interface BaseResult {
  id: number;
  title: string;
  resourceURI: string;
  urls: Uri[];
  modified: string;
  thumbnail: Thumbnail;
}

interface Comic extends BaseResult {
  format: string;
  issueNumber: number;
  series: Item;
  dates: MarvelDate[];
}

export class Config {
  public readonly publicKey: string;
  public readonly privateKey: string;
  public readonly host: string;
  public readonly attributionLink: string;

  constructor() {
    this.publicKey = process.env.MARVEL_API_PUBLIC_KEY ?? '';
    if (this.publicKey === '')
      throw new Error('required environment MARVEL_API_PUBLIC_KEY not set');

    this.privateKey = process.env.MARVEL_API_PRIVATE_KEY ?? '';
    if (this.privateKey === '')
      throw new Error('required environment MARVEL_API_PRIVATE_KEY not set');

    this.host = process.env.MARVEL_API_HOST ?? '';
    if (this.host === '')
      throw new Error('required environment MARVEL_API_HOST not set');

    this.attributionLink = process.env.MARVEL_ATTRIBUTION_LINK ?? '';
    if (this.attributionLink === '')
      throw new Error('required environment MARVEL_ATTRIBUTION_LINK not set');
  }
}

export class Client {
  private config: Config;

  constructor(config: Config) {
    this.config = config;
  }

  public async weeklyComics(date: Date): Promise<DataWrapper<Comic>> {
    const first = new Date(date);
    if (first.getDay() === 0) first.setDate(first.getDate() - 1);

    first.setMonth(first.getMonth() - 3);

    while (first.getDay() !== 0) {
      first.setDate(first.getDate() - 1);
    }

    first.setDate(first.getDate() - 7);

    const last = new Date(first);
    last.setDate(last.getDate() + 6);

    const ts = Date.now();

    const hash = crypto
      .createHash('md5')
      .update(`${ts}${this.config.privateKey}${this.config.publicKey}`)
      .digest('hex');

    const url = `${this.config.host}/v1/public/comics?format=comic&formatType=comic&noVariants=true&dateRange=${first.toISOString().substring(0, 10)},${last.toISOString().substring(0, 10)}&hasDigitalIssue=true&orderBy=issueNumber&limit=100&ts=${ts}&hash=${hash}&apikey=${this.config.publicKey}`;

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        Accept: 'application/json',
      },
    });

    if (!response.ok) throw new Error(response.statusText);

    return (await response.json()) as DataWrapper<Comic>;
  }
}
