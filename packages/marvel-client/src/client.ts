import type { Config } from './config.ts';
import { BaseClient, type ComicDataWrapper } from './base.ts';

export class Client extends BaseClient {
  constructor(config: Config) {
    super(config);
  }

  public async weeklyComics(date: Date): Promise<ComicDataWrapper> {
    const first = new Date(date);
    if (first.getDay() === 0) first.setDate(first.getDate() - 1);

    first.setMonth(first.getMonth() - 3);

    while (first.getDay() !== 0) {
      first.setDate(first.getDate() - 1);
    }

    first.setDate(first.getDate() - 7);

    const last = new Date(first);
    last.setDate(last.getDate() + 6);

    return await this.getV1PublicComics({
      format: 'comic',
      formatType: 'comic',
      noVariants: true,
      dateRange: [
        first.toISOString().substring(0, 10),
        last.toISOString().substring(0, 10),
      ],
      hasDigitalIssue: true,
      orderBy: ['issueNumber'],
      limit: 100,
    });
  }
}
