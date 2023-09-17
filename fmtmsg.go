package main

func FmtGreen(s string) string {
	return string("\033[32m") + s + string("\033[0m")
}

func FmtRed(s string) string {
	return string("\033[31m") + s + string("\033[0m")
}
