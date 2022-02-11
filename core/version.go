package core

const VERSION = "v0.0.5"

func GetVersion() string {
	return "netspy: " + VERSION
}

func PrintVersion() {
	print(GetVersion())
}
