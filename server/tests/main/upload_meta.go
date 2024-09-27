package main

import (
	"bufio"
	"fmt"
	"onij/infra/mysql"
	"onij/logic"
	"os"
	"strconv"
	"strings"
	"time"
)

const MetaPath = "./meta.txt"

func uploadMeta() {
	file, err := os.Open(MetaPath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // 按空格拆分行数据
		if len(parts) < 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		// 解析MetaEnumCode
		metaEnumCode, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Failed to parse MetaEnumCode:", parts[0])
			continue
		}

		// 解析后续的Value和Name
		var metas []*mysql.Meta
		for i := 1; i < len(parts); i++ {
			valueNamePairs := strings.Split(parts[i], ",")
			if len(valueNamePairs) != 2 {
				fmt.Println("Invalid value,name format:", parts[i])
				continue
			}

			value, err := strconv.Atoi(valueNamePairs[0])
			if err != nil {
				fmt.Println("Failed to parse Value:", valueNamePairs[0])
				continue
			}
			// Name直接取
			name := valueNamePairs[1]

			// 构建Meta对象
			meta := mysql.Meta{
				MetaEnumCode: metaEnumCode,
				Value:        value,
				Name:         name,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			metas = append(metas, &meta)
		}
		err = logic.NewLocalLogic().SaveMeta(metas)
		if err != nil {
			fmt.Println("Failed to save meta:", err)
			return
		}
	}

	// 检查读取过程是否出错
	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading file:", err)
	}
}
