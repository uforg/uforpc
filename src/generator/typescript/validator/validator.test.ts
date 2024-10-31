import { assertEquals } from "@std/assert";
import { SchemaValidator, schValidator } from "./validator.ts";

Deno.test("String validation", async (t) => {
  await t.step("should validate required string", () => {
    const schema = schValidator.string().required();
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, true);

    const result2 = schema.validate("");
    assertEquals(result2.isValid, true);

    const result3 = schema.validate(" ");
    assertEquals(result3.isValid, true);
  });

  await t.step("should invalidate undefined when string is required", () => {
    const schema = schValidator.string().required();

    const result1 = schema.validate(undefined);
    assertEquals(result1.isValid, false);

    const result2 = schema.validate(null);
    assertEquals(result2.isValid, false);
  });

  await t.step("should invalidate required with custom message", () => {
    const schema = schValidator.string().required("My custom error message");

    const result1 = schema.validate(undefined);
    assertEquals(result1.isValid, false);
    assertEquals(result1.error, "My custom error message");
  });

  await t.step("should validate string matching regex pattern", () => {
    const schema = schValidator.string().regex(/^[A-Z]+$/);
    const result = schema.validate("HOLA");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not matching regex pattern", () => {
    const schema = schValidator.string().regex(/^[A-Z]+$/);
    const result = schema.validate("Hola");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String does not match pattern");
  });

  await t.step("should validate string with length", () => {
    const schema = schValidator.string().length(5);
    const result = schema.validate("Hola!");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string different than length", () => {
    const schema = schValidator.string().length(5);
    const result = schema.validate("Hola Mundo!!");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String length does not match expected length");
  });

  await t.step("should validate string with minLength", () => {
    const schema = schValidator.string().minLength(5);
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string shorter than minLength", () => {
    const schema = schValidator.string().minLength(10);
    const result = schema.validate("Hola");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String is shorter than minimum length");
  });

  await t.step("should validate string with maxLength", () => {
    const schema = schValidator.string().maxLength(10);
    const result = schema.validate("Hola");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string longer than maxLength", () => {
    const schema = schValidator.string().maxLength(5);
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String is longer than maximum length");
  });

  await t.step("should validate string equals specific value", () => {
    const schema = schValidator.string().equals("Hola");
    const result = schema.validate("Hola");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not equal to specific value", () => {
    const schema = schValidator.string().equals("Hola");
    const result = schema.validate("Mundo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Value does not equal expected value");
  });

  await t.step("should validate string contains substring", () => {
    const schema = schValidator.string().contains("Mundo");
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not containing substring", () => {
    const schema = schValidator.string().contains("Adiós");
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, false);
    assertEquals(
      result.error,
      "String does not contain the required substring",
    );
  });

  await t.step("should validate email format - basic case", () => {
    const schema = schValidator.string().email();
    const tests: string[] = [
      "test@example.com",
      "user.name@domain.com",
      "user+label@domain.co.uk",
      "first.last@subdomain.domain.org",
    ];

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, true, `Failed for ${test}`);
    }
  });

  await t.step("should invalidate incorrect email formats", () => {
    const schema = schValidator.string().email();
    const tests: string[] = [
      "test@example",
      "test@.com",
      "@domain.com",
      "test@domain..com",
      "test@dom ain.com",
      "te st@domain.com",
    ];

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, false, `Failed for ${test}`);
      assertEquals(
        result.error,
        "String is not a valid email address",
        `Failed for ${test}`,
      );
    }
  });

  await t.step("should validate ISO8601 date formats", () => {
    const schema = schValidator.string().iso8601();
    const tests: string[] = [
      "2023-10-31T13:45:00Z",
      "2023-10-31T13:45:00.123Z",
      "2023-10-31T13:45:00+01:00",
      "2023-10-31T13:45:00",
    ];

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, true, `Failed for ${test}`);
    }
  });

  await t.step("should invalidate incorrect ISO8601 date formats", () => {
    const schema = schValidator.string().iso8601();
    const tests: string[] = [
      "31-10-2023",
      "2023-10-31T13:45",
      "2023-10-31T13:45:00+0100",
      "2023/10/31T13:45:00Z",
      "2023-10-31 13:45:00Z",
    ];

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, false, `Failed for ${test}`);
      assertEquals(
        result.error,
        "String is not a valid ISO8601 date",
        `Failed for ${test}`,
      );
    }
  });

  await t.step("should validate UUID strings", () => {
    const schema = schValidator.string().uuid();
    const tests: string[] = [
      "a8e636cc-97b2-11ef-b864-0242ac120002", // v1
      "000003e8-97b2-21ef-8900-325096b39f47", // v2
      "f57d6300-2b9d-3f7a-87e1-79c6ac813803", // v3
      "0ce3a2ea-1ad9-4694-9d58-13b78a7a9359", // v4
      "97767cc1-d7d2-5352-951e-0311a8530ca5", // v5
    ];

    // Add other random generated UUIDs
    for (let i = 0; i < 20; i++) {
      tests.push(crypto.randomUUID());
    }

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, true, `Failed for ${test}`);
    }
  });

  await t.step("should invalidate incorrect UUID strings", () => {
    const schema = schValidator.string().uuid();
    const tests: string[] = [
      "",
      "12345",
      "other-text",
      "a8e636cc-97b2-11ef-b864-0242ac12000",
      "000003e8-97b2-21ef-8900-325096b39f4",
      "f57d6300-2b9d-3f7a-87e1-79c6ac81380",
      "97767cc1-d7d2-5352-951e-0311a8530ca5a",
    ];

    for (const test of tests) {
      const result = schema.validate(test);
      assertEquals(result.isValid, false, `Failed for ${test}`);
      assertEquals(
        result.error,
        "String is not a valid UUID",
        `Failed for ${test}`,
      );
    }
  });

  await t.step("should validate lowercase string", () => {
    const schema = schValidator.string().lowercase();
    const result = schema.validate("hola mundo");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not in lowercase", () => {
    const schema = schValidator.string().lowercase();
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String is not in lowercase");
  });

  await t.step("should validate uppercase string", () => {
    const schema = schValidator.string().uppercase();
    const result = schema.validate("HOLA MUNDO");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not in uppercase", () => {
    const schema = schValidator.string().uppercase();
    const result = schema.validate("Hola Mundo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String is not in uppercase");
  });

  await t.step("should validate string is valid JSON", () => {
    const schema = schValidator.string().json();
    const result = schema.validate('{"nombre": "Juan", "edad": 30}');
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate string not valid JSON", () => {
    const schema = schValidator.string().json();
    const result = schema.validate('{"nombre": "Juan", "edad": }');
    assertEquals(result.isValid, false);
    assertEquals(result.error, "String is not valid JSON");
  });
});

