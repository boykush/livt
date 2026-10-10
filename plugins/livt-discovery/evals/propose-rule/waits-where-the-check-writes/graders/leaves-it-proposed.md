---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-04/example/EX-04
type: regex
pattern: '^(?![\s\S]*status: accepted)[\s\S]*status: proposed'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
