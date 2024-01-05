package del

import (
	"fmt"
	"github.com/kloudlite/kl/domain/client"
	common_util "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/fzf"

	"github.com/kloudlite/kl/constants"
	"github.com/spf13/cobra"
)

var deleteMountCommand = &cobra.Command{
	Use:   "mount",
	Short: "remove one mount from your " + constants.CmdName + "-config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(_ *cobra.Command, _ []string) {
		removeConfigMount()
	},
}

func removeConfigMount() {

	klFile, err := client.GetKlFile(nil)

	if err != nil {
		common_util.PrintError(err)
		return
	}

	selectedMount, err := fzf.FindOne(
		klFile.FileMount.Mounts,
		func(item client.FileEntry) string {
			return fmt.Sprintf("%s | %s | %s", item.Type, item.Path, item.Name)
		},
		fzf.WithPrompt("Select Config Group >"),
	)

	if err != nil {
		common_util.PrintError(err)
		return
	}

	newMounts := make([]client.FileEntry, 0)
	for _, fe := range klFile.FileMount.Mounts {
		if fe.Name == selectedMount.Name {
			continue
		}
		newMounts = append(newMounts, fe)
	}

	klFile.FileMount.Mounts = newMounts

	err = client.WriteKLFile(*klFile)

	if err != nil {
		common_util.PrintError(err)
	}

	fmt.Println("mount removed")
}
