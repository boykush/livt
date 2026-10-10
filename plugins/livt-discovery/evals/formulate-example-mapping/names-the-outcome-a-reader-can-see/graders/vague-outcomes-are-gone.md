---
type: regex
pattern: 'name: not shipped, no due date$|name: 12,000 EUR with a purchase order number, accepted$|name: 9,000 EUR, no purchase order number needed$'
flags: m
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
