package colors

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func RedString(str string) string {
	return red + str + reset
}
func GreenString(str string) string {
	return green + str + reset
}
func YellowString(str string) string {
	return yellow + str + reset
}
func BlueString(str string) string {
	return blue + str + reset
}
func MagentaString(str string) string {
	return magenta + str + reset
}
func CyanString(str string) string {
	return cyan + str + reset
}
func GrayString(str string) string {
	return gray + str + reset
}
func WhiteString(str string) string {
	return white + str + reset
}
