package util

import (
	"regexp"
	"strings"
)

// 一个简化的正则：匹配 4 类 XML "标签或声明"
// 1) 注释：       <!-- ... -->
// 2) Doctype：    <!DOCTYPE ...>
// 3) 处理指令：    <? ... ?>
// 4) 普通标签(含自闭合)：<...>或<.../>
var reTag = regexp.MustCompile(`(?s)<!--.*?-->|<!DOCTYPE.*?>|<\?.*?\?>|<[^>]+>`)

// removeUnmatchedClosingTags 会：
//   - 原样保留注释 <!--...-->、<!DOCTYPE ...>、<?...?> 声明
//   - 对普通标签：
//   - 如果是自闭合，如 <tag .../>，直接写回，不入栈
//   - 如果是开标签，如 <tag ...>，压栈等待匹配
//   - 如果是闭标签，如 </tag>，与栈顶比对，匹配则弹栈，否则删除
func removeUnmatchedClosingTags(xmlStr string) string {
	var stack []string
	var sb strings.Builder

	lastIndex := 0
	matches := reTag.FindAllStringIndex(xmlStr, -1)
	if matches == nil {
		// 没有任何 "<...>" 类片段，说明基本没标签可处理
		return xmlStr
	}

	for _, m := range matches {
		start, end := m[0], m[1]
		// 先把 [lastIndex, start) 的纯文本加入输出
		sb.WriteString(xmlStr[lastIndex:start])

		tagText := xmlStr[start:end]
		// 注释/DOCTYPE/处理指令 => 原样保留
		if strings.HasPrefix(tagText, "<!--") ||
			strings.HasPrefix(tagText, "<!DOCTYPE") ||
			strings.HasPrefix(tagText, "<?") {
			sb.WriteString(tagText)
			lastIndex = end
			continue
		}

		// 普通标签(可能 <tag>, <tag/>, 或 </tag>)
		content := strings.TrimSpace(tagText[1 : len(tagText)-1]) // 去掉 < 和 >
		if len(content) == 0 {
			// 极罕见: <> 空标签
			lastIndex = end
			continue
		}

		// 判断是否自闭合，如 <tag .../>
		isSelfClose := strings.HasSuffix(content, "/")
		if isSelfClose {
			// 去掉尾部 "/", 只用来识别标签名
			content = strings.TrimSpace(content[:len(content)-1])
		}

		if content[0] == '/' {
			// 闭合标签，如 </tagName>
			name := strings.TrimSpace(content[1:])
			// 找出真正的标签名(去掉可能的空格、属性——标准闭合标签一般没属性，这里兼容处理)
			if idx := strings.IndexAny(name, " \t\r\n"); idx != -1 {
				name = name[:idx]
			}
			// 和栈顶匹配
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == name {
					// 栈顶匹配 => 弹栈 + 保留该闭合标签
					stack = stack[:len(stack)-1]
					sb.WriteString(tagText)
				} else {
					// 不匹配 => 多余闭合标签 => 不写入 sb
				}
			} else {
				// 栈空 => 多余闭合标签 => 跳过
			}
		} else {
			// 开标签
			// 取第一个空格前内容作为标签名(可能含属性, 这里只截取前面部分)
			name := content
			if idx := strings.IndexAny(name, " \t\r\n"); idx != -1 {
				name = name[:idx]
			}
			// 自闭合标签不入栈
			if !isSelfClose {
				stack = append(stack, name)
			}
			sb.WriteString(tagText)
		}
		lastIndex = end
	}

	// 最后拼接 剩余的纯文本
	sb.WriteString(xmlStr[lastIndex:])
	return sb.String()
}
