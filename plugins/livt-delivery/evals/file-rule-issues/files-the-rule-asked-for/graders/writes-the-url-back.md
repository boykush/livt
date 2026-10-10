---
# livt:automates livt://mapping/file-automation-issues-to-impl-repos/rule/R-02/example/EX-01
type: regex
pattern: '- id: R-02\n(?:(?!  - id: )[\s\S])*issues:\s*\n\s*- \S*acme/shop-backend/issues/101'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
