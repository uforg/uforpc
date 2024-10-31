// -----------------------------------------------------------------------------
// Schema validator
// -----------------------------------------------------------------------------

/** Available schema types for validation */
type SchemaValidatorType =
  | "string"
  | "number"
  | "int"
  | "float"
  | "boolean"
  | "array"
  | "object";

/** Result of schema validation containing validity status and optional error message */
type SchemaValidatorResult = {
  isValid: boolean;
  error?: string;
};

/**
 * Schema class for type-safe validation
 * @template T - The type of value being validated
 */
export class SchemaValidator<T> {
  private regexes = {
    email: /^((?!\.)[\w\-+_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$/,
    iso8601:
      /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[\+\-]\d{2}:\d{2})?$/,
    uuid:
      /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$/i,
  };

  private type: SchemaValidatorType;
  private arraySchema?: SchemaValidator<unknown>;
  private objectSchema?: Record<string, SchemaValidator<unknown>>;

  private errorMessage?: string;
  private isRequired = false;
  private pattern?: RegExp;
  private equalsValue?: T;
  private containsValue?: string;
  private lengthValue?: number;
  private minLengthValue?: number;
  private maxLengthValue?: number;
  private enumValues?: T[];
  private minValue?: number;
  private maxValue?: number;
  private isEmail?: boolean;
  private isIso8601?: boolean;
  private isUUID?: boolean;
  private isJSON?: boolean;
  private isLowercase?: boolean;
  private isUppercase?: boolean;

  /**
   * Creates a new validation schema
   * @param type - The type of value to validate
   * @param errorMessage - Optional custom error message
   */
  constructor(type: SchemaValidatorType, errorMessage?: string) {
    this.type = type;
    this.errorMessage = errorMessage;
  }

  /**
   * Makes the schema required (non-nullable/undefined)
   * @param errorMessage - Optional custom error message for required validation
   */
  required(errorMessage?: string): SchemaValidator<T> {
    this.isRequired = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Adds regex validation for string schemas
   * @param pattern - RegExp to test against string values
   * @param errorMessage - Optional custom error message for pattern validation
   */
  regex(pattern: RegExp, errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.pattern = pattern;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the value equals the specified value
   * @param value - Value to compare against
   * @param errorMessage - Optional custom error message for equals validation
   */
  equals(value: T, errorMessage?: string): SchemaValidator<T> {
    this.equalsValue = value;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string contains the specified substring
   * @param substring - Substring to search for
   * @param errorMessage - Optional custom error message for contains validation
   */
  contains(substring: string, errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.containsValue = substring;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string has the exact specified length
   * @param length - Exact length the string should have
   * @param errorMessage - Optional custom error message for length validation
   */
  length(length: number, errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.lengthValue = length;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string has at least the specified length
   * @param length - Minimum length the string should have
   * @param errorMessage - Optional custom error message for minLength validation
   */
  minLength(length: number, errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.minLengthValue = length;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string does not exceed the specified length
   * @param length - Maximum length the string should have
   * @param errorMessage - Optional custom error message for maxLength validation
   */
  maxLength(length: number, errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.maxLengthValue = length;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the value is one of the specified enum values
   * @param values - Array of allowed values
   * @param errorMessage - Optional custom error message for enum validation
   */
  enum(values: T[], errorMessage?: string): SchemaValidator<T> {
    this.enumValues = values;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string is a valid email address
   * @param errorMessage - Optional custom error message for email validation
   */
  email(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.isEmail = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string is a valid ISO8601 date
   * @param errorMessage - Optional custom error message for ISO8601 validation
   */
  iso8601(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.isIso8601 = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string is a valid UUIDv4
   * @param errorMessage - Optional custom error message for UUIDv4 validation
   */
  uuid(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.isUUID = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the value is valid JSON
   * @param errorMessage - Optional custom error message for JSON validation
   */
  json(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string" && this.type !== "object") return this;
    this.isJSON = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string is lowercase
   * @param errorMessage - Optional custom error message for lowercase validation
   */
  lowercase(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.isLowercase = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the string is uppercase
   * @param errorMessage - Optional custom error message for uppercase validation
   */
  uppercase(errorMessage?: string): SchemaValidator<T> {
    if (this.type !== "string") return this;
    this.isUppercase = true;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the number meets the minimum value
   * @param value - Minimum value allowed
   * @param errorMessage - Optional custom error message for min validation
   */
  min(value: number, errorMessage?: string): SchemaValidator<T> {
    if (
      this.type !== "number" && this.type !== "int" && this.type !== "float"
    ) {
      return this;
    }
    this.minValue = value;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Checks if the number does not exceed the maximum value
   * @param value - Maximum value allowed
   * @param errorMessage - Optional custom error message for max validation
   */
  max(value: number, errorMessage?: string): SchemaValidator<T> {
    if (
      this.type !== "number" && this.type !== "int" && this.type !== "float"
    ) {
      return this;
    }
    this.maxValue = value;
    this.errorMessage = errorMessage || this.errorMessage;
    return this;
  }

  /**
   * Creates an array schema
   * @param schema - Schema for array elements
   * @param errorMessage - Optional custom error message for array validation
   */
  array<U>(
    schema: SchemaValidator<U>,
    errorMessage?: string,
  ): SchemaValidator<U[]> {
    const newSchema = new SchemaValidator<U[]>("array", errorMessage);
    newSchema.arraySchema = schema;
    return newSchema;
  }

  /**
   * Creates an object schema
   * @param schema - Record of property schemas
   * @param errorMessage - Optional custom error message for object validation
   */
  object<U extends Record<string, unknown>>(
    schema: { [K in keyof U]: SchemaValidator<U[K]> },
    errorMessage?: string,
  ): SchemaValidator<U> {
    const newSchema = new SchemaValidator<U>("object", errorMessage);
    newSchema.objectSchema = schema as Record<
      string,
      SchemaValidator<unknown>
    >;
    return newSchema;
  }

  /**
   * Creates a lazy schema for recursive validation
   * @param schema - Function returning the validation schema
   * @param errorMessage - Optional custom error message for lazy validation
   */
  static lazy<T>(
    schema: () => SchemaValidator<T>,
    errorMessage?: string,
  ): SchemaValidator<T> {
    const lazySchema = new SchemaValidator<T>("object", errorMessage);
    lazySchema.validate = (value: unknown): SchemaValidatorResult =>
      schema().validate(value);
    return lazySchema;
  }

  /**
   * Validates a value against the schema
   * @param value - Value to validate
   * @returns Validation result with boolean and optional error message
   */
  validate(value: unknown): SchemaValidatorResult {
    if (value === undefined || value === null) {
      if (this.isRequired) {
        return {
          isValid: false,
          error: this.errorMessage || "Field is required",
        };
      }
      return { isValid: true };
    }

    if (!this.validateType(value)) {
      return {
        isValid: false,
        error: this.errorMessage || `Invalid type, expected ${this.type}`,
      };
    }

    if (this.equalsValue !== undefined) {
      if (value !== this.equalsValue) {
        return {
          isValid: false,
          error: this.errorMessage || `Value does not equal expected value`,
        };
      }
    }

    if (this.enumValues !== undefined) {
      if (!this.enumValues.includes(value as T)) {
        return {
          isValid: false,
          error: this.errorMessage || "Value is not in the allowed enumeration",
        };
      }
    }

    if (this.type === "string" && typeof value === "string") {
      if (this.pattern && !this.pattern.test(value)) {
        return {
          isValid: false,
          error: this.errorMessage || "String does not match pattern",
        };
      }
      if (this.containsValue !== undefined) {
        if (!value.includes(this.containsValue)) {
          return {
            isValid: false,
            error: this.errorMessage ||
              "String does not contain the required substring",
          };
        }
      }
      if (this.lengthValue !== undefined) {
        if (value.length !== this.lengthValue) {
          return {
            isValid: false,
            error: this.errorMessage ||
              "String length does not match expected length",
          };
        }
      }
      if (this.minLengthValue !== undefined) {
        if (value.length < this.minLengthValue) {
          return {
            isValid: false,
            error: this.errorMessage || "String is shorter than minimum length",
          };
        }
      }
      if (this.maxLengthValue !== undefined) {
        if (value.length > this.maxLengthValue) {
          return {
            isValid: false,
            error: this.errorMessage || "String is longer than maximum length",
          };
        }
      }
      if (this.isEmail) {
        if (!this.regexes.email.test(value)) {
          return {
            isValid: false,
            error: this.errorMessage || "String is not a valid email address",
          };
        }
      }
      if (this.isIso8601) {
        if (!this.regexes.iso8601.test(value)) {
          return {
            isValid: false,
            error: this.errorMessage || "String is not a valid ISO8601 date",
          };
        }
      }
      if (this.isUUID) {
        if (!this.regexes.uuid.test(value)) {
          return {
            isValid: false,
            error: this.errorMessage || "String is not a valid UUID",
          };
        }
      }
      if (this.isJSON) {
        try {
          JSON.parse(value);
        } catch {
          return {
            isValid: false,
            error: this.errorMessage || "String is not valid JSON",
          };
        }
      }
      if (this.isLowercase && value !== value.toLowerCase()) {
        return {
          isValid: false,
          error: this.errorMessage || "String is not in lowercase",
        };
      }
      if (this.isUppercase && value !== value.toUpperCase()) {
        return {
          isValid: false,
          error: this.errorMessage || "String is not in uppercase",
        };
      }
    }

    if (
      (this.type === "number" ||
        this.type === "int" ||
        this.type === "float") &&
      typeof value === "number"
    ) {
      if (this.minValue !== undefined && value < this.minValue) {
        return {
          isValid: false,
          error: this.errorMessage ||
            "Number is less than the minimum allowed value",
        };
      }
      if (this.maxValue !== undefined && value > this.maxValue) {
        return {
          isValid: false,
          error: this.errorMessage ||
            "Number is greater than the maximum allowed value",
        };
      }
      if (this.type === "int" && !Number.isInteger(value)) {
        return {
          isValid: false,
          error: this.errorMessage || "Number is not an integer",
        };
      }
    }

    if (this.type === "boolean" && typeof value === "boolean") {
      // No additional boolean validations implemented for the moment
    }

    if (this.isJSON && this.type === "object" && typeof value === "object") {
      try {
        JSON.stringify(value);
      } catch {
        return {
          isValid: false,
          error: this.errorMessage || "Object is not valid JSON",
        };
      }
    }

    if (this.type === "array" && Array.isArray(value)) {
      if (this.arraySchema) {
        for (const item of value) {
          const result = this.arraySchema.validate(item);
          if (!result.isValid) return result;
        }
      }
    }

    if (
      this.type === "object" &&
      typeof value === "object" &&
      !Array.isArray(value)
    ) {
      if (this.objectSchema) {
        for (const [key, schema] of Object.entries(this.objectSchema)) {
          const result = schema.validate(
            (value as Record<string, unknown>)[key],
          );
          if (!result.isValid) return result;
        }
      }
    }

    return { isValid: true };
  }

  /**
   * Validates the type of a value
   * @param value - Value to validate type of
   * @returns Whether the value matches the schema type
   */
  private validateType(value: unknown): boolean {
    switch (this.type) {
      case "string":
        return typeof value === "string";
      case "number":
        return typeof value === "number";
      case "int":
        return typeof value === "number" && Number.isInteger(value);
      case "float":
        return typeof value === "number";
      case "boolean":
        return typeof value === "boolean";
      case "array":
        return Array.isArray(value);
      case "object":
        return typeof value === "object" && !Array.isArray(value);
      default:
        return false;
    }
  }
}

/** Factory object for creating validation schemas */
export const schValidator = {
  /** Creates a string validation schema
   * @param errorMessage - Optional custom error message
   */
  string: (errorMessage?: string) =>
    new SchemaValidator<string>("string", errorMessage),

  /** Creates a number validation schema
   * @param errorMessage - Optional custom error message
   */
  number: (errorMessage?: string) =>
    new SchemaValidator<number>("number", errorMessage),

  /** Creates an integer validation schema
   * @param errorMessage - Optional custom error message
   */
  int: (errorMessage?: string) =>
    new SchemaValidator<number>("int", errorMessage),

  /** Creates a float validation schema
   * @param errorMessage - Optional custom error message
   */
  float: (errorMessage?: string) =>
    new SchemaValidator<number>("float", errorMessage),

  /** Creates a boolean validation schema
   * @param errorMessage - Optional custom error message
   */
  boolean: (errorMessage?: string) =>
    new SchemaValidator<boolean>("boolean", errorMessage),

  /** Creates an array validation schema
   * @param schema - Schema for array elements
   * @param errorMessage - Optional custom error message
   */
  array: <T>(schema: SchemaValidator<T>, errorMessage?: string) =>
    new SchemaValidator<T[]>("array", errorMessage).array(schema),

  /** Creates an object validation schema
   * @param schema - Record of property schemas
   * @param errorMessage - Optional custom error message
   */
  object: <T extends Record<string, unknown>>(
    schema: { [K in keyof T]: SchemaValidator<T[K]> },
    errorMessage?: string,
  ) => new SchemaValidator<T>("object", errorMessage).object(schema),

  /** Creates a lazy validation schema for recursive validation
   * @param schema - Function returning the validation schema
   * @param errorMessage - Optional custom error message
   */
  lazy: <T>(
    schema: () => SchemaValidator<T>,
    errorMessage?: string,
  ): SchemaValidator<T> => SchemaValidator.lazy(schema, errorMessage),
};
