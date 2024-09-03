package conn

import "github.com/pion/webrtc/v4"

var iceServers = []webrtc.ICEServer{
	{URLs: []string{"stun:stun.l.google.com:19302"}},
	{URLs: []string{"stun:stun.l.google.com:5349"}},
	{URLs: []string{"stun:stun1.l.google.com:3478"}},
	{URLs: []string{"stun:stun1.l.google.com:5349"}},
	{URLs: []string{"stun:stun2.l.google.com:19302"}},
	{URLs: []string{"stun:stun2.l.google.com:5349"}},
	{URLs: []string{"stun:stun3.l.google.com:3478"}},
	{URLs: []string{"stun:stun3.l.google.com:5349"}},
	{URLs: []string{"stun:stun4.l.google.com:19302"}},
	{URLs: []string{"stun:stun4.l.google.com:5349"}},
}