Deno.test("Number validation", async (t) => {
  await t.step("should validate required number", () => {
    const schema = schValidator.number().required();
    const result = schema.validate(100);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate undefined when number is required", () => {
    const schema = schValidator.number().required();
    const result = schema.validate(undefined);
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Field is required");
  });

  await t.step("should validate number with min value", () => {
    const schema = schValidator.number().min(50);
    const result = schema.validate(75);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate number less than min value", () => {
    const schema = schValidator.number().min(50);
    const result = schema.validate(25);
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Number is less than the minimum allowed value");
  });

  await t.step("should validate number with max value", () => {
    const schema = schValidator.number().max(100);
    const result = schema.validate(75);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate number greater than max value", () => {
    const schema = schValidator.number().max(50);
    const result = schema.validate(75);
    assertEquals(result.isValid, false);
    assertEquals(
      result.error,
      "Number is greater than the maximum allowed value",
    );
  });

  await t.step("should validate integer number", () => {
    const schema = schValidator.int();
    const result = schema.validate(42);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate non-integer number for int schema", () => {
    const schema = schValidator.int();
    const result = schema.validate(42.5);
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Invalid type, expected int");
  });

  await t.step("should validate float number", () => {
    const schema = schValidator.float();
    const result = schema.validate(42.5);
    assertEquals(result.isValid, true);
  });

  await t.step("should validate number equals specific value", () => {
    const schema = schValidator.number().equals(100);
    const result = schema.validate(100);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate number not equal to specific value", () => {
    const schema = schValidator.number().equals(100);
    const result = schema.validate(50);
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Value does not equal expected value");
  });
});

Deno.test("Boolean validation", async (t) => {
  await t.step("should validate boolean true", () => {
    const schema = schValidator.boolean();
    const result = schema.validate(true);
    assertEquals(result.isValid, true);
  });

  await t.step("should validate boolean false", () => {
    const schema = schValidator.boolean();
    const result = schema.validate(false);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate non-boolean value", () => {
    const schema = schValidator.boolean();
    const result = schema.validate("true");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Invalid type, expected boolean");
  });
});

