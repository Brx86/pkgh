package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	ARCHIVE_URL = "https://asia.archive.pkgbuild.com/packages"
	PATTERN     = `-([^-]+-[0-9.]+)-(any|x86_64).(pkg.tar.[0-9a-z]*)<\/a>\s*(.*\d)\s{1,}(\w*)`
)

func FetchPkgHistory(pkgName string) [][]string {
	pattern := regexp.MustCompile(pkgName + PATTERN)
	realUrl := fmt.Sprintf("%s/%s/%s/", ARCHIVE_URL, string(pkgName[0]), pkgName)
	resp, _ := http.Get(realUrl)
	defer resp.Body.Close()
	text, _ := io.ReadAll(resp.Body)
	return pattern.FindAllStringSubmatch(string(text), -1)
}

type PkgInfo struct {
	ts      int64
	fTime   string
	version string
	size    string
}

func MatchPkg(pkgName string) []PkgInfo {
	result := make([]PkgInfo, 0)
	for _, p := range FetchPkgHistory(pkgName) {
		pTime, _ := time.Parse("02-Jan-2006 15:04", p[4])
		fTime := pTime.Format("2006-01-02 15:04")
		result = append(result, PkgInfo{
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

func MakeTable(pkgName string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"Version", "Update Time", "Size"})
	for _, p := range MatchPkg(pkgName) {
		t.AppendRow(table.Row{p.version, p.fTime, p.size})
	}
	t.Render()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<package>")
		return
	}
	MakeTable(os.Args[1])
}
