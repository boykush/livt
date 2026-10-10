---
max_turns: 40
timeout_seconds: 600
allowed_tools: [Read, Glob, Grep, Skill, Edit, Write, Bash]
---

Working through the pay-by-invoice story on my own I found a rule the others have not seen yet: an invoice over 10,000 EUR needs a purchase order number. A 12,000 EUR invoice without a purchase order number is refused, the same invoice with one is accepted, and a 9,000 EUR invoice needs none.

And a second one: an invoice is refused for a customer with an invoice more than 60 days overdue. A customer 61 days overdue is refused, one 59 days overdue is not.

Put both on the mapping as proposals. Alice (@alice) decides both.
