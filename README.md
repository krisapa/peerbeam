# PeerBeam

`PeerBeam` is a CLI tool that allows two computers to quickly and securely transfer files. 

- enables **direct file transfer** between two computers
- uses **WebRTC** for secure, p2p communication
- **cross-platform**: works on Windows, Linux, macOS, your toaster
- **no port-forwarding** or network config needed
- supports **ipv6** and **ipv4**

> **Note:** Some side channel is required to exchange the initial connection info. (RDP, SSH, text, email, etc.)

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

You can also query a STUN server for ICE candidates by running:
```bash
peerbeam stun
```

## References
* [pion/webrtc](https://github.com/pion/webrtc)
* [pion/sctp](https://github.com/pion/sctp)
* [schollz/croc](https://github.com/schollz/croc)

## License

PeerBeam is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
