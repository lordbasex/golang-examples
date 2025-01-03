
# GO Asterisk Core Show Channels - AMI Client

Lord BaseX (c) 2014-2025
 Federico Pereira lord.basex@gmail.com

This Go application connects to an Asterisk server via the Asterisk Manager Interface (AMI), executes the `core show channels concise` command every second, and outputs detailed information about active channels.

### Requirements
- Go 1.18 or later
- A running Asterisk server
- AMI (Asterisk Manager Interface) enabled and configured on the Asterisk server

### Features
- Fetches and processes channel information from the Asterisk server.
- Allows configuration through command-line arguments or environment variables.
- Continuously outputs information about active channels in a structured format.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/lordbasex/golang-examples.git
   cd golang-examples/go-asterisk-core-show-channels
   ```

2. Build the application:
   ```bash
   make
   ```

### Configuration

The application can be configured using environment variables or command-line arguments.

#### Environment Variables

Set the following environment variables before running the application:

- `AMI_HOST` – The IP address or hostname of the Asterisk server (default: `127.0.0.1`)
- `AMI_PORT` – The AMI port (default: `5038`)
- `AMI_USERNAME` – The AMI username (default: `admin`)
- `AMI_PASSWORD` – The AMI password (default: `password`)

Example:

```bash
export AMI_HOST=127.0.0.1
export AMI_PORT=8080
export AMI_USERNAME=my_user
export AMI_PASSWORD=my_password
```

#### Command-Line Arguments

Alternatively, you can configure the application via command-line arguments:

```bash
./asterisk-core-show-channels-aarch64 -host=127.0.0.1 -port=5038 -username=my_user -password=my_password
```
or 

```bash
./asterisk-core-show-channels-amd64 -host=127.0.0.1 -port=5038 -username=my_user -password=my_password
```

### Usage

After building the application, you can run it with the following command:

```bash
./asterisk-core-show-channels-xxx
```

By default, it will connect to `127.0.0.1:5038` with the username `admin` and password `password`. You can override these values using environment variables or command-line arguments.

Once connected, the application will continuously retrieve information from the Asterisk server every second, processing and displaying channel information as follows:

Example output:

```
Connected: Asterisk Call Manager/7.0.3
--------------->>
Channel: SIP/ANTEL-CELU-0000001e
    Context: demo-incoming
    Exten:
    Priority: 1
    ChannelState: Down
    Application: AppDial2
    ApplicationData: (Outgoing Line)
    CallerIDNum: 1111111122222
    AccountCode: 3
    PeerAccount: 3
    Duration: 0
    BridgeId: 1735937566.30

    Dump: SIP/ANTEL-CELU-0000001e!demo-incoming!!1!Down!AppDial2!(Outgoing Line)!1111111122222!3!3!3!0!!1735937566.30
<<---------------
```

### Explanation of Output

Each line represents a channel. The data is split into 14 parts by the `!` separator. The relevant fields are displayed for each channel, such as:

- **Channel**: The channel identifier (e.g., `SIP/ANTEL-CELU-0000001e`)
- **Context**: The context associated with the channel (e.g., `demo-incoming`)
- **ChannelState**: The current state of the channel (e.g., `Down`)
- **Application**: The application currently being used on the channel (e.g., `AppDial2`)
- **Duration**: The duration of the call or activity on the channel (e.g., `0`)
- **BridgeId**: The bridge ID for the call (e.g., `1735937566.30`)

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
