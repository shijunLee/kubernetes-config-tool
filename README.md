# kubernetes-config-merge-tool

* this is a tool to merge kubernetes config files

* use this tool method,please use `kubernetes-config-tool` to see the cmd params
```bash

    # use `-h` to see the tool help info
    kubernetes-config-tool -h

    # merge to config file ,if the source file not in <system user home>/.kube/config please set -s command
    kubernetes-config-tool m -i ./test.config -o ./use.config -n test

    # use l or list to see the config file context list
    kubernetes-config-tool l <-s "the config file">

    # use d or delete to delete a cluster form the config file,like: 
    kubernetes-config-tool d --sourcefile ./use.config  --cluster-name=kubernetes

    # use change to change current-context for config file
    kubernetes-config-tool change --sourcefile ./use.config  --context-name=dev
```
