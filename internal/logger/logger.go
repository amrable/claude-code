package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var colors = struct {
	Reset   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	Gray    string
	White   string
	Bold    string
}{
	Reset:   "\033[0m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	Gray:    "\033[90m",
	White:   "\033[97m",
	Bold:    "\033[1m",
}

type StyledFormatter struct {
	TimestampFormat string
}

func (f *StyledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	levelColor := f.getLevelColor(entry.Level)
	levelIcon := f.getLevelIcon(entry.Level)

	b.WriteString(levelColor)
	b.WriteString(levelIcon)
	b.WriteString(" ")

	if entry.Message != "" {
		b.WriteString(entry.Message)
	}

	b.WriteString(colors.Reset)
	b.WriteString("\n")

	return b.Bytes(), nil
}

func (f *StyledFormatter) getLevelColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return colors.Gray
	case logrus.InfoLevel:
		return colors.Cyan
	case logrus.WarnLevel:
		return colors.Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colors.Red
	default:
		return colors.Reset
	}
}

func (f *StyledFormatter) getLevelIcon(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "◆"
	case logrus.InfoLevel:
		return "●"
	case logrus.WarnLevel:
		return "▲"
	case logrus.ErrorLevel:
		return "✖"
	case logrus.FatalLevel:
		return "✖"
	default:
		return "○"
	}
}

func Setup() {
	logrus.SetFormatter(&StyledFormatter{})
	logrus.SetOutput(os.Stdout)

	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func PrintBanner() {
	banner := `
` + colors.Cyan + colors.Bold + `╔════════════════════════════════════════╗
║      ` + colors.White + `Claude Code CLI - Go Edition` + colors.Cyan + `      ║
╚════════════════════════════════════════╝` + colors.Reset + `

Type your message and press Enter.
Press Ctrl+C to exit.
`
	fmt.Println(banner)
}

func Prompt() {
	fmt.Print(colors.Green + "> " + colors.Reset)
}

func AssistantResponse(msg string) {
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		fmt.Printf(colors.Cyan+"%s"+colors.Reset+"\n", line)
	}
}

func ToolCall(name string) {
	fmt.Printf(colors.Yellow+"⟳ Calling tool: %s"+colors.Reset+"\n", name)
}

func Error(err error) {
	fmt.Printf(colors.Red+"✖ Error: %v"+colors.Reset+"\n", err)
}
