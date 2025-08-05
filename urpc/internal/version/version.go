package version

import "strings"

const Version = "0.3.2"
const VersionWithPrefix = "v" + Version

var AsciiArt = strings.TrimSpace(`
+------------------------------------------------+
|              ╦ ╦╔═╗╔═╗  ╦═╗╔═╗╔═╗              |
|              ║ ║╠╣ ║ ║  ╠╦╝╠═╝║                |
|              ╚═╝╚  ╚═╝  ╩╚═╩  ╚═╝              |
| Star the repo: https://github.com/uforg/uforpc |
| Show usage:    urpc --help                     |
| Show version:  urpc --version                  |
+------------------------------------------------+
`)
