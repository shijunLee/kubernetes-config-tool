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

	"github.com/shijunLee/kubernetes-config-tool/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	inputFile   string
	outputFile  string
	clusterName string
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:     "merge",
	Aliases: []string{"m"},
	Short:   "merge a kubernetes config file to another",
	Long: `merge kubernete config file to the source kubernetes kubernetes config file :

kct merge --input=./test.config --output=./config --cluster-name=test [--sourcefile=~/.kube/config]
or
kct m -i ./test.config -o  ./config -n test [-s ~/.kube/config]`,
	Run: func(cmd *cobra.Command, args []string) {
		runMergeCommand()
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	mergeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "the kubernetes config file will merge file path ")
	mergeCmd.Flags().StringVarP(&outputFile, "output", "o", "./config", "the merged kubernetes config file output path")
	mergeCmd.Flags().StringVarP(&clusterName, "cluster-name", "n", "", "the merge config cluster name")
}

func runMergeCommand() {
	inputFile, err := filepath.Abs(inputFile)
	if err != nil {
		fmt.Println("the input path error")
		return
	}
	_, err = os.Stat(inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("the merge config source file not exist")
			return
		}
	}
	config, err := clientcmd.LoadFromFile(defaultPath)
	if err != nil {
		fmt.Println("load the default kubernetes config file error", err)
		return
	}
	if clientcmdapi.IsConfigEmpty(config) {
		fmt.Println("load the default kubernetes config file is empty")
		return
	}
	mergeConfig, err := clientcmd.LoadFromFile(inputFile)
	if err != nil {
		fmt.Println("load the default kubernetes config file error", err)
		return
	}
	if clientcmdapi.IsConfigEmpty(mergeConfig) {
		fmt.Println("load the default kubernetes config file is empty")
		return
	}
	outputFile, err = filepath.Abs(outputFile)
	if err != nil {
		fmt.Println("the output path error")
		return
	}
	fileInfo, err := os.Stat(outputFile)
	var f *os.File
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("the output file is exist")
			return
		} else if os.IsNotExist(err) {
			f, err = os.Create(outputFile)
			defer f.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err)
			return
		}
	}
	if fileInfo != nil && fileInfo.IsDir() {
		outputFile = path.Join(outputFile, "config")
		_, err = os.Stat(outputFile)
		if os.IsNotExist(err) {
			f, err = os.Create(outputFile)
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			f, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0600)
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		f, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0600)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if clusterName == "" {
		suffix := utils.String(5)
		clusterName = suffix
	}

	for key, value := range mergeConfig.Clusters {
		mergeKey := fmt.Sprintf("%s-%s", key, clusterName)
		config.Clusters[mergeKey] = value
	}

	for key, value := range mergeConfig.Contexts {
		mergeKey := fmt.Sprintf("%s-%s", key, clusterName)
		value.Cluster = fmt.Sprintf("%s-%s", value.Cluster, clusterName)
		value.AuthInfo = fmt.Sprintf("%s-%s", value.AuthInfo, clusterName)
		config.Contexts[mergeKey] = value
	}
	for key, value := range mergeConfig.AuthInfos {
		mergeKey := fmt.Sprintf("%s-%s", key, clusterName)
		if value.Username != "" {
			value.Username = fmt.Sprintf("%s-%s", value.Username, clusterName)
		}
		config.AuthInfos[mergeKey] = value
	}
	data, err := clientcmd.Write(*config)
	if err != nil {
		fmt.Println("get merged config error")
		return
	}

	_, err = f.Write(data)
	if err != nil {
		fmt.Println("write file error", err)
	}

}
