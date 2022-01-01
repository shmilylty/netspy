package core

const VERSION = "0.0.1"

func GetVersion() string {
	return "version: " + VERSION
}

func PrintVersion() {
	print(GetVersion())
}
