package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wd, _ := os.Getwd()
	fmt.Println(wd)

	md := filepath.Dir(wd)
	fmt.Println(md)

	sqls, err := findSQLs(wd)
	if err != nil {
		fmt.Println(err)
		return
	}

	excludedSQLs := map[string]struct{}{}
	for name, path := range sqls {
		if _, ok := excludedSQLs[name]; ok {
			continue
		}

		g := &generator{DestPath: md, DdlPath: path}
		g.Run()
	}
}

func findSQLs(dir string) (map[string]string, error) {
	var kvs = make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), "sql") {
			kvs[info.Name()] = path
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return kvs, err
}
