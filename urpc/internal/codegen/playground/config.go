package playground

import "fmt"

// Config is the configuration for the playground generator.
type Config struct {
	// OutputDir is the directory to output the generated playground to.
	OutputDir string `toml:"output_dir"`
}

func (c Config) Validate() error {
	if c.OutputDir == "" {
		return fmt.Errorf(`"output_dir" is required`)
	}
	return nil
}
