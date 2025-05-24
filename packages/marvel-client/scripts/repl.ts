import { configFromEnv } from '../src/config.ts';
import { BaseClient } from '../src/base.ts';
import repl from 'node:repl';

const config = configFromEnv();

const client = new BaseClient(config);

const context = {
  config,
  client,
};

const r = repl.start();

Object.assign(r.context, context);
