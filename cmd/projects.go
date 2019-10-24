/*
Copyright © 2019 Doppler <support@doppler.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	api "doppler-cli/api"
	configuration "doppler-cli/config"
	"doppler-cli/utils"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List projects",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")

		localConfig := configuration.LocalConfig(cmd)
		_, info := api.GetAPIProjects(cmd, localConfig.Key.Value)

		printProjectsInfo(info, jsonFlag)
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get [project_id]",
	Short: "Get info for a project",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		localConfig := configuration.LocalConfig(cmd)

		project := localConfig.Project.Value
		if len(args) > 0 {
			project = args[0]
		}

		_, info := api.GetAPIProject(cmd, localConfig.Key.Value, project)

		printProjectInfo(info, jsonFlag)
	},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a project",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		description := cmd.Flag("description").Value.String()

		name := cmd.Flag("name").Value.String()
		if len(args) > 0 {
			name = args[0]
		}

		localConfig := configuration.LocalConfig(cmd)
		_, info := api.CreateAPIProject(cmd, localConfig.Key.Value, name, description)

		if !silent {
			printProjectInfo(info, jsonFlag)
		}
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete [project_id]",
	Short: "Delete a project",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO prompt user with a confirmation before proceeding (and add a --yes flag to skip it)
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		localConfig := configuration.LocalConfig(cmd)

		project := localConfig.Project.Value
		if len(args) > 0 {
			project = args[0]
		}

		api.DeleteAPIProject(cmd, localConfig.Key.Value, project)

		// fetch and display projects
		if !silent {
			_, info := api.GetAPIProjects(cmd, localConfig.Key.Value)
			printProjectsInfo(info, jsonFlag)
		}
	},
}

var projectsUpdateCmd = &cobra.Command{
	Use:   "update [project_id]",
	Short: "Update a project",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		localConfig := configuration.LocalConfig(cmd)

		project := localConfig.Project.Value
		if len(args) > 0 {
			project = args[0]
		}

		name := cmd.Flag("name").Value.String()
		description := cmd.Flag("description").Value.String()

		_, info := api.UpdateAPIProject(cmd, localConfig.Key.Value, project, name, description)

		if !silent {
			printProjectInfo(info, jsonFlag)
		}
	},
}

func init() {
	projectsGetCmd.Flags().Bool("json", false, "output json")
	projectsGetCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	projectsCmd.AddCommand(projectsGetCmd)

	projectsCreateCmd.Flags().Bool("json", false, "output json")
	projectsCreateCmd.Flags().Bool("silent", false, "don't output the response")
	projectsCreateCmd.Flags().String("name", "", "project name")
	projectsCreateCmd.Flags().String("description", "", "project description")
	projectsCmd.AddCommand(projectsCreateCmd)

	projectsDeleteCmd.Flags().Bool("json", false, "output json")
	projectsDeleteCmd.Flags().Bool("silent", false, "don't output the response")
	projectsCmd.AddCommand(projectsDeleteCmd)

	projectsUpdateCmd.Flags().Bool("json", false, "output json")
	projectsUpdateCmd.Flags().Bool("silent", false, "don't output the response")
	projectsUpdateCmd.Flags().String("name", "", "project name")
	projectsUpdateCmd.Flags().String("description", "", "project description")
	projectsUpdateCmd.MarkFlagRequired("name")
	projectsUpdateCmd.MarkFlagRequired("description")
	projectsCmd.AddCommand(projectsUpdateCmd)

	projectsCmd.Flags().Bool("json", false, "output json")
	rootCmd.AddCommand(projectsCmd)
}

func printProjectsInfo(info []api.ProjectInfo, jsonFlag bool) {
	if jsonFlag {
		resp, err := json.Marshal(info)
		if err != nil {
			utils.Err(err)
		}

		fmt.Println(string(resp))
		return
	}

	var rows [][]string
	for _, projectInfo := range info {
		rows = append(rows, []string{projectInfo.ID, projectInfo.Name, projectInfo.Description, projectInfo.SetupAt, projectInfo.CreatedAt})
	}
	utils.PrintTable([]string{"id", "name", "description", "setup_at", "created_at"}, rows)
}

func printProjectInfo(info api.ProjectInfo, jsonFlag bool) {
	if jsonFlag {
		resp, err := json.Marshal(info)
		if err != nil {
			utils.Err(err)
		}

		fmt.Println(string(resp))
		return
	}

	rows := [][]string{{info.ID, info.Name, info.Description, info.SetupAt, info.CreatedAt}}
	utils.PrintTable([]string{"id", "name", "description", "setup_at", "created_at"}, rows)
}