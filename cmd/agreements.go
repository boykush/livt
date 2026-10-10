package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/agreement"
	"github.com/spf13/cobra"
)

var (
	agreementsRoot      string
	agreementsApprovals string
	agreementsWrite     bool
)

func init() {
	agreementsCheckCmd.Flags().StringVar(&agreementsRoot, "root", ".", "the livt repository to check")
	agreementsCheckCmd.Flags().StringVar(&agreementsApprovals, "approvals", "", "file saying who wrote and who approved the change (default: nobody)")
	agreementsCheckCmd.Flags().BoolVar(&agreementsWrite, "write", false, "record who agreed, and accept proposed rules everyone has agreed to")
	agreementsCmd.AddCommand(agreementsCheckCmd)
	rootCmd.AddCommand(agreementsCmd)
}

var agreementsCmd = &cobra.Command{
	Use:   "agreements",
	Short: "Hold accepted rules to the people they name",
}

var agreementsCheckCmd = &cobra.Command{
	Use:   "check <base> [head]",
	Short: "Check that every rule a change accepts has its decision makers' agreement",
	Long: `Check a change against the people its rules name.

A rule may name the people whose agreement it needs as decision_makers. A
change that makes such a rule accepted — by turning a proposal accepted, or by
adding the rule already accepted — passes only when every one of them stands
behind it: as the author of the change, as someone who approved it, or as
someone already on record in agreements/{story-key}.json.

Who wrote and who approved the change is not something this command can know.
It is read from --approvals, a file whatever speaks to your forge writes:

    {"change": "<url>", "author": "@alice", "approved_by": ["@bob"]}

The change is the difference between <base> and [head], or between <base> and
the working tree when no head is given. With --write, which needs the working
tree, the people who stand behind a rule the change touches are added to the
record, and a proposed rule everyone has agreed to is rewritten accepted.

Exits non-zero when a rule is accepted without its people.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		head := ""
		if len(args) == 2 {
			head = args[1]
		}
		return checkAgreements(cmd.OutOrStdout(), agreementsRoot, args[0], head, agreementsApprovals, agreementsWrite)
	},
}

func loadSide(root, rev string) (map[string]agreement.Side, error) {
	if rev == "" {
		return agreement.LoadTree(root)
	}
	return agreement.LoadRev(root, rev)
}

// checkAgreements writes before it checks, so what is judged is the state the
// run leaves behind: a rewrite that did not pass its own check would be a
// second opinion from the same command.
func checkAgreements(out io.Writer, root, base, head, approvals string, write bool) error {
	if write && head != "" {
		return fmt.Errorf("--write changes the working tree, so it cannot be used with a head revision")
	}
	facts, err := agreement.LoadFacts(approvals)
	if err != nil {
		return err
	}
	baseSide, err := agreement.LoadRev(root, base)
	if err != nil {
		return err
	}
	headSide, err := loadSide(root, head)
	if err != nil {
		return err
	}
	if write {
		mappings := filepath.Join(root, filepath.FromSlash(agreement.MappingsDir))
		written, err := agreement.Write(root, mappings, baseSide, headSide, facts)
		if err != nil {
			return err
		}
		for _, u := range written.Recorded {
			fmt.Fprintln(out, "recorded: "+u)
		}
		for _, u := range written.Accepted {
			fmt.Fprintln(out, "accepted: "+u)
		}
		if headSide, err = agreement.LoadTree(root); err != nil {
			return err
		}
	}

	res := agreement.Check(baseSide, headSide, facts)
	if res.Checked == 0 {
		fmt.Fprintln(out, "nothing to check: the change accepts no rule that names decision makers")
		return nil
	}
	for _, f := range res.Findings {
		fmt.Fprintln(out, f.URI)
		if len(f.Missing) > 0 {
			fmt.Fprintln(out, "  waiting for: "+strings.Join(f.Missing, ", "))
		}
		if len(f.Teams) > 0 {
			fmt.Fprintln(out, "  cannot match a team, name people instead: "+strings.Join(f.Teams, ", "))
		}
		if len(f.Unbacked) > 0 {
			fmt.Fprintln(out, "  on record without having written or approved this change: "+strings.Join(f.Unbacked, ", "))
		}
	}
	if !res.OK() {
		return fmt.Errorf("%d of %d accepted rule(s) lack the agreement of the people they name", len(res.Findings), res.Checked)
	}
	fmt.Fprintf(out, "ok: %d accepted rule(s) have the agreement of everyone they name\n", res.Checked)
	return nil
}
