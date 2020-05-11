package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	inputFile   string
	outputFile  string
	sourceFile  string
	clusterName string
	help        bool
)

const (
	// We omit vowels from the set of available characters to reduce the chances
	// of "bad words" being formed.
	alphanums = "bcdfghjklmnpqrstvwxz2456789"
	// No. of bits required to index into alphanums string.
	alphanumsIdxBits = 5
	// Mask used to extract last alphanumsIdxBits of an int.
	alphanumsIdxMask = 1<<alphanumsIdxBits - 1
	// No. of random letters we can extract from a single int63.
	maxAlphanumsPerInt = 63 / alphanumsIdxBits
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&inputFile, "i", "", "the kubernetes config file will merge file path ")
	flag.StringVar(&outputFile, "o", "./config", "the merged kubernetes config file output path")
	flag.StringVar(&sourceFile, "s", "", "the source config file will merged default is user dir ~/.kube/config")
	flag.StringVar(&clusterName, "n", "", "the merge config cluster name")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	var defaultPath = ""
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
	_, err := os.Stat(inputFile)
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
	if fileInfo.IsDir() {
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
	}else{
		f, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0600)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if clusterName == "" {
		suffix := String(5)
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
		value.Username = fmt.Sprintf("%s-%s", value.Username, clusterName)
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

var rng = struct {
	sync.Mutex
	rand *rand.Rand
}{
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
}

func String(n int) string {
	b := make([]byte, n)
	rng.Lock()
	defer rng.Unlock()

	randomInt63 := rng.rand.Int63()
	remaining := maxAlphanumsPerInt
	for i := 0; i < n; {
		if remaining == 0 {
			randomInt63, remaining = rng.rand.Int63(), maxAlphanumsPerInt
		}
		if idx := int(randomInt63 & alphanumsIdxMask); idx < len(alphanums) {
			b[i] = alphanums[idx]
			i++
		}
		randomInt63 >>= alphanumsIdxBits
		remaining--
	}
	return string(b)
}
