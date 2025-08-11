package dart

import (
	"fmt"
	"strings"
)

// Config is the configuration for the Dart code generator.
type Config struct {
	// OutputFile is the file to output the generated code to.
	OutputFile string `toml:"output_file"`
}

func (c Config) Validate() error {
	if c.OutputFile == "" {
		return fmt.Errorf(`"output_file" is required`)
	}
	if !strings.HasSuffix(c.OutputFile, ".dart") {
		return fmt.Errorf(`"output_file" must end with ".dart"`)
	}
	return nil
}
