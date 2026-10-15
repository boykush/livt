---
type: regex
pattern: '- id: R-01\n(?:(?!\n  - id: )[\s\S])*issues:\s*\n\s*- \S*acme/shop-backend/issues/101'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
