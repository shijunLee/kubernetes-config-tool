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
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:     "set-current",
	Aliases: []string{"c", "change"},
	Short:   "a commend to change current context",
	Long: `a commend to change current context. For example:

kubeconfig set-current test

or

kubeconfig c test

you can use 'kubeconfig list' to see the context list.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		changeCurrentContext(args[0])
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func changeCurrentContext(contextName string) {
	config, err := clientcmd.LoadFromFile(defaultPath)
	if err != nil {
		fmt.Println("load the default kubernetes config file error", err)
		return
	}
	if contextName == "" {
		fmt.Println("not set the context you want to change to,if you don't know all context see below")
		listCurrentKubernetesContext()
		return
	}
	var contexts []string
	for key := range config.Contexts {
		contexts = append(contexts, key)
	}
	isContain := false

	for _, item := range contexts {
		if contextName == item {
			isContain = true
		}
	}
	if !isContain {
		fmt.Println("not found the context you want to set,current config info please see below")
		listCurrentKubernetesContext()
		return
	}
	config.CurrentContext = contextName
	err = clientcmd.WriteToFile(*config, defaultPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("write file to %s error", defaultPath), err)
		return
	}
	fmt.Println(fmt.Sprintf("change current context to %s success", contextName))
}
