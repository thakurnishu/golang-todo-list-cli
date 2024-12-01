package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thakurnishu/golang-todo-list-cli/util"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pending task",
	Long:  `list all pending task`,
	Args:  cobra.ExactArgs(0),
	Run:   listCmdRun,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolP("all", "a", true, "list all task")
}
func listCmdRun(cmd *cobra.Command, args []string) {
	allFlagPassed := cmd.Flags().Changed("all")

	err := util.ListTasks(taskFilename, allFlagPassed)
	if err != nil {
		fmt.Println("Error Listing Task: ", err)
	}
}
