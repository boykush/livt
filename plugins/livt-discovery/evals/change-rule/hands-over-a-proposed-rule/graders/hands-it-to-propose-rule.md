---
# livt:automates livt://mapping/carry-proposals-to-agreement-with-skill/rule/R-01/example/EX-03
type: llm
---

PASS if the reply says R-02 is still a proposal and that changing it belongs to propose-rule, or the run went on to handle it through propose-rule.
FAIL if it changed the rule as an agreed one without mentioning that it is still proposed.
