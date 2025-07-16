package embedplayground

import "embed"

//go:embed all:build/**
var BuildFS embed.FS
