# PeerBeam

https://github.com/user-attachments/assets/9579b274-1067-42c8-b8c9-816e3168f6b4

`PeerBeam` is a tool for fast and secure file transfer between computers.

- **direct file transfer** between two computers
- **WebRTC** for secure, p2p communication
- **cross-platform**: works on Windows, Linux, macOS
- **no port-forwarding** or network config needed
- supports **ipv6** and **ipv4**

## Installation

[Install Go](https://golang.org/dl/) then run:
```
go install github.com/6b70/peerbeam@latest
```
This will install the `peerbeam` binary to your `$GOPATH/bin`.

## Usage
To send a file run:
```bash
peerbeam send <file1> <file2> ...
```

To receive files run:
```bash
peerbeam receive
```

You can also query a STUN server:
```bash
peerbeam stun
```

## References
* [pion/webrtc](https://github.com/pion/webrtc)
* [pion/sctp](https://github.com/pion/sctp)
* [schollz/croc](https://github.com/schollz/croc)

