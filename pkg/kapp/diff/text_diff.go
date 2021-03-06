package diff

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/aryann/difflib"
)

type TextDiff []difflib.DiffRecord

func NewTextDiff(existingLines, newLines []string) TextDiff {
	return TextDiff(difflib.Diff(existingLines, newLines))
}

func (l TextDiff) HasChanges() bool {
	for _, diff := range l {
		if diff.Delta != difflib.Common {
			return true
		}
	}
	return false
}

func (l TextDiff) MinimalMD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(l.MinimalString())))
}

func (l TextDiff) MinimalString() string { return l.String(false) }
func (l TextDiff) FullString() string    { return l.String(true) }

func (l TextDiff) String(full bool) string {
	var sb strings.Builder

	for _, diff := range l {
		var mark string

		switch diff.Delta {
		case difflib.RightOnly:
			mark = " + "
		case difflib.LeftOnly:
			mark = " - "
		case difflib.Common:
			if !full {
				continue
			}
			mark = "   "
		}

		// make sure to have line numbers to make sure diff is truly unique
		sb.WriteString(fmt.Sprintf("%3d,%3d%s%s\n", diff.LineLeft, diff.LineRight, mark, diff.Payload))
	}

	return sb.String()
}
