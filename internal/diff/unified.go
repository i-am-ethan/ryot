package diff

import (
	"fmt"
	"strings"
)

type OpKind int

const (
	Equal OpKind = iota
	Delete
	Insert
)

type Op struct {
	Kind OpKind
	Line string
}

// UnifiedDiff は最小実装の unified diff を生成する。
// - Myers等の最適化はせず、LCS(DP)で編集列を作る（LT用・小さめファイル前提）
// - ハンクは1つにまとめる（読みやすさ優先）
func UnifiedDiff(fromFile, toFile string, aText, bText string) string {
	a := splitLinesKeepNL(aText)
	b := splitLinesKeepNL(bText)
	ops := lcsOps(a, b)

	oldCount := len(a)
	newCount := len(b)

	var sb strings.Builder
	sb.WriteString("--- ")
	sb.WriteString(fromFile)
	sb.WriteString("\n")
	sb.WriteString("+++ ")
	sb.WriteString(toFile)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("@@ -1,%d +1,%d @@\n", oldCount, newCount))

	for _, op := range ops {
		switch op.Kind {
		case Equal:
			sb.WriteString(" ")
		case Delete:
			sb.WriteString("-")
		case Insert:
			sb.WriteString("+")
		}
		// op.Line already includes newline when present
		sb.WriteString(op.Line)
		if !strings.HasSuffix(op.Line, "\n") {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func splitLinesKeepNL(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.SplitAfter(s, "\n")
	// If input doesn't end with \n, SplitAfter returns last chunk without \n (OK)
	// If input ends with \n, last part is "" (remove)
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func lcsOps(a, b []string) []Op {
	n := len(a)
	m := len(b)

	// dp[i][j] = LCS length of a[i:] and b[j:]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else {
				if dp[i+1][j] >= dp[i][j+1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j+1]
				}
			}
		}
	}

	ops := make([]Op, 0, n+m)
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			ops = append(ops, Op{Kind: Equal, Line: a[i]})
			i++
			j++
			continue
		}
		if dp[i+1][j] >= dp[i][j+1] {
			ops = append(ops, Op{Kind: Delete, Line: a[i]})
			i++
		} else {
			ops = append(ops, Op{Kind: Insert, Line: b[j]})
			j++
		}
	}
	for i < n {
		ops = append(ops, Op{Kind: Delete, Line: a[i]})
		i++
	}
	for j < m {
		ops = append(ops, Op{Kind: Insert, Line: b[j]})
		j++
	}
	return ops
}
