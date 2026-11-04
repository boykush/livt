package diff

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/boykush/livt/internal/uri"
)

// mappingLines is the line each item of a mapping file is written on: the line
// of its own words — a name, or a question's text — since that is the line a
// review comment on it is about. The mapping itself answers to its name, or to
// the top of the file when it has none. A file that will not read gives no
// lines, and the items it would have placed simply have nowhere to point.
func mappingLines(path, mappingURI string) map[string]int {
	lines := map[string]int{mappingURI: 1}
	data, err := os.ReadFile(path)
	if err != nil {
		return lines
	}
	var doc yaml.Node
	if yaml.Unmarshal(data, &doc) != nil || len(doc.Content) == 0 {
		return lines
	}
	root := doc.Content[0]
	if n := valueOf(root, "name"); n != nil {
		lines[mappingURI] = n.Line
	}
	parsed, _ := uri.Parse(mappingURI)
	key := parsed.StoryKey
	for _, r := range items(valueOf(root, "rules")) {
		ruleID := scalar(valueOf(r, "id"))
		lines[uri.Rule(key, ruleID)] = lineOf(r, "name")
		for _, ex := range items(valueOf(r, "examples")) {
			lines[uri.Example(key, ruleID, scalar(valueOf(ex, "id")))] = lineOf(ex, "name")
		}
	}
	for _, q := range items(valueOf(root, "questions")) {
		lines[uri.Question(key, scalar(valueOf(q, "id")))] = lineOf(q, "text")
	}
	return lines
}

func valueOf(m *yaml.Node, key string) *yaml.Node {
	if m == nil || m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

func items(seq *yaml.Node) []*yaml.Node {
	if seq == nil || seq.Kind != yaml.SequenceNode {
		return nil
	}
	return seq.Content
}

func scalar(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	return n.Value
}

// lineOf is the line of the item's own words, or of the item itself when it
// has none written.
func lineOf(item *yaml.Node, key string) int {
	if n := valueOf(item, key); n != nil {
		return n.Line
	}
	return item.Line
}
