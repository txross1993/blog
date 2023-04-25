// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

type Tree struct {
	Nodes []*Node
}

type Node struct {
	Left, Right *Node
	Data        map[string]interface{}
}

func main() {
	var (
		left = &Node{
			Data: map[string]interface{}{
				"key1": "left1",
				"key2": "left2",
			},
		}
		right = &Node{
			Data: map[string]interface{}{
				"key1": "right1",
			},
		}
		root = &Node{Left: left, Right: right}
		tree = &Tree{Nodes: []*Node{root, left, right}}
	)

	tplStr := `First left key: {{ range $node := .Nodes -}}
    {{- range $key := $node.Data -}}
        {{- with $keyPrefix := (slice $key 0 4) -}}
            {{- if eq $keyPrefix "left" }} {{- $key }} {{ break }} {{ end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}`

	templateParsed, err := template.New("tree").Parse(tplStr)
	if err != nil {
		fmt.Println("failed to parse template: ", err)
		os.Exit(1)
	}

	dest := bytes.NewBuffer([]byte{})
	err = templateParsed.Execute(dest, tree)
	if err != nil {
		fmt.Println("failed to execute template: ", err)
		os.Exit(1)
	}

	fmt.Println(dest.String())
}
