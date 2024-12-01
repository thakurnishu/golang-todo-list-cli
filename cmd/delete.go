package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thakurnishu/golang-todo-list-cli/util"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete task with provide taskId",
	Long:  `delete task with provide taskId`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := util.DeleteTaskFromCSV(taskFilename, args[0])
		if err != nil {
			fmt.Println("Error Deleting Task: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
