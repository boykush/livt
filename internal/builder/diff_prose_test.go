package builder

import (
	"strings"
	"testing"

	"github.com/boykush/livt/internal/diff"
	"github.com/boykush/livt/internal/i18n"
)

func prose(op diff.Op, value string, parts ...diff.Part) diff.Line {
	return diff.Line{Op: op, Field: diff.Field{Value: value, Prose: true}, Parts: parts}
}

func TestDiffProse(t *testing.T) {
	ctx := func(v string) diff.Line { return prose(diff.OpContext, v) }
	tests := []struct {
		name    string
		lines   []diff.Line
		want    []string
		notWant []string
	}{
		{
			name:  "markdown renders as html",
			lines: []diff.Line{ctx("## 背景"), ctx(""), ctx("- **強調**")},
			want:  []string{"<h2>背景</h2>", "<li><strong>強調</strong></li>"},
		},
		{
			name:    "only the changed line of a paragraph is marked",
			lines:   []diff.Line{ctx("一行目"), prose(diff.OpAdd, "二行目")},
			want:    []string{"一行目<br>", `<span class="diff-add">二行目</span>`},
			notWant: []string{`<p class="`},
		},
		{
			name: "a reworded line marks only the words that moved",
			lines: []diff.Line{prose(diff.OpDel, "目的の成果",
				diff.Part{Text: "目的の"}, diff.Part{Changed: true, Text: "成果"})},
			want: []string{`<span class="diff-del reworded">目的の<span class="moved">成果</span></span>`},
		},
		{
			name:  "a list item keeps its marker when its text changes",
			lines: []diff.Line{ctx("1. 一つ目"), prose(diff.OpAdd, "2. 二つ目")},
			want:  []string{`<li><span class="diff-add">二つ目</span></li>`},
		},
		{
			name:  "a code span is marked whole",
			lines: []diff.Line{prose(diff.OpAdd, "`livt://` で引ける")},
			want:  []string{`<span class="diff-add"><code>livt://</code></span>`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string((&Builder{}).diffProse(tt.lines, ""))
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("got %q, want it to contain %q", got, w)
				}
			}
			for _, w := range tt.notWant {
				if strings.Contains(got, w) {
					t.Errorf("got %q, want it not to contain %q", got, w)
				}
			}
		})
	}
}

func TestDiffProseFolds(t *testing.T) {
	long := func() []diff.Line {
		lines := []diff.Line{prose(diff.OpAdd, "変わった段落"), prose(diff.OpContext, "")}
		for range foldMin {
			lines = append(lines, prose(diff.OpContext, "- 変わらない項目"))
		}
		return lines
	}
	b := &Builder{Lang: i18n.Lang("ja")}

	t.Run("a long unchanged run folds", func(t *testing.T) {
		got := string(b.diffProse(long(), ""))
		if !strings.Contains(got, `<details class="diff-fold"><summary>変更のない`) {
			t.Errorf("got %q, want the unchanged list folded", got)
		}
		if strings.Index(got, "変わった段落") > strings.Index(got, "<details") {
			t.Errorf("got %q, want the change outside the fold", got)
		}
	})

	t.Run("a short unchanged run stays open", func(t *testing.T) {
		got := string(b.diffProse(long()[:foldMin], ""))
		if strings.Contains(got, "<details") {
			t.Errorf("got %q, want no fold", got)
		}
	})

	t.Run("a body with no change stays open", func(t *testing.T) {
		lines := long()
		lines[0].Op = diff.OpContext
		if got := string(b.diffProse(lines, "")); strings.Contains(got, "<details") {
			t.Errorf("got %q, want an unchanged body read whole", got)
		}
	})
}
