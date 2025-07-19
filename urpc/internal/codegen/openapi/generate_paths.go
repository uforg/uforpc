package openapi

import (
	"github.com/uforg/uforpc/urpc/internal/schema"
)

func generatePaths(sch schema.Schema) (Paths, error) {
	paths := Paths{}

	for _, procNode := range sch.GetProcNodes() {
		name := procNode.Name
		// inputName := fmt.Sprintf("%sInput", name)
		// responseName := fmt.Sprintf("%sResponse", name)

		doc := ""
		if procNode.Doc != nil {
			doc = *procNode.Doc
		}

		paths["/"+name] = map[string]any{
			"post": map[string]any{
				"tags":        []string{"procedures"},
				"description": doc,
			},
		}
	}

	for _, streamNode := range sch.GetStreamNodes() {
		name := streamNode.Name
		// inputName := fmt.Sprintf("%sInput", name)
		// responseName := fmt.Sprintf("%sResponse", name)

		doc := ""
		if streamNode.Doc != nil {
			doc = *streamNode.Doc
		}

		paths["/"+name] = map[string]any{
			"post": map[string]any{
				"tags":        []string{"streams"},
				"description": doc,
			},
		}
	}

	return paths, nil
}
