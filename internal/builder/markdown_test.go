package builder

import (
	"strings"
	"testing"
)

func TestMarkdown(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{"emphasis", "**強調** と `code`", []string{"<strong>強調</strong>", "<code>code</code>"}},
		{"list", "- 一つ目\n- 二つ目", []string{"<ul>", "<li>一つ目</li>"}},
		{"link", "[livt](https://example.com)", []string{`<a href="https://example.com">livt</a>`}},
		{"single newline stays a break", "一行目\n二行目", []string{"一行目<br>"}},
		{"raw html is not passed through", "<script>alert(1)</script>", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(markdown(tt.src))
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("markdown(%q) = %q, want it to contain %q", tt.src, got, w)
				}
			}
			if strings.Contains(got, "<script>") {
				t.Errorf("markdown(%q) = %q, raw HTML leaked", tt.src, got)
			}
		})
	}
}
