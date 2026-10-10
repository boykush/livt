---
# livt:automates livt://mapping/file-story-issues-to-impl-repos/rule/R-02/example/EX-03
type: regex
pattern: '\nissue: https://board\.acme\.example/cards/42\n'
target: { source: file, path: stories/pay-by-invoice.md }
---
