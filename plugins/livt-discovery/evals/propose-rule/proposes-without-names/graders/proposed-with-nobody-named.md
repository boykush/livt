---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-03/example/EX-05
type: regex
pattern: '^(?![\s\S]*decision_makers)[\s\S]*status: proposed'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
