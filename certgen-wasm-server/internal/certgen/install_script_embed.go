package certgen

import _ "embed"

// installCACertScript is embedded so WASM can include the helper script in the downloaded zip
// without reading from a local filesystem at runtime.
//
//go:embed assets/install_ca_cert.sh
var installCACertScript []byte

