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
    return await request(this.config, 'GET', `/v1/public/characters`, query);
  }

  /**
   * Fetches a single character by id.
   */
  async getV1PublicCharactersByCharacterId(path: {
    characterId: number;
  }): Promise<CharacterDataWrapper> {
    const query = {};
    const { characterId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/characters/${characterId}`,
      query,
    );
  }

  /**
   * Fetches lists of comics filtered by a character id.
   */
  async getV1PublicCharactersByCharacterIdComics(
    path: { characterId: number },
    query: {
      format?:
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic';
      formatType?: 'comic' | 'collection';
      noVariants?: boolean;
      dateDescriptor?: 'lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth';
      dateRange?: string[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      issueNumber?: number;
      diamondCode?: string;
      digitalId?: number;
      upc?: string;
      isbn?: string;
      ean?: string;
      issn?: string;
      hasDigitalIssue?: boolean;
      modifiedSince?: string;
      creators?: number[];
      series?: number[];
      events?: number[];
      stories?: number[];
      sharedAppearances?: number[];
      collaborators?: number[];
      orderBy?: (
        | 'focDate'
        | 'onsaleDate'
        | 'title'
        | 'issueNumber'
        | 'modified'
        | '-focDate'
        | '-onsaleDate'
        | '-title'
        | '-issueNumber'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<ComicDataWrapper> {
    const { characterId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/characters/${characterId}/comics`,
      query,
    );
  }

  /**
   * Fetches lists of events filtered by a character id.
   */
  async getV1PublicCharactersByCharacterIdEvents(
    path: { characterId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      creators?: number[];
      series?: number[];
      comics?: number[];
      stories?: number[];
      orderBy?: (
        | 'name'
        | 'startDate'
        | 'modified'
        | '-name'
        | '-startDate'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<EventDataWrapper> {
    const { characterId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/characters/${characterId}/events`,
      query,
    );
  }

  /**
   * Fetches lists of series filtered by a character id.
   */
  async getV1PublicCharactersByCharacterIdSeries(
    path: { characterId: number },
    query: {
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      modifiedSince?: string;
      comics?: number[];
      stories?: number[];
      events?: number[];
      creators?: number[];
      seriesType?: 'collection' | 'one shot' | 'limited' | 'ongoing';
      contains?: (
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic'
      )[];
      orderBy?: (
        | 'title'
        | 'modified'
        | 'startYear'
        | '-title'
        | '-modified'
        | '-startYear'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<SeriesDataWrapper> {
    const { characterId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/characters/${characterId}/series`,
      query,
    );
  }

  /**
   * Fetches lists of stories filtered by a character id.
   */
  async getV1PublicCharactersByCharacterIdStories(
    path: { characterId: number },
    query: {
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      events?: number[];
      creators?: number[];
      orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<StoryDataWrapper> {
    const { characterId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/characters/${characterId}/stories`,
      query,
    );
  }

  /**
   * Fetches lists of comics.
   */
  async getV1PublicComics(query: {
    format?:
      | 'comic'
      | 'magazine'
      | 'trade paperback'
      | 'hardcover'
      | 'digest'
      | 'graphic novel'
      | 'digital comic'
      | 'infinite comic';
    formatType?: 'comic' | 'collection';
    noVariants?: boolean;
    dateDescriptor?: 'lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth';
    dateRange?: string[];
    title?: string;
    titleStartsWith?: string;
    startYear?: number;
    issueNumber?: number;
    diamondCode?: string;
    digitalId?: number;
    upc?: string;
    isbn?: string;
    ean?: string;
    issn?: string;
    hasDigitalIssue?: boolean;
    modifiedSince?: string;
    creators?: number[];
    characters?: number[];
    series?: number[];
    events?: number[];
    stories?: number[];
    sharedAppearances?: number[];
    collaborators?: number[];
    orderBy?: (
      | 'focDate'
      | 'onsaleDate'
      | 'title'
      | 'issueNumber'
      | 'modified'
      | '-focDate'
      | '-onsaleDate'
      | '-title'
      | '-issueNumber'
      | '-modified'
    )[];
    limit?: number;
    offset?: number;
  }): Promise<ComicDataWrapper> {
    return await request(this.config, 'GET', `/v1/public/comics`, query);
  }

  /**
   * Fetches a single comic by id.
   */
  async getV1PublicComicsByComicId(path: {
    comicId: number;
  }): Promise<ComicDataWrapper> {
    const query = {};
    const { comicId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/comics/${comicId}`,
      query,
    );
  }

  /**
   * Fetches lists of characters filtered by a comic id.
   */
  async getV1PublicComicsByComicIdCharacters(
    path: { comicId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      series?: number[];
      events?: number[];
      stories?: number[];
      orderBy?: ('name' | 'modified' | '-name' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CharacterDataWrapper> {
    const { comicId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/comics/${comicId}/characters`,
      query,
    );
  }

  /**
   * Fetches lists of creators filtered by a comic id.
   */
  async getV1PublicComicsByComicIdCreators(
    path: { comicId: number },
    query: {
      firstName?: string;
      middleName?: string;
      lastName?: string;
      suffix?: string;
      nameStartsWith?: string;
      firstNameStartsWith?: string;
      middleNameStartsWith?: string;
      lastNameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      stories?: number[];
      orderBy?: (
        | 'lastName'
        | 'firstName'
        | 'middleName'
        | 'suffix'
        | 'modified'
        | '-lastName'
        | '-firstName'
        | '-middleName'
        | '-suffix'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CreatorDataWrapper> {
    const { comicId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/comics/${comicId}/creators`,
      query,
    );
  }

  /**
   * Fetches lists of events filtered by a comic id.
   */
  async getV1PublicComicsByComicIdEvents(
    path: { comicId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      series?: number[];
      stories?: number[];
      orderBy?: (
        | 'name'
        | 'startDate'
        | 'modified'
        | '-name'
        | '-startDate'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<EventDataWrapper> {
    const { comicId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/comics/${comicId}/events`,
      query,
    );
  }

  /**
   * Fetches lists of stories filtered by a comic id.
   */
  async getV1PublicComicsByComicIdStories(
    path: { comicId: number },
    query: {
      modifiedSince?: string;
      series?: number[];
      events?: number[];
      creators?: number[];
      characters?: number[];
      orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<StoryDataWrapper> {
    const { comicId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/comics/${comicId}/stories`,
      query,
    );
  }

  /**
   * Fetches lists of creators.
   */
  async getV1PublicCreators(query: {
    firstName?: string;
    middleName?: string;
    lastName?: string;
    suffix?: string;
    nameStartsWith?: string;
    firstNameStartsWith?: string;
    middleNameStartsWith?: string;
    lastNameStartsWith?: string;
    modifiedSince?: string;
    comics?: number[];
    series?: number[];
    events?: number[];
    stories?: number[];
    orderBy?: (
      | 'lastName'
      | 'firstName'
      | 'middleName'
      | 'suffix'
      | 'modified'
      | '-lastName'
      | '-firstName'
      | '-middleName'
      | '-suffix'
      | '-modified'
    )[];
    limit?: number;
    offset?: number;
  }): Promise<CreatorDataWrapper> {
    return await request(this.config, 'GET', `/v1/public/creators`, query);
  }

  /**
   * Fetches a single creator by id.
   */
  async getV1PublicCreatorsByCreatorId(path: {
    creatorId: number;
  }): Promise<CreatorDataWrapper> {
    const query = {};
    const { creatorId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/creators/${creatorId}`,
      query,
    );
  }

  /**
   * Fetches lists of comics filtered by a creator id.
   */
  async getV1PublicCreatorsByCreatorIdComics(
    path: { creatorId: number },
    query: {
      format?:
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic';
      formatType?: 'comic' | 'collection';
      noVariants?: boolean;
      dateDescriptor?: 'lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth';
      dateRange?: string[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      issueNumber?: number;
      diamondCode?: string;
      digitalId?: number;
      upc?: string;
      isbn?: string;
      ean?: string;
      issn?: string;
      hasDigitalIssue?: boolean[];
      modifiedSince?: string;
      characters?: number[];
      series?: number[];
      events?: number[];
      stories?: number[];
      sharedAppearances?: number[];
      collaborators?: number[];
      orderBy?: (
        | 'focDate'
        | 'onsaleDate'
        | 'title'
        | 'issueNumber'
        | 'modified'
        | '-focDate'
        | '-onsaleDate'
        | '-title'
        | '-issueNumber'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<ComicDataWrapper> {
    const { creatorId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/creators/${creatorId}/comics`,
      query,
    );
  }

  /**
   * Fetches lists of events filtered by a creator id.
   */
  async getV1PublicCreatorsByCreatorIdEvents(
    path: { creatorId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      characters?: number[];
      series?: number[];
      comics?: number[];
      stories?: number[];
      orderBy?: (
        | 'name'
        | 'startDate'
        | 'modified'
        | '-name'
        | '-startDate'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<EventDataWrapper> {
    const { creatorId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/creators/${creatorId}/events`,
      query,
    );
  }

  /**
   * Fetches lists of series filtered by a creator id.
   */
  async getV1PublicCreatorsByCreatorIdSeries(
    path: { creatorId: number },
    query: {
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      modifiedSince?: string;
      comics?: number[];
      stories?: number[];
      events?: number[];
      characters?: number[];
      seriesType?: 'collection' | 'one shot' | 'limited' | 'ongoing';
      contains?: (
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic'
      )[];
      orderBy?: (
        | 'title'
        | 'modified'
        | 'startYear'
        | '-title'
        | '-modified'
        | '-startYear'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<SeriesDataWrapper> {
    const { creatorId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/creators/${creatorId}/series`,
      query,
    );
  }

  /**
   * Fetches lists of stories filtered by a creator id.
   */
  async getV1PublicCreatorsByCreatorIdStories(
    path: { creatorId: number },
    query: {
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      events?: number[];
      characters?: number[];
      orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<StoryDataWrapper> {
    const { creatorId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/creators/${creatorId}/stories`,
      query,
    );
  }

  /**
   * Fetches lists of events.
   */
  async getV1PublicEvents(query: {
    name?: string;
    nameStartsWith?: string;
    modifiedSince?: string;
    creators?: number[];
    characters?: number[];
    series?: number[];
    comics?: number[];
    stories?: number[];
    orderBy?: (
      | 'name'
      | 'startDate'
      | 'modified'
      | '-name'
      | '-startDate'
      | '-modified'
    )[];
    limit?: number;
    offset?: number;
  }): Promise<EventDataWrapper> {
    return await request(this.config, 'GET', `/v1/public/events`, query);
  }

  /**
   * Fetches a single event by id.
   */
  async getV1PublicEventsByEventId(path: {
    eventId: number;
  }): Promise<EventDataWrapper> {
    const query = {};
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}`,
      query,
    );
  }

  /**
   * Fetches lists of characters filtered by an event id.
   */
  async getV1PublicEventsByEventIdCharacters(
    path: { eventId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      stories?: number[];
      orderBy?: ('name' | 'modified' | '-name' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CharacterDataWrapper> {
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}/characters`,
      query,
    );
  }

  /**
   * Fetches lists of comics filtered by an event id.
   */
  async getV1PublicEventsByEventIdComics(
    path: { eventId: number },
    query: {
      format?:
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic';
      formatType?: 'comic' | 'collection';
      noVariants?: boolean[];
      dateDescriptor?: ('lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth')[];
      dateRange?: string[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      issueNumber?: number;
      diamondCode?: string;
      digitalId?: number;
      upc?: string;
      isbn?: string;
      ean?: string;
      issn?: string;
      hasDigitalIssue?: boolean[];
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      series?: number[];
      events?: number[];
      stories?: number[];
      sharedAppearances?: number[];
      collaborators?: number[];
      orderBy?: (
        | 'focDate'
        | 'onsaleDate'
        | 'title'
        | 'issueNumber'
        | 'modified'
        | '-focDate'
        | '-onsaleDate'
        | '-title'
        | '-issueNumber'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<ComicDataWrapper> {
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}/comics`,
      query,
    );
  }

  /**
   * Fetches lists of creators filtered by an event id.
   */
  async getV1PublicEventsByEventIdCreators(
    path: { eventId: number },
    query: {
      firstName?: string;
      middleName?: string;
      lastName?: string;
      suffix?: string;
      nameStartsWith?: string;
      firstNameStartsWith?: string;
      middleNameStartsWith?: string;
      lastNameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      stories?: number[];
      orderBy?: (
        | 'lastName'
        | 'firstName'
        | 'middleName'
        | 'suffix'
        | 'modified'
        | '-lastName'
        | '-firstName'
        | '-middleName'
        | '-suffix'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CreatorDataWrapper> {
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}/creators`,
      query,
    );
  }

  /**
   * Fetches lists of series filtered by an event id.
   */
  async getV1PublicEventsByEventIdSeries(
    path: { eventId: number },
    query: {
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      modifiedSince?: string;
      comics?: number[];
      stories?: number[];
      creators?: number[];
      characters?: number[];
      seriesType?: 'collection' | 'one shot' | 'limited' | 'ongoing';
      contains?: (
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic'
      )[];
      orderBy?: (
        | 'title'
        | 'modified'
        | 'startYear'
        | '-title'
        | '-modified'
        | '-startYear'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<SeriesDataWrapper> {
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}/series`,
      query,
    );
  }

  /**
   * Fetches lists of stories filtered by an event id.
   */
  async getV1PublicEventsByEventIdStories(
    path: { eventId: number },
    query: {
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      creators?: number[];
      characters?: number[];
      orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<StoryDataWrapper> {
    const { eventId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/events/${eventId}/stories`,
      query,
    );
  }

  /**
   * Fetches lists of series.
   */
  async getV1PublicSeries(query: {
    title?: string;
    titleStartsWith?: string;
    startYear?: number;
    modifiedSince?: string;
    comics?: number[];
    stories?: number[];
    events?: number[];
    creators?: number[];
    characters?: number[];
    seriesType?: 'collection' | 'one shot' | 'limited' | 'ongoing';
    contains?: (
      | 'comic'
      | 'magazine'
      | 'trade paperback'
      | 'hardcover'
      | 'digest'
      | 'graphic novel'
      | 'digital comic'
      | 'infinite comic'
    )[];
    orderBy?: (
      | 'title'
      | 'modified'
      | 'startYear'
      | '-title'
      | '-modified'
      | '-startYear'
    )[];
    limit?: number;
    offset?: number;
  }): Promise<SeriesDataWrapper> {
    return await request(this.config, 'GET', `/v1/public/series`, query);
  }

  /**
   * Fetches a single comic series by id.
   */
  async getV1PublicSeriesBySeriesId(path: {
    seriesId: number;
  }): Promise<SeriesDataWrapper> {
    const query = {};
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}`,
      query,
    );
  }

  /**
   * Fetches lists of characters filtered by a series id.
   */
  async getV1PublicSeriesBySeriesIdCharacters(
    path: { seriesId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      events?: number[];
      stories?: number[];
      orderBy?: ('name' | 'modified' | '-name' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CharacterDataWrapper> {
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}/characters`,
      query,
    );
  }

  /**
   * Fetches lists of comics filtered by a series id.
   */
  async getV1PublicSeriesBySeriesIdComics(
    path: { seriesId: number },
    query: {
      format?:
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic';
      formatType?: 'comic' | 'collection';
      noVariants?: boolean[];
      dateDescriptor?: ('lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth')[];
      dateRange?: string[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      issueNumber?: number;
      diamondCode?: string;
      digitalId?: number;
      upc?: string;
      isbn?: string;
      ean?: string;
      issn?: string;
      hasDigitalIssue?: boolean[];
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      events?: number[];
      stories?: number[];
      sharedAppearances?: number[];
      collaborators?: number[];
      orderBy?: (
        | 'focDate'
        | 'onsaleDate'
        | 'title'
        | 'issueNumber'
        | 'modified'
        | '-focDate'
        | '-onsaleDate'
        | '-title'
        | '-issueNumber'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<ComicDataWrapper> {
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}/comics`,
      query,
    );
  }

  /**
   * Fetches lists of creators filtered by a series id.
   */
  async getV1PublicSeriesBySeriesIdCreators(
    path: { seriesId: number },
    query: {
      firstName?: string;
      middleName?: string;
      lastName?: string;
      suffix?: string;
      nameStartsWith?: string;
      firstNameStartsWith?: string;
      middleNameStartsWith?: string;
      lastNameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      events?: number[];
      stories?: number[];
      orderBy?: (
        | 'lastName'
        | 'firstName'
        | 'middleName'
        | 'suffix'
        | 'modified'
        | '-lastName'
        | '-firstName'
        | '-middleName'
        | '-suffix'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CreatorDataWrapper> {
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}/creators`,
      query,
    );
  }

  /**
   * Fetches lists of events filtered by a series id.
   */
  async getV1PublicSeriesBySeriesIdEvents(
    path: { seriesId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      comics?: number[];
      stories?: number[];
      orderBy?: (
        | 'name'
        | 'startDate'
        | 'modified'
        | '-name'
        | '-startDate'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<EventDataWrapper> {
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}/events`,
      query,
    );
  }

  /**
   * Fetches lists of stories filtered by a series id.
   */
  async getV1PublicSeriesBySeriesIdStories(
    path: { seriesId: number },
    query: {
      modifiedSince?: string;
      comics?: number[];
      events?: number[];
      creators?: number[];
      characters?: number[];
      orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<StoryDataWrapper> {
    const { seriesId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/series/${seriesId}/stories`,
      query,
    );
  }

  /**
   * Fetches lists of stories.
   */
  async getV1PublicStories(query: {
    modifiedSince?: string;
    comics?: number[];
    series?: number[];
    events?: number[];
    creators?: number[];
    characters?: number[];
    orderBy?: ('id' | 'modified' | '-id' | '-modified')[];
    limit?: number;
    offset?: number;
  }): Promise<StoryDataWrapper> {
    return await request(this.config, 'GET', `/v1/public/stories`, query);
  }

  /**
   * Fetches a single comic story by id.
   */
  async getV1PublicStoriesByStoryId(path: {
    storyId: number;
  }): Promise<StoryDataWrapper> {
    const query = {};
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}`,
      query,
    );
  }

  /**
   * Fetches lists of characters filtered by a story id.
   */
  async getV1PublicStoriesByStoryIdCharacters(
    path: { storyId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      events?: number[];
      orderBy?: ('name' | 'modified' | '-name' | '-modified')[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CharacterDataWrapper> {
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}/characters`,
      query,
    );
  }

  /**
   * Fetches lists of comics filtered by a story id.
   */
  async getV1PublicStoriesByStoryIdComics(
    path: { storyId: number },
    query: {
      format?:
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic';
      formatType?: 'comic' | 'collection';
      noVariants?: boolean[];
      dateDescriptor?: ('lastWeek' | 'thisWeek' | 'nextWeek' | 'thisMonth')[];
      dateRange?: string[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      issueNumber?: number;
      diamondCode?: string;
      digitalId?: number;
      upc?: string;
      isbn?: string;
      ean?: string;
      issn?: string;
      hasDigitalIssue?: boolean[];
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      series?: number[];
      events?: number[];
      sharedAppearances?: number[];
      collaborators?: number[];
      orderBy?: (
        | 'focDate'
        | 'onsaleDate'
        | 'title'
        | 'issueNumber'
        | 'modified'
        | '-focDate'
        | '-onsaleDate'
        | '-title'
        | '-issueNumber'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<ComicDataWrapper> {
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}/comics`,
      query,
    );
  }

  /**
   * Fetches lists of creators filtered by a story id.
   */
  async getV1PublicStoriesByStoryIdCreators(
    path: { storyId: number },
    query: {
      firstName?: string;
      middleName?: string;
      lastName?: string;
      suffix?: string;
      nameStartsWith?: string;
      firstNameStartsWith?: string;
      middleNameStartsWith?: string;
      lastNameStartsWith?: string;
      modifiedSince?: string;
      comics?: number[];
      series?: number[];
      events?: number[];
      orderBy?: (
        | 'lastName'
        | 'firstName'
        | 'middleName'
        | 'suffix'
        | 'modified'
        | '-lastName'
        | '-firstName'
        | '-middleName'
        | '-suffix'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<CreatorDataWrapper> {
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}/creators`,
      query,
    );
  }

  /**
   * Fetches lists of events filtered by a story id.
   */
  async getV1PublicStoriesByStoryIdEvents(
    path: { storyId: number },
    query: {
      name?: string;
      nameStartsWith?: string;
      modifiedSince?: string;
      creators?: number[];
      characters?: number[];
      series?: number[];
      comics?: number[];
      orderBy?: (
        | 'name'
        | 'startDate'
        | 'modified'
        | '-name'
        | '-startDate'
        | '-modified'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<EventDataWrapper> {
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}/events`,
      query,
    );
  }

  /**
   * Fetches lists of series filtered by a story id.
   */
  async getV1PublicStoriesByStoryIdSeries(
    path: { storyId: number },
    query: {
      events?: number[];
      title?: string;
      titleStartsWith?: string;
      startYear?: number;
      modifiedSince?: string;
      comics?: number[];
      creators?: number[];
      characters?: number[];
      seriesType?: 'collection' | 'one shot' | 'limited' | 'ongoing';
      contains?: (
        | 'comic'
        | 'magazine'
        | 'trade paperback'
        | 'hardcover'
        | 'digest'
        | 'graphic novel'
        | 'digital comic'
        | 'infinite comic'
      )[];
      orderBy?: (
        | 'title'
        | 'modified'
        | 'startYear'
        | '-title'
        | '-modified'
        | '-startYear'
      )[];
      limit?: number;
      offset?: number;
    },
  ): Promise<SeriesDataWrapper> {
    const { storyId } = path;
    return await request(
      this.config,
      'GET',
      `/v1/public/stories/${storyId}/series`,
      query,
    );
  }
}

export interface ComicList {
  available: number;
  returned: number;
  collectionURI: string;
  items: ComicSummary[];
}

export interface EventList {
  available: number;
  returned: number;
  collectionURI: string;
  items: EventSummary[];
}

export interface CreatorList {
  available: number;
  returned: number;
  collectionURI: string;
  items: CreatorSummary[];
}

export interface CharacterList {
  available: number;
  returned: number;
  collectionURI: string;
  items: CharacterSummary[];
}

export interface SeriesList {
  available: number;
  returned: number;
  collectionURI: string;
  items: SeriesSummary[];
}

export interface StoryList {
  available: number;
  returned: number;
  collectionURI: string;
  items: StorySummary[];
}

export interface CharacterSummary {
  resourceURI: string;
  name: string;
  role: string;
}

export interface EventSummary {
  resourceURI: string;
  name: string;
}

export interface SeriesSummary {
  resourceURI: string;
  name: string;
}

export interface ComicSummary {
  resourceURI: string;
  name: string;
}

export interface Url {
  type: string;
  url: string;
}

export interface CreatorSummary {
  resourceURI: string;
  name: string;
  role: string;
}

export interface StorySummary {
  resourceURI: string;
  name: string;
  type: string;
}

export interface Image {
  path: string;
  extension: string;
}

export interface ComicDate {
  type: string;
  date: string;
}

export interface CharacterDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Character[];
}

export interface EventDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Event[];
}

export interface ComicPrice {
  type: string;
  price: number;
}

export interface EventDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: EventDataContainer;
  etag: string;
}

export interface Creator {
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
}

export interface Event {
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
}

export interface ComicDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Comic[];
}

export interface TextObject {
  type: string;
  language: string;
  text: string;
}

export interface CreatorDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: CreatorDataContainer;
  etag: string;
}

export interface StoryDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: StoryDataContainer;
  etag: string;
}

export interface Character {
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
}

export interface CharacterDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: CharacterDataContainer;
  etag: string;
}

export interface ComicDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: ComicDataContainer;
  etag: string;
}

export interface Series {
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
}

export interface SeriesDataWrapper {
  code: number;
  status: string;
  copyright: string;
  attributionText: string;
  attributionHTML: string;
  data: SeriesDataContainer;
  etag: string;
}

export interface SeriesDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Series[];
}

export interface StoryDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Story[];
}

export interface Comic {
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
}

export interface CreatorDataContainer {
  offset: number;
  limit: number;
  total: number;
  count: number;
  results: Creator[];
}

export interface Story {
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
}
