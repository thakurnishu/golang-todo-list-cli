package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thakurnishu/golang-todo-list-cli/util"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add to task",
	Long:  `add to task`,
	Args:  cobra.ExactArgs(1),
	Run:   addCmdRun,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCmdRun(cmd *cobra.Command, args []string) {
	err := util.AddTask(taskFilename, args[0])
	if err != nil {
		fmt.Println("Error Creating Task: ", err)
	}
}
