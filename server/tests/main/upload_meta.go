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
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// 将行拆分为 MetaEnumCode 和剩余的部分
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		metaEnumCodeStr := parts[0]
		valueNameStr := parts[1]

		// 解析 MetaEnumCode
		metaEnumCode, err := strconv.Atoi(metaEnumCodeStr)
		if err != nil {
			fmt.Println("Failed to parse MetaEnumCode:", metaEnumCodeStr)
			continue
		}

		// 将剩余的部分按逗号分隔
		valueNameList := strings.Split(valueNameStr, ",")
		if len(valueNameList)%2 != 0 {
			fmt.Println("Invalid value,name pairs in line:", line)
			continue
		}

		var metas []*mysql.Meta
		for i := 0; i < len(valueNameList); i += 2 {
			valueStr := valueNameList[i]
			name := valueNameList[i+1]

			value, err := strconv.Atoi(valueStr)
			if err != nil {
				fmt.Println("Failed to parse Value:", valueStr)
				continue
			}

			// 构建 Meta 对象
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
