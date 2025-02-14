package sdp_transform

import (
    "github.com/seamory/mediasou-client-go/sdp_transform/pointer"
    "log"
    "testing"
)

func TestWrite(t *testing.T) {
    sdp := `v=0
o=- 6186858436061843296 1739497432 IN IP4 0.0.0.0
s=-
t=0 0
a=fingerprint:sha-256 45:A7:FA:D6:EE:39:58:CD:77:4E:DD:26:C7:06:42:20:EB:34:E8:83:B8:26:41:E1:EE:63:27:DA:01:72:40:04
a=extmap-allow-mixed
a=group:BUNDLE video audio

m=video 9 UDP/TLS/RTP/SAVPF 96 97 102 103 104 105 106 107 108 109 127 125 39 40 45 46 98 99 100 101 112 113
c=IN IP4 0.0.0.0
a=setup:actpass
a=mid:video
a=ice-ufrag:GhPdTZXCgcDZEmuy
a=ice-pwd:IfMtQsMZMGDEUrhUfiIiJMAQCvYFlwup
a=rtcp-mux
a=rtcp-rsize
a=rtpmap:96 VP8/90000
a=rtcp-fb:96 goog-remb
a=rtcp-fb:96 ccm fir
a=rtcp-fb:96 nack
a=rtcp-fb:96 nack pli
a=rtcp-fb:96 nack
a=rtcp-fb:96 nack pli
a=rtcp-fb:96 transport-cc

a=rtpmap:97 rtx/90000
a=fmtp:97 apt=96
a=rtcp-fb:97 nack
a=rtcp-fb:97 nack pli
a=rtcp-fb:97 transport-cc

a=rtpmap:102 H264/90000
a=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f
a=rtcp-fb:102 goog-remb
a=rtcp-fb:102 ccm fir
a=rtcp-fb:102 nack
a=rtcp-fb:102 nack pli
a=rtcp-fb:102 nack
a=rtcp-fb:102 nack pli
a=rtcp-fb:102 transport-cc
a=rtpmap:103 rtx/90000
a=fmtp:103 apt=102
a=rtcp-fb:103 nack
a=rtcp-fb:103 nack pli
a=rtcp-fb:103 transport-cc
a=rtpmap:104 H264/90000
a=fmtp:104 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f
a=rtcp-fb:104 goog-remb
a=rtcp-fb:104 ccm fir
a=rtcp-fb:104 nack
a=rtcp-fb:104 nack pli
a=rtcp-fb:104 nack
a=rtcp-fb:104 nack pli
a=rtcp-fb:104 transport-cc

a=rtpmap:105 rtx/90000
a=fmtp:105 apt=104
a=rtcp-fb:105 nack
a=rtcp-fb:105 nack pli
a=rtcp-fb:105 transport-cc
a=rtpmap:106 H264/90000
a=fmtp:106 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f
a=rtcp-fb:106 goog-remb
a=rtcp-fb:106 ccm fir
a=rtcp-fb:106 nack
a=rtcp-fb:106 nack pli
a=rtcp-fb:106 nack
a=rtcp-fb:106 nack pli
a=rtcp-fb:106 transport-cc
a=rtpmap:107 rtx/90000
a=fmtp:107 apt=106
a=rtcp-fb:107 nack
a=rtcp-fb:107 nack pli
a=rtcp-fb:107 transport-cc
a=rtpmap:108 H264/90000
a=fmtp:108 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f
a=rtcp-fb:108 goog-remb
a=rtcp-fb:108 ccm fir
a=rtcp-fb:108 nack
a=rtcp-fb:108 nack pli
a=rtcp-fb:108 nack
a=rtcp-fb:108 nack pli
a=rtcp-fb:108 transport-cc
a=rtpmap:109 rtx/90000
a=fmtp:109 apt=108
a=rtcp-fb:109 nack
a=rtcp-fb:109 nack pli
a=rtcp-fb:109 transport-cc
a=rtpmap:127 H264/90000
a=fmtp:127 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=4d001f
a=rtcp-fb:127 goog-remb
a=rtcp-fb:127 ccm fir
a=rtcp-fb:127 nack
a=rtcp-fb:127 nack pli
a=rtcp-fb:127 nack
a=rtcp-fb:127 nack pli
a=rtcp-fb:127 transport-cc

a=rtpmap:125 rtx/90000
a=fmtp:125 apt=127
a=rtcp-fb:125 nack
a=rtcp-fb:125 nack pli
a=rtcp-fb:125 transport-cc
a=rtpmap:39 H264/90000
a=fmtp:39 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=4d001f
a=rtcp-fb:39 goog-remb
a=rtcp-fb:39 ccm fir
a=rtcp-fb:39 nack
a=rtcp-fb:39 nack pli
a=rtcp-fb:39 nack
a=rtcp-fb:39 nack pli
a=rtcp-fb:39 transport-cc

a=rtpmap:40 rtx/90000
a=fmtp:40 apt=39
a=rtcp-fb:40 nack
a=rtcp-fb:40 nack pli
a=rtcp-fb:40 transport-cc

a=rtpmap:45 AV1/90000
a=rtcp-fb:45 goog-remb
a=rtcp-fb:45 ccm fir
a=rtcp-fb:45 nack
a=rtcp-fb:45 nack pli
a=rtcp-fb:45 nack
a=rtcp-fb:45 nack pli
a=rtcp-fb:45 transport-cc

a=rtpmap:46 rtx/90000
a=fmtp:46 apt=45
a=rtcp-fb:46 nack
a=rtcp-fb:46 nack pli
a=rtcp-fb:46 transport-cc

a=rtpmap:98 VP9/90000
a=fmtp:98 profile-id=0
a=rtcp-fb:98 goog-remb
a=rtcp-fb:98 ccm fir
a=rtcp-fb:98 nack
a=rtcp-fb:98 nack pli
a=rtcp-fb:98 nack
a=rtcp-fb:98 nack pli
a=rtcp-fb:98 transport-cc

a=rtpmap:99 rtx/90000
a=fmtp:99 apt=98
a=rtcp-fb:99 nack
a=rtcp-fb:99 nack pli
a=rtcp-fb:99 transport-cc
a=rtpmap:100 VP9/90000
a=fmtp:100 profile-id=2
a=rtcp-fb:100 goog-remb
a=rtcp-fb:100 ccm fir
a=rtcp-fb:100 nack
a=rtcp-fb:100 nack pli
a=rtcp-fb:100 nack
a=rtcp-fb:100 nack pli
a=rtcp-fb:100 transport-cc
a=rtpmap:101 rtx/90000
a=fmtp:101 apt=100
a=rtcp-fb:101 nack
a=rtcp-fb:101 nack pli
a=rtcp-fb:101 transport-cc
a=rtpmap:112 H264/90000
a=fmtp:112 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=64001f
a=rtcp-fb:112 goog-remb
a=rtcp-fb:112 ccm fir
a=rtcp-fb:112 nack
a=rtcp-fb:112 nack pli
a=rtcp-fb:112 nack
a=rtcp-fb:112 nack pli
a=rtcp-fb:112 transport-cc
a=rtpmap:113 rtx/90000
a=fmtp:113 apt=112
a=rtcp-fb:113 nack
a=rtcp-fb:113 nack pli
a=rtcp-fb:113 transport-cc
a=extmap:1 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
a=ssrc:1903097929 cname:EGVljhFhXQWilIua
a=ssrc:1903097929 msid:EGVljhFhXQWilIua leamIBPccZyBZDay
a=ssrc:1903097929 mslabel:EGVljhFhXQWilIua
a=ssrc:1903097929 label:leamIBPccZyBZDay
a=sendrecv
m=audio 9 UDP/TLS/RTP/SAVPF 111 9 0 8
c=IN IP4 0.0.0.0
a=setup:actpass
a=mid:audio
a=ice-ufrag:GhPdTZXCgcDZEmuy
a=ice-pwd:IfMtQsMZMGDEUrhUfiIiJMAQCvYFlwup
a=rtcp-mux
a=rtcp-rsize
a=rtpmap:111 opus/48000/2
a=fmtp:111 minptime=10;useinbandfec=1
a=rtcp-fb:111 transport-cc
a=rtpmap:9 G722/8000
a=rtcp-fb:9 transport-cc
a=rtpmap:0 PCMU/8000
a=rtcp-fb:0 transport-cc
a=rtpmap:8 PCMA/8000
a=rtcp-fb:8 transport-cc
a=extmap:1 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
a=ssrc:2678770010 cname:TSpmnuMCCrSOjrwt
a=ssrc:2678770010 msid:TSpmnuMCCrSOjrwt VJDuDrhFhEVLApBq
a=ssrc:2678770010 mslabel:TSpmnuMCCrSOjrwt
a=ssrc:2678770010 label:VJDuDrhFhEVLApBq
a=sendrecv
`

    description, err := Parse(sdp)
    if err != nil {
        panic(err)
    }

    description.Version = pointer.String("123")

    log.Println(Write(*description, nil))
}
