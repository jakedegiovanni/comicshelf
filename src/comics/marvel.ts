import crypto from 'node:crypto';
import { Comic, Page, Url } from './comics.js';
import fetchJson from '../http/fetch.js';

interface dataWrapper<T> {
  code: number;
  status: string;
  copyright: string;
  attribution_text: string;
  attribution_html: string;
  etag: string;
  data: dataContainer<T>;
}

interface dataContainer<T> {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: T[];
}

interface item {
  name: string;
  resourceURI: string;
}

interface uri {
  type: string;
  url: string;
}

interface date {
  type: string;
  date: string;
}

interface thumbnail {
  path: string;
  extension: string;
}

interface baseResult {
  id: number;
  title: string;
  resourceURI: string;
  urls: uri[];
  modified: string;
  thumbnail: thumbnail;
}

interface comic extends baseResult {
  format: string;
  issueNumber: number;
  series: item;
  dates: date[];
}

export class MarvelConfig {
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

export class MarvelClient {
  #config: MarvelConfig;

  constructor(config: MarvelConfig) {
    this.#config = config;
  }

  public async weeklyComics(): Promise<Page<Comic>> {
    const first = new Date(); // todo - this should be a form on the ui
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
      .update(`${ts}${this.#config.privateKey}${this.#config.publicKey}`)
      .digest('hex');

    const url = `${this.#config.host}/v1/public/comics?format=comic&formatType=comic&noVariants=true&dateRange=${first.toISOString().substring(0, 10)},${last.toISOString().substring(0, 10)}&hasDigitalIssue=true&orderBy=issueNumber&limit=100&ts=${ts}&hash=${hash}&apikey=${this.#config.publicKey}`;

    const marvelComics = await fetchJson<dataWrapper<comic>>(url);

    return new Page<Comic>(
      marvelComics.data.limit,
      marvelComics.data.total,
      marvelComics.data.count,
      marvelComics.data.offset,
      marvelComics.data.results.map(
        (comic: comic): Comic =>
          new Comic(
            comic.id,
            comic.title,
            comic.urls.map((url: uri): Url => new Url(url.type, url.url)),
            `${comic.thumbnail.path}/portrait_uncanny.${comic.thumbnail.extension}`,
            comic.format,
            comic.issueNumber,
            comic.dates.find((date: date): boolean => {
              return date.type.toLowerCase().trim() === 'onsaledate';
            })?.date ?? '',
            marvelComics.attribution_text,
            this.#config.attributionLink,
            comic.series.resourceURI.match(/\/(?<seriesId>[0-9]+)$/)?.groups
              ?.seriesId ?? '',
          ),
      ),
    );
  }
}