Deno.test("Array validation", async (t) => {
  await t.step("should validate array of integers", () => {
    const schema = schValidator.array(schValidator.int());
    const result = schema.validate([1, 2, 3]);
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate array with invalid element", () => {
    const schema = schValidator.array(schValidator.int());
    const result = schema.validate([1, 2, 3.5]);
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Invalid type, expected int");
  });

  await t.step("should invalidate non-array value", () => {
    const schema = schValidator.array(schValidator.int());
    const result = schema.validate("no es un arreglo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Invalid type, expected array");
  });
});

Deno.test("Object validation", async (t) => {
  await t.step("should validate object with defined schema", () => {
    const schema = schValidator.object({
      nombre: schValidator.string().required(),
      edad: schValidator.int().min(0),
    });
    const result = schema.validate({ nombre: "Juan", edad: 30 });
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate object with missing required field", () => {
    const schema = schValidator.object({
      nombre: schValidator.string().required(),
      edad: schValidator.int().min(0),
    });
    const result = schema.validate({ edad: 30 });
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Field is required");
  });

  await t.step("should invalidate object with invalid field", () => {
    const schema = schValidator.object({
      nombre: schValidator.string().required(),
      edad: schValidator.int().min(0),
    });
    const result = schema.validate({ nombre: "Juan", edad: -5 });
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Number is less than the minimum allowed value");
  });

  await t.step("should validate nested object", () => {
    const direccionSchema = schValidator.object({
      calle: schValidator.string().required(),
      ciudad: schValidator.string().required(),
    });

    const schema = schValidator.object({
      nombre: schValidator.string().required(),
      direccion: direccionSchema,
    });

    const result = schema.validate({
      nombre: "Juan",
      direccion: { calle: "Calle Principal 123", ciudad: "Ciudad" },
    });

    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate nested object with invalid field", () => {
    const direccionSchema = schValidator.object({
      calle: schValidator.string().required(),
      ciudad: schValidator.string().required(),
    });

    const schema = schValidator.object({
      nombre: schValidator.string().required(),
      direccion: direccionSchema,
    });

    const result = schema.validate({
      nombre: "Juan",
      direccion: { calle: "Calle Principal 123" },
    });

    assertEquals(result.isValid, false);
    assertEquals(result.error, "Field is required");
  });
});

Deno.test("Enum validation", async (t) => {
  await t.step("should validate value in enum", () => {
    const schema = schValidator.string().enum(["rojo", "verde", "azul"]);
    const result = schema.validate("rojo");
    assertEquals(result.isValid, true);
  });

  await t.step("should invalidate value not in enum", () => {
    const schema = schValidator.string().enum(["rojo", "verde", "azul"]);
    const result = schema.validate("amarillo");
    assertEquals(result.isValid, false);
    assertEquals(result.error, "Value is not in the allowed enumeration");
  });
});

Deno.test("Lazy validation", async (t) => {
  await t.step("should validate recursive structures using lazy schema", () => {
    type Nodo = {
      valor: number;
      siguiente?: Nodo;
    };

    const nodoSchema: SchemaValidator<Nodo> = schValidator.lazy(() =>
      schValidator.object({
        valor: schValidator.number().required(),
        siguiente: schValidator.lazy(() => nodoSchema),
      })
    );

    const data: Nodo = {
      valor: 1,
      siguiente: {
        valor: 2,
        siguiente: {
          valor: 3,
        },
      },
    };

    const result = nodoSchema.validate(data);
    assertEquals(result.isValid, true);
  });

  await t.step(
    "should invalidate recursive structure with invalid data",
    () => {
      type Nodo = {
        valor: number;
        siguiente?: Nodo;
      };

      const nodoSchema: SchemaValidator<Nodo> = schValidator.lazy(() =>
        schValidator.object({
          valor: schValidator.number().required(),
          siguiente: schValidator.lazy(() => nodoSchema),
        })
      );

      const data = {
        valor: 1,
        siguiente: {
          valor: "invalido",
        },
      };

      const result = nodoSchema.validate(data);
      assertEquals(result.isValid, false);
      assertEquals(result.error, "Invalid type, expected number");
    },
  );
});
