#!/usr/bin/env bash
# A livt repository with one story and its mapping, committed on main, so a
# run has a revision to report and a clean tree to leave its write-back in.
set -euo pipefail
mkdir -p stories discoveries/example-mappings
cat > stories/pay-by-invoice.md <<'EOF'
---
name: Pay by invoice
repos:
  - acme/shop-backend
issues:
  - https://github.com/acme/shop-backend/issues/3
---

As a business customer
I want to pay an order by invoice
So that I can settle it inside my company's payment process
EOF
cat > discoveries/example-mappings/pay-by-invoice.yaml <<'EOF'
rules:
  - id: R-01
    name: An invoice is due 30 days after the order ships
    issues:
      - https://github.com/acme/shop-backend/issues/7
    examples:
      - id: EX-01
        name: An order shipped on 1 March is due on 31 March
      - id: EX-02
        name: An order not yet shipped has no due date
  - id: R-02
    name: An invoice over 10,000 EUR needs a purchase order number
    examples:
      - id: EX-01
        name: A 12,000 EUR invoice without a purchase order number is refused
      - id: EX-02
        name: A 12,000 EUR invoice from a trusted customer needs none
        retired: true
      - id: EX-03
        name: A 9,000 EUR invoice needs no purchase order number
  - id: R-03
    name: An invoice is refused for a customer more than 60 days overdue
    status: proposed
    examples:
      - id: EX-01
        name: A customer 61 days overdue is refused
EOF
git init --quiet --initial-branch=main
git config user.name 'Eval Fixture'
git config user.email 'eval@example.com'
git add --all
git commit --quiet --message 'docs: record pay-by-invoice'
