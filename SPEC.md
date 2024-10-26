# UFO RPC Specification

## Table of Contents

1. [Introduction](#introduction)
2. [Design Goals](#design-goals)
3. [Architecture Overview](#architecture-overview)
4. [Schema Definition Language](#schema-definition-language)
   - [Data Types](#data-types)
     - [Primitive Types](#primitive-types)
     - [Custom Types](#custom-types)
     - [Arrays and Nested Objects](#arrays-and-nested-objects)
   - [Procedures](#procedures)
     - [Procedure Definition](#procedure-definition)
     - [Procedure Metadata](#procedure-metadata)
     - [Queries vs. Mutations](#queries-vs-mutations)
     - [Input and Output Structure](#input-and-output-structure)
5. [Data Serialization](#data-serialization)
6. [Error Handling](#error-handling)
   - [Error Response Format](#error-response-format)
7. [Supported Languages](#supported-languages)
8. [Example Schema](#example-schema)
9. [Conclusion](#conclusion)
10. [Appendix](#appendix)
    - [Implementing Arrays and Nested Objects](#implementing-arrays-and-nested-objects)
    - [Handling Custom Types](#handling-custom-types)

---

## Introduction

**UFO RPC** is a cross-language Remote Procedure Call (RPC) framework designed
to facilitate communication between clients and servers, as well as between
microservices. It emphasizes a schema-first approach, strong type safety, and
exceptional ease of use. The framework enables developers to define APIs using a
simple schema language, from which both client and server code can be
automatically generated.

Initially, UFO RPC provides support for **TypeScript** and **Go**, with plans to
extend support to additional programming languages in the future.

### Schema Format

UFO RPC uses **JSON** as the primary format for schema definitions. Internally,
the framework processes schemas written in JSON. However, to enhance the
development experience, developers can write their schemas in **YAML** or other
formats, which are then transpiled into JSON before being processed by UFO RPC.
This provides syntactic sugar and improved readability without affecting the
formal specification or internal processing. For the purposes of this
specification, all examples and definitions will be presented in **JSON**
format.

You can find the UFO RPC json schema in the
[src/schema/schema.json](./src/schema/schema.json) file and the validation
schema in the [src/schema/schema.ts](./src/schema/schema.ts) file

---

## Design Goals

- **Simplicity**: Offer a straightforward and minimalistic framework to reduce
  overhead in API development.
- **Type Safety**: Ensure strong typing across client and server code to catch
  errors at compile time.
- **Schema-First Approach**: Use a schema as the single source of truth for API
  definitions.
- **Cross-Language Compatibility**: Enable seamless integration between
  different programming languages.
- **Standard HTTP Methods**: Leverage HTTP GET and POST methods for efficient
  communication and caching.
- **JSON Data Serialization**: Utilize JSON for data transport to ensure wide
  compatibility and ease of debugging.
- **Extensibility**: Allow future enhancements without breaking existing
  implementations.

---

## Architecture Overview

UFO RPC uses a schema definition language written in **JSON** format. Developers
define their APIs in this schema, specifying data types and procedures (queries
and mutations). The schema is then used to generate both client and server code,
ensuring consistency and reducing manual coding efforts.

To enhance the development experience, UFO RPC allows schemas to be written in
**YAML** or other formats, which are transpiled into JSON before processing.
This provides syntactic sugar and improved readability without affecting the
formal specification or internal processing.

---

## Schema Definition Language

### Data Types

#### Primitive Types

UFO RPC supports the following primitive data types:

- `string`
- `int`
- `float`
- `boolean`

#### Custom Types

Custom types are user-defined data structures composed of primitive types and
other custom types. They allow for modeling complex data structures in a
type-safe manner. All custom type names must start with a capital letter.

**Syntax:**

```json
{
  "types": [
    {
      "name": "TypeName",
      "desc": "Description of the type.",
      "fields": {
        "fieldName": {
          "type": "fieldType",
          "desc": "Description of the field."
        },
        "anotherField": {
          "type": "anotherType",
          "desc": "Description of another field."
        }
      }
    }
  ]
}
```

**Example:**

```json
{
  "types": [
    {
      "name": "User",
      "desc": "Represents a user in the system.",
      "fields": {
        "id": {
          "type": "string",
          "desc": "Unique identifier of the user."
        },
        "username": {
          "type": "string",
          "desc": "Username of the user."
        },
        "email": {
          "type": "string",
          "desc": "Email address of the user."
        }
      }
    }
  ]
}
```

#### Arrays and Nested Objects

##### Arrays

Arrays are denoted by appending `[]` to the type of the elements they contain.

**Syntax:**

```json
{
  "fieldName": {
    "type": "elementType[]",
    "desc": "Description of the array field."
  }
}
```

**Example:**

```json
{
  "roles": {
    "type": "string[]",
    "desc": "List of roles assigned to the user."
  },
  "users": {
    "type": "User[]",
    "desc": "Array of user objects."
  }
}
```

##### Nested Objects

Nested objects can be defined inline or by referencing custom types. When
defining inline nested objects, you must specify `"type": "object"` and include
a `"fields"` section. Similarly, you can use `"type": "object[]"` for arrays of
objects.

- **Inline Definition:**

  ```json
  {
    "profile": {
      "type": "object",
      "desc": "User profile information.",
      "fields": {
        "age": {
          "type": "int",
          "desc": "Age of the user."
        },
        "address": {
          "type": "object",
          "desc": "Address details.",
          "fields": {
            "street": {
              "type": "string",
              "desc": "Street name."
            },
            "city": {
              "type": "string",
              "desc": "City name."
            },
            "zipCode": {
              "type": "string",
              "desc": "Postal code."
            }
          }
        }
      }
    }
  }
  ```

- **Reference to Custom Types:**

  If `Address` is defined as a custom type:

  ```json
  {
    "types": [
      {
        "name": "Address",
        "desc": "Represents a physical address.",
        "fields": {
          "street": {
            "type": "string",
            "desc": "Street name."
          },
          "city": {
            "type": "string",
            "desc": "City name."
          },
          "zipCode": {
            "type": "string",
            "desc": "Postal code."
          }
        }
      }
    ]
  }
  ```

  Then, use it in another type:

  ```json
  {
    "address": {
      "type": "Address",
      "desc": "Address details."
    }
  }
  ```

### Procedures

#### Procedure Definition

Procedures represent the remote functions that can be called. All procedure
names must start with a capital letter. They are defined with the following
attributes:

- `name`: A unique identifier for the procedure (should start with a capital
  letter).
- `type`: Specifies whether it is a `query` or a `mutation`.
- `desc`: An optional description of the procedure.
- `input`: The input parameters of the procedure.
- `output`: The output returned by the procedure.
- `meta`: Optional metadata for the procedure.

If a procedure does not require input parameters or does not return any output,
the `input` or `output` fields can be omitted.

**Syntax:**

```json
{
  "procedures": [
    {
      "name": "ProcedureName",
      "type": "query | mutation",
      "desc": "Description of the procedure.",
      "input": {
        "parameterName": {
          "type": "parameterType",
          "desc": "Description of the parameter."
        },
        "anotherParameter": {
          "type": "anotherType",
          "desc": "Description of another parameter."
        }
      },
      "output": {
        "fieldName": {
          "type": "fieldType",
          "desc": "Description of the output field."
        },
        "anotherField": {
          "type": "anotherType",
          "desc": "Description of another output field."
        }
      },
      "meta": {
        "key": "value"
      }
    }
  ]
}
```

#### Procedure Metadata

Procedure metadata provides additional information about procedures, which can
be utilized when **using the code generated by UFO RPC**. These metadata are not
used during the code generation process itself but can be leveraged in the
generated code for various purposes such as authentication, logging, or custom
annotations.

The metadata values can only be [primitive types](#primitive-types)

**Syntax:**

```json
"meta": {
  "key": "value"
}
```

**Example:**

```json
{
  "procedures": [
    {
      "name": "CreateUser",
      "type": "mutation",
      "desc": "Creates a new user in the system.",
      "meta": {
        "requiresAuth": true,
        "logLevel": "debug"
      },
      "input": {
        "user": {
          "type": "User",
          "desc": "User information to create."
        }
      },
      "output": {
        "id": {
          "type": "string",
          "desc": "Unique identifier of the created user."
        },
        "username": {
          "type": "string",
          "desc": "Username of the created user."
        }
      }
    }
  ]
}
```

#### Queries vs. Mutations

The distinction between **queries** and **mutations** is fundamental for both
semantic clarity and operational behavior:

- **Queries**:
  - Represent **read-only** operations that do not modify server state.
  - Utilize the HTTP **GET** method.
  - Can be easily cached by clients and intermediaries due to their idempotent
    nature.
  - Enhance reasoning about the codebase by clearly indicating non-mutating
    operations.

- **Mutations**:
  - Represent operations that **modify server state**.
  - Utilize the HTTP **POST** method.
  - Not cached, ensuring that state changes are always executed.
  - Provide semantic clarity by distinguishing state-changing operations.

This separation allows developers to reason more effectively about their
codebase and leverages HTTP semantics for efficient network communication and
caching strategies.

#### Input and Output Structure

The `input` and `output` fields in a procedure definition can be:

- **An object** defining the parameters (fields) of the input/output. Each field
  can have:
  - `type`: The type of the field.
  - `desc` (optional): A description of the field.

- **A string representing a type**. In this case, the input/output is of that
  type.

If the `input` or `output` is assigned only a string, it is interpreted as the
type of the input/output.

**Important Note:**

- **Fields in `input` and `output` are treated as an anonymous custom type**
  with exactly the same fields. This means they have the same structure and
  behavior as a named custom type but without a specific name.

**Examples:**

- **Object Definition:**

  ```json
  {
    "input": {
      "userId": {
        "type": "string",
        "desc": "Unique identifier of the user."
      },
      "profile": {
        "type": "Profile",
        "desc": "Profile information to update."
      }
    }
  }
  ```

- **Type Reference:**

  ```json
  {
    "input": "User"
  }
  ```

This flexibility allows for concise definitions when the input or output
corresponds directly to a custom type.

---

## Data Serialization

All data exchanged between clients and servers using UFO RPC is serialized in
**JSON** format. JSON was chosen due to its ubiquity, readability, and ease of
use across multiple programming languages. This ensures:

- **Wide Compatibility**: Almost all programming languages have built-in support
  for JSON parsing and serialization.
- **Ease of Debugging**: Human-readable format simplifies the debugging process
  during development and troubleshooting.

---

## Error Handling

### Error Response Format

When errors occur during the execution of a procedure, UFO RPC returns a
standardized JSON error response. This consistent format allows client
applications to handle errors efficiently and uniformly.

**Error Response Structure:**

```json
{
  "error": {
    "message": "The user with the specified ID does not exist.",
    "details": {
      "field1": "field value",
      "userId": "12345"
    }
  }
}
```

**Fields:**

- `error`: An object containing error details.
  - `message` (string): A human-readable description of the error, intended to
    be displayed to the end user.
  - `details` (object, optional): An object containing additional information
    relevant to the error.
    - Additional fields: Any other relevant data for debugging or error
      handling.

**Usage:**

- **Client-Side Handling**: Clients can parse the `error` object to determine
  the cause of the failure and implement appropriate error handling strategies.
- **User Feedback**: The `message` field is intended for display to the end
  user, providing a clear explanation of what went wrong.
- **Versatility**: The `details` field allows servers to provide additional
  context, aiding in debugging and error handling without exposing sensitive
  information to the end user.

---

## Supported Languages

UFO RPC initially supports code generation for:

- **TypeScript**
- **Go**

Support for additional languages is planned for future releases, enabling
broader adoption and integration into diverse technology stacks.

---

## Example Schema

Below is a comprehensive example of a UFO RPC schema incorporating all the
discussed features.

```json
{
  "types": [
    {
      "name": "User",
      "desc": "Represents a user in the system.",
      "fields": {
        "id": "string",
        "username": "string",
        "email": "string",
        "roles": {
          "type": "string[]",
          "desc": "List of roles assigned to the user."
        },
        "profile": {
          "type": "Profile",
          "desc": "User profile information."
        }
      }
    },
    {
      "name": "Address",
      "desc": "Represents a physical address.",
      "fields": {
        "street": {
          "type": "string",
          "desc": "Street name."
        },
        "city": {
          "type": "string",
          "desc": "City name."
        },
        "zipCode": {
          "type": "string",
          "desc": "Postal code."
        }
      }
    },
    {
      "name": "Profile",
      "desc": "User profile information.",
      "fields": {
        "age": {
          "type": "int",
          "desc": "Age of the user."
        },
        "address": {
          "type": "Address",
          "desc": "Address details."
        }
      }
    }
  ],
  "procedures": [
    {
      "name": "GetUser",
      "type": "query",
      "desc": "Retrieves user information by ID.",
      "input": {
        "userId": {
          "type": "string",
          "desc": "The user's unique identifier."
        }
      },
      "output": {
        "user": {
          "type": "User",
          "desc": "The user details."
        }
      }
    },
    {
      "name": "CreateUser",
      "type": "mutation",
      "desc": "Creates a new user in the system.",
      "meta": {
        "requiresAuth": true,
        "logLevel": "debug"
      },
      "input": {
        "user": {
          "type": "User",
          "desc": "User information to create."
        }
      },
      "output": {
        "userId": {
          "type": "string",
          "desc": "Unique identifier of the created user."
        }
      }
    },
    {
      "name": "ListUsers",
      "type": "query",
      "desc": "Retrieves a paginated list of users.",
      "input": {
        "page": {
          "type": "int",
          "desc": "Page number."
        },
        "pageSize": {
          "type": "int",
          "desc": "Number of users per page."
        }
      },
      "output": {
        "users": {
          "type": "User[]",
          "desc": "Array of user objects."
        },
        "totalCount": {
          "type": "int",
          "desc": "Total number of users."
        }
      }
    },
    {
      "name": "DeleteUser",
      "type": "mutation",
      "desc": "Deletes a user from the system.",
      "meta": {
        "requiresAuth": true,
        "logLevel": "info"
      },
      "input": {
        "userId": {
          "type": "string",
          "desc": "Unique identifier of the user to delete."
        }
      },
      "output": {
        "success": {
          "type": "boolean",
          "desc": "Indicates if the deletion was successful."
        }
      }
    },
    {
      "name": "UpdateUserProfile",
      "type": "mutation",
      "desc": "Updates the user's profile information.",
      "input": {
        "userId": {
          "type": "string",
          "desc": "The user's unique identifier."
        },
        "profile": {
          "type": "object",
          "desc": "Profile information to update.",
          "fields": {
            "age": {
              "type": "int",
              "desc": "Age of the user."
            },
            "address": {
              "type": "object",
              "desc": "Address details.",
              "fields": {
                "street": {
                  "type": "string",
                  "desc": "Street name."
                },
                "city": {
                  "type": "string",
                  "desc": "City name."
                },
                "zipCode": {
                  "type": "string",
                  "desc": "Postal code."
                }
              }
            }
          }
        }
      },
      "output": {
        "success": {
          "type": "boolean",
          "desc": "Indicates if the update was successful."
        }
      }
    }
  ]
}
```

In this example:

- **Primitive Types**: Includes `string`, `int`, `float`, and `boolean`.
- **Arrays and Nested Objects**: Demonstrates the use of arrays
  (`roles: string[]`) and nested objects with `"type": "object"` and `"fields"`.
- **Custom Types**: `Address` and `Profile` are custom types used within other
  types and procedures.
- **Input and Output Structure**:
  - Inputs and outputs can be objects with fields or assigned directly to a
    type.
  - Fields in `input` and `output` are treated as anonymous custom types with
    exactly the same fields.
- **Descriptions (`desc`)**: Added to types, procedures, and fields to enhance
  documentation and understanding.
- **Procedure Metadata**: Uses the `meta` field for additional procedure
  metadata, which are useful when using the code generated by UFO RPC.

---

## Conclusion

UFO RPC provides a simple yet powerful framework for defining and implementing
RPC APIs. By focusing on a minimalistic schema definition language and
leveraging standard HTTP methods along with JSON for data transport, UFO RPC
ensures ease of use, type safety, and cross-language compatibility.

The framework's initial support for TypeScript and Go, combined with its
extensible design, positions UFO RPC as a valuable tool for developers seeking
to streamline API development and focus on building their products without
unnecessary complexity.

---

**Note**: This specification is designed to be extensible and may evolve over
time. Future versions may introduce additional features such as advanced type
definitions, enhanced error handling mechanisms, or integration with
authentication frameworks, while maintaining backward compatibility and the core
principle of simplicity.

---

## Appendix

### Implementing Arrays and Nested Objects

#### Arrays

Arrays are denoted by appending `[]` to the type of the elements they contain.

**Syntax:**

```json
{
  "fieldName": {
    "type": "elementType[]",
    "desc": "Description of the array field."
  }
}
```

**Examples:**

```json
{
  "roles": {
    "type": "string[]",
    "desc": "List of roles assigned to the user."
  },
  "userIds": {
    "type": "int[]",
    "desc": "Array of user IDs."
  },
  "scores": {
    "type": "float[]",
    "desc": "List of score values."
  },
  "users": {
    "type": "User[]",
    "desc": "Array of user objects."
  },
  "addresses": {
    "type": "object[]",
    "desc": "List of addresses.",
    "fields": {
      "street": {
        "type": "string",
        "desc": "Street name."
      },
      "city": {
        "type": "string",
        "desc": "City name."
      },
      "zipCode": {
        "type": "string",
        "desc": "Postal code."
      }
    }
  }
}
```

#### Nested Objects

Nested objects can be defined in two ways:

1. **Inline Definition with `"type": "object"`**:

   Allows defining objects directly within the parent type or procedure
   input/output. You must specify `"type": "object"` and include a `"fields"`
   section.

   **Example:**

   ```json
   {
     "profile": {
       "type": "object",
       "desc": "User profile information.",
       "fields": {
         "age": {
           "type": "int",
           "desc": "Age of the user."
         },
         "address": {
           "type": "object",
           "desc": "Address details.",
           "fields": {
             "street": {
               "type": "string",
               "desc": "Street name."
             },
             "city": {
               "type": "string",
               "desc": "City name."
             },
             "zipCode": {
               "type": "string",
               "desc": "Postal code."
             }
           }
         }
       }
     }
   }
   ```

2. **Reference to Custom Types**:

   Define the nested object as a custom type and reference it by name.

   **Example:**

   ```json
   {
     "types": [
       {
         "name": "Address",
         "desc": "Represents a physical address.",
         "fields": {
           "street": {
             "type": "string",
             "desc": "Street name."
           },
           "city": {
             "type": "string",
             "desc": "City name."
           },
           "zipCode": {
             "type": "string",
             "desc": "Postal code."
           }
         }
       },
       {
         "name": "Profile",
         "desc": "User profile information.",
         "fields": {
           "age": {
             "type": "int",
             "desc": "Age of the user."
           },
           "address": {
             "type": "Address",
             "desc": "Address details."
           }
         }
       }
     ],
     "profile": {
       "type": "Profile",
       "desc": "User profile information."
     }
   }
   ```

### Handling Custom Types

Custom types are used directly by assigning them as the type of a field, just
like primitive types. Inputs and outputs can be assigned directly to a type by
specifying the type as a string.

**Defining Custom Types:**

```json
{
  "types": [
    {
      "name": "TypeName",
      "desc": "Description of the type.",
      "fields": {
        "fieldName": {
          "type": "fieldType",
          "desc": "Description of the field."
        }
      }
    }
  ]
}
```

**Using Custom Types:**

```json
{
  "fields": {
    "profile": {
      "type": "Profile",
      "desc": "User profile information."
    }
  },
  "input": {
    "user": {
      "type": "User",
      "desc": "User information."
    }
  },
  "input": "User"
}
```

**Fields in Input/Output as Anonymous Custom Types:**

- When you define fields directly under `input` or `output`, they are treated as
  an **anonymous custom type** with exactly the same fields.
- This means they behave like a custom type but without a specific name.

**Example:**

```json
{
  "input": {
    "userId": {
      "type": "string",
      "desc": "Unique identifier of the user."
    },
    "profile": {
      "type": "object",
      "desc": "Profile information to update.",
      "fields": {
        "age": {
          "type": "int",
          "desc": "Age of the user."
        },
        "address": {
          "type": "object",
          "desc": "Address details.",
          "fields": {
            "street": {
              "type": "string",
              "desc": "Street name."
            },
            "city": {
              "type": "string",
              "desc": "City name."
            },
            "zipCode": {
              "type": "string",
              "desc": "Postal code."
            }
          }
        }
      }
    }
  }
}
```

By allowing inputs and outputs to be assigned directly to a type or defined as
anonymous custom types with fields, the schema remains flexible, concise, and
easy to understand.

---

**Data Transport Format:**

All data is serialized and transported in **JSON** format. This includes:

- **Procedure Inputs**: Parameters sent from the client to the server.
- **Procedure Outputs**: Responses sent from the server to the client.
- **Error Responses**: Standardized error messages in JSON format.

Using JSON ensures compatibility across different programming languages and
simplifies the development process.
