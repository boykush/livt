---
type: regex
pattern: '[0-9][0-9,.]*\s*(EUR|USD|GBP|JPY|€|\$|£|¥)|(EUR|USD|GBP|JPY|€|\$|£|¥)\s*[0-9]'
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
