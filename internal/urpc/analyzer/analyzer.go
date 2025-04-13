// Package analyzer provides semantic analysis for URPC schemas.
// This implementation prioritizes simplicity and maintainability over performance
// by performing a full analysis without caching results between calls.
package analyzer

// Analyzer manages the analysis process for URPC schemas without caching.
type Analyzer struct {
	fileProvider FileProvider
	resolver     *resolver
}

// NewAnalyzer creates a new cache-less Analyzer instance.
func NewAnalyzer(fileProvider FileProvider) (*Analyzer, error) {
	return &Analyzer{
		fileProvider: fileProvider,
		resolver:     newResolver(fileProvider),
	}, nil
}

// Analyze performs semantic analysis on a URPC schema starting from the given entry point.
// It resolves all imports, combines the schemas, and returns the combined schema along with
// any diagnostics encountered during the resolution and analysis phases.
//
// It consists of two phases:
//   - Resolution phase: Resolves all imports and combines the schemas.
//   - Semantic analysis phase: Performs semantic analysis on the combined schema.
func (a *Analyzer) Analyze(entryPointFilePath string) (CombinedSchema, []Diagnostic, error) {
	combinedSchema, resolverDiagnostics, _ := a.resolver.resolve(entryPointFilePath)

	semanalyzer := newSemanalyzer(combinedSchema)
	semanalyzerDiagnostics, _ := semanalyzer.analyze()

	diagnostics := []Diagnostic{}
	if len(resolverDiagnostics) > 0 {
		diagnostics = append(diagnostics, resolverDiagnostics...)
	}
	if len(semanalyzerDiagnostics) > 0 {
		diagnostics = append(diagnostics, semanalyzerDiagnostics...)
	}

	if len(diagnostics) > 0 {
		return combinedSchema, diagnostics, diagnostics[0]
	}
	return combinedSchema, nil, nil
}
