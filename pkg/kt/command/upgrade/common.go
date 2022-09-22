package upgrade

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Handle(serverUrl, currentVersion string) error {
	log.Info().Msgf("Visit server %s for check version", serverUrl)
	// 查看服务端文件列表
	res, err := http.Get(serverUrl + "/index")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err := validateStatusCode(res.StatusCode); err != nil {
		log.Error().Msgf(string(body))
		return err
	}

	forwardVersion := analyse(currentVersion, body)

	if forwardVersion == nil {
		log.Info().Msgf("Your system installed et-%s is the latest version, nothing would be done...", currentVersion)
		return nil
	}

	// 用户选择版本
	userSelectVersion := ""
	for {
		fmt.Println("You have installed et-" + currentVersion + ", please choose version upgrade:\n" + strings.Join(forwardVersion, "\n"))
		fmt.Print("Input version(default:" + forwardVersion[len(forwardVersion)-1] + "):")
		fmt.Scanln(&userSelectVersion)

		if userSelectVersion == "" {
			userSelectVersion = forwardVersion[len(forwardVersion)-1]
			log.Info().Msgf("you have choosed version %s, starting download now...", userSelectVersion)
			break
		} else {
			if !IsContainInArray(userSelectVersion, forwardVersion) {
				fmt.Println("Warning: You entered a wrong version!!!")
				userSelectVersion = ""
			} else {
				break
			}
		}
	}
	if err := DownloadAndRename(serverUrl, userSelectVersion); err != nil {
		return err
	}

	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" Upgrade ET package success, Enjoy!")
	log.Info().Msg("---------------------------------------------------------------")

	return nil
}

func analyse(currentVersion string, body []byte) []string {
	// 去掉最后一个回车，并分割成数组
	versionArray := strings.Split(strings.TrimRight(string(body), "\n"), "\n")
	// 统计要升级的版本
	var forwardVersion []string
	for i := 0; i < len(versionArray); i++ {
		if versionArray[i] == currentVersion && i+1 != len(versionArray) {
			forwardVersion = versionArray[i+1:]
			break
		}
	}
	// 如果当前版本不在服务端列表中，则全部展示给用户选择
	if !IsContainInArray(currentVersion, versionArray) {
		forwardVersion = versionArray
	}

	return forwardVersion
}

func DownloadAndRename(serverUrl, userSelectVersion string) error {
	resp, err := http.Get(serverUrl + "/" + userSelectVersion + "/et.exe")
	if err != nil {
		return err
	}
	if err := validateStatusCode(resp.StatusCode); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Error().Msgf(string(body))
		return err
	}

	base := os.Getenv("ET_HOME") + "\\et.exe"
	out, err := os.Create(base + ".temp")
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	log.Info().Msgf("Download package complete.")
	log.Info().Msgf("Starting install new package.")
	out.Close()

	os.Rename(base, base+".old")
	os.Rename(base+".temp", base)
	return nil
}

func IsContainInArray(item string, items []string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func validateStatusCode(code int) error {
	if code < 200 || code > 300 {
		return fmt.Errorf("HTTP error code:%d", code)
	}
	return nil
}
