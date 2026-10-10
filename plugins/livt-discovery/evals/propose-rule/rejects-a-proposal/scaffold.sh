#!/usr/bin/env bash
# A livt repository with one story and its mapping, committed on main: the
# skills read the default branch's copy to tell a draft from a filed rule.
set -euo pipefail
mkdir -p stories discoveries/example-mappings
cat > stories/pay-by-invoice.md <<'EOF'
---
name: Pay by invoice
---

As a business customer
I want to pay an order by invoice
So that I can settle it inside my company's payment process
EOF
cat > discoveries/example-mappings/pay-by-invoice.yaml <<'EOF'
rules:
  - id: R-01
    name: An invoice is due 30 days after the order ships
    examples:
      - id: EX-01
        name: An order shipped on 1 March is due on 31 March
      - id: EX-02
        name: An order not yet shipped has no due date
  - id: R-02
    name: An invoice over 10,000 EUR needs a purchase order number
    status: proposed
    decision_makers:
      - "@alice"
      - "@bob"
    examples:
      - id: EX-01
        name: A 12,000 EUR invoice without a purchase order number is refused
      - id: EX-02
        name: A 12,000 EUR invoice with a purchase order number is accepted
      - id: EX-03
        name: A 9,000 EUR invoice needs no purchase order number
EOF
git init --quiet --initial-branch=main
git config user.name 'Eval Fixture'
git config user.email 'eval@example.com'
git add --all
git commit --quiet --message 'docs: record pay-by-invoice'
