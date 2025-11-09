import type {
  Contact,
  Schema,
  XAEPResourceRef,
} from "../openapi/types.js";

export interface API {
  serverURL: string;
  name: string;
  contact: Contact | null;
  schemas: Record<string, Schema>;
  resources: Record<string, Resource>;
}

export interface Resource {
  singular: string;
  plural: string;
  parents: Resource[];
  children: Resource[];
  patternElems: string[];
  schema: Schema;
  getMethod?: GetMethod;
  listMethod?: ListMethod;
  createMethod?: CreateMethod;
  updateMethod?: UpdateMethod;
  deleteMethod?: DeleteMethod;
  customMethods: CustomMethod[];
}

export interface GetMethod {}

export interface ListMethod {
  hasUnreachableResources: boolean;
  supportsFilter: boolean;
  supportsSkip: boolean;
}

export interface CreateMethod {
  supportsUserSettableCreate: boolean;
}

export interface UpdateMethod {}

export interface DeleteMethod {}

export interface CustomMethod {
  name: string;
  method: string;
  request: Schema | null;
  response: Schema | null;
}

export interface PatternInfo {
  isResourcePattern: boolean;
  customMethodName: string;
}

export interface PathWithParams {
  pattern: string;
  params: Array<{
    in: string;
    name: string;
    required: boolean;
    schema: Schema;
    "x-aep-resource-reference"?: XAEPResourceRef;
  }>;
}
