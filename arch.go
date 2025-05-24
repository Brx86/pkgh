package main

import (
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	ARCHIVE_URL = "https://asia.archive.pkgbuild.com/packages/"
	PATTERN     = `-([^-]+-[0-9.]+)-(any|x86_64).(pkg.tar.[0-9a-z]*)<\/a>\s*(.*\d)\s{1,}(\w*)`
)

type pkgInfo struct {
	ts      int64
	fTime   string
	version string
	size    string
}

func fetchPkgHistory(pkgName string) [][]string {
	pattern := regexp.MustCompile(pkgName + PATTERN)
	realUrl := ARCHIVE_URL + pkgName[:1] + "/" + pkgName
	resp, err := resty.New().R().Get(realUrl)
	if err != nil {
		return nil
	}
	text := resp.String()
	return pattern.FindAllStringSubmatch(text, -1)
}

func matchPkg(pkgName string) []pkgInfo {
	historyData := fetchPkgHistory(pkgName)
	result := make([]pkgInfo, 0, len(historyData))
	for _, p := range historyData {
		pTime, _ := time.Parse("02-Jan-2006 15:04", p[4])
		fTime := pTime.Format("2006-01-02 15:04")
		result = append(result, pkgInfo{
			ts:      pTime.Unix(),
			fTime:   "\033[36m" + fTime + "\033[0m",
			version: "\033[32m" + p[1] + "\033[0m",
			size:    "\033[90m" + p[5] + "\033[0m",
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ts < result[j].ts
	})
	return result
}

func MakeTable(pkgName string, render bool) string {
	data := matchPkg(pkgName)
	if len(data) == 0 {
		return ""
	}
	println("https://archive.aya1.de/?p=" + pkgName)
	t := table.NewWriter()
	if render {
		t.SetOutputMirror(os.Stdout)
	}
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"Version", "Update Time", "Size"})
	for _, p := range data {
		t.AppendRow(table.Row{p.version, p.fTime, p.size})
	}
	return t.Render()
}
