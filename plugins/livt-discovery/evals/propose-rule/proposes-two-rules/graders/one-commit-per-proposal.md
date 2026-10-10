---
type: regex
pattern: '\tcommit: '
match: 'count:2'
target: { source: file, path: .git/logs/HEAD }
---
