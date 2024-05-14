package cmd

import (
	"fmt"
	"github.com/EscanBE/node-setup-check/constants"
	"github.com/EscanBE/node-setup-check/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

const (
	flagType = "type"
)

func GetCheckCmd() *cobra.Command {
	validTargetValues := strings.Join(types.AllNodeTypeNames(), "/")

	var cmd = &cobra.Command{
		Use:     "check [home]",
		Aliases: []string{},
		Args:    cobra.ExactArgs(1),
		Short:   "Check node setup",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("App version", constants.VERSION)
			fmt.Println("NOTICE: always update to latest version for accurate check")
			time.Sleep(2 * time.Second)

			typeName, _ := cmd.Flags().GetString(flagType)
			nodeType := types.NodeTypeFromString(typeName)
			if nodeType == types.UnspecifiedNodeType {
				exitWithErrorMsgf("ERR: Invalid node type, can be either %s\n", validTargetValues)
			}

			defer func() {
				if len(checkRecords) == 0 {
					fmt.Println("All checks passed")
					return
				}

				printCheckRecords()
				os.Exit(1)
			}()

			home := args[0]
			checkHome(home)

			checkHomeKeyring(home, nodeType == types.ValidatorNode)
			checkHomeConfig(home, nodeType)
			checkHomeData(home, nodeType)
		},
	}

	cmd.Flags().String(flagType, "", fmt.Sprintf("type of node to check, can be: %s", validTargetValues))

	return cmd
}

func init() {
	rootCmd.AddCommand(GetCheckCmd())
}
