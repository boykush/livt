---
# livt:automates livt://mapping/propose-rule-before-agreement/rule/R-05/example/EX-02
type: regex
pattern: '- id: R-02(?=[\s\S]*status: rejected)(?=[\s\S]*purchase order number is refused)'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
