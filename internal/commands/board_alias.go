package commands

import (
	"fmt"
	"sort"

	"github.com/robzolkos/fizzy-cli/internal/config"
	"github.com/robzolkos/fizzy-cli/internal/errors"
	"github.com/robzolkos/fizzy-cli/internal/response"
	"github.com/spf13/cobra"
)

var boardAliasCmd = &cobra.Command{
	Use:   "board-alias",
	Short: "Manage board aliases",
	Long:  "Commands for managing board aliases. Aliases let you use short names instead of board IDs.",
}

var boardAliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all board aliases",
	Long:  "Lists all configured board aliases.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadGlobal()
		aliases := cfg.BoardAliases

		if len(aliases) == 0 {
			printSuccessWithBreadcrumbs(
				[]interface{}{},
				"No aliases configured",
				[]response.Breadcrumb{
					breadcrumb("add", "fizzy board-alias add <alias> <board_id>", "Add an alias"),
					breadcrumb("boards", "fizzy board list", "List boards to find IDs"),
				},
			)
			return
		}

		// Sort aliases alphabetically
		keys := make([]string, 0, len(aliases))
		for k := range aliases {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		data := make([]map[string]string, 0, len(aliases))
		for _, alias := range keys {
			data = append(data, map[string]string{
				"alias":    alias,
				"board_id": aliases[alias],
			})
		}

		summary := fmt.Sprintf("%d alias(es)", len(aliases))
		breadcrumbs := []response.Breadcrumb{
			breadcrumb("add", "fizzy board-alias add <alias> <board_id>", "Add an alias"),
			breadcrumb("remove", "fizzy board-alias remove <alias>", "Remove an alias"),
		}

		printSuccessWithBreadcrumbs(data, summary, breadcrumbs)
	},
}

var boardAliasAddCmd = &cobra.Command{
	Use:   "add ALIAS BOARD_ID",
	Short: "Add a board alias",
	Long: `Adds a board alias. The alias can then be used anywhere a board ID is expected.

Examples:
  fizzy board-alias add visa 03flgarh3nklbrb8uthxhvsoa
  fizzy card list --board visa
  fizzy board show visa`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		boardID := args[1]

		cfg := config.LoadGlobal()
		cfg.SetBoardAlias(alias, boardID)

		if err := cfg.Save(); err != nil {
			exitWithError(errors.NewError(fmt.Sprintf("Failed to save config: %s", err)))
		}

		data := map[string]string{
			"alias":    alias,
			"board_id": boardID,
		}

		breadcrumbs := []response.Breadcrumb{
			breadcrumb("list", "fizzy board-alias list", "List all aliases"),
			breadcrumb("use", fmt.Sprintf("fizzy card list --board %s", alias), "Use the alias"),
			breadcrumb("remove", fmt.Sprintf("fizzy board-alias remove %s", alias), "Remove this alias"),
		}

		printSuccessWithBreadcrumbs(data, fmt.Sprintf("Alias '%s' → %s", alias, boardID), breadcrumbs)
	},
}

var boardAliasRemoveCmd = &cobra.Command{
	Use:   "remove ALIAS",
	Short: "Remove a board alias",
	Long:  "Removes a board alias.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

		cfg := config.LoadGlobal()
		if !cfg.RemoveBoardAlias(alias) {
			exitWithError(errors.NewNotFoundError(fmt.Sprintf("Alias '%s' not found", alias)))
		}

		if err := cfg.Save(); err != nil {
			exitWithError(errors.NewError(fmt.Sprintf("Failed to save config: %s", err)))
		}

		data := map[string]string{
			"alias":   alias,
			"removed": "true",
		}

		breadcrumbs := []response.Breadcrumb{
			breadcrumb("list", "fizzy board-alias list", "List remaining aliases"),
		}

		printSuccessWithBreadcrumbs(data, fmt.Sprintf("Removed alias '%s'", alias), breadcrumbs)
	},
}

func init() {
	rootCmd.AddCommand(boardAliasCmd)
	boardAliasCmd.AddCommand(boardAliasListCmd)
	boardAliasCmd.AddCommand(boardAliasAddCmd)
	boardAliasCmd.AddCommand(boardAliasRemoveCmd)
}
