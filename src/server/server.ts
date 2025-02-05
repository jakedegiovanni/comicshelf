import express, { Express, Request, Response } from 'express';
import cors from 'cors';
import { MarvelClient, MarvelConfig } from '../comics/marvel.js';

const port: number = Number(process.env.SERVER_PORT ?? -1);
if (isNaN(port) || port < 0) {
  throw new Error("couldn't parse supplied port into a number");
}

const marvelConfig = new MarvelConfig();
const marvelClient = new MarvelClient(marvelConfig);

const app: Express = express();

app.use(cors());

app.get('/api/v1/comics', async (req: Request, resp: Response) => {
  try {
    resp.json(await marvelClient.weeklyComics());
  } catch (e) {
    console.warn(`Got an exception retrieving comics: ${e}`);
    resp.sendStatus(500);
  }
});

app.listen(port, () => {
  console.log(`server is running: ${port}`);
});
