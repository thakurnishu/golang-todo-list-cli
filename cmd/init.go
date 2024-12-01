/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thakurnishu/golang-todo-list-cli/util"
)

var taskHeaders = []string{"ID", "Description", "CreatedAt", "IsComplete"}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create Initial CSV File for Database",
	Long:  `Create Initial CSV File for Database`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CreateInitialDatabase(taskFilename, taskHeaders)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
