package sdp_transform

import (
    "github.com/dlclark/regexp2"
)

func split(str, pattern string, limit int) []string {
    reg := regexp2.MustCompile(pattern, regexp2.None)
    var result []string
    lastIndex := 0

    // 用于存储匹配到的索引
    match, _ := reg.FindStringMatch(string(str))

    for match != nil {
        start, end := match.Index, match.Index+match.Length
        result = append(result, str[lastIndex:start])

        for _, g := range match.Groups()[1:] {
            if len(g.Captures) > 0 {
                result = append(result, g.Captures[0].String())
            }
        }

        lastIndex = end
        // 查找下一个匹配项
        match, _ = reg.FindNextMatch(match)
    }

    // 添加最后一部分
    result = append(result, str[lastIndex:])
    if limit > 0 {
        if limit > len(result) {
            limit = len(result)
        }
        return result[0:limit]
    }
    return result
}
