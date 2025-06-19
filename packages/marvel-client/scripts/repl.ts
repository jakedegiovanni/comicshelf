import { configFromEnv } from '../src/config.ts';
import { Client, MARVEL_UNLIMITED_OFFSET } from '../src/index.ts';
import repl from 'node:repl';

const config = configFromEnv();

const client = new Client(config);

const context = {
  config,
  client,
  MARVEL_UNLIMITED_OFFSET,
};

const r = repl.start();

Object.assign(r.context, context);
