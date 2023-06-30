package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/VichyGopher/hcmonitor/internal/filesystem"
	"github.com/VichyGopher/hcmonitor/internal/utils"
)

var (
	Config      *filesystem.ConfigStruct
	hcRegex     = regexp.MustCompile(`/captcha/v1/[A-Za-z0-9]+/static/images`)
	versionList []*HcVersion
)

func checkForUpdate() (string, error) {
	response, err := http.Get("https://hcaptcha.com/1/api.js?render=explicit&onload=hcaptchaOnLoad")
	if utils.HandleError(err) {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if utils.HandleError(err) {
		return "", err
	}

	var found string
	for _, match := range hcRegex.FindAllStringSubmatch(string(body), -1) {
		found = match[0]
		break
	}

	if found == "" {
		return "", fmt.Errorf("cant find version")
	}

	version := strings.Split(strings.Split(found, "v1/")[1], "/static")[0]

	return version, nil
}

func PushToGithub(v *HcVersion) {
	// soon
}

func DownloadVersionFiles(v *HcVersion) {
	files := []string{"hsw.js", "hsj.js", "hsl.js"}
	var urls []string
	output := ""

	os.Mkdir(fmt.Sprintf("%s/builds/%s", filesystem.BasePath, v.Version), os.ModePerm)

	for _, file := range files {

		url := fmt.Sprintf("https://newassets.hcaptcha.com/c/%s/%s", v.HswVersion, file)
		urls = append(urls, url)

		if err := filesystem.DownloadFile(url, fmt.Sprintf("builds/%s/%s", v.Version, file)); err != nil {
			fmt.Printf("[-] Failed download: %s (%s)\n", url, err.Error())
			return
		}

		fmt.Printf("[+] Success download: %s\n", url)
		output += fmt.Sprintf("\n- %s", url)
	}

	// download asset files
	assets := map[string]string{
		"api.js":        "https://hcaptcha.com/1/api.js?render=explicit&onload=hcaptchaOnLoad",
		"hcaptcha.js":   fmt.Sprintf("https://newassets.hcaptcha.com/captcha/v1/%s/hcaptcha.js", v.Version),
		"challenge.js":  fmt.Sprintf("https://newassets.hcaptcha.com/captcha/challenge/image_label_binary/%s/challenge.js", v.Version),
		"hcaptcha.html": fmt.Sprintf("https://newassets.hcaptcha.com/captcha/v1/%s/static/hcaptcha.html", v.Version),
	}

	for name, url := range assets {
		if err := filesystem.DownloadFile(url, fmt.Sprintf("builds/%s/%s", v.Version, name)); err != nil {
			fmt.Printf("[-] Failed download: %s (%s)\n", url, err.Error())
			return
		}

		fmt.Printf("[+] Success download: %s\n", url)
		output += fmt.Sprintf("\n- %s", url)
	}

	webhook, _ := disgohook.NewWebhookClientByToken(nil, nil, Config.Webhooks.Version)
	_, er := webhook.SendEmbeds(api.NewEmbedBuilder().
		SetDescription(fmt.Sprintf("### New version found `hcaptcha: %s, hsw: %s`.\n\n%s", v.Version, v.HswVersion, output)).
		SetEmbedFooter(&api.EmbedFooter{Text: fmt.Sprintf("%s - %s", time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"))}).
		Build(),
	)

	if er != nil {
		panic(er)
	}

	filesystem.AppendFile("versions.csv", fmt.Sprintf("%s,%s,%d", v.Version, v.HswVersion, v.ReleaseDate))
	go PushToGithub(v)
}

func loadVersions() {
	lines, err := filesystem.ReadFile("versions.csv")
	if utils.HandleError(err) {
		panic(err)
	}

	for i, line := range lines {
		parsed := strings.Split(line, ",")

		timestampInt, err := strconv.ParseInt(parsed[2], 10, 64)
		if err != nil {
			panic(err)
		}

		versionList = append(versionList, &HcVersion{
			Version:     parsed[0],
			HswVersion:  parsed[1],
			ReleaseDate: timestampInt,
		})

		fmt.Printf("  #%d version: %s, release date: %s\n", i, parsed[0], time.Unix(timestampInt, 0).Format("02/01/2006 15:04:05"))
	}

	fmt.Printf("\n[+] Loaded %d versions from versions.csv !\n\n", len(versionList))
}

func main() {
	var err error
	Config, err = filesystem.LoadConfigFile()

	if utils.HandleError(err) {
		panic(err)
	}

	// Load already saved versions
	loadVersions()

	for {
		version, err := checkForUpdate()
		if utils.HandleError(err) {
			continue
		}

		fmt.Printf("Version: %s\n", version)

		found := false
		for _, v := range versionList {
			if v.Version == version {
				found = true
				break
			}
		}

		if !found {
			v := &HcVersion{
				Version:     version,
				ReleaseDate: time.Now().Unix(),
			}

			jwt, err := ScrapeJwt(v)

			if utils.HandleError(err) {
				panic(err)
			}

			hc_version := strings.Split(jwt.VersionBaseUrl, "https://newassets.hcaptcha.com/c/")[1]

			v.HswVersion = hc_version

			versionList = append(versionList, v)
			go DownloadVersionFiles(v)
		}

		time.Sleep(1 * time.Minute)
	}
}
