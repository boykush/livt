---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-02/example/EX-02
type: regex
pattern: '^rules:\n  - id: R-01\n    name: An invoice is due 30 days after the order ships\n    examples:\n      - id: EX-01\n        name: An order shipped on 1 March is due on 31 March\n      - id: EX-02\n        name: An order not yet shipped has no due date\n  - id: R-02\n    name: An invoice over 10,000 EUR needs a purchase order number\n    status: accepted\n    decision_makers:\n      - "@alice"\n      - "@bob"\n    examples:\n      - id: EX-01\n        name: A 12,000 EUR invoice without a purchase order number is refused\n      - id: EX-02\n        name: A 12,000 EUR invoice with a purchase order number is accepted\n      - id: EX-03\n        name: A 9,000 EUR invoice needs no purchase order number\n\s*$'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
