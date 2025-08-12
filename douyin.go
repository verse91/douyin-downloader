package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

const cVUrl = "https://www.iesdouyin.com/aweme/v1/play/?video_id=%s&ratio=1080p&line=0"
// add more info here if you want music, thumbnail, ...
type VideoInfo struct {
	Type         string   `json:"type"`
	VideoID      string   `json:"video_id,omitempty"`
	Like         int      `json:"like,omitempty"`
	Comment      int      `json:"comment,omitempty"`
	Collect      int      `json:"save,omitempty"`
	Share        int      `json:"share,omitempty"`
	VideoDesc    string   `json:"video_desc,omitempty"`
	CreateAt     string   `json:"create_at,omitempty"`
	VideoURL     string   `json:"video_url,omitempty"`
	ImageURLList []string `json:"img_url,omitempty"`
}

type UserInfo struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func formatDate(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func doGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseImgList(body string) []string {
	content := strings.ReplaceAll(strings.ReplaceAll(body, "\\u002F", "/"), "\\/", "/")
	reg := regexp.MustCompile(`{"uri":"[^\s"]+","url_list":\["(https://p\d{1,2}-sign\.douyinpic\.com/.*?)"`)
	urlRet := regexp.MustCompile(`"uri":"([^\s"]+)","url_list":`)
	firstUrls := reg.FindAllStringSubmatch(content, -1)
	urlMatches := urlRet.FindAllStringSubmatch(content, -1)

	var firstURLs, urlList []string
	for _, match := range firstUrls {
		if len(match) > 1 {
			firstURLs = append(firstURLs, match[1])
		}
	}
	for _, match := range urlMatches {
		if len(match) > 1 {
			urlList = append(urlList, match[1])
		}
	}

	urlSet := make(map[string]bool)
	for _, url := range urlList {
		urlSet[url] = true
	}

	var rList []string
	for urlSetKey := range urlSet {
		for _, firstURL := range firstURLs {
			if strings.Contains(firstURL, urlSetKey) {
				rList = append(rList, firstURL)
				break
			}
		}
	}

	var filteredRList []string
	for _, url := range rList {
		if !strings.Contains(url, "/obj/") {
			filteredRList = append(filteredRList, url)
		}
	}
	return filteredRList
}

func getVideoInfo(url string) (*VideoInfo, *UserInfo, error) {
	typeStr := "video"
	var imgList []string
	var videoURL string

	body, err := doGet(url)
	if err != nil {
		return nil, nil, err
	}

	pattern := regexp.MustCompile(`"video":{"play_addr":{"uri":"([a-z0-9]+)"`)
	match := pattern.FindStringSubmatch(body)
	if len(match) < 2 {
		typeStr = "img"
	} else {
		videoURL = fmt.Sprintf(cVUrl, match[1])
	}
	if typeStr == "img" {
		imgList = parseImgList(body)
	}

	statsRegex := regexp.MustCompile(`"statistics"\s*:\s*\{([\s\S]*?)\},`)
	statsMatch := statsRegex.FindStringSubmatch(body)
	if len(statsMatch) < 2 {
		return nil, nil, fmt.Errorf("no stats found in the response")
	}
	innerContent := statsMatch[1]

	awemeIDRegex := regexp.MustCompile(`"aweme_id"\s*:\s*"([^"]+)"`)
	commentCountRegex := regexp.MustCompile(`"comment_count"\s*:\s*(\d+)`)
	diggCountRegex := regexp.MustCompile(`"digg_count"\s*:\s*(\d+)`)
	shareCountRegex := regexp.MustCompile(`"share_count"\s*:\s*(\d+)`)
	collectCountRegex := regexp.MustCompile(`"collect_count"\s*:\s*(\d+)`)
	authorRegex := regexp.MustCompile(`"nickname":\s*"([^"]+)",\s*"signature":\s*"([^"]+)"`)
	usernameRegex := regexp.MustCompile(`"unique_id":\s*"([^"]+)"`)
	createTimeRegex := regexp.MustCompile(`"create_time":\s*(\d+)`)
	descRegex := regexp.MustCompile(`"desc":\s*"([^"]+)"`)

	videoInfo := &VideoInfo{
		Type:         typeStr,
		VideoURL:     videoURL,
		ImageURLList: imgList,
	}
	userInfo := &UserInfo{}

	if m := awemeIDRegex.FindStringSubmatch(innerContent); len(m) > 1 {
		videoInfo.VideoID = m[1]
	}
	if m := commentCountRegex.FindStringSubmatch(innerContent); len(m) > 1 {
		if count, err := strconv.Atoi(m[1]); err == nil {
			videoInfo.Comment = count
		}
	}
	if m := diggCountRegex.FindStringSubmatch(innerContent); len(m) > 1 {
		if count, err := strconv.Atoi(m[1]); err == nil {
			videoInfo.Like = count
		}
	}
	if m := shareCountRegex.FindStringSubmatch(innerContent); len(m) > 1 {
		if count, err := strconv.Atoi(m[1]); err == nil {
			videoInfo.Share = count
		}
	}
	if m := collectCountRegex.FindStringSubmatch(innerContent); len(m) > 1 {
		if count, err := strconv.Atoi(m[1]); err == nil {
			videoInfo.Collect = count
		}
	}
	if m := authorRegex.FindStringSubmatch(body); len(m) > 2 {
		userInfo.Name = m[1]
		userInfo.Bio = m[2]
	}
	if m := usernameRegex.FindStringSubmatch(body); len(m) > 1 {
		userInfo.Username = m[1]
	}
	if m := createTimeRegex.FindStringSubmatch(body); len(m) > 1 {
		if ts, err := strconv.ParseInt(m[1], 10, 64); err == nil {
			videoInfo.CreateAt = formatDate(ts)
		}
	}
	if m := descRegex.FindStringSubmatch(body); len(m) > 1 {
		videoInfo.VideoDesc = m[1]
	}

	return videoInfo, userInfo, nil
}

func main() {
	// url := "https://v.douyin.com/AUVa3G-5QUE/"
	url := "https://v.douyin.com/9N2HGwrYB70/"
	videoInfo, userInfo, err := getVideoInfo(url)
	if err != nil {
		fmt.Printf("Error getting video info: %v\n", err)
		return
	}

	output := struct {
		Media interface{} `json:"media"`
		User  interface{} `json:"user"`
	}{
		Media: videoInfo,
		User:  userInfo,
	}

	resultJson, err := sonic.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling output JSON: %v\n", err)
		return
	}
	fmt.Println(string(resultJson))
}
