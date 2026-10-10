---
type: regex
pattern: 'name: not shipped, no due date$'
flags: m
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
