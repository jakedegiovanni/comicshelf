import { request } from './request.ts';
import type { Config } from './config.ts';

export class BaseClient {
  protected config: Config;

  constructor(config: Config) {
    this.config = config;
  }

  /**
   * Fetches lists of characters.
   */
  async getV1PublicCharacters(query: {
    name?: string;
    nameStartsWith?: string;
    modifiedSince?: string;
    comics?: number[];
    series?: number[];
    events?: number[];
    stories?: number[];
    orderBy?: ('name' | 'modified' | '-name' | '-modified')[];
    limit?: number;
    offset?: number;
  }): Promise<void> {
    return await request(this.config, 'GET', '/v1/public/characters', query);
  }
}
