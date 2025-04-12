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
	if content, ok := m.files[path]; ok {
		return content, "mock-hash", nil
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
		require.Empty(t, combined.TypeDefs)
		require.Empty(t, combined.ProcDefs)
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
		require.Empty(t, combined.TypeDefs)
		require.Empty(t, combined.ProcDefs)
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
		require.Contains(t, err.Error(), "version statement already defined for this schema")

		// Verify that error diagnostics were generated
		require.NotEmpty(t, diagnostics)
		require.Contains(t, diagnostics[0].Message, "version statement already defined for this schema")

		// Verify that no schema was processed
		require.NotNil(t, combined.Schema)
		require.Empty(t, combined.TypeDefs)
		require.Empty(t, combined.ProcDefs)
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
		require.Empty(t, combined.TypeDefs)  // No types should be defined
		require.Len(t, combined.ProcDefs, 1) // Only the proc from main file
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
		require.Empty(t, combined.TypeDefs)  // No types should be defined
		require.Len(t, combined.ProcDefs, 1) // Only the proc from main file
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

		// Verify the definition maps
		require.Len(t, combined.TypeDefs, 1)
		require.Len(t, combined.ProcDefs, 1)
		require.Contains(t, combined.TypeDefs, "User")
		require.Contains(t, combined.ProcDefs, "GetUser")
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

		// Verify the definition maps
		require.Len(t, combined.TypeDefs, 2)
		require.Len(t, combined.ProcDefs, 1)
		require.Len(t, combined.RuleDefs, 1)
		require.Contains(t, combined.TypeDefs, "User")
		require.Contains(t, combined.TypeDefs, "Profile")
		require.Contains(t, combined.ProcDefs, "GetUser")
		require.Contains(t, combined.RuleDefs, "validBio")
	})

	t.Run("Circular imports", func(t *testing.T) {
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
		combined, diagnostics, err := analyzer.Analyze("/main.urpc")
		// Expect an error for circular imports
		require.Error(t, err)
		require.Contains(t, err.Error(), "circular import detected")

		// Verify that circular import diagnostics were generated
		require.NotEmpty(t, diagnostics, "Expected circular import diagnostics")
		require.Contains(t, diagnostics[0].Message, "circular import detected")

		// Verify that the schema was still combined despite the circular imports
		require.NotNil(t, combined.Schema)
		require.Len(t, combined.TypeDefs, 2)
		require.Contains(t, combined.TypeDefs, "User")
		require.Contains(t, combined.TypeDefs, "Profile")
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
		require.Len(t, combined.ProcDefs, 1)
		require.Contains(t, combined.ProcDefs, "GetUser")
	})
}
