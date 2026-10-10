---
# livt:automates livt://mapping/propose-rule-before-agreement/rule/R-05/example/EX-04
type: regex
pattern: 'status: accepted[\s\S]*- id: Q-01(?=[\s\S]*retired: true)(?=[\s\S]*superseded_by:[\s\S]*livt://mapping/pay-by-invoice/rule/R-02)'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
