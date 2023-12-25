package use

import (
	"fmt"

	"github.com/kloudlite/kl/cmd/util"
	"github.com/kloudlite/kl/lib/common"
	"github.com/kloudlite/kl/lib/common/ui/text"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "select cluster to use later with all commands",
	Long: `Select cluster
Examples:
  # select cluster
  kl use cluster
	# this will open selector where you can select one of the cluster accessible to you.

  # select account with cluster id
  kl use cluster <clusterId>
	`,
	Run: func(_ *cobra.Command, args []string) {
		clusterName, err := util.SelectCluster(args)
		if err != nil {
			common.PrintError(err)
			return
		}

		fmt.Println(text.Bold(text.Green("\nSelected cluster:")), text.Blue(clusterName))
	},
}
