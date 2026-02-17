import { AxiosInstance, AxiosRequestConfig, isAxiosError } from "axios";
import { Resource } from "../api/types.js";

type RequestLoggingFunction = (
  ctx: unknown,
  req: AxiosRequestConfig,
  ...args: unknown[]
) => void;
type ResponseLoggingFunction = (
  ctx: unknown,
  resp: unknown,
  ...args: unknown[]
) => void;

export class Client {
  private headers: Record<string, string>;
  private client: AxiosInstance;
  private requestLoggingFunction: RequestLoggingFunction;
  private responseLoggingFunction: ResponseLoggingFunction;

  constructor(
    client: AxiosInstance,
    headers: Record<string, string>,
    requestLoggingFunction: RequestLoggingFunction,
    responseLoggingFunction: ResponseLoggingFunction,
  ) {
    this.client = client;
    this.headers = headers;
    this.requestLoggingFunction = requestLoggingFunction;
    this.responseLoggingFunction = responseLoggingFunction;
  }

  async create(
    ctx: unknown,
    resource: Resource,
    serverUrl: string,
    body: Record<string, unknown>,
    parameters: Record<string, string>,
  ): Promise<Record<string, unknown>> {
    let suffix = "";
    if (resource.createMethod?.supportsUserSettableCreate) {
      const id = body.id;
      if (!id) {
        throw new Error(`id field not found in ${JSON.stringify(body)}`);
      }
      if (typeof id === "string") {
        suffix = `?id=${id}`;
      }
    }

    const url = this.basePath(ctx, resource, serverUrl, parameters, suffix);
    const response = await this.makeRequest(ctx, "POST", url, body);
    return response;
  }

  async list(
    ctx: unknown,
    resource: Resource,
    serverUrl: string,
    parameters: Record<string, string>,
  ): Promise<Record<string, unknown>[]> {
    const url = this.basePath(ctx, resource, serverUrl, parameters, "");
    const response = await this.makeRequest(ctx, "GET", url);

    const kebab = this.kebabToCamelCase(resource.plural);
    const lowerKebab =
      kebab.length > 1 ? kebab.charAt(0).toLowerCase() + kebab.slice(1) : "";

    const possibleKeys = ["results", resource.plural, kebab, lowerKebab];

    for (const key of possibleKeys) {
      if (response[key] && Array.isArray(response[key])) {
        return response[key].filter(
          (item: unknown) => typeof item === "object",
        ) as Record<string, unknown>[];
      }
    }

    throw new Error("No valid list key was found");
  }

  async get(
    ctx: unknown,
    serverUrl: string,
    path: string,
  ): Promise<Record<string, unknown>> {
    const url = `${serverUrl}/${path.replace(/^\//, "")}`;
    return this.makeRequest(ctx, "GET", url);
  }

  async getWithFullUrl(
    ctx: unknown,
    url: string,
  ): Promise<Record<string, unknown>> {
    return this.makeRequest(ctx, "GET", url);
  }

  async delete(ctx: unknown, serverUrl: string, path: string): Promise<void> {
    const url = `${serverUrl}/${path.replace(/^\//, "")}`;
    await this.makeRequest(ctx, "DELETE", url);
  }

  async update(
    ctx: unknown,
    serverUrl: string,
    path: string,
    body: Record<string, unknown>,
  ): Promise<Record<string, unknown>> {
    const url = `${serverUrl}/${path.replace(/^\//, "")}`;
    return this.makeRequest(ctx, "PATCH", url, body);
  }

  private async makeRequest(
    ctx: unknown,
    method: string,
    url: string,
    body?: Record<string, unknown>,
  ): Promise<Record<string, unknown>> {
    if (body) {
      Object.keys(body).forEach((key) => {
        if (body[key] === null) {
          delete body[key];
        }
      });
    }
    const config: AxiosRequestConfig = {
      method,
      url,
      headers: this.headers,
      data: body,
    };

    this.requestLoggingFunction(ctx, config);

    try {
      const response = await this.client.request(config);
      this.responseLoggingFunction(ctx, response);

      const data = response.data;
      this.checkErrors(data);
      return data;
    } catch (error: unknown) {
      if (isAxiosError(error) && error.response) {
        this.responseLoggingFunction(ctx, error.response);
        throw new Error(
          `Request failed: ${JSON.stringify(error.response.data)} for request ${JSON.stringify(config)}`,
        );
      }
      throw error;
    }
  }

  private checkErrors(response: Record<string, unknown>): void {
    if (response.error) {
      throw new Error(`Returned errors: ${JSON.stringify(response.error)}`);
    }
  }

  private basePath(
    ctx: unknown,
    resource: Resource,
    serverUrl: string,
    parameters: Record<string, string>,
    suffix: string,
  ): string {
    serverUrl = serverUrl.replace(/\/$/, "");
    const urlElems = [serverUrl];

    for (let i = 0; i < resource.patternElems.length - 1; i++) {
      const elem = resource.patternElems[i];
      if (i % 2 === 0) {
        urlElems.push(elem);
      } else {
        const paramName = elem.slice(1, -1);
        const value = parameters[paramName];
        if (!value) {
          throw new Error(
            `Parameter ${paramName} not found in parameters ${JSON.stringify(parameters)}`,
          );
        }

        const lastValue = value.split("/").pop() || value;
        urlElems.push(lastValue);
      }
    }

    let result = urlElems.join("/");
    if (suffix) {
      result += suffix;
    }
    return result;
  }

  private kebabToCamelCase(str: string): string {
    return str.replace(/-([a-z])/g, (_, letter) => letter.toUpperCase());
  }
}
