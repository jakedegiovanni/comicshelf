import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { Client, Config } from '@jakedegiovanni/marvel-client';

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
    console.log(new Client(new Config()).weeklyComics(new Date()));
    console.log(`Server is running on http://localhost:${info.port}`);
  },
);
