package analyzer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockFileProvider is a mock implementation of FileProvider for testing
type mockFileProvider struct {
	files map[string]string
}

func (m *mockFileProvider) GetFileAndHash(relativeTo string, path string) (string, string, error) {
	// Try with the path as is
	if content, ok := m.files[path]; ok {
		return content, "mock-hash", nil
	}

	// If relativeTo is provided, try with the path relative to it
	if relativeTo != "" {
		relativePath := relativeTo + "/" + path
		if content, ok := m.files[relativePath]; ok {
			return content, "mock-hash", nil
		}
	}

	return "", "", os.ErrNotExist
}

func TestResolver(t *testing.T) {
	t.Run("Schema without version statement", func(t *testing.T) {
		// Create a mock file provider with a schema missing version statement
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						type User {
							id: string
							name: string
						}

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for missing version statement
		require.Error(t, err)
		require.Contains(t, err.Error(), "the first statement must be a version statement")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "the first statement must be a version statement")

		// Verify that no schema was processed
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDecls)
		require.Empty(t, combined.ProcDecls)
	})

	t.Run("Schema with unsupported version", func(t *testing.T) {
		// Create a mock file provider with a schema having unsupported version
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 2

						type User {
							id: string
							name: string
						}

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for unsupported version
		require.Error(t, err)
		require.Contains(t, err.Error(), "at the moment, the only supported version is 1")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "at the moment, the only supported version is 1")

		// Verify that no schema was processed
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDecls)
		require.Empty(t, combined.ProcDecls)
	})

	t.Run("Schema with multiple version statements", func(t *testing.T) {
		// Create a mock file provider with a schema having multiple version statements
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						type User {
							id: string
							name: string
						}

						version 1

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for multiple version statements
		require.Error(t, err)
		require.Contains(t, err.Error(), "version statement already declared for this schema")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "version statement already declared for this schema")

		// Verify that no schema was processed
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDecls)
		require.Empty(t, combined.ProcDecls)
	})

	t.Run("Schema with imported file missing version", func(t *testing.T) {
		// Create a mock file provider with a main file that imports a file without version
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						import "/types.urpc"

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
				"/types.urpc": `
						// Missing version statement
						type User {
							id: string
							name: string
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for missing version in imported file
		require.Error(t, err)
		require.Contains(t, err.Error(), "the first statement must be a version statement")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "the first statement must be a version statement")

		// Verify that the main schema was processed but not the imported one
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDecls)  // No types should be declared
		require.Len(t, combined.ProcDecls, 1) // Only the proc from main file
	})

	t.Run("Schema with external markdown", func(t *testing.T) {
		// Create a mock file provider with a schema that references an external markdown
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						""" docs/user.md """

						type User {
							id: string
							name: string
						}

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
				"/main.urpc/docs/user.md": "# User Documentation\n\nThis is the documentation for the User type.\n",
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)
		require.Empty(t, diagnostics)

		// Verify that the docstring was resolved
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.Schema.GetDocstrings(), 1)
		require.Equal(t, "# User Documentation\n\nThis is the documentation for the User type.\n", combined.Schema.GetDocstrings()[0].Value)
	})

	t.Run("Schema with missing external markdown", func(t *testing.T) {
		// Create a mock file provider with a schema that references a non-existent external markdown
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						""" docs/missing.md """

						type User {
							id: string
							name: string
						}

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for missing external markdown
		require.Error(t, err)
		require.Contains(t, err.Error(), "external markdown file not found")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "external markdown file not found")

		// Verify that the schema was still processed
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.TypeDecls, 1)
		require.Len(t, combined.ProcDecls, 1)
	})

	t.Run("Schema with external markdowns in different nodes", func(t *testing.T) {
		// Create a mock file provider with a schema that has external markdowns in different types of nodes
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						""" docs/overview.md """

						"""docs/rule.md"""
						rule @validEmail {
							for: string
							error: "Invalid email format"
						}

						"""docs/type.md"""
						type User {
							id: string
							email: string
								@validEmail
						}

						""" docs/proc.md """
						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
				"/main.urpc/docs/overview.md": "# API Overview\n\nThis is the main API documentation.\n",
				"/main.urpc/docs/rule.md":     "# Email Validation Rule\n\nValidates email format.\n",
				"/main.urpc/docs/type.md":     "# User Type\n\nRepresents a user in the system.\n",
				"/main.urpc/docs/proc.md":     "# GetUser Procedure\n\nRetrieves a user by ID.\n",
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)
		require.Empty(t, diagnostics)

		// Verify that all docstrings were resolved
		require.NotNil(t, combined.Schema)

		// Check standalone docstring
		require.Len(t, combined.Schema.GetDocstrings(), 1)
		require.Equal(t, "# API Overview\n\nThis is the main API documentation.\n", combined.Schema.GetDocstrings()[0].Value)

		// Check rule docstring
		require.Len(t, combined.Schema.GetRules(), 1)
		require.Equal(t, "# Email Validation Rule\n\nValidates email format.\n", combined.Schema.GetRules()[0].Docstring.Value)

		// Check type docstring
		require.Len(t, combined.Schema.GetTypes(), 1)
		require.Equal(t, "# User Type\n\nRepresents a user in the system.\n", combined.Schema.GetTypes()[0].Docstring.Value)

		// Check proc docstring
		require.Len(t, combined.Schema.GetProcs(), 1)
		require.Equal(t, "# GetUser Procedure\n\nRetrieves a user by ID.\n", combined.Schema.GetProcs()[0].Docstring.Value)
	})

	t.Run("Schema with imported file having unsupported version", func(t *testing.T) {
		// Create a mock file provider with a main file that imports a file with unsupported version
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
						version 1

						import "/types.urpc"

						proc GetUser {
							input {
								id: string
							}
							output {
								user: User
							}
						}
					`,
				"/types.urpc": `
						version 2

						type User {
							id: string
							name: string
						}
					`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for unsupported version in imported file
		require.Error(t, err)
		require.Contains(t, err.Error(), "at the moment, the only supported version is 1")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "at the moment, the only supported version is 1")

		// Verify that the main schema was processed but not the imported one
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDecls)  // No types should be declared
		require.Len(t, combined.ProcDecls, 1) // Only the proc from main file
	})

	t.Run("Basic schema with no imports", func(t *testing.T) {
		// Create a mock file provider with a single file
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
					version 1

					type User {
						id: string
						name: string
					}

					proc GetUser {
						input {
							id: string
						}
						output {
							user: User
						}
					}
				`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)
		require.Empty(t, diagnostics)

		// Verify the combined schema
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.Schema.GetTypes(), 1)
		require.Len(t, combined.Schema.GetProcs(), 1)
		require.Equal(t, "User", combined.Schema.GetTypes()[0].Name)
		require.Equal(t, "GetUser", combined.Schema.GetProcs()[0].Name)

		// Verify the declaration maps
		require.Len(t, combined.TypeDecls, 1)
		require.Len(t, combined.ProcDecls, 1)
		require.Contains(t, combined.TypeDecls, "User")
		require.Contains(t, combined.ProcDecls, "GetUser")
	})

	t.Run("Schema with imports", func(t *testing.T) {
		// Create a mock file provider with multiple files
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
					version 1

					import "/types.urpc"

					proc GetUser {
						input {
							id: string
						}
						output {
							user: User
						}
					}
				`,
				"/types.urpc": `
					version 1

					type User {
						id: string
						name: string
						profile: Profile
					}

					import "/profile.urpc"
				`,
				"/profile.urpc": `
					version 1

					type Profile {
						bio: string
						avatar: string
					}

					rule @validBio {
						for: string
						error: "Invalid bio"
					}
				`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)
		require.Empty(t, diagnostics)

		// Verify the combined schema
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.Schema.GetTypes(), 2)
		require.Len(t, combined.Schema.GetProcs(), 1)
		require.Len(t, combined.Schema.GetRules(), 1)

		// Verify the declaration maps
		require.Len(t, combined.TypeDecls, 2)
		require.Len(t, combined.ProcDecls, 1)
		require.Len(t, combined.RuleDecls, 1)
		require.Contains(t, combined.TypeDecls, "User")
		require.Contains(t, combined.TypeDecls, "Profile")
		require.Contains(t, combined.ProcDecls, "GetUser")
		require.Contains(t, combined.RuleDecls, "validBio")
	})

	t.Run("Circular imports allowed but resolved only once", func(t *testing.T) {
		// Create a mock file provider with circular imports
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
					version 1

					import "/types.urpc"

					proc GetUser {
						input {
							id: string
						}
						output {
							user: User
						}
					}
				`,
				"/types.urpc": `
					version 1

					import "/profile.urpc"

					type User {
						id: string
						name: string
						profile: Profile
					}
				`,
				"/profile.urpc": `
					version 1

					import "/main.urpc"

					type Profile {
						bio: string
						avatar: string
					}
				`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, _, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)

		// Verify that the schema was combined despite the circular import
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.TypeDecls, 2)
		require.Len(t, combined.ProcDecls, 1)
		require.Contains(t, combined.TypeDecls, "User")
		require.Contains(t, combined.TypeDecls, "Profile")
		require.Contains(t, combined.ProcDecls, "GetUser")
	})

	t.Run("Missing import file", func(t *testing.T) {
		// Create a mock file provider with a missing import
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
					version 1

					import "/missing.urpc"

					proc GetUser {
						input {
							id: string
						}
						output {
							user: User
						}
					}
				`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for missing import file
		require.Error(t, err)
		require.Contains(t, err.Error(), "error reading file")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "error reading file")

		// Verify that the main schema was still processed
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.ProcDecls, 1)
		require.Contains(t, combined.ProcDecls, "GetUser")
	})

	t.Run("Schema with deduplication", func(t *testing.T) {
		// Create a mock file provider with multiple files
		provider := &mockFileProvider{
			files: map[string]string{
				"/main.urpc": `
					version 1

					import "/schema_users.urpc"
					import "/schema_posts.urpc"
				`,
				"/schema_common.urpc": `
					version 1

					rule @uuid {
						for: string
						error: "Invalid UUID format"
					}

					rule @email {
						for: string
						error: "Invalid email format"
					}

					type WithId {
						id: string
							@uuid
					}

					proc Ping {
						input {
							id: string
						}
						output {
							success: boolean
						}
					}
				`,
				"/schema_users.urpc": `
					version 1

					import "/schema_common.urpc"

					type User extends WithId {
						name: string
						email: string
							@minlen(6) // a@a.aa
							@email
						age: int
							@min(18)
						isActive: boolean
						createdAt: datetime
						updatedAt: datetime
					}
				`,
				"/schema_posts.urpc": `
					version 1

					import "/schema_common.urpc"

					type Post extends WithId {
						title: string
						content: string
						authorId: string
						createdAt: datetime
						updatedAt: datetime
					}
				`,
			},
		}

		// Create an analyzer with the mock provider
		analyzer, err := NewAnalyzer(provider)
		require.NoError(t, err)

		// Analyze the schema
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		require.NoError(t, err)
		require.Empty(t, diagnostics)

		// Verify the combined schema
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.Schema.GetTypes(), 3)
		require.Len(t, combined.Schema.GetProcs(), 1)
		require.Len(t, combined.Schema.GetRules(), 2)

		// Verify the declaration maps
		require.Len(t, combined.TypeDecls, 3)
		require.Len(t, combined.ProcDecls, 1)
		require.Len(t, combined.RuleDecls, 2)
		require.Contains(t, combined.RuleDecls, "uuid")
		require.Contains(t, combined.RuleDecls, "email")
		require.Contains(t, combined.TypeDecls, "WithId")
		require.Contains(t, combined.TypeDecls, "User")
		require.Contains(t, combined.TypeDecls, "Post")
		require.Contains(t, combined.ProcDecls, "Ping")
	})
}
