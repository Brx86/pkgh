package main

import (
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 || strings.HasPrefix(os.Args[1], "-") {
		println("Search package history in arch and archlinuxcn\n")
		println("Usage:", os.Args[0], "<package>")
		return
	}
	pkgName := os.Args[1]
	archChan := make(chan string)
	archChanCN := make(chan string)
	resultChan := make(chan string)
	go func() {
		archChanCN <- MakeTableCN(pkgName, false)
	}()
	go func() {
		archChan <- MakeTable(pkgName, false)
	}()
	go func() {
		for range 2 {
			select {
			case result := <-archChanCN:
				if result != "" {
					resultChan <- result
					return
				}
			case result := <-archChan:
				if result != "" {
					resultChan <- result
					return
				}
			}
		}
		resultChan <- "错误：Arch 官方仓库与 archlinuxcn 仓库中未找到此包。"
	}()
	select {
	case result := <-resultChan:
		println(result)
	case <-time.After(10 * time.Second):
		println("错误：网络超时")
	}
}
