package dto

// ZebraParams defines the zebra printer parameters
type ZebraParams struct {
	Device    string      `json:"device"`
	IpPrinter string      `json:"body"`
	Data      interface{} `json:"data"`
}

// PosParams defines the pos printer parameters
type PosParams struct {
	Device    string      `json:"device"`
	IpPrinter string      `json:"body"`
	Data      interface{} `json:"data"`
}

// DotMatrixParams defines the dot matrix printer parameters
type DotMatrixParams struct {
	Device    string      `json:"device"`
	IpPrinter string      `json:"body"`
	Data      interface{} `json:"data"`
}
