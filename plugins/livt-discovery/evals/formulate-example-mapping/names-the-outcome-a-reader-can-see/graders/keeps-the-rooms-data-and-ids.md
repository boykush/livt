---
type: regex
pattern: '- id: R-01\n    name: An invoice is due 30 days after the order ships[\s\S]*- id: EX-01\n        name: An order shipped on 1 March is due on 31 March[\s\S]*- id: EX-02[\s\S]*- id: R-02\n    name: An invoice over 10,000 EUR needs a purchase order number[\s\S]*- id: EX-01\n        name: A 12,000 EUR invoice without a purchase order number is refused[\s\S]*- id: EX-02\n        name: [^\n]*12,000 EUR[^\n]*purchase order number[^\n]*[\s\S]*- id: EX-03\n        name: [^\n]*9,000 EUR[^\n]*'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
