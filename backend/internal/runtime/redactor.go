package runtime

import "regexp"

var redactionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(token|secret|password)\s*=\s*([^\s]+)`),
	regexp.MustCompile(`(?i)(bearer)\s+([^\s]+)`),
	regexp.MustCompile(`sk-[A-Za-z0-9_-]+`),
}

// Redact 移除给定内容中的敏感值。
func Redact(input string) string {
	output := input
	for _, pattern := range redactionPatterns {
		output = pattern.ReplaceAllStringFunc(output, func(_ string) string {
			return "[REDACTED]"
		})
	}
	return output
}
