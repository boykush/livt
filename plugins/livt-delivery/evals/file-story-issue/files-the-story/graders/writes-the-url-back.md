---
# livt:automates livt://mapping/file-story-issues-to-impl-repos/rule/R-02/example/EX-01
type: regex
pattern: '\nissues:\s*\n\s*- \S*acme/shop-backend/issues/101'
target: { source: file, path: stories/pay-by-invoice.md }
---
