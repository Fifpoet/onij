// 列出并删除桶根目录下七牛 HLS 默认哈希文件夹，以及源片旁边的旧 *.mp4.m3u8。
// 默认 dry-run；加 -apply 才删除。
// 运行：cd server && go run ./cmd/cleanup_hls_hash_dirs [-apply]
package main

import (
	"flag"
	"fmt"
	"os"

	"onij/util"
)

func main() {
	apply := flag.Bool("apply", false, "真正删除；默认只打印")
	flag.Parse()

	roots, err := util.ListCommonPrefixes("", "/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "list root prefixes: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("root prefixes:")
	var hashDirs []string
	for _, p := range roots {
		keep := !util.IsHlsHashRootPrefix(p)
		mark := "keep"
		if !keep {
			mark = "HASH"
			hashDirs = append(hashDirs, p)
		}
		fmt.Printf("  [%s] %s\n", mark, p)
	}

	var toDelete []string
	for _, p := range hashDirs {
		keys, err := util.ListKeys(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "list %s: %v\n", p, err)
			os.Exit(1)
		}
		toDelete = append(toDelete, keys...)
		fmt.Printf("hash dir %s: %d objects\n", p, len(keys))
	}

	for _, prefix := range []string{"cloud/", "tran/"} {
		keys, err := util.ListKeys(prefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "list %s: %v\n", prefix, err)
			os.Exit(1)
		}
		n := 0
		for _, k := range keys {
			if util.IsLegacyHlsPlaylistKey(k) {
				toDelete = append(toDelete, k)
				n++
				fmt.Printf("legacy m3u8 %s\n", k)
			}
		}
		if n == 0 {
			fmt.Printf("no legacy m3u8 under %s\n", prefix)
		}
	}

	fmt.Printf("total delete candidates: %d\n", len(toDelete))
	if !*apply {
		fmt.Println("dry-run; pass -apply to delete")
		return
	}
	if err := util.DeleteKeys(toDelete); err != nil {
		fmt.Fprintf(os.Stderr, "delete: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("deleted %d objects\n", len(toDelete))
}
