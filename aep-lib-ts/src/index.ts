// Main exports for the aep-lib-ts package
export * from './cases/cases.js';
export * from './api/api.js';
export * from './api/resource.js';
export * from './openapi/openapi.js';
export * from './client/client.js';
export * from './constants/constants.js';
export { logger } from './utils/logger.js';

// Re-export types with explicit naming to avoid conflicts
export type {
  API,
  Contact,
  CustomMethod,
  CreateMethod,
  DeleteMethod,
  GetMethod,
  ListMethod,
  UpdateMethod,
  Resource,
  PatternInfo,
  APISchema
} from './api/types.js';

export type {
  Components,
  Info,
  MediaType,
  OpenAPI,
  Operation,
  Parameter,
  PathItem,
  RequestBody,
  Response,
  Server,
  ServerVariable,
  XAEPResource,
  XAEPResourceRef
} from './openapi/types.js';