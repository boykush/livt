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
        name: An order not yet shipped has no due date
  - id: R-02
    name: A large invoice needs a purchase order number
    examples:
      - id: EX-01
        name: Over the threshold, no purchase order number, refused
      - id: EX-02
        name: Over the threshold with a purchase order number is accepted
EOS
git init --quiet --initial-branch=main
git config user.name 'Eval Fixture'
git config user.email 'eval@example.com'
git add --all
git commit --quiet --message 'docs: record pay-by-invoice'
