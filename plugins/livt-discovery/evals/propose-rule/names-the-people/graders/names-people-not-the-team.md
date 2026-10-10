---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-03/example/EX-02
type: regex
pattern: '^(?![\s\S]*@acme/billing)[\s\S]*decision_makers:[\s\S]*@alice'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
