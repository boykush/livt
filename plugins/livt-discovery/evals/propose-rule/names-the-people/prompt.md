---
max_turns: 40
timeout_seconds: 600
allowed_tools: [Read, Glob, Grep, Skill, Edit, Write, Bash]
---

Working through the pay-by-invoice story on my own I found a rule the others have not seen yet: an invoice over 10,000 EUR needs a purchase order number. A 12,000 EUR invoice without a purchase order number is refused, the same invoice with one is accepted, and a 9,000 EUR invoice needs none.

Put it on the mapping as a proposal. I'm @dana and it is mine to decide together with Alice (@alice) and the billing team (@acme/billing). Carol (@carol) only needs to hear about it.
