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
  }): Promise<CharacterDataWrapper> {
    return await request(this.config, 'GET', '/v1/public/characters', query);
  }
}

export type ComicList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: ComicSummary[];
};

export type EventList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: EventSummary[];
};

export type CreatorList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: CreatorSummary[];
};

export type CharacterList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: CharacterSummary[];
};

export type SeriesList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: SeriesSummary[];
};

export type StoryList = {
  available: number;
  returned: number;
  collectionURI: string;
  items: StorySummary[];
};

export type CharacterSummary = {
  resourceURI: string;
  name: string;
  role: string;
};

export type EventSummary = { resourceURI: string; name: string };

export type SeriesSummary = { resourceURI: string; name: string };

export type ComicSummary = { resourceURI: string; name: string };

export type Url = { type: string; url: string };

export type CreatorSummary = {
  resourceURI: string;
  name: string;
  role: string;
};

export type StorySummary = { resourceURI: string; name: string; type: string };

export type Image = { path: string; extension: string };

export type ComicDate = { type: string; date: string };

export type CharacterDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Character[];
};

export type EventDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Event[];
};

export type ComicPrice = { type: string; price: number };

export type EventDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: EventDataContainer;
  etag: string;
};

export type Creator = {
  id: number;
  firstName: string;
  middleName: string;
  lastName: string;
  suffix: string;
  fullName: string;
  modified: string;
  resourceURI: string;
  urls: Url[];
  thumbnail: Image;
  series: SeriesList;
  stories: StoryList;
  comics: ComicList;
  events: EventList;
};

export type Event = {
  id: number;
  title: string;
  description: string;
  resourceURI: string;
  urls: Url[];
  modified: string;
  start: string;
  end: string;
  thumbnail: Image;
  comics: ComicList;
  stories: StoryList;
  series: SeriesList;
  characters: CharacterList;
  creators: CreatorList;
  next: EventSummary;
  previous: EventSummary;
};

export type ComicDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Comic[];
};

export type TextObject = { type: string; language: string; text: string };

export type CreatorDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: CreatorDataContainer;
  etag: string;
};

export type StoryDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: StoryDataContainer;
  etag: string;
};

export type Character = {
  id: number;
  name: string;
  description: string;
  modified: string;
  resourceURI: string;
  urls: Url[];
  thumbnail: Image;
  comics: ComicList;
  stories: StoryList;
  events: EventList;
  series: SeriesList;
};

export type CharacterDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: CharacterDataContainer;
  etag: string;
};

export type ComicDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: ComicDataContainer;
  etag: string;
};

export type Series = {
  id: number;
  title: string;
  description: string;
  resourceURI: string;
  urls: Url[];
  startYear: number;
  endYear: number;
  rating: string;
  modified: string;
  thumbnail: Image;
  comics: ComicList;
  stories: StoryList;
  events: EventList;
  characters: CharacterList;
  creators: CreatorList;
  next: SeriesSummary;
  previous: SeriesSummary;
};

export type SeriesDataWrapper = {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: SeriesDataContainer;
  etag: string;
};

export type SeriesDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Series[];
};

export type StoryDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Story[];
};

export type Comic = {
  id: number;
  digitalId: number;
  title: string;
  issueNumber: number;
  variantDescription: string;
  description: string;
  modified: string;
  isbn: string;
  upc: string;
  diamondCode: string;
  ean: string;
  issn: string;
  format: string;
  pageCount: number;
  textObjects: TextObject[];
  resourceURI: string;
  urls: Url[];
  series: SeriesSummary;
  variants: ComicSummary[];
  collections: ComicSummary[];
  collectedIssues: ComicSummary[];
  dates: ComicDate[];
  prices: ComicPrice[];
  thumbnail: Image;
  images: Image[];
  creators: CreatorList;
  characters: CharacterList;
  stories: StoryList;
  events: EventList;
};

export type CreatorDataContainer = {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Creator[];
};

export type Story = {
  id: number;
  title: string;
  description: string;
  resourceURI: string;
  type: string;
  modified: string;
  thumbnail: Image;
  comics: ComicList;
  series: SeriesList;
  events: EventList;
  characters: CharacterList;
  creators: CreatorList;
  originalissue: ComicSummary;
};
