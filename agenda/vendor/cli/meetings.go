// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"cmd"
	"strings"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var createMeetingsCmd = &cobra.Command{
	Use:   "createMeetings",
	Short: "Create meetings.",
	Long: `If you want to create a meeting, you should declare the title name, 
		which can't be the same as others' title name, some participators(at least one),
		start time as the format of (yyyy-mm-dd), and end time as the format of (yyyy-mm-dd).`,
	Run: func(comd *cobra.Command, args []string) {
		title, _ := comd.Flags().GetString("title")
		checkEmpty("title", title)
		
		participatorStr, _ := comd.Flags().GetString("participators")
		checkEmpty("participators", participatorStr)
		participators := strings.Split(participatorStr, " ")
		
		startTime, _ := comd.Flags().GetString("start")
		checkEmpty("Start Time", startTime)

		endTime, _ := comd.Flags().GetString("end")
		checkEmpty("End Time", endTime)

		cmd.HostMeeting(title, participators, startTime, endTime)
	},
}

var changeParticipatorsCmd = &cobra.Command{
	Use: "changeParticipators",
	Short: "Change your own meetings' participators.",
	Long: `You can append or remove some participators from your own meeting
	by specifying the title name.`,
	Run: func(comd *cobra.Command, args []string) {
		title, _ := comd.Flags().GetString("title")
		checkEmpty("title", title)

		participatorStr, _ := comd.Flags().GetString("participators")
		checkEmpty("participators", participatorStr)
		participators := strings.Split(participatorStr, " ")
		
		deleteOrNot, _ := comd.Flags().GetBool("delete")
		if deleteOrNot {
			var error int
			for _, each := range participators {
				error = cmd.AddParticipant(title, each)
 				if error != 0 {
 					return
 				}
			}
		} else {
			var error int
			for _, each := range participators {
				error = cmd.RemoveParticipant(title, each)
				if error != 0 {
					return
				}
			}
		}
	},
}

var listMeetingsCmd = &cobra.Command{
	Use: "listMeetingsCmd",
	Short: "List all of your own meetings during a time interval.",
	Long: `You can see the detail information of all of meetings,
	which you attended, during a time interval.`,
	Run: func(comd *cobra.Command, args []string) {
		startTime, _ := comd.Flags().GetString("start")
		checkEmpty("Start Time", startTime)

		endTime, _ := comd.Flags().GetString("end")
		checkEmpty("End Time", endTime)

		cmd.QueryMeeting(startTime, endTime)
	},
}

var cancelCmd = &cobra.Command{
	Use: "cancel",
	Short: "Cancel your own meeting by specifying title name.",
	Long: `Using this command, you are able to cancel the meetings, which are created by you.`,
	Run: func(comd *cobra.Command, args []string) {
		title, _ := comd.Flags().GetString("title")
		checkEmpty("Title", title)

		cmd.CancelMeeting(title)
	},
}

var quitCmd = &cobra.Command{
	Use: "quit",
	Short: "Quit meetings.",
	Long: `You can quit any meetings you want, which are you attended, not created.`,
	Run: func(comd *cobra.Command, args []string) {
		title, _ := comd.Flags().GetString("title")
		checkEmpty("Title", title)

		cmd.QuitMeeting(title)
	},
}

var clearCmd = &cobra.Command{
	Use: "clear",
	Short: "Clear all meetings you attended or created.",
	Long: `Using this command, you can clear all of the meetings you attended or created.`,
	Run: func(comd *cobra.Command, args []string) {
		cmd.ClearMeetings()
	},
}


func init() {
	createMeetingsCmd.Flags().StringP("title", "t", "", "Input title name.")
	createMeetingsCmd.Flags().StringP("participators", "p", "", "Input participator name.")
	createMeetingsCmd.Flags().StringP("start", "s", "", "Input start time as the format of (yyyy-mm-dd).")
	createMeetingsCmd.Flags().StringP("end", "e", "", "Input end time as the format of (yyyy-mm-dd).")

	changeParticipatorsCmd.Flags().BoolP("delete", "d", false, "If true, delete participators, otherwise append participators.")
	changeParticipatorsCmd.Flags().StringP("title", "t", "", "Input the title name.")
	changeParticipatorsCmd.Flags().StringP("participators", "p", "", "Input the participators.")

	listMeetingsCmd.Flags().StringP("start", "s", "", "Input the start time.")
	listMeetingsCmd.Flags().StringP("end", "e", "", "Input the end time.")

	cancelCmd.Flags().StringP("title", "t", "", "Input the title.")

	quitCmd.Flags().StringP("title", "t", "", "Input the title.")

	RootCmd.AddCommand(createMeetingsCmd)
	RootCmd.AddCommand(changeParticipatorsCmd)
	RootCmd.AddCommand(listMeetingsCmd)
	RootCmd.AddCommand(cancelCmd)
	RootCmd.AddCommand(quitCmd)
	RootCmd.AddCommand(clearCmd)


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
