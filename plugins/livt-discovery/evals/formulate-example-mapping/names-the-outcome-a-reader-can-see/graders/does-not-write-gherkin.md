---
type: regex
pattern: '^\s*(Given|When|Then|Scenario|Feature):?\s'
flags: m
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
