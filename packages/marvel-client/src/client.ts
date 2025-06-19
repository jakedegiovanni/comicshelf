import { BaseClient, type ComicDataWrapper } from './base.ts';

export declare const DAY: {
  SUNDAY: 0;
  MONDAY: 1;
  TUESDAY: 2;
  WEDNESDAY: 3;
  THURSDAY: 4;
  FRIDAY: 5;
  SATURDAY: 6;
};

export type DAY = (typeof DAY)[keyof typeof DAY];

export const MARVEL_UNLIMITED_OFFSET = { year: 0, month: -3, day: -7 };

export class Client extends BaseClient {
  // todo - dayJs for easier manipulation
  // offset of {year: 0, moneth: -3}
  public async weeklyComics(
    date: Date,
    offset: { year: number; month: number; day: number } = {
      year: 0,
      month: 0,
      day: 0,
    },
    weekStart: DAY = 0,
  ): Promise<ComicDataWrapper> {
    const first = new Date(date);
    // if (first.getDay() === 0) first.setDate(first.getDate() - 1);

    while (first.getDay() !== weekStart) {
      first.setUTCDate(first.getUTCDate() - 1);
    }

    first.setFullYear(first.getFullYear() + offset.year);
    first.setMonth(first.getMonth() + offset.month);
    first.setDate(first.getDate() + offset.day);

    const last = new Date(first);
    last.setUTCDate(last.getUTCDate() + 6);

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
