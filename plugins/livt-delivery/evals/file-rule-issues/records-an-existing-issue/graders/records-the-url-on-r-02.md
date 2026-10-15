---
# livt:automates livt://mapping/file-automation-issues-to-impl-repos/rule/R-09/example/EX-02
type: regex
pattern: '- id: R-02\n(?:(?!\n  - id: )[\s\S])*issues:\s*\n\s*- \S*acme/shop-backend/issues/55'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
