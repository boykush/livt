---
# livt:automates livt://mapping/file-automation-issues-to-impl-repos/rule/R-03/example/EX-01
type: regex
pattern: 'livt://mapping/pay-by-invoice/rule/R-01[^/\w]'
match: not_contains
target: mock_calls
---
