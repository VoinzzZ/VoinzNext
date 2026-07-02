package style

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Color definitions
	Cyan    = color.New(color.FgCyan, color.Bold)
	Green   = color.New(color.FgGreen, color.Bold)
	Yellow  = color.New(color.FgYellow, color.Bold)
	Red     = color.New(color.FgRed, color.Bold)
	White   = color.New(color.FgHiWhite, color.Bold)
	Gray    = color.New(color.FgHiBlack)
	Dim     = color.New(color.FgWhite)
	Magenta = color.New(color.FgMagenta, color.Bold)
	Blue    = color.New(color.FgBlue, color.Bold)

	// Styled Sprint functions
	SprintCyan    = Cyan.SprintFunc()
	SprintGreen   = Green.SprintFunc()
	SprintYellow  = Yellow.SprintFunc()
	SprintRed     = Red.SprintFunc()
	SprintWhite   = White.SprintFunc()
	SprintGray    = Gray.SprintFunc()
	SprintDim     = Dim.SprintFunc()
	SprintMagenta = Magenta.SprintFunc()
	SprintBlue    = Blue.SprintFunc()
)

// ── Logo / Header ──

func Logo() string {
	return SprintCyan("VoinzNext")
}

func Banner(title, subtitle string) {
	width := 50
	line := strings.Repeat("─", width)

	fmt.Println()
	fmt.Printf("  %s%s%s\n", SprintCyan("╭"), line, SprintCyan("╮"))
	fmt.Printf("  %s  %s%s\n", SprintCyan("│"), centerText(title, width-4), SprintCyan("│"))
	fmt.Printf("  %s  %s%s\n", SprintCyan("│"), centerText(subtitle, width-4), SprintCyan("│"))
	fmt.Printf("  %s%s%s\n", SprintCyan("╰"), line, SprintCyan("╯"))
	fmt.Println()
}

// ── Progress steps ──

func StepRunning(label string) {
	fmt.Printf("  %s %s...\n", SprintCyan("●"), label)
}

func StepDone(label string) {
	fmt.Printf("  %s %s\n", SprintGreen("✔"), SprintDim(label))
}

func StepWarn(label string, detail string) {
	fmt.Printf("  %s %s %s\n", SprintYellow("⚠"), label, SprintDim(detail))
}

func StepError(label string, err error) {
	fmt.Printf("  %s %s: %v\n", SprintRed("✘"), label, err)
}

// ── Success / Error banners ──

func SuccessBanner(projectName string) {
	fmt.Println()
	fmt.Printf("  %s\n", SprintGreen(strings.Repeat("─", 50)))
	fmt.Printf("  %s  %s%s\n", SprintGreen("◆"), centerText("Project created successfully!", 44), SprintGreen("◆"))
	fmt.Printf("  %s\n", SprintGreen(strings.Repeat("─", 50)))
	fmt.Println()
	fmt.Printf("  %s %s\n", SprintGreen("✔"), fmt.Sprintf("Project %q has been generated.", SprintWhite(projectName)))
	fmt.Println()
}

func ErrorBanner(err error) {
	fmt.Println()
	fmt.Printf("  %s\n", SprintRed(strings.Repeat("─", 40)))
	fmt.Printf("  %s  %s\n", SprintRed("◆"), centerText("Error", 36))
	fmt.Printf("  %s\n", SprintRed(strings.Repeat("─", 40)))
	fmt.Printf("  %s\n\n", SprintRed(fmt.Sprintf("  %v", err)))
}

// ── Next steps ──

func NextSteps(steps [][2]string) {
	fmt.Printf("  %s\n", SprintCyan("┌──────────────────────────────────────────┐"))
	fmt.Printf("  %s\n", SprintCyan("│")+"  "+SprintWhite("Next Steps")+strings.Repeat(" ", 33)+SprintCyan("│"))
	fmt.Printf("  %s\n", SprintCyan("├──────────────────────────────────────────┤"))
	for _, step := range steps {
		label := step[0]
		cmd := step[1]
		fmt.Printf("  %s  %s  %s\n", SprintCyan("│"), SprintDim(label), SprintWhite(cmd))
	}
	fmt.Printf("  %s\n", SprintCyan("└──────────────────────────────────────────┘"))
	fmt.Println()
}

// ── Misc helpers ──

func Label(text string) string {
	return SprintCyan(text)
}

func Value(text string) string {
	return SprintWhite(text)
}

func Dimmed(text string) string {
	return SprintDim(text)
}

func Check() string {
	return SprintGreen("✔")
}

func Bullet(text string) string {
	return fmt.Sprintf("  %s %s", SprintCyan("●"), text)
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	padding := width - len(text)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

// FeatureList displays a feature with active/inactive status
func FeatureList(feature string, active bool) string {
	if active {
		return fmt.Sprintf("  %s  %s", SprintGreen("✔"), SprintWhite(feature))
	}
	return fmt.Sprintf("  %s  %s", SprintDim("·"), SprintDim(feature))
}

// Separator draws a dim line
func Separator() {
	fmt.Printf("  %s\n", SprintDim(strings.Repeat("─", 50)))
}
