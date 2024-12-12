package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// 从文件中读取以 "||" 或 "@@" 开头的ADG过滤规则，并返回字符串切片。
func readRulesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rules []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && (strings.HasPrefix(line, "||") || strings.HasPrefix(line, "@@")) {
			rules = append(rules, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// 从两个ADG过滤规则切片中移除重复项。
func removeDuplicates(rules1, rules2 []string) []string {
	uniqueRules := make(map[string]bool)
	for _, rule := range append(rules1, rules2...) {
		uniqueRules[rule] = true
	}

	// Convert map keys to slice
	var uniqueRulesSlice []string
	for rule := range uniqueRules {
		uniqueRulesSlice = append(uniqueRulesSlice, rule)
	}

	fmt.Printf("去重后剩余 %d 条规则\n", len(uniqueRulesSlice))
	return uniqueRulesSlice
}

// 将ADG过滤规则写入文件。
func writeRulesToFile(filename string, rules []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 添加注释头部
	header := fmt.Sprintf("! Title: AdGuard DNS filter & CHN: anti-AD\n! Last modified: %s\n! Total lines: %d\n", time.Now().Format("2006-01-02 15:04:05"), len(rules))
	_, err = writer.WriteString(header)
	if err != nil {
		return err
	}

	// 写入规则
	for _, rule := range rules {
		_, err := writer.WriteString(rule + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("已成功写入 %d 条规则到文件 %s\n", len(rules), filename)
	return nil
}

func main() {
	// 规则文件
	filePath1 := "../Rules/easylist.txt"
	filePath2 := "../Rules/filter.txt"
	// 输出文件
	outputFilePath := "../Rules/rules.txt"

	// 从文件中读取ADG过滤规则
	fmt.Println("开始读取规则文件...")
	rules1, err := readRulesFromFile(filePath1)
	if err != nil {
		fmt.Println("读取文件", filePath1, "时出错:", err)
		return
	}
	fmt.Printf("AdGuard DNS filter：共 %d 条规则\n", len(rules1))

	rules2, err := readRulesFromFile(filePath2)
	if err != nil {
		fmt.Println("读取文件", filePath2, "时出错:", err)
		return
	}
	fmt.Printf("CHN: anti-AD：共 %d 条规则\n", len(rules2))

	// 移除重复规则
	fmt.Println("开始去重...")
	uniqueRules := removeDuplicates(rules1, rules2)

	// 将有效的唯一规则写入新文件
	fmt.Println("开始写入最终规则文件...")
	err = writeRulesToFile(outputFilePath, uniqueRules)
	if err != nil {
		fmt.Println("写入文件", outputFilePath, "时出错:", err)
		return
	}
}
