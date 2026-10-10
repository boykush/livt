#!/usr/bin/env bash
# A livt repository with one story and its mapping, committed on main: the
# record is formulation's baseline, and the skill checks it is committed.
set -euo pipefail
mkdir -p stories discoveries/example-mappings
cat > stories/pay-by-invoice.md <<'EOS'
---
name: Pay by invoice
---

As a business customer
I want to pay an order by invoice
So that I can settle it inside my company's payment process
EOS
cat > discoveries/example-mappings/pay-by-invoice.yaml <<'EOS'
rules:
  - id: R-01
    name: An invoice is due 30 days after the order ships
    examples:
      - id: EX-01
        name: An order shipped on 1 March is due on 31 March
      - id: EX-02
        name: not shipped, no due date
  - id: R-02
    name: An invoice over 10,000 EUR needs a purchase order number
    examples:
      - id: EX-01
        name: A 12,000 EUR invoice without a purchase order number is refused
      - id: EX-02
        name: 12,000 EUR with a purchase order number, accepted
      - id: EX-03
        name: 9,000 EUR, no purchase order number needed
EOS
git init --quiet --initial-branch=main
git config user.name 'Eval Fixture'
git config user.email 'eval@example.com'
git add --all
git commit --quiet --message 'docs: record pay-by-invoice'
