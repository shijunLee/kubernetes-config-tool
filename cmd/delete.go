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

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var deleteClusterName = ""

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteCluster()
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
	deleteCmd.Flags().StringVarP(&deleteClusterName, "cluster-name", "", "", "the cluster name which you want to delete")
}

func deleteCluster() {
	config, err := clientcmd.LoadFromFile(defaultPath)
	if err != nil {
		fmt.Println("load the default kubernetes config file error", err)
		return
	}

	if deleteClusterName == "" {
		fmt.Println("not set the `cluster-name` you want to delete ,if you don't know all clusters see below")
		listCurrentKubernetesContext()
		return
	}
	isContain := false
	for name := range config.Clusters {
		if name == deleteClusterName {
			isContain = true
			break
		}
	}
	if !isContain {
		fmt.Println("the `cluster-name` you want to delete is not found,if you don't know all clusters see below")
		listCurrentKubernetesContext()
		return
	}

	var newConfig = clientcmdapi.NewConfig()
	for name, cluster := range config.Clusters {
		if name != deleteClusterName {
			newConfig.Clusters[name] = cluster
		}
	}
	var deleteContext []string
	var deleteContextUsers []string
	for name, context := range config.Contexts {
		if context.Cluster != deleteClusterName {
			newConfig.Contexts[name] = context
		} else {
			deleteContext = append(deleteContext, name)
			deleteContextUsers = append(deleteContextUsers, context.AuthInfo)
		}
	}
	for name, authInfo := range config.AuthInfos {
		isDelete := false
		for _, item := range deleteContextUsers {
			if item == name {
				isDelete = true
				break
			}
		}
		if !isDelete {
			newConfig.AuthInfos[name] = authInfo
		}
	}
	newConfig.APIVersion = config.APIVersion
	newConfig.Extensions = config.Extensions
	newConfig.Preferences = config.Preferences
	isCurrentContextDelete := false
	for _, item := range deleteContext {
		if item == config.CurrentContext {
			isCurrentContextDelete = true
		}
	}
	if isCurrentContextDelete {
		var contexts []string
		for contextName := range newConfig.Contexts {
			contexts = append(contexts, contextName)
		}
		if len(contexts) > 0 {
			newConfig.CurrentContext = contexts[0]
		}
	} else {
		newConfig.CurrentContext = config.CurrentContext
	}

	err = clientcmd.WriteToFile(*newConfig, defaultPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("write file to %s error", defaultPath), err)
		return
	}

}
