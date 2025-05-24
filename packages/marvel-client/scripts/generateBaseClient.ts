import { writeFileSync } from 'node:fs';
import swagger from './api.swagger.json' with { type: 'json' };
import path from 'node:path';

const packageDir = new URL('..', import.meta.url).pathname;

const mapDataType = (param: {
  dataType: string;
  allowMultiple: boolean;
  allowableValues?: { values: string[] };
}): string => {
  let dataType = param.dataType;
  if (dataType === 'int') dataType = 'number';
  if (dataType === 'Date') dataType = 'string';

  if (param.allowableValues?.values?.length) {
    dataType = param.allowableValues.values.map(v => `"${v}"`).join(' | ');
    if (param.allowMultiple) dataType = `(${dataType})`;
  }

  if (param.allowMultiple) dataType = `${dataType}[]`;

  return dataType;
};

const createMethodName = (
  path: string,
  operation: { httpMethod: string },
): string => {
  return `${operation.httpMethod.toLowerCase()}${api.path
    .split('/')
    .map(e => `${e.slice(0, 1).toUpperCase()}${e.slice(1)}`)
    .join('')}`;
};

const createQueryMap = (operation: {
  parameters: {
    allowMultiple: boolean;
    required: boolean;
    name: string;
    dataType: string;
    paramType: string;
    allowableValues?: { values: string[] };
  }[];
}): string => {
  const params = operation.parameters
    .filter(param => param.paramType === 'query')
    .map(
      param =>
        `${param.name}${param.required ? '' : '?'}: ${mapDataType(param)}`,
    )
    .join(',');

  if (!params) return '';

  return `query: {${params}}`;
};

const createInputs = (operation: {
  parameters: {
    allowMultiple: boolean;
    required: boolean;
    name: string;
    dataType: string;
    paramType: string;
    allowableValues?: { values: string[] };
  }[];
}): string => {
  const queries = createQueryMap(operation);

  return queries;
};

const api = swagger.apis[0];
const methods = api.operations.map(operation => {
  return `/**
  * ${operation.summary}
  */
  async ${createMethodName(api.path, operation)}(${createInputs(operation)}): Promise<void> {
        return await request(this.config, "${operation.httpMethod}", "${api.path}", query);
    }`;
});

const template = `import { request } from "./request.ts";
import type { Config } from "./config.ts";

export class BaseClient {

    protected config: Config;

    constructor(config: Config) {
        this.config = config;
    }

    ${methods.join('\n')}
}`;

writeFileSync(path.join(packageDir, 'src', 'base.ts'), template);

console.log('base client generated');
