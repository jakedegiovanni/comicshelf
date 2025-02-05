export class Page<T> {
  constructor(
    readonly limit: number,
    readonly total: number,
    readonly count: number,
    readonly offset: number,
    readonly results: T[],
  ) {}
}

export class Url {
  constructor(
    readonly type: string,
    readonly url: string,
  ) {}
}

export class Comic {
  constructor(
    readonly id: number,
    readonly title: string,
    readonly urls: Url[],
    readonly thumbnail: string,
    readonly format: string,
    readonly issuer_number: number,
    readonly on_sale_date: string,
    readonly attribution: string,
    readonly attribution_link: string,
    readonly series_id: string,
  ) {}
}
