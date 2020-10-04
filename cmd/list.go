/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"
	"sort"

	"github.com/shijunLee/kubernetes-config-tool/pkg/utils"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listCurrentKubernetesContext()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type contextResult struct {
	Name        string
	IsCurrent   string
	Server      string
	Namespace   string
	ClusterName string
}

func listCurrentKubernetesContext() {
	config, err := clientcmd.LoadFromFile(defaultPath)
	if err != nil {
		fmt.Println("load the default kubernetes config file error", err)
		return
	}
	var contexts []contextResult
	for name, context := range config.Contexts {
		isCurrent := ""
		if name == config.CurrentContext {
			isCurrent = "√"
		}
		contexts = append(contexts, contextResult{
			Name:        name,
			Namespace:   context.Namespace,
			ClusterName: context.Cluster,
			IsCurrent:   isCurrent,
		})
	}
	var result []contextResult
	for _, item := range contexts {
		for name, server := range config.Clusters {
			if item.ClusterName == name {
				item.Server = server.Server
			}
		}
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ClusterName > result[j].ClusterName
	})
	utils.PrintObjectTable(result, os.Stdout)
}
