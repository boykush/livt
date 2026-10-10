---
# livt:automates livt://mapping/file-automation-issues-to-impl-repos/rule/R-07/example/EX-02
type: regex
pattern: '\tcommit'
match: 'count:1'
target: { source: file, path: .git/logs/HEAD }
---
