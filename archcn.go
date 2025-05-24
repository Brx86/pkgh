package main

import (
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"
)

const REPO_API string = "https://api.github.com/repos/archlinuxcn/repo/commits"
const REPO_URL string = "https://github.com/archlinuxcn/repo/tree/master/archlinuxcn/"

type commitInfo struct {
	Commit struct {
		Committer struct {
			Date time.Time
		}
		Message string
	}
	// HtmlUrl string `json:"html_url"`
}

func extractVersion(message string) string {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return ""
	}
	parts := strings.Split(lines[0], "auto updated to ")
	if len(parts) == 2 {
		return parts[1]
	}
	parts = strings.Split(lines[0], ": ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func getPkgCommits(pkgName string) []commitInfo {
	client := resty.New()
	client.SetQueryParam("per_page", "100").
		SetQueryParam("path", "archlinuxcn/"+pkgName+"/PKGBUILD").
		SetHeader("Accept", "application/vnd.github+json")
	if os.Getenv("GITHUB_TOKEN") == "" {
		println("警告: 未设置 GITHUB_TOKEN 环境变量。对于公开仓库可能仍然可以访问，但速率限制会更低。")
	} else {
		client.SetAuthToken(os.Getenv("GITHUB_TOKEN"))
	}
	result := []commitInfo{}
	for page := 1; ; page++ {
		client.SetQueryParam("page", strconv.Itoa(page))
		resp := &[]commitInfo{}
		_, err := client.R().SetResult(resp).Get(REPO_API)
		if err != nil {
			println("错误：请求 GITHUB API 出错，请检查网络或重试。")
			return nil
		}
		if len(*resp) == 0 {
			slices.Reverse(result)
			return result
		}
		result = append(result, *resp...)
	}
}

func MakeTableCN(pkgName string, render bool) string {
	data := getPkgCommits(pkgName)
	if len(data) == 0 {
		return ""
	}
	println(REPO_URL + pkgName)
	t := table.NewWriter()
	if render {
		t.SetOutputMirror(os.Stdout)
	}
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"Version", "Commit Time"})
	for _, c := range data {
		fmtTime := "\033[36m" + c.Commit.Committer.Date.Format("2006-01-02 15:04") + "\033[0m"
		version := "\033[32m" + extractVersion(c.Commit.Message) + "\033[0m"
		// url := "\033[90m" + c.HtmlUrl[:len(c.HtmlUrl)-34] + "\033[0m"
		t.AppendRow(table.Row{version, fmtTime})
	}
	return t.Render()
}
