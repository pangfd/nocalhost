/*
Copyright 2020 The Nocalhost Authors.
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

package cmds

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/user"
	"strings"
)

func GetK8sRestClientConfig() (*restclient.Config, error) {
	home := GetHomePath()
	kubeconfigPath := fmt.Sprintf("%s/.kube/config", home) // default kubeconfig
	if settings.KubeConfig != "" {
		kubeconfigPath = settings.KubeConfig
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getClientSet() (*kubernetes.Clientset, error) {
	k8sConfig, err := GetK8sRestClientConfig()
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	return clientSet, nil
}

//func GetDeploymentClient(nameSpace string) (appsV1.DeploymentInterface, error) {
//	clientSet, err := getClientSet()
//	if err != nil {
//		fmt.Printf("%v", err)
//		return nil, err
//	}
//
//	deploymentsClient := clientSet.AppsV1().Deployments(nameSpace)
//	return deploymentsClient, nil
//}

func GetRestClient() (*restclient.RESTClient, error) {
	k8sConfig, err := GetK8sRestClientConfig()
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	return restclient.RESTClientFor(k8sConfig)
}

func GetPodClient(nameSpace string) (coreV1.PodInterface, error) {
	clientSet, err := getClientSet()
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	podClient := clientSet.CoreV1().Pods(nameSpace)
	return podClient, nil
}

//func GetReplicaSetsControlledByDeployment(deployment string) (map[int]*v1.ReplicaSet, error) {
//	var rsList *v1.ReplicaSetList
//	clientSet, err := getClientSet()
//	if err == nil {
//		replicaSetsClient := clientSet.AppsV1().ReplicaSets(nameSpace)
//		rsList, err = replicaSetsClient.List(context.TODO(), metav1.ListOptions{})
//	}
//	if err != nil {
//		//fmt.Printf("failed to get rs: %v\n", err)
//		return nil, err
//	}
//
//	rsMap := make(map[int]*v1.ReplicaSet)
//	for _, item := range rsList.Items {
//		if item.OwnerReferences != nil {
//			for _, owner := range item.OwnerReferences {
//				if owner.Name == deployment && item.Annotations["deployment.kubernetes.io/revision"] != "" {
//					revision, err := strconv.Atoi(item.Annotations["deployment.kubernetes.io/revision"])
//					if err == nil {
//						rsMap[revision] = item.DeepCopy()
//					}
//				}
//			}
//		}
//	}
//	return rsMap, nil
//}

func printlnErr(info string, err error) {
	fmt.Printf("%s, err: %v\n", info, err)
}

func GetHomePath() string {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir
	}
	return ""
}

func GetNocalhostConfigPath(appName string) {

}

func GetApplicationHomeDir(appName string) string {
	// GetHomePath() + "/.nhctl/" + "application/" + applicationName
	return fmt.Sprintf("%s%c%s%c%s%c%s", GetHomePath(), os.PathSeparator, ".nhctl", os.PathSeparator, "application", os.PathSeparator, appName)
}

//func GetAppPortForwardDir(appName string) string {
//	dir := fmt.Sprintf("%s%c%s", GetApplicationHomeDir(appName), os.PathSeparator, DefaultPortForwardDir)
//	return dir
//}

// .nhctl/application/xxx/port-forward/pid/
//func GetAppPortForwardPidDir(appName string, pid int) string {
//	return fmt.Sprintf("%s%c%d", GetAppPortForwardDir(appName), os.PathSeparator, pid)
//}

// .nhctl/application/xxx/port-forward/{pid}/{local_port}_{remote_port}
//func SavePortForwardInfo(appName string, localPort string, remotePort string) {
//	pid := os.Getpid()
//	pidDir := GetAppPortForwardPidDir(appName, pid)
//	fileName := fmt.Sprintf("%s%c%s_%s", pidDir, os.PathSeparator, localPort, remotePort)
//	f, err := os.Create(fileName)
//	utils.Mush(err)
//	utils.Mush(f.Close())
//}

//func GetAppConfigPath(appName string) string {
//	return fmt.Sprintf("%s%c%s%c%s", GetApplicationHomeDir(appName), os.PathSeparator, ".nocalhost", os.PathSeparator, "config.yaml")
//}

func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			fs, ds, err := GetFilesAndDirs(dirPth + PthSep + fi.Name())
			if err != nil {
				return files, dirs, err
			}
			dirs = append(dirs, ds...)
			files = append(files, fs...)
		} else {
			ok := strings.HasSuffix(fi.Name(), ".yaml")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			} else if strings.HasSuffix(fi.Name(), ".yml") {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}