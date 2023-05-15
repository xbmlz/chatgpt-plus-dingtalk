package utils

import "strings"

func ExtractStringBetween(content, start, end string) (ret string) {
	startIndex := strings.Index(content, start)
	endIndex := strings.LastIndex(content, end)
	if startIndex != -1 || endIndex != -1 {
		ret = content[startIndex+len(start) : endIndex]
	}
	return
}

func AfterString(after, content string) (ret string) {
	afterIndex := strings.LastIndex(content, after)
	if afterIndex != -1 {
		// 从后往前一直到index的内容
		ret = content[afterIndex+len(after):]
	}
	return
}
