import { configFromEnv } from '../src/config.ts';
import { BaseClient } from '../src/base.ts';

const config = configFromEnv();

const client = new BaseClient(config);

const response = await client.getV1PublicCharacters({ limit: 50 });
console.log(JSON.stringify(response, null, 2));
