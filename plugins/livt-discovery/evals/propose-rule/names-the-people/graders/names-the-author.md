---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-03/example/EX-03
type: regex
pattern: 'decision_makers:[\s\S]*@dana'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
