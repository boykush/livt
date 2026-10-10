---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-02/example/EX-01
type: regex
pattern: '- id: R-02(?=[\s\S]*status: proposed)(?=[\s\S]*purchase order number is refused)(?=[\s\S]*needs no)'
flags: i
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
