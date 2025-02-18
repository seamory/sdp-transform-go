package sdp_transform

type MediaDescription struct {
    SharedDescriptionFields
    MediaAttributes
}

type Origin struct {
    Username       string `json:"username"`
    SessionID      string `json:"sessionId"`
    SessionVersion string `json:"sessionVersion"`
    NetType        string `json:"netType"`
    IPVer          string `json:"ipVer"`
    Address        string `json:"address"`
}

type Timing struct {
    Start string `json:"start"`
    Stop  string `json:"stop"`
}

type Media struct {
    Type     string  `json:"type"`
    Port     string  `json:"port"`
    Protocol string  `json:"protocol"`
    Payloads *string `json:"payloads,omitempty"`
    MediaDescription
    MediaExtensionAttributes
}

// SessionDescription
// Descriptor fields that exist only at the session level (before an m= block).
// See the SDP grammar for more details: https://tools.ietf.org/html/rfc4566#section-9
type SessionDescription struct {
    SharedDescriptionFields
    SessionAttributes
    Version          *string  `json:"version,omitempty"`
    Origin           *Origin  `json:"origin,omitempty"`
    Name             *string  `json:"name,omitempty"`
    URI              *string  `json:"uri,omitempty"`
    Email            *string  `json:"email,omitempty"`
    Phone            *string  `json:"phone,omitempty"`
    Timing           *Timing  `json:"timing,omitempty"`
    Timezones        *string  `json:"timezones,omitempty"`
    Repeats          *string  `json:"repeats,omitempty"`
    Media            []*Media `json:"media,omitempty"`
    ExtmapAllowMixed *string  `json:"extmapAllowMixed,omitempty"`
}

type Ext struct {
    Value      string  `json:"value"`
    Direction  *string `json:"direction,omitempty"`
    EncryptUri *string `json:"encrypt_uri,omitempty"`
    URI        string  `json:"uri"`
    Config     *string `json:"config,omitempty"`
}

type Fingerprint struct {
    Type string `json:"type"`
    Hash string `json:"hash"`
}

type SourceFilter struct {
    FilterMode   string `json:"filterMode"`
    NetType      string `json:"netType"`
    AddressTypes string `json:"addressTypes"`
    DestAddress  string `json:"destAddress"`
    SrcList      string `json:"srcList"`
}

type Invalid struct {
    Value string `json:"value"`
}

// SharedAttributes
// These attributes can exist on both the session level and the media level.
// https://www.iana.org/assignments/sdp-parameters/sdp-parameters.xhtml#sdp-parameters-8
type SharedAttributes struct {
    Direction    *string       `json:"direction,omitempty"`
    Control      *string       `json:"control,omitempty"`
    Ext          []*Ext        `json:"ext,omitempty"`
    Setup        *string       `json:"setup,omitempty"`
    IceUfrag     *string       `json:"iceUfrag,omitempty"`
    IcePwd       *string       `json:"icePwd,omitempty"`
    Fingerprint  *Fingerprint  `json:"fingerprint,omitempty"`
    SourceFilter *SourceFilter `json:"sourceFilter,omitempty"`
    Invalid      []*Invalid    `json:"invalid,omitempty"`
}

type MsidSemantic struct {
    Semantic string `json:"semantic"`
    Token    string `json:"token"`
}

type Group struct {
    Type string `json:"type"`
    Mids string `json:"mids"`
}

// SessionAttributes
// Attributes that only exist at the session level (before an m= block).
// https://www.iana.org/assignments/sdp-parameters/sdp-parameters.xhtml#sdp-parameters-7
type SessionAttributes struct {
    SharedAttributes
    IceLite      *string       `json:"icelite,omitempty"`
    IceOptions   *string       `json:"iceOptions,omitempty"`
    MsidSemantic *MsidSemantic `json:"msidSemantic,omitempty"`
    Groups       []*Group      `json:"groups,omitempty"`
}

type RTP struct {
    Payload  string  `json:"payload"`
    Codec    string  `json:"codec"`
    Rate     *string `json:"rate,omitempty"`
    Encoding *string `json:"encoding,omitempty"`
}

type RTCP struct {
    Port    string  `json:"port"`
    NetType *string `json:"netType,omitempty"`
    IPVer   *string `json:"ipVer,omitempty"`
    Address *string `json:"address,omitempty"`
}

type RTCPFB struct {
    Payload string  `json:"payload"`
    Type    string  `json:"type"`
    SubType *string `json:"subtype,omitempty"`
}

type RTCPFBTrrInt struct {
    Payload string `json:"payload"`
    Value   string `json:"value"`
}

type FMTP struct {
    Payload string `json:"payload"`
    Config  string `json:"config"`
}

type Crypto struct {
    ID            string  `json:"id"`
    Suite         string  `json:"suite"`
    Config        string  `json:"config"`
    SessionConfig *string `json:"sessionConfig,omitempty"`
}

