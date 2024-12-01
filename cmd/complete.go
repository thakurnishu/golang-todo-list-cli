package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thakurnishu/golang-todo-list-cli/util"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "mark task as complete",
	Long:  `mark task as complete`,
	Args:  cobra.ExactArgs(1),
	Run:   completeCmdRun,
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func completeCmdRun(cmd *cobra.Command, args []string) {
	taskId := args[0]
	if num, err := strconv.Atoi(taskId); err != nil || num < 0 {
		fmt.Println("Passed TaskID should be positive integer")
		return
	}
	if err := util.MarksTaskAsComplete(taskFilename, taskId); err != nil {
		fmt.Println("Error Marking Task As Complete: ", err)
	}
}
