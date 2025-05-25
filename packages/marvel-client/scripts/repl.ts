import { configFromEnv } from '../src/config.ts';
import { Client } from '../src/index.ts';
import repl from 'node:repl';

const config = configFromEnv();

const client = new Client(config);

const context = {
  config,
  client,
};

const r = repl.start();

Object.assign(r.context, context);
