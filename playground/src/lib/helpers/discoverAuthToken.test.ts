import { describe, expect, it } from "vitest";

import {
  discoverAuthToken,
  discoverFirstAuthToken,
} from "./discoverAuthToken.ts";

describe("discoverAuthToken", () => {
  describe("basic token detection", () => {
    it("finds simple token at root level", () => {
      const data = { token: "abc123" };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual({
        key: "token",
        value: "abc123",
        path: "token",
        depth: 0,
      });
    });

    it("finds multiple tokens at root level", () => {
      const data = {
        token: "abc123",
        authToken: "def456",
        jwt: "ghi789",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(3);
      expect(result.map((t) => t.key)).toEqual(["authToken", "jwt", "token"]);
    });

    it("returns empty array when no tokens found", () => {
      const data = { username: "john", email: "john@example.com" };
      const result = discoverAuthToken(data);

      expect(result).toEqual([]);
    });

    it("handles null and undefined input", () => {
      expect(discoverAuthToken(null)).toEqual([]);
      expect(discoverAuthToken(undefined)).toEqual([]);
    });

    it("handles empty object", () => {
      const result = discoverAuthToken({});
      expect(result).toEqual([]);
    });
  });

  describe("token pattern matching", () => {
    it("detects various token naming conventions", () => {
      const data = {
        token: "value1",
        authToken: "value2",
        auth_token: "value3",
        userToken: "value4",
        sessionToken: "value5",
        refreshToken: "value6",
        ACCESS_TOKEN: "value7",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(7);
      expect(result.every((t) => t.value.startsWith("value"))).toBe(true);
    });

    it("detects JWT tokens", () => {
      const data = {
        jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
        jwtToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(2);
      expect(result.map((t) => t.key)).toEqual(["jwt", "jwtToken"]);
    });

    it("detects API key patterns", () => {
      const data = {
        apiKey: "key123",
        api_key: "key456",
        authKey: "key789",
        API_KEY: "keyABC",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(4);
      expect(result.every((t) => t.key.toLowerCase().includes("key"))).toBe(
        true,
      );
    });

    it("detects bearer tokens", () => {
      const data = {
        bearer: "bearer123",
        bearerToken: "bearer456",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(2);
      expect(result.map((t) => t.key)).toEqual(["bearer", "bearerToken"]);
    });

    it("detects session and access tokens", () => {
      const data = {
        access: "access123",
        session: "session456",
        refresh: "refresh789",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(3);
      expect(result.map((t) => t.key)).toEqual([
        "access",
        "refresh",
        "session",
      ]);
    });
  });

  describe("nested object detection", () => {
    it("finds tokens in nested objects", () => {
      const data = {
        user: {
          profile: {
            token: "nested123",
          },
        },
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual({
        key: "token",
        value: "nested123",
        path: "user.profile.token",
        depth: 2,
      });
    });

    it("finds tokens at multiple depths", () => {
      const data = {
        token: "root123",
        auth: {
          token: "auth123",
          nested: {
            sessionToken: "session123",
          },
        },
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(3);
      expect(result[0].depth).toBe(0); // root token comes first
      expect(result[1].depth).toBe(1); // auth.token
      expect(result[2].depth).toBe(2); // auth.nested.sessionToken
    });

    it("maintains correct paths for nested tokens", () => {
      const data = {
        response: {
          data: {
            authentication: {
              accessToken: "access123",
            },
          },
        },
      };
      const result = discoverAuthToken(data);

      expect(result[0].path).toBe("response.data.authentication.accessToken");
    });
  });

  describe("array handling", () => {
    it("finds tokens in arrays", () => {
      const data = {
        users: [{ token: "user1token" }, { token: "user2token" }],
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(2);
      expect(result[0].path).toBe("users[0].token");
      expect(result[1].path).toBe("users[1].token");
    });

    it("handles nested arrays", () => {
      const data = {
        sessions: [
          {
            tokens: [{ jwt: "jwt1" }, { jwt: "jwt2" }],
          },
        ],
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(2);
      expect(result[0].path).toBe("sessions[0].tokens[0].jwt");
      expect(result[1].path).toBe("sessions[0].tokens[1].jwt");
    });

    it("skips arrays when includeArrays is false", () => {
      const data = {
        token: "root123",
        users: [{ token: "array123" }],
      };
      const result = discoverAuthToken(data, { includeArrays: false });

      expect(result).toHaveLength(1);
      expect(result[0].key).toBe("token");
      expect(result[0].value).toBe("root123");
    });
  });

  describe("value validation", () => {
    it("ignores null and undefined token values", () => {
      const data = {
        token: null,
        authToken: undefined,
        validToken: "valid123",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(1);
      expect(result[0].key).toBe("validToken");
    });

    it("ignores empty string tokens", () => {
      const data = {
        token: "",
        authToken: "   ",
        validToken: "valid123",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(1);
      expect(result[0].key).toBe("validToken");
    });

    it("ignores non-string token values", () => {
      const data = {
        token: 123,
        authToken: true,
        objectToken: { nested: "value" },
        validToken: "valid123",
      };
      const result = discoverAuthToken(data);

      expect(result).toHaveLength(1);
      expect(result[0].key).toBe("validToken");
    });

    it("respects minimum token length", () => {
      const data = {
        token: "ab",
        authToken: "abc",
        longToken: "abcd",
      };
      const result = discoverAuthToken(data, { minTokenLength: 3 });

      expect(result).toHaveLength(2);
      expect(result.map((t) => t.key)).toEqual(["authToken", "longToken"]);
    });
  });

  describe("configuration options", () => {
    it("respects maxDepth option", () => {
      const data = {
        level1: {
          level2: {
            level3: {
              token: "deep123",
            },
          },
        },
      };
      const result = discoverAuthToken(data, { maxDepth: 1 });

      expect(result).toEqual([]);
    });

    it("uses custom regex pattern", () => {
      const data = {
        token: "ignored",
        customAuth: "found123",
        customSecret: "found456",
      };
      const customPattern = /^custom.*/i;
      const result = discoverAuthToken(data, { customPattern });

      expect(result).toHaveLength(2);
      expect(result.map((t) => t.key)).toEqual(["customAuth", "customSecret"]);
    });

    it("sorts results by depth then path", () => {
      const data = {
        zebraToken: "token1",
        auth: {
          betaToken: "token2",
          alphaToken: "token3",
        },
      };
      const result = discoverAuthToken(data);

      expect(result.map((t) => t.path)).toEqual([
        "zebraToken", // depth 0
        "auth.alphaToken", // depth 1, alphabetically first
        "auth.betaToken", // depth 1, alphabetically second
      ]);
    });
  });

  describe("complex real-world scenarios", () => {
    it("handles typical API authentication response", () => {
      const apiResponse = {
        success: true,
        data: {
          user: {
            id: 123,
            email: "user@example.com",
          },
          tokens: {
            accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
            refreshToken: "refresh_token_here",
            expiresIn: 3600,
          },
        },
      };
      const result = discoverAuthToken(apiResponse);

      expect(result).toHaveLength(2);
      expect(result.map((t) => t.key)).toEqual(["accessToken", "refreshToken"]);
      expect(result[0].path).toBe("data.tokens.accessToken");
    });

    it("handles OAuth provider response", () => {
      const oauthResponse = {
        access_token: "ya29.a0AfH6SMC...",
        token_type: "Bearer",
        expires_in: 3599,
        refresh_token: "1//04_token...",
        scope: "https://www.googleapis.com/auth/userinfo.email",
      };
      const result = discoverAuthToken(oauthResponse);

      expect(result).toHaveLength(3); // access_token, refresh_token, and token_type
      expect(result.map((t) => t.key)).toEqual([
        "access_token",
        "refresh_token",
        "token_type",
      ]);
    });

    it("handles mixed authentication data", () => {
      const mixedData = {
        jwt: "header.payload.signature",
        user: {
          apiKey: "ak_test_12345",
          profile: {
            sessionToken: "sess_abcdef",
          },
        },
        metadata: {
          bearerToken: "bearer_xyz789",
        },
      };
      const result = discoverAuthToken(mixedData);

      expect(result).toHaveLength(4);
      expect(result[0].key).toBe("jwt"); // depth 0
      expect(result[1].key).toBe("bearerToken"); // depth 1, alphabetically first
      expect(result[2].key).toBe("apiKey"); // depth 1, alphabetically second
      expect(result[3].key).toBe("sessionToken"); // depth 2
    });
  });
});

describe("discoverFirstAuthToken", () => {
  it("returns the first token found", () => {
    const data = {
      secondToken: "second",
      firstToken: "first",
    };
    const result = discoverFirstAuthToken(data);

    expect(result).toEqual({
      key: "firstToken",
      value: "first",
      path: "firstToken",
      depth: 0,
    });
  });

  it("returns null when no tokens found", () => {
    const data = { username: "john", email: "john@example.com" };
    const result = discoverFirstAuthToken(data);

    expect(result).toBeNull();
  });

  it("returns shallowest token when multiple depths exist", () => {
    const data = {
      nested: {
        deep: {
          token: "deep_token",
        },
      },
      token: "shallow_token",
    };
    const result = discoverFirstAuthToken(data);

    expect(result?.value).toBe("shallow_token");
    expect(result?.depth).toBe(0);
  });

  it("accepts same options as discoverAuthToken", () => {
    const data = {
      token: "short",
      authToken: "longer_token",
    };
    const result = discoverFirstAuthToken(data, { minTokenLength: 10 });

    expect(result?.key).toBe("authToken");
    expect(result?.value).toBe("longer_token");
  });
});