type Candidate struct {
    Foundation  string  `json:"foundation"`
    Component   string  `json:"component"`
    Transport   string  `json:"transport"`
    Priority    string  `json:"priority"`
    IP          string  `json:"ip"`
    Port        string  `json:"port"`
    Type        string  `json:"type"`
    Raddr       *string `json:"raddr,omitempty"`
    Rport       *string `json:"rport,omitempty"`
    TCPType     *string `json:"tcptype,omitempty"`
    Generation  *string `json:"generation,omitempty"`
    NetworkID   *string `json:"network-id,omitempty"`
    NetworkCost *string `json:"network-cost,omitempty"`
}

type SSRC struct {
    ID        string  `json:"id"`
    Attribute string  `json:"attribute"`
    Value     *string `json:"value,omitempty"`
}

type SSRCGroup struct {
    Semantics string `json:"semantics"`
    SSRCs     string `json:"ssrcs"`
}

type SCTPMap struct {
    SCTPMapNumber  string `json:"sctpmapNumber"`
    App            string `json:"app"`
    MaxMessageSize string `json:"maxMessageSize"`
}

type RID struct {
    ID        string  `json:"id"`
    Direction string  `json:"direction"`
    Params    *string `json:"params,omitempty"`
}

type ImageAttr struct {
    PT     interface{} `json:"pt"`
    Dir1   string      `json:"dir1"`
    Attrs1 string      `json:"attrs1"`
    Dir2   *string     `json:"dir2,omitempty"`
    Attrs2 *string     `json:"attrs2,omitempty"`
}

type Simulcast struct {
    Dir1  string  `json:"dir1"`
    List1 string  `json:"list1"`
    Dir2  *string `json:"dir2,omitempty"`
    List2 *string `json:"list2,omitempty"`
}

type Simulcast03 struct {
    Value string `json:"value"`
}

// MediaAttributes
// Attributes that only exist at the media level (within an m= block).
// https://www.iana.org/assignments/sdp-parameters/sdp-parameters.xhtml#sdp-parameters-9
type MediaAttributes struct {
    SharedAttributes
    RTP              []*RTP          `json:"rtp"`
    RTCP             *RTCP           `json:"rtcp,omitempty"`
    RTCPFB           []*RTCPFB       `json:"rtcpFb,omitempty"`
    RTCPFBTrrInt     []*RTCPFBTrrInt `json:"rtcpFbTrrInt,omitempty"`
    FMTP             []*FMTP         `json:"fmtp"`
    MID              *string         `json:"mid,omitempty"`
    MSID             *string         `json:"msid,omitempty"`
    PTIME            *string         `json:"ptime,omitempty"`
    MaxPTIME         *string         `json:"maxptime,omitempty"`
    Crypto           *Crypto         `json:"crypto,omitempty"`
    Candidates       []*Candidate    `json:"candidates,omitempty"`
    EndOfCandidates  *string         `json:"endOfCandidates,omitempty"`
    RemoteCandidates *string         `json:"remoteCandidates,omitempty"`
    SSRCs            []*SSRC         `json:"ssrcs,omitempty"`
    SSRCGroups       []*SSRCGroup    `json:"ssrcGroups,omitempty"`
    RTCPMux          *string         `json:"rtcpMux,omitempty"`
    RTCPRsize        *string         `json:"rtcpRsize,omitempty"`
    SCTPMap          *SCTPMap        `json:"sctpmap,omitempty"`
    XGoogleFlag      *string         `json:"xGoogleFlag,omitempty"`
    RIDs             []*RID          `json:"rids,omitempty"`
    ImageAttrs       []*ImageAttr    `json:"imageattrs,omitempty"`
    Simulcast        *Simulcast      `json:"simulcast,omitempty"`
    Simulcast03      *Simulcast03    `json:"simulcast_03,omitempty"`
    Framerate        *string         `json:"framerate,omitempty"`
}

type Connection struct {
    Version string `json:"version"`
    IP      string `json:"ip"`
}

type Bandwidth struct {
    Type  string `json:"type"`
    Limit string `json:"limit"` // string or number
}

// SharedDescriptionFields
// Descriptor fields that exist at both the session level and media level.
// See the SDP grammar for more details: https://tools.ietf.org/html/rfc4566#section-9
type SharedDescriptionFields struct {
    Description *string      `json:"description,omitempty"`
    Connection  *Connection  `json:"connection,omitempty"`
    Bandwidth   []*Bandwidth `json:"bandwidth,omitempty"`
}

// MediaExtensionAttributes mediasoup used.
type MediaExtensionAttributes struct {
    IceOptions       *string `json:"iceOptions,omitempty"`
    ExtmapAllowMixed *string `json:"extmapAllowMixed,omitempty"`
    SctpPort         *string `json:"sctpPort,omitempty"`
    MaxMessageSize   *string `json:"maxMessageSize,omitempty"`
}
