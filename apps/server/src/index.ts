import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { Client, configFromEnv } from '@jakedegiovanni/marvel-client';

const app = new Hono();

app.get('/', c => {
  return c.text('Hello Hono!');
});

serve(
  {
    fetch: app.fetch,
    port: 3000,
  },
  info => {
    console.log(new Client(configFromEnv()).weeklyComics(new Date()));
    console.log(`Server is running on http://localhost:${info.port}`);
  },
);
