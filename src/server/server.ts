import express, { Express, Request, Response } from 'express';
import crypto from 'node:crypto';
import cors from 'cors';
import { Comic, Page, Url } from '../comics/comics.js';

const port: number = Number(process.env.SERVER_PORT);
if (isNaN(port)) {
  throw new Error("couldn't parse supplied port into a number");
}

const app: Express = express();

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

app.use(cors());

app.get('/api/v1/comics', async (req: Request, resp: Response) => {
  try {
    const comics: Page<Comic> = await getWeeklyComics();
    resp.json(comics);
  } catch (e) {
    console.warn(`Got an exception retrieving comics: ${e}`);
    resp.sendStatus(500);
  }
});

app.listen(port, () => {
  console.log(`server is running: ${port}`);
});

async function getWeeklyComics(): Promise<Page<Comic>> {
  const apiKey = process.env.MARVEL_API_PUBLIC_KEY;
  const ts = Date.now();
  const hash = crypto
    .createHash('md5')
    .update(`${ts}${process.env.MARVEL_API_PRIVATE_KEY}${apiKey}`)
    .digest('hex');
  const url = `${process.env.MARVEL_API_HOST}/v1/public/comics?format=comic&formatType=comic&noVariants=true&dateRange=2024-12-01,2025-01-01&hasDigitalIssue=true&orderBy=issueNumber&limit=100&ts=${ts}&hash=${hash}&apikey=${apiKey}`;

  const response: globalThis.Response = await fetch(url);

  if (!response.ok) {
    console.log(response.headers);
    throw new Error(
      `Response not ok: ${response.status} ${response.statusText}: ${JSON.stringify(await response.json())}`,
    );
  }

  const marvelComics: dataWrapper<comic> = await response.json();

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
          process.env.MARVEL_ATTRIBUTION_LINK ?? '', // todo - should be guaranteed at this stage
          comic.series.resourceURI.match(/\/(?<seriesId>[0-9]+)$/)?.groups
            ?.seriesId ?? '',
        ),
    ),
  );
}
