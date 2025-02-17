package sdp_transform

import (
    "bufio"
    "encoding/json"
    "github.com/dlclark/regexp2"
    "strconv"
    "strings"
)

func attachProperties(match *regexp2.Match, location map[string]interface{}, names []string, rawName string) {
    if rawName != "" && len(names) == 0 {
        if len(match.Groups()) != 0 {
            location[rawName] = match.Groups()[1].String()
        } else {
            location[rawName] = match.String()
        }
    } else {
        for i, v := range names {
            val := match.Groups()[i+1].String()
            if val != "" {
                location[v] = val
            }
        }
    }
}

func parseReg(rule Rule, location map[string]interface{}, content string) {
    if rule.Push != "" {
        _, ok := location[rule.Push]
        if !ok {
            location[rule.Push] = []map[string]interface{}{}
        }
    }
    if rule.Name != "" {
        _, ok := location[rule.Name]
        if !ok {
            location[rule.Name] = map[string]interface{}{}
        }
    }
    keyLocatin := map[string]interface{}{}
    if rule.Push != "" {
        //
    } else if rule.Name != "" && len(rule.Names) != 0 {
        keyLocatin = location[rule.Name].(map[string]interface{})
    } else {
        keyLocatin = location
    }

    matched, err := rule.Reg.FindStringMatch(content)
    if err != nil {
        panic(err)
    }

    //if rule.Name != "" && len(rule.Names) == 0 {
    //    if len(matched.Groups()) != 0 {
    //        keyLocatin[rule.Name] = matched.Groups()[1].String()
    //    } else {
    //        keyLocatin[rule.Name] = matched.String()
    //    }
    //} else {
    //    for i, v := range rule.Names {
    //        val := matched.Groups()[i+1].String()
    //        if val != "" {
    //            keyLocatin[v] = val
    //        }
    //    }
    //}

    attachProperties(matched, keyLocatin, rule.Names, rule.Name)

    if rule.Push != "" {
        arr := location[rule.Push].([]map[string]interface{})
        arr = append(arr, keyLocatin)
        location[rule.Push] = arr
    }
}

func Parse(description string) (*SessionDescription, error) {
    session := map[string]interface{}{}
    media := make([]map[string]interface{}, 0)

    var location map[string]interface{}
    location = session

    validLine := regexp2.MustCompile(`^([a-z])=(.*)`, regexp2.None)
    scanner := bufio.NewScanner(strings.NewReader(description))
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue
        }
        if ok, _ := validLine.MatchString(line); !ok {
            continue
        }

        match, err := validLine.FindStringMatch(line)
        if err != nil {
            return nil, err
        }
        typ := match.Groups()[1].String()
        content := match.Groups()[2].String()

        if typ == "m" {
            media = append(media, map[string]interface{}{
                "rtp":  []map[string]interface{}{},
                "fmtp": []map[string]interface{}{},
            })
            location = media[len(media)-1]
        }

        for _, rule := range grammarMap[typ] {
            if ok, _ := rule.Reg.MatchString(content); ok {
                parseReg(*rule, location, content)
                break
            }
        }
    }
    session["media"] = media
    marshal, err := json.Marshal(session)
    if err != nil {
        return nil, err
    }
    var s SessionDescription
    err = json.Unmarshal(marshal, &s)
    if err != nil {
        return nil, err
    }
    return &s, nil
}

type ParamMap map[string]*string

func ParseParams(str string) ParamMap {
    paramMap := ParamMap{}
    params := split(str, `;\s?`, -1)
    for _, param := range params {
        s := split(param, `=(.+)`, 2)
        if len(s) == 2 {
            paramMap[s[0]] = &s[1]
        } else {
            paramMap[s[0]] = nil
        }
    }
    return ParamMap{}
}

var ParseFmtpConfig = ParseParams

func ParsePayloads(payloads string) []int {
    payloadsNums := make([]int, 0)
    for _, s := range strings.Split(payloads, " ") {
        i, err := strconv.ParseInt(s, 10, 32)
        if err != nil {
            panic(err)
        }
        payloadsNums = append(payloadsNums, int(i))
    }
    return payloadsNums
}

type RemoteCandidate struct {
    Component string `json:"component,omitempty"`
    IP        string `json:"ip,omitempty"`
    Port      string `json:"port,omitempty"`
}

func ParseRemoteCandidates(str string) []RemoteCandidate {
    candidates := make([]RemoteCandidate, 0)
    parts := strings.Split(str, " ")
    for i := 0; i < len(parts); i = i + 3 {
        candidates = append(candidates, RemoteCandidate{
            Component: parts[i],
            IP:        parts[i+1],
            Port:      parts[i+2],
        })
    }
    return candidates
}

func ParseImageAttributes(str string) []ParamMap {
    // attrs1: '[x=1280,y=720] [x=320,y=180]'
    paramMaps := make([]ParamMap, 0)
    for _, item := range strings.Split(str, " ") {
        item, _ = strings.CutPrefix(item, "[")
        item, _ = strings.CutSuffix(item, "]")
        paramMap := ParamMap{}
        for _, param := range strings.Split(item, ",") {
            s := split(param, `=(.+)`, 2)
            if len(s) == 2 {
                paramMap[s[0]] = &s[1]
            } else {
                paramMap[s[0]] = nil
            }
        }
        paramMaps = append(paramMaps, paramMap)
    }
    return paramMaps
}

type SimulcastStream struct {
    SCID   string `json:"scid,omitempty"`
    Paused bool   `json:"paused,omitempty"`
}

func ParseSimulcastStreamList(str string) [][]SimulcastStream {
    //  list1: '1,~4;2;3',
    list := make([][]SimulcastStream, 0)
    for _, stream := range strings.Split(str, ";") {
        simulcasts := make([]SimulcastStream, 0)
        for _, format := range strings.Split(stream, ",") {
            var scid string
            var paused bool
            if !strings.HasPrefix(format, "~") {
                scid = format
            } else {
                scid, _ = strings.CutPrefix(format, "~")
                paused = true
            }
            simulcasts = append(simulcasts, SimulcastStream{
                SCID:   scid,
                Paused: paused,
            })
        }
        list = append(list, simulcasts)
    }
    return list
}
