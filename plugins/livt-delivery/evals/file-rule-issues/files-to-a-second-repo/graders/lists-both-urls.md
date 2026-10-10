---
# livt:automates livt://mapping/file-automation-issues-to-impl-repos/rule/R-08/example/EX-01
type: regex
pattern: 'issues:\s*\n\s*- \S*acme/shop-backend/issues/7\s*\n\s*- \S*acme/shop-web/issues/101'
target: { source: file, path: discoveries/example-mappings/pay-by-invoice.yaml }
---
