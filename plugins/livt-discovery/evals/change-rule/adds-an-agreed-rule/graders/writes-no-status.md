---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-01/example/EX-02
type: regex
pattern: '^(?![\s\S]*status:)[\s\S]*- id: R-02'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
