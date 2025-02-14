package sdp_transform

import (
    "github.com/dlclark/regexp2"
    "strings"
)

type MapStringString map[string]string

func (m MapStringString) has(key string) bool {
    v, ok := m[key]
    if ok && v != "" {
        return true
    }
    return false
}

type Rule struct {
    Name   string
    Push   string
    Reg    *regexp2.Regexp
    Names  []string
    Format func(m map[string]string) string
}

// https://github.com/clux/sdp-transform/blob/master/lib/grammar.js

type GrammarMap map[string][]*Rule

func regexp2MustCompile(str string) *regexp2.Regexp {
    return regexp2.MustCompile(str, regexp2.None)
}

var grammarMap GrammarMap = map[string][]*Rule{
    "v": {
        {
            Name: "version",
            Reg:  regexp2MustCompile(`^(\d*)$`),
        },
    },
    "o": {
        {
            Name:  "origin",
            Reg:   regexp2MustCompile(`^(\S*) (\d*) (\d*) (\S*) IP(\d) (\S*)`),
            Names: []string{"username", "sessionId", "sessionVersion", "netType", "ipVer", "address"},
            Format: func(m map[string]string) string {
                return "%s %s %d %s IP%d %s"
            },
        },
    },
    "s": {
        {
            Name: "name",
        },
    },
    "i": {
        {
            Name: "description",
        },
    },
    "u": {
        {
            Name: "uri",
        },
    },
    "e": {
        {
            Name: "email",
        },
    },
    "p": {
        {
            Name: "phone",
        },
    },
    "z": {
        {
            Name: "timezones",
        },
    },
    "r": {
        {
            Name: "repeats",
        },
    },
    "t": {
        {
            // t=0 0
            Name:  "timing",
            Reg:   regexp2MustCompile(`^(\d*) (\d*)`),
            Names: []string{"start", "stop"},
            Format: func(m map[string]string) string {
                return "%d %d"
            },
        },
    },
    "c": {
        {
            // c=IN IP4 10.47.197.26
            Name:  "connection",
            Reg:   regexp2MustCompile(`^IN IP(\d) (\S*)`),
            Names: []string{"version", "ip"},
            Format: func(m map[string]string) string {
                return "IN IP%d %s"
            },
        },
    },
    "b": {
        {
            // b=AS:4000
            Push:  "bandwidth",
            Reg:   regexp2MustCompile(`^(TIAS|AS|CT|RR|RS):(\d*)`),
            Names: []string{"type", "limit"},
            Format: func(m map[string]string) string {
                return "%s:%s"
            },
        },
    },
    "m": {
        {
            // m=video 51744 RTP/AVP 126 97 98 34 31
            // NB: special - pushes to session
            // TODO: rtp/fmtp should be filtered by the payloads found here?
            Reg:   regexp2MustCompile(`^(\w*) (\d*) ([\w/]*)(?: (.*))?`),
            Names: []string{"type", "port", "protocol", "payloads"},
            Format: func(m map[string]string) string {
                return "%s %d %s %s"
            },
        },
    },
    "a": {
        {
            // a=rtpmap:110 opus/48000/2
            Push:  "rtp",
            Reg:   regexp2MustCompile(`^rtpmap:(\d*) ([\w\-.]*)(?:\s*\/(\d*)(?:\s*\/(\S*))?)?`),
            Names: []string{"payload", "codec", "rate", "encoding"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("encoding") {
                    return "rtpmap:%d %s/%s/%s"
                }
                if mss.has("rate") {
                    return "rtpmap:%d %s/%s"
                }
                return "rtpmap:%d %s"
            },
        },
        {
            // a=fmtp:108 profile-level-id=24;object=23;bitrate=64000
            // a=fmtp:111 minptime=10; useinbandfec=1
            Push:  "fmtp",
            Reg:   regexp2MustCompile(`^fmtp:(\d*) ([\S| ]*)`),
            Names: []string{"payload", "config"},
            Format: func(m map[string]string) string {
                return "fmtp:%d %s"
            },
        },
        {
            // a=control:streamid=0
            Name: "control",
            Reg:  regexp2MustCompile(`^control:(.*)`),
            Format: func(m map[string]string) string {
                return "control:%s"
            },
        },
        {
            // a=rtcp:65179 IN IP4 193.84.77.194
            Name:  "rtcp",
            Reg:   regexp2MustCompile(`^rtcp:(\d*)(?: (\S*) IP(\d) (\S*))?`),
            Names: []string{"port", "netType", "ipVer", "address"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("address") {
                    return "rtcp:%d %s IP%d %s"
                }
                return "rtcp:%d"
            },
        },
        {
            // a=rtcp-fb:98 trr-int 100
            Push:  "rtcpFbTrrInt",
            Reg:   regexp2MustCompile(`^rtcp-fb:(\*|\d*) trr-int (\d*)`),
            Names: []string{"payload", "value"},
            Format: func(m map[string]string) string {
                return "rtcp-fb:%s trr-int %d"
            },
        },
        {
            // a=rtcp-fb:98 nack rpsi
            Push:  "rtcpFb",
            Reg:   regexp2MustCompile(`^rtcp-fb:(\*|\d*) ([\w-_]*)(?: ([\w-_]*))?`),
            Names: []string{"payload", "type", "subtype"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("subtype") {
                    return "rtcp-fb:%s %s %s"
                }
                return "rtcp-fb:%s %s"
            },
        },
        {
            // a=extmap:2 urn:ietf:params:rtp-hdrext:toffset
            // a=extmap:1/recvonly URI-gps-string
            // a=extmap:3 urn:ietf:params:rtp-hdrext:encrypt urn:ietf:params:rtp-hdrext:smpte-tc 25@600/24
            Push:  "ext",
            Reg:   regexp2MustCompile(`^extmap:(\d+)(?:\/(\w+))?(?: (urn:ietf:params:rtp-hdrext:encrypt))? (\S*)(?: (\S*))?`),
            Names: []string{"value", "direction", "encrypt-uri", "uri", "config"},
            Format: func(m map[string]string) string {
                sb := strings.Builder{}
                sb.WriteString("extmap:%d")

                mss := MapStringString(m)

                if mss.has("direction") {
                    sb.WriteString("/%s")
                } else {
                    sb.WriteString("%v")
                }

                if mss.has("encrypt-uri") {
                    sb.WriteString(" %s")
                } else {
                    sb.WriteString("%v")
                }

                sb.WriteString(" %s")

                if mss.has("config") {
                    sb.WriteString("%s")
                }

                return sb.String()
            },
        },
        {
            // a=extmap-allow-mixed
            Name: "extmapAllowMixed",
            Reg:  regexp2MustCompile(`^(extmap-allow-mixed)`),
        },
        {
            Push:  "crypto",
            Reg:   regexp2MustCompile(`^crypto:(\d*) ([\w_]*) (\S*)(?: (\S*))?`),
            Names: []string{"id", "suite", "config", "sessionConfig"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("sessionConfig") {
                    return "crypto:%d %s %s %s"
                } else {
                    return "crypto:%d %s %s"
                }
            },
        },
        {
            // a=setup:actpass
            Name: "setup",
            Reg:  regexp2MustCompile(`^setup:(\w*)`),
            Format: func(m map[string]string) string {
                return "setup:%s"
            },
        },
        {
            // a=connection:new
            Name: "connectionType",
            Reg:  regexp2MustCompile(`^connection:(new|existing)`),
            Format: func(m map[string]string) string {
                return "connection:%s"
            },
        },
        {
            // a=mid:1
            Name: "mid",
            Reg:  regexp2MustCompile(`^mid:([^\s]*)`),
            Format: func(m map[string]string) string {
                return "mid:%s"
            },
        },
        {
            // a=msid:0c8b064d-d807-43b4-b434-f92a889d8587 98178685-d409-46e0-8e16-7ef0db0db64a
            Name: "msid",
            Reg:  regexp2MustCompile(`^msid:(.*)`),
            Format: func(m map[string]string) string {
                return "msid:%s"
            },
        },
        {
            // a=ptime:20
            Name: "ptime",
            Reg:  regexp2MustCompile(`^ptime:(\d*(?:\.\d*)*)`),
            Format: func(m map[string]string) string {
                return "ptime:%d"
            },
        },
        {
            // a=maxptime:60
            Name: "maxptime",
            Reg:  regexp2MustCompile(`^maxptime:(\d*(?:\.\d*)*)`),
            Format: func(m map[string]string) string {
                return "maxptime:%d"
            },
        },
        {
            // a=sendrecv
            Name: "direction",
            Reg:  regexp2MustCompile(`^(sendrecv|recvonly|sendonly|inactive)`),
        },
        {
            // a=ice-lite
            Name: "icelite",
            Reg:  regexp2MustCompile(`^(ice-lite)`),
        },
        {
            // a=ice-ufrag:F7gI
            Name: "iceUfrag",
            Reg:  regexp2MustCompile(`^ice-ufrag:(\S*)`),
            Format: func(m map[string]string) string {
                return "ice-ufrag:%s"
            },
        },
        {
            // a=ice-pwd:x9cml/YzichV2+XlhiMu8g
            Name: "icePwd",
            Reg:  regexp2MustCompile(`^ice-pwd:(\S*)`),
            Format: func(m map[string]string) string {
                return "ice-pwd:%s"
            },
        },
        {
            // a=fingerprint:SHA-1 00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33
            Name:  "fingerprint",
            Reg:   regexp2MustCompile(`^fingerprint:(\S*) (\S*)`),
            Names: []string{"type", "hash"},
            Format: func(m map[string]string) string {
                return "fingerprint:%s %s"
            },
        },
        {
            // a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host
            // a=candidate:1162875081 1 udp 2113937151 192.168.34.75 60017 typ host generation 0 network-id 3 network-cost 10
            // a=candidate:3289912957 2 udp 1845501695 193.84.77.194 60017 typ srflx raddr 192.168.34.75 rport 60017 generation 0 network-id 3 network-cost 10
            // a=candidate:229815620 1 tcp 1518280447 192.168.150.19 60017 typ host tcptype active generation 0 network-id 3 network-cost 10
            // a=candidate:3289912957 2 tcp 1845501695 193.84.77.194 60017 typ srflx raddr 192.168.34.75 rport 60017 tcptype passive generation 0 network-id 3 network-cost 10
            Push:  "candidates",
            Reg:   regexp2MustCompile(`^candidate:(\S*) (\d*) (\S*) (\d*) (\S*) (\d*) typ (\S*)(?: raddr (\S*) rport (\d*))?(?: tcptype (\S*))?(?: generation (\d*))?(?: network-id (\d*))?(?: network-cost (\d*))?`),
            Names: []string{"foundation", "component", "transport", "priority", "ip", "port", "type", "raddr", "rport", "tcptype", "generation", "network-id", "network-cost"},
            Format: func(m map[string]string) string {
                sb := strings.Builder{}
                sb.WriteString("candidate:%s %d %s %d %s %d typ %s")

                mss := MapStringString(m)
                if mss.has("raddr") {
                    sb.WriteString(" raddr %s rport %d")
                } else {
                    sb.WriteString("%v%v")
                }

                if mss.has("tcptype") {
                    sb.WriteString(" tcptype %s")
                } else {
                    sb.WriteString("%v")
                }

                if mss.has("generation") {
                    sb.WriteString(" generation %d")
                }

                if mss.has("network-id") {
                    sb.WriteString(" network-id %d")
                } else {
                    sb.WriteString("%v")
                }

                if mss.has("network-cost") {
                    sb.WriteString(" network-cost %d")
                } else {
                    sb.WriteString("%v")
                }

                return sb.String()
            },
        },
        {
            // a=end-of-candidates (keep after the candidates line for readability)
            Name: "endOfCandidates",
            Reg:  regexp2MustCompile(`^(end-of-candidates)`),
        },
        {
            // a=remote-candidates:1 203.0.113.1 54400 2 203.0.113.1 54401 ...
            Name: "remoteCandidates",
            Reg:  regexp2MustCompile(`^remote-candidates:(.*)`),
            Format: func(m map[string]string) string {
                return "remote-candidates:%s"
            },
        },
        {
            // a=ice-options:google-ice
            Name: "iceOptions",
            Reg:  regexp2MustCompile(`^ice-options:(\S*)`),
            Format: func(m map[string]string) string {
                return "ice-options:%s"
            },
        },
        {
            // a=ssrc:2566107569 cname:t9YU8M1UxTF8Y1A1
            Push:  "ssrcs",
            Reg:   regexp2MustCompile(`^ssrc:(\d*) ([\w_-]*)(?::(.*))?`),
            Names: []string{"id", "attribute", "value"},
            Format: func(m map[string]string) string {
                sb := strings.Builder{}
                sb.WriteString("ssrc:%d")

                mss := MapStringString(m)
                if mss.has("attribute") {
                    sb.WriteString(" %s")
                    if mss.has("value") {
                        sb.WriteString(":%s")
                    }
                }
                return sb.String()
            },
        },
        {
            // a=ssrc-group:FEC 1 2
            // a=ssrc-group:FEC-FR 3004364195 1080772241
            Push: "ssrcGroups",
            // token-char = %x21 / %x23-27 / %x2A-2B / %x2D-2E / %x30-39 / %x41-5A / %x5E-7E
            Reg:   regexp2MustCompile(`^ssrc-group:([\x21\x23\x24\x25\x26\x27\x2A\x2B\x2D\x2E\w]*) (.*)`),
            Names: []string{"semantics", "ssrcs"},
            Format: func(m map[string]string) string {
                return "ssrc-group:%s %s"
            },
        },
        {
            // a=msid-semantic: WMS Jvlam5X3SX1OP6pn20zWogvaKJz5Hjf9OnlV
            Name:  "msidSemantic",
            Reg:   regexp2MustCompile(`^msid-semantic:\s?(\w*) (\S*)`),
            Names: []string{"semantic", "token"},
            Format: func(m map[string]string) string {
                return "msid-semantic: %s %s" // space after ":" is not accidental
            },
        },
        {
            // a=group:BUNDLE audio video
            Push:  "groups",
            Reg:   regexp2MustCompile(`^group:(\w*) (.*)`),
            Names: []string{"type", "mids"},
            Format: func(m map[string]string) string {
                return "group:%s %s"
            },
        },
        {
            // a=rtcp-mux
            Name: "rtcpMux",
            Reg:  regexp2MustCompile(`^(rtcp-mux)`),
        },
        {
            // a=rtcp-rsize
            Name: "rtcpRsize",
            Reg:  regexp2MustCompile(`^(rtcp-rsize)`),
        },
        {
            // a=sctpmap:5000 webrtc-datachannel 1024
            Name:  "sctpmap",
            Reg:   regexp2MustCompile(`^sctpmap:([\w_/]*) (\S*)(?: (\S*))?`),
            Names: []string{"sctpmapNumber", "app", "maxMessageSize"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("maxMessageSize") {
                    return "sctpmap:%s %s %s"
                } else {
                    return "sctpmap:%s %s"
                }
            },
        },
        {
            // a=x-google-flag:conference
            Name: "xGoogleFlag",
            Reg:  regexp2MustCompile(`^x-google-flag:([^\s]*)`),
            Format: func(m map[string]string) string {
                return "x-google-flag:%s"
            },
        },
        {
            // a=rid:1 send max-width=1280;max-height=720;max-fps=30;depend=0
            Push:  "rids",
            Reg:   regexp2MustCompile(`^rid:([\d\w]+) (\w+)(?: ([\S| ]*))?`),
            Names: []string{"id", "direction", "params"},
            Format: func(m map[string]string) string {
                mss := MapStringString(m)
                if mss.has("params") {
                    return "rid:%s %s %s"
                } else {
                    return "rid:%s %s"
                }
            },
        },
        {
            // a=imageattr:97 send [x=800,y=640,sar=1.1,q=0.6] [x=480,y=320] recv [x=330,y=250]
            // a=imageattr:* send [x=800,y=640] recv *
            // a=imageattr:100 recv [x=320,y=240]
            Push: "imageattrs",
            Reg: regexp2MustCompile(
                `^imageattr:(\d+|\*)` +
                    `[\s\t]+(send|recv)[\s\t]+(\*|\[\S+\](?:[\s\t]+\[\S+\])*)` +
                    `(?:[\s\t]+(recv|send)[\s\t]+(\*|\[\S+\](?:[\s\t]+\[\S+\])*))?`,
            ),
            Names: []string{"pt", "dir1", "attrs1", "dir2", "attrs2"},
            Format: func(m map[string]string) string {
                return "imageattr:%s %s %s' + (o.dir2 ? ' %s %s' : '')"
            },
        },
        {
            // a=simulcast:send 1,2,3;~4,~5 recv 6;~7,~8
            // a=simulcast:recv 1;4,5 send 6;7
            Name: "simulcast",
            Reg: regexp2MustCompile(
                `^simulcast:` +
                    `(send|recv) ([a-zA-Z0-9\-_~;,]+)` +
                    `(?:\s?(send|recv) ([a-zA-Z0-9\-_~;,]+))?` +
                    `$`,
            ),
            Names: []string{"dir1", "list1", "dir2", "list2"},
            Format: func(m map[string]string) string {
                return "simulcast:%s %s' + (o.dir2 ? ' %s %s' : '')"
            },
        },
        {
            // old simulcast draft 03 (implemented by Firefox)
            //   https://tools.ietf.org/html/draft-ietf-mmusic-sdp-simulcast-03
            // a=simulcast: recv pt=97;98 send pt=97
            // a=simulcast: send rid=5;6;7 paused=6,7
            Name:  "simulcast_03",
            Reg:   regexp2MustCompile(`^simulcast:[\s\t]+([\S+\s\t]+)$`),
            Names: []string{"value"},
            Format: func(m map[string]string) string {
                return "simulcast: %s"
            },
        },
        {
            // a=framerate:25
            // a=framerate:29.97
            Name: "framerate",
            Reg:  regexp2MustCompile(`^framerate:(\d+(?:$|\.\d+))`),
            Format: func(m map[string]string) string {
                return "framerate:%s"
            },
        },
        {
            // RFC4570
            // a=source-filter: incl IN IP4 239.5.2.31 10.1.15.5
            Name:  "sourceFilter",
            Reg:   regexp2MustCompile(`^source-filter: *(excl|incl) (\S*) (IP4|IP6|\*) (\S*) (.*)`),
            Names: []string{"filterMode", "netType", "addressTypes", "destAddress", "srcList"},
            Format: func(m map[string]string) string {
                return "source-filter: %s %s %s %s %s"
            },
        },
        {
            // a=bundle-only
            Name: "bundleOnly",
            Reg:  regexp2MustCompile(`^(bundle-only)`),
        },
        {
            // a=label:1
            Name: "label",
            Reg:  regexp2MustCompile(`^label:(.+)`),
            Format: func(m map[string]string) string {
                return "label:%s"
            },
        },
        {
            // RFC version 26 for SCTP over DTLS
            // https://tools.ietf.org/html/draft-ietf-mmusic-sctp-sdp-26#section-5
            Name: "sctpPort",
            Reg:  regexp2MustCompile(`^sctp-port:(\d+)$`),
            Format: func(m map[string]string) string {
                return "sctp-port:%s"
            },
        },
        {
            // RFC version 26 for SCTP over DTLS
            // https://tools.ietf.org/html/draft-ietf-mmusic-sctp-sdp-26#section-6
            Name: "maxMessageSize",
            Reg:  regexp2MustCompile(`^max-message-size:(\d+)$`),
            Format: func(m map[string]string) string {
                return "max-message-size:%s"
            },
        },
        {
            // RFC7273
            // a=ts-refclk:ptp=IEEE1588-2008:39-A7-94-FF-FE-07-CB-D0:37
            Push:  "tsRefClocks",
            Reg:   regexp2MustCompile(`^ts-refclk:([^\s=]*)(?:=(\S*))?`),
            Names: []string{"clksrc", "clksrcExt"},
            Format: func(m map[string]string) string {
                sb := strings.Builder{}
                sb.WriteString("ts-refclk:%s")

                mss := MapStringString(m)
                if mss.has("clksrcExt") {
                    sb.WriteString("=%s")
                }
                return sb.String()
            },
        },
        {
            // RFC7273
            // a=mediaclk:direct=963214424
            Name:  "mediaClk",
            Reg:   regexp2MustCompile(`^mediaclk:(?:id=(\S*))? *([^\s=]*)(?:=(\S*))?(?: *rate=(\d+)\/(\d+))?`),
            Names: []string{"id", "mediaClockName", "mediaClockValue", "rateNumerator", "rateDenominator"},
            Format: func(m map[string]string) string {
                sb := strings.Builder{}
                sb.WriteString("mediaclk:")
                mss := MapStringString(m)
                if mss.has("id") {
                    sb.WriteString("id=%s %s")
                } else {
                    sb.WriteString("%v%s")
                }

                if mss.has("mediaClockValue") {
                    sb.WriteString("=%s")
                }

                if mss.has("rateNumerator") {
                    sb.WriteString(" rate=%s")
                }

                if mss.has("rateDenominator") {
                    sb.WriteString("/%s")
                }

                return sb.String()
            },
        },
        {
            // a=keywds:keywords
            Name: "keywords",
            Reg:  regexp2MustCompile(`^keywds:(.+)$`),
            Format: func(m map[string]string) string {
                return "keywds:%s"
            },
        },
        {
            // a=content:main
            Name: "content",
            Reg:  regexp2MustCompile(`^content:(.+)`),
            Format: func(m map[string]string) string {
                return "content:%s"
            },
        },
        // BFCP https://tools.ietf.org/html/rfc4583
        {
            // a=floorctrl:c-s
            Name: "bfcpFloorCtrl",
            Reg:  regexp2MustCompile(`^floorctrl:(c-only|s-only|c-s)`),
            Format: func(m map[string]string) string {
                return "floorctrl:%s"
            },
        },
        {
            // a=confid:1
            Name: "bfcpConfId",
            Reg:  regexp2MustCompile(`^confid:(\d+)`),
            Format: func(m map[string]string) string {
                return "confid:%s"
            },
        },
        {
            // a=userid:1
            Name: "bfcpUserId",
            Reg:  regexp2MustCompile(`^userid:(\d+)`),
            Format: func(m map[string]string) string {
                return "userid:%s"
            },
        },
        {
            // a=floorid:1
            Name:  "bfcpFloorId",
            Reg:   regexp2MustCompile(`^floorid:(.+) (?:m-stream|mstrm):(.+)`),
            Names: []string{"id", "mStream"},
            Format: func(m map[string]string) string {
                return "floorid:%s mstrm:%s"
            },
        },
        {
            // any a= that we don"t understand is kept verbatim on media.invalid
            Push:  "invalid",
            Names: []string{"value"},
        },
    },
}

func init() {
    for _, rules := range grammarMap {
        for _, rule := range rules {
            if rule.Reg == nil {
                rule.Reg = regexp2MustCompile(`(.*)`)
            }
            if rule.Format == nil {
                rule.Format = func(m map[string]string) string {
                    return "%s"
                }
            }
        }
    }
}
