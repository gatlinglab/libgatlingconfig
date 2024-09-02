package libgatlingconfig

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CGatlingConfig struct {
	kValue  map[string]string
	appPath string
}

const c_Key_ConfigServerUrl = "CONFIGSERVERURL"
const c_Key_ConfigServerToken = "CONFIGSERVERTOKEN" // config server token is for X-API-KEY
const c_Key_ConfigAPPToken = "APPTOKEN"             // app token for the kv key;

var g_singleGatlingConfig *CGatlingConfig = &CGatlingConfig{kValue: map[string]string{}}

func GetSingleGatlingConfig() *CGatlingConfig {
	return g_singleGatlingConfig
}

func (pInst *CGatlingConfig) Initialize(appName string) error {
	if appName == "" {
		return errors.New("appname is empty")
	}
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	pInst.appPath = path[:index+1]

	pInst.listEnv()
	pInst.loadAppConfig(appName)

	servertoken := pInst.kValue[c_Key_ConfigServerToken]
	appToken := pInst.kValue[c_Key_ConfigAPPToken]
	serverUrl := pInst.kValue[c_Key_ConfigServerUrl]

	if serverUrl != "" {
		pInst.loadServerConfig(serverUrl, servertoken, appName)
	}

	if serverUrl != "" && appToken != "" {
		pInst.loadServerConfig(serverUrl, servertoken, appName+"-"+appToken)
	}

	return nil
}

func (pInst *CGatlingConfig) listEnv() int {
	var iCount = 0
	for i, env := range os.Environ() {
		// env is
		envPair := strings.SplitN(env, "=", 2)
		key := envPair[0]
		value := envPair[1]
		if key != "" {
			pInst.kValue[key] = value
			iCount = i
		}

	}

	return iCount
}

func (pInst *CGatlingConfig) loadAppConfig(appName string) int {
	var iCount = 0
	f, err := os.Open(appName + ".cfg")
	if err != nil {
		return -1
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		iRet := pInst.analyseConfig(scanner.Text())
		if iRet > 0 {
			iCount++
		}
	}

	return iCount
}
func (pInst *CGatlingConfig) loadServerConfig(serverUrl, serverToken, appToken string) int {
	if serverUrl[len(serverUrl)-1] != '/' {
		serverUrl += "/"
	}
	serverUrl += appToken
	req, err := http.NewRequest(http.MethodGet, serverUrl, nil)
	if err != nil {
		return -1
	}
	if serverToken != "" {
		req.Header.Add("X-API-KEY", serverToken)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -2
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return -3
	}
	strConect := strings.ReplaceAll(string(content), "\r", "\n")
	strConect = strings.ReplaceAll(strConect, "█", "\n")
	strConect = strings.ReplaceAll(strConect, "\n\n", "\n")
	strConect = strings.ReplaceAll(strConect, "\n\n", "\n")
	strList := strings.Split(strConect, "\n")
	var iCount = 0
	for _, line := range strList {
		iRet := pInst.analyseConfig(line)
		if iRet > 0 {
			iCount++
		}
	}
	return iCount
}

func (pInst *CGatlingConfig) analyseConfig(line string) int {
	strPair := strings.SplitN(line, "=", 2)
	if len(strPair) < 2 {
		return -1
	}
	if strPair[0] != "" {
		pInst.kValue[strPair[0]] = strPair[1]
	}
	return 1
}

func (pInst *CGatlingConfig) Get(key string) string {
	return pInst.kValue[key]
}

func (pInst *CGatlingConfig) Set(key string, value string) {
	pInst.kValue[key] = value
}
