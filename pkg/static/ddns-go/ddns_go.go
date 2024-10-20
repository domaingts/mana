package ddnsgo

import _ "embed"

//go:embed ddns-go.service
var Service []byte

//go:embed config.yaml
var Config []byte
