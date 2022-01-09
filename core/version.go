package core

const VERSION = "v0.0.3"

func GetVersion() string {
	return "netspy: " + VERSION
}

func PrintVersion() {
	print(GetVersion())
}
