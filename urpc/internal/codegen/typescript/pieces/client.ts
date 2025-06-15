import { Response, UfoError } from "./core_types";

/**
 * Mocks for the parts that are generated but not exported
 */

function asError(err: unknown): UfoError {
  return err as UfoError;
}

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Internal Client
// -----------------------------------------------------------------------------

/**
 * internalClient is the engine used by the generated façade. All identifiers
 * are deliberately un-exported because user code should interact only with the
 * generated wrappers.
 */
class internalClient {
  private baseURL: string;
  private fetchFn: typeof fetch;
  private globalHeaders: Record<string, string> = {};
  private procSet: Set<string>;
  private streamSet: Set<string>;

  constructor(
    baseURL: string,
    procNames: string[],
    streamNames: string[],
    opts: internalClientOption[]
  ) {
    this.baseURL = baseURL.replace(/\/+$/, "");
    this.procSet = new Set(procNames);
    this.streamSet = new Set(streamNames);
    this.fetchFn = globalThis.fetch?.bind(globalThis) as typeof fetch;

    opts.forEach((o) => o(this));

    if (!this.fetchFn) {
      throw new Error(
        "globalThis.fetch is undefined - please supply a custom fetch using WithFetch()"
      );
    }
  }

  async callProc(
    name: string,
    input: unknown,
    headers: Record<string, string>
  ): Promise<Response<any>> {
    if (!this.procSet.has(name)) {
      return {
        ok: false,
        error: new UfoError({
          message: `${name} procedure not found in schema`,
          category: "ClientError",
          code: "INVALID_PROC",
        }),
      };
    }

    let payload: string;
    try {
      payload = input == null ? "{}" : JSON.stringify(input);
    } catch (err) {
      return {
        ok: false,
        error: asError(err),
      };
    }

    const url = `${this.baseURL}/${name}`;
    const hdrs: Record<string, string> = {
      "content-type": "application/json",
      accept: "application/json",
      ...this.globalHeaders,
      ...headers,
    };

    try {
      const fetchResp = await this.fetchFn(url, {
        method: "POST",
        headers: hdrs,
        body: payload,
      });

      if (!fetchResp.ok) {
        return {
          ok: false,
          error: new UfoError({
            message: `Unexpected HTTP status: ${fetchResp.status}`,
            category: "HTTPError",
            code: "BAD_STATUS",
            details: { status: fetchResp.status },
          }),
        };
      }

      return await fetchResp.json();
    } catch (err) {
      return { ok: false, error: asError(err) };
    }
  }

  stream(
    name: string,
    input: unknown,
    headers: Record<string, string>
  ): {
    generator: AsyncGenerator<Response<any>, void, unknown>;
    abortController: AbortController;
  } {
    const self = this;
    const abortController = new AbortController();

    async function* generator() {
      if (!self.streamSet.has(name)) {
        yield {
          ok: false,
          error: new UfoError({
            message: `${name} stream not found in schema`,
            category: "ClientError",
            code: "INVALID_STREAM",
          }),
        } as Response<any>;
        return;
      }

      let payload: string;
      try {
        payload = input == null ? "{}" : JSON.stringify(input);
      } catch (err) {
        yield { ok: false, error: asError(err) } as Response<any>;
        return;
      }

      const url = `${self.baseURL}/${name}`;
      const hdrs: Record<string, string> = {
        "content-type": "application/json",
        accept: "text/event-stream",
        ...self.globalHeaders,
        ...headers,
      };

      let fetchResp: globalThis.Response;
      try {
        fetchResp = await self.fetchFn(url, {
          method: "POST",
          headers: hdrs,
          body: payload,
          signal: abortController.signal,
        });
      } catch (err) {
        yield { ok: false, error: asError(err) } as Response<any>;
        return;
      }

      if (!fetchResp.ok || !fetchResp.body) {
        yield {
          ok: false,
          error: new UfoError({
            message: `Unexpected HTTP status: ${fetchResp.status}`,
            category: "HTTPError",
            code: "BAD_STATUS",
            details: { status: fetchResp.status },
          }),
        } as Response<any>;
        return;
      }

      const reader = fetchResp.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });

        // Process lines
        let idx: number;
        while ((idx = buffer.indexOf("\n\n")) !== -1) {
          const line = buffer.slice(0, idx).trimEnd();
          buffer = buffer.slice(idx + 1);

          if (line === "") {
            // ignore
            continue;
          }
          if (line.startsWith("data:")) {
            const jsonStr = line.slice(5).trim();
            try {
              const evt = JSON.parse(jsonStr) as Response<any>;
              yield evt;
            } catch (err) {
              yield { ok: false, error: asError(err) } as Response<any>;
              return;
            }
          }
        }
      }
    }

    return { generator: generator(), abortController };
  }

  // Exposed mutators from builder
  setFetch(fn: typeof fetch) {
    this.fetchFn = fn.bind(globalThis) as typeof fetch;
  }

  addGlobalHeader(k: string, v: string) {
    this.globalHeaders[k] = v;
  }
}

// -----------------------------------------------------------------------------
// Builder Helpers
// -----------------------------------------------------------------------------

type internalClientOption = (c: internalClient) => void;

function withFetch(fn: typeof fetch): internalClientOption {
  return (c) => c.setFetch(fn);
}

function withGlobalHeader(key: string, value: string): internalClientOption {
  return (c) => c.addGlobalHeader(key, value);
}

// -----------------------------------------------------------------------------
// Fluent Builders exposed to generated wrappers
// -----------------------------------------------------------------------------

class clientBuilder {
  private baseURL: string;
  private opts: internalClientOption[] = [];

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  withFetch(fn: typeof fetch): clientBuilder {
    this.opts.push(withFetch(fn));
    return this;
  }

  withGlobalHeader(key: string, value: string): clientBuilder {
    this.opts.push(withGlobalHeader(key, value));
    return this;
  }

  build(procNames: string[], streamNames: string[]): internalClient {
    return new internalClient(this.baseURL, procNames, streamNames, this.opts);
  }
}
