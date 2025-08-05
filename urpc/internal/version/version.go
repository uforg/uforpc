package version

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const Version = "0.3.2"
const VersionWithPrefix = "v" + Version

// asciiArtRaw is used to generate AsciiArt
var asciiArtRaw = strings.Join([]string{
	"╦ ╦╔═╗╔═╗  ╦═╗╔═╗╔═╗",
	"║ ║╠╣ ║ ║  ╠╦╝╠═╝║  ",
	"╚═╝╚  ╚═╝  ╩╚═╩  ╚═╝",
}, "\n")

// basicInfoRaw is used to generate AsciiArt
var basicInfoRaw = strings.Join([]string{
	"Star the repo: https://github.com/uforg/uforpc",
	"Show usage:    urpc --help",
	"Show version:  urpc --version",
}, "\n")

// AsciiArt is the ASCII art for the UFO RPC logo.
// It is generated dynamically to ensure that the logo is always
// centered and the lines are always the same width.
var AsciiArt = func() string {
	maxWidth := 0
	for line := range strings.SplitSeq(basicInfoRaw, "\n") {
		if utf8.RuneCountInString(line) > maxWidth {
			maxWidth = utf8.RuneCountInString(line)
		}
	}
	dashes := strings.Repeat("-", maxWidth)

	combined := strings.Join([]string{
		centerText(asciiArtRaw, maxWidth),
		centerText(VersionWithPrefix, maxWidth),
		"",
		basicInfoRaw,
	}, "\n")

	combinedWithLines := ""
	for line := range strings.SplitSeq(combined, "\n") {
		spaces := maxWidth - utf8.RuneCountInString(line)
		combinedWithLines += fmt.Sprintf("| %s%s |\n", line, strings.Repeat(" ", spaces))
	}

	lines := []string{
		"+-" + dashes + "-+",
		strings.TrimSpace(combinedWithLines),
		"+-" + dashes + "-+",
	}

	return strings.Join(lines, "\n")
}()

// centerText centers text within a given width.
// It handles both single and multi-line strings, treating multi-line
// strings as a single block. It prevents panics by not adding
// padding if the text exceeds the desired width.
func centerText(text string, desiredWidth int) string {
	lines := strings.Split(text, "\n")

	// Find the widest line to determine the block's width
	var longestLineWidth int
	for _, line := range lines {
		lineWidth := utf8.RuneCountInString(line)
		if lineWidth > longestLineWidth {
			longestLineWidth = lineWidth
		}
	}

	// Calculate the left padding for the entire block
	blockLeftPaddingCount := 0
	if longestLineWidth < desiredWidth {
		blockLeftPaddingCount = (desiredWidth - longestLineWidth) / 2
	}
	blockLeftPadding := strings.Repeat(" ", blockLeftPaddingCount)

	// Build the result line by line
	var resultBuilder strings.Builder
	for i, line := range lines {
		// Add block padding and the line itself
		resultBuilder.WriteString(blockLeftPadding)
		resultBuilder.WriteString(line)

		// Calculate and add right padding to fill the remaining space
		currentWidth := blockLeftPaddingCount + utf8.RuneCountInString(line)
		rightPaddingCount := desiredWidth - currentWidth

		// Only add padding if the count is positive
		if rightPaddingCount > 0 {
			resultBuilder.WriteString(strings.Repeat(" ", rightPaddingCount))
		}

		// Add a newline if it's not the last line
		if i < len(lines)-1 {
			resultBuilder.WriteString("\n")
		}
	}

	return resultBuilder.String()
}
