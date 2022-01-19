package core

const VERSION = "v0.0.4"

func GetVersion() string {
	return "netspy: " + VERSION
}

func PrintVersion() {
	print(GetVersion())
}
