---
type: regex
pattern: '- id: R-02(?=[\s\S]*?status: rejected)(?=[\s\S]*?superseded_by:[\s\S]*?/rule/R-03)[\s\S]*- id: R-03(?=[\s\S]*status: proposed)'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
