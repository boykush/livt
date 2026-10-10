# Turns what the forge says about a pull request into the facts livt reads.
# Reviews that only comment take no side, so the latest one that did is what
# a person stands on.
($pr[0]) as $p
| ($reviews[0] // []) as $r
| {
    change: $p.html_url,
    author: (if $p.user.type == "Bot" then "" else "@" + $p.user.login end),
    approved_by: (
      $r
      | map(select(.user != null and (.state == "APPROVED" or .state == "CHANGES_REQUESTED" or .state == "DISMISSED")))
      | group_by(.user.login)
      | map(max_by(.submitted_at))
      | map(select(.state == "APPROVED") | "@" + .user.login)
      | sort
    )
  }
