package golang

import (
	"strings"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
)

var packageHeader = strings.TrimSpace(`
// Code generated by UFO RPC. DO NOT EDIT.
// If you edit this file, it will be overwritten the next time it is generated.
// 
// This file is licensed under the MIT License.
// See https://github.com/uforg/uforpc for more information.
//
// Copyright (c) [Generated by UFO RPC - User retains copyright]
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//nolint:all
`)

func generatePackage(_ schema.Schema, config Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line(packageHeader)
	g.Break()

	g.Linef("// Package %s contains the generated code for the UFO RPC schema", config.PackageName)
	g.Linef("package %s", config.PackageName)
	g.Break()

	imports := []string{
		"bufio",
		"bytes",
		"context",
		"encoding/json",
		"fmt",
		"io",
		"net/http",
		"slices",
		"strings",
		"sync",
		"time",
	}

	g.Line("import (")
	g.Block(func() {
		for _, imp := range imports {
			g.Linef(`"%s"`, imp)
		}
	})
	g.Line(")")
	g.Break()

	return g.String(), nil
}
