# PeerBeam

PeerBeam is a CLI file transfer app powered by WebRTC. It establishes a direct, peer-to-peer connection without relying on intermediary data relays. 

- **Direct P2P File Transfer**: WebRTC NAT traversal can connect two clients that are both behind NATs and firewalls without any port forwarding
- **STUN Server**: Uses Google's public STUN server to retrieve public IP address, port, and NAT type information

⚠️ **Warning:** PeerBeam is under active development, and the API is not yet stable. Breaking changes may occur, so use with caution.

[xkcd: File Transfer](https://xkcd.com/949/)

## Installation

To install PeerBeam, you'll need to have Go installed. You can download and install PeerBeam using the following commands:

```bash
go install github.com/ksp237/peerbeam@latest
```

This will install the `peerbeam` binary to your `$GOPATH/bin`.

## Usage

PeerBeam provides simple commands for sending and receiving files:

### Send Files

To send files:

```bash
peerbeam send <file1> <file2> ...
```

This will initiate the file transfer process. The receiving peer will need to execute the `receive` command to accept the files.

### Receive Files

To receive files:

```bash
peerbeam receive
```

This will listen for incoming file transfers and prompt for acceptance before the transfer begins.

### Find Your Public IP

To fetch srflx candidate (NAT IP:port mappings)

```bash
peerbeam stun
```

## Development

To build PeerBeam from source:

```bash
git clone https://github.com/yourusername/peerbeam.git
cd peerbeam
go build
```

## License

PeerBeam is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## References
* [pion/webrtc](https://github.com/pion/webrtc)
* [pion/sctp](https://github.com/pion/sctp)
