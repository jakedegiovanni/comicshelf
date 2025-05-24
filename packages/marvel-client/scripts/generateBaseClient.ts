import { writeFileSync } from 'node:fs';
import swagger from './api.swagger.json' with { type: 'json' };
import path from 'node:path';

const packageDir = new URL('..', import.meta.url).pathname;

const mapDataType = (
  dt: string,
  allowMultiple: boolean,
  values: string[] | undefined,
  itemsRef: string | undefined,
): string => {
  let dataType = dt;

  if (dataType === 'Array') {
    if (!itemsRef)
      throw new Error(`itemsRef must be defined when type is of Array`);

    return `${itemsRef}[]`;
  }

  if (dataType === 'int' || dataType === 'float' || dataType === 'double')
    dataType = 'number';
  if (dataType === 'Date') dataType = 'string';

  if (values?.length) {
    dataType = values.map(v => `"${v}"`).join(' | ');
    if (allowMultiple) dataType = `(${dataType})`;
  }

  if (allowMultiple) dataType = `${dataType}[]`;

  return dataType;
};

const createMethodName = (
  path: string,
  operation: { httpMethod: string },
): string => {
  const methodPart = operation.httpMethod.toLowerCase();
  const endpointParts = path
    .split('/')
    .map(e => {
      if (e.includes('{')) {
        e = e.replaceAll('{', '').replaceAll('}', '');
        return `By${e.slice(0, 1).toUpperCase()}${e.slice(1)}`;
      }
      return `${e.slice(0, 1).toUpperCase()}${e.slice(1)}`;
    })
    .join('');

  return `${methodPart}${endpointParts}`;
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
        `${param.name}${param.required ? '' : '?'}: ${mapDataType(param.dataType, param.allowMultiple, param.allowableValues?.values, undefined)}`,
    )
    .join(',');

  if (!params) return '';

  return `query: {${params}}`;
};

const createPathParamMap = (operation: {
  parameters: {
    allowMultiple: boolean;
    required: boolean;
    name: string;
    dataType: string;
    paramType: string;
    allowableValues?: { values: string[] };
  }[];
}): { ty: string; names: string[] } => {
  const names: string[] = [];
  const params = operation.parameters
    .filter(param => param.paramType === 'path')
    .map(param => {
      names.push(param.name);
      return `${param.name}${param.required ? '' : '?'}: ${mapDataType(param.dataType, param.allowMultiple, param.allowableValues?.values, undefined)}`;
    })
    .join(',');

  if (!params) return { ty: '', names };

  return { ty: `path: {${params}}`, names };
};

const methods = swagger.apis.flatMap(api =>
  api.operations.map(operation => {
    const { ty: pathTy, names: pathNames } = createPathParamMap(operation);
    const queries = createQueryMap(operation);

    return `/**
  * ${operation.summary}
  */
  async ${createMethodName(api.path, operation)}(${[pathTy, queries].filter(i => i).join(',')}): Promise<${operation.responseClass}> {
        ${queries ? '' : 'const query = {}'}
        ${pathNames.length ? `const {${pathNames.join(',')}} = path` : ''}
        return await request(this.config, "${operation.httpMethod}", \`${api.path.replaceAll('{', '${')}\`, query);
    }`;
  }),
);

const models = Object.values(swagger.models).map(model => {
  const fields = Object.entries(model.properties).map(([k, v]) => {
    //@ts-expect-error(2339)
    return `${k}: ${mapDataType(v.type, false, [], v.items?.$ref)}`;
  });
  return `export type ${model.id} = {${fields.join(',')}}`;
});

const template = `import { request } from "./request.ts";
import type { Config } from "./config.ts";

export class BaseClient {

    protected config: Config;

    constructor(config: Config) {
        this.config = config;
    }

    ${methods.join('\n\n')}
}

${models.join('\n\n')}
`;

writeFileSync(path.join(packageDir, 'src', 'base.ts'), template);

console.log('base client generated');
