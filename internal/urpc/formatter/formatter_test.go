package formatter

import (
	"embed"
	"path"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/util/strutil"
)

//go:embed tests/*.urpc
var testFiles embed.FS

func TestFormatEmptySchema(t *testing.T) {
	input := ""
	expected := ""

	formatted, err := Format("schema.urpc", input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat(t *testing.T) {
	files, err := testFiles.ReadDir("tests")
	require.NoError(t, err)

	for _, file := range files {
		content, err := testFiles.ReadFile(path.Join("tests", file.Name()))
		require.NoError(t, err)

		separator := "\n// >>>>\n\n"
		input := strutil.GetStrBefore(string(content), separator)
		expected := strutil.GetStrAfter(string(content), separator)

		formatted, err := Format(file.Name(), input)
		require.NoError(t, err, "error formatting %s", file.Name())
		require.Equal(t, expected, formatted, "incorrect formatting for %s", file.Name())
	}
}
