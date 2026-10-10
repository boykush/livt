---
type: regex
pattern: 'decision_makers:\s*\S'
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
