package ui

import (
	"fmt"
	"strings"
)

type Dividers struct{}

func NewDividers() *Dividers { return &Dividers{} }

func (d *Dividers) RenderPaginationFooter(currentPage, totalPages int) string {
	divider := strings.Repeat("─", 40)
	return fmt.Sprintf("\n%s\n Page [%d / %d] | [N]ext | [P]rev | [Q]uit\n%s", 
		divider, currentPage, totalPages, divider)
}
