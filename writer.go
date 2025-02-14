package sdp_transform

import (
    "encoding/json"
    "fmt"
    "github.com/seamory/sdp-transform-go/pointer"
    "regexp"
    "strings"
)

// customized util.format - discards excess arguments and can void middle ones
func format(formatStr string, args ...string) string {
    i := 0
    l := len(args)
    reg := regexp.MustCompile("%[sdv%]")
    return reg.ReplaceAllStringFunc(formatStr, func(x string) string {
        if i > l-1 {
            return x
        }
        arg := args[i]
        i += 1
        switch x {
        case "%%":
            return "%"
        case "%s":
            return arg
        case "%d":
            return arg
        case "%v":
            return ""
        }
        return x
    })
}

func convertMapStringString(msi map[string]interface{}) map[string]string {
    m := make(map[string]string)
    for k, v := range msi {
        m[k] = v.(string)
    }
    return m
}

func makeLine(typ string, obj Rule, location map[string]interface{}) string {
    var str string
    if obj.Push != "" {
        str = obj.Format(convertMapStringString(location))
    } else {
        m, ok := location[obj.Name].(map[string]interface{})
        if !ok {
            str = obj.Format(nil)
        } else {
            str = obj.Format(convertMapStringString(m))
        }
    }

    args := make([]string, 0)
    args = append(args, fmt.Sprintf("%s=%s", typ, str))

    if len(obj.Names) != 0 {
        for i, name := range obj.Names {
            if obj.Name != "" {
                m, ok := location[obj.Name].(map[string]interface{})
                if !ok {
                    args = append(args, "")
                } else {
                    args = append(args, m[name].(string))
                }
            } else {
                s, ok := location[obj.Names[i]].(string)
                if !ok {
                    args = append(args, "")
                } else {
                    args = append(args, s)
                }

            }
        }
    } else {
        s, ok := location[obj.Name].(string)
        if !ok {
            args = append(args, "")
        } else {
            args = append(args, s)
        }
    }

    return format(args[0], args[1:]...)
}

// RFC specified order
// TODO: extend this with all the rest
var DefaultOuterOrder = []string{
    "v", "o", "s", "i",
    "u", "e", "p", "c",
    "b", "t", "r", "z", "a",
}
var DefaultInnerOrder = []string{"i", "c", "b", "a"}

type WriteOptions struct {
    OuterOrder []string
    InnerOrder []string
}

func Write(session SessionDescription, options *WriteOptions) string {

    if session.Version == nil {
        session.Version = pointer.String("")
    }
    for _, mLine := range session.Media {
        if mLine.Payloads == nil {
            mLine.Payloads = pointer.String("")
        }
    }

    marshal, err := json.Marshal(session)
    if err != nil {
        return ""
    }
    var s map[string]interface{}
    err = json.Unmarshal(marshal, &s)
    if err != nil {
        return ""
    }

    outerOrder := DefaultOuterOrder
    innerOrder := DefaultInnerOrder
    if options != nil && len(options.OuterOrder) != 0 {
        outerOrder = options.OuterOrder
    }
    if options != nil && len(options.InnerOrder) != 0 {
        innerOrder = options.InnerOrder
    }

    sdp := make([]string, 0)

    for _, typ := range outerOrder {
        for _, obj := range grammarMap[typ] {
            if v, ok := s[obj.Name]; ok && v != nil {
                sdp = append(sdp, makeLine(typ, *obj, s))
            } else if v, ok = s[obj.Push]; ok && v != nil {
                m := s[obj.Push].([]interface{})
                for _, val := range m {
                    sdp = append(sdp, makeLine(typ, *obj, val.(map[string]interface{})))
                }
            }
        }
    }

    medias := s["media"].([]interface{})
    for _, media := range medias {
        mLine := media.(map[string]interface{})
        sdp = append(sdp, makeLine("m", *grammarMap["m"][0], mLine))

        for _, typ := range innerOrder {
            for _, obj := range grammarMap[typ] {
                if v, ok := mLine[obj.Name]; ok && v != nil {
                    sdp = append(sdp, makeLine(typ, *obj, mLine))
                } else if v, ok = mLine[obj.Push]; ok && v != nil {
                    m := mLine[obj.Push].([]interface{})
                    for _, el := range m {
                        sdp = append(sdp, makeLine(typ, *obj, el.(map[string]interface{})))
                    }
                }
            }
        }
    }

    return fmt.Sprintf("%s\r\n", strings.Join(sdp, "\r\n"))
}
