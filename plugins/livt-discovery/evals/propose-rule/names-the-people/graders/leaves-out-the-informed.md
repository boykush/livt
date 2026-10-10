---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-03/example/EX-04
type: regex
pattern: '@carol'
match: not_contains
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
