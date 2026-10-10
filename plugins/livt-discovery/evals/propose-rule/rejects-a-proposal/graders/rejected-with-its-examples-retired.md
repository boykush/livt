---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-02/example/EX-03
type: regex
pattern: '- id: R-02(?=[\s\S]*status: rejected)(?=(?:[\s\S]*?retired: true){3})'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
