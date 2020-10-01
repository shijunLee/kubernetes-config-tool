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
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var sourceFile string
var defaultPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubernetes-config-tool",
	Short: "a tool for manager kubernetes config file",
	Long:  `this is command app for manager kubernetes config files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&sourceFile, "sourcefile", "s", "", "the source config file will merged default is user dir ~/.kube/config")
	initDefaultKuberConfig()
}

func initDefaultKuberConfig() {
	if sourceFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir, err = filepath.Abs("~")
			if err != nil {
				fmt.Println("get user home path error")
				panic(err)
			}
		}
		defaultPath = path.Join(homeDir, ".kube", "config")
	} else {
		defaultPath = sourceFile
		_, err := os.Stat(sourceFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("the config source file not exist")
				return
			}
		}
	}
}
