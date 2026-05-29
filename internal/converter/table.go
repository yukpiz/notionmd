package converter

import (
	"github.com/brittonhayes/notionmd/chunk"
	"github.com/dstotijn/go-notion"
	"github.com/gomarkdown/markdown/ast"
)

func isTable(node ast.Node) bool {
	_, ok := node.(*ast.Table)
	return ok
}

func convertTable(node *ast.Table) notion.Block {
	if node == nil {
		return nil
	}

	var rows []notion.Block
	hasColumnHeader := false

	for _, child := range node.GetChildren() {
		switch c := child.(type) {
		case *ast.TableHeader:
			hasColumnHeader = true
			for _, row := range c.GetChildren() {
				if tableRow, ok := row.(*ast.TableRow); ok {
					rows = append(rows, convertTableRow(tableRow))
				}
			}
		case *ast.TableBody:
			for _, row := range c.GetChildren() {
				if tableRow, ok := row.(*ast.TableRow); ok {
					rows = append(rows, convertTableRow(tableRow))
				}
			}
		}
	}

	if len(rows) == 0 {
		return nil
	}

	tableWidth := 0
	if len(rows) > 0 {
		if firstRow, ok := rows[0].(notion.TableRowBlock); ok {
			tableWidth = len(firstRow.Cells)
		}
	}

	return notion.TableBlock{
		TableWidth:      tableWidth,
		HasColumnHeader: hasColumnHeader,
		HasRowHeader:    false,
		Children:        rows,
	}
}

func convertTableRow(node *ast.TableRow) notion.Block {
	var cells [][]notion.RichText

	for _, child := range node.GetChildren() {
		cell, ok := child.(*ast.TableCell)
		if !ok {
			continue
		}
		cells = append(cells, convertTableCell(cell))
	}

	return notion.TableRowBlock{
		Cells: cells,
	}
}

func convertTableCell(node *ast.TableCell) []notion.RichText {
	if node == nil {
		return nil
	}

	var richText []notion.RichText
	for _, child := range node.GetChildren() {
		if isLink(child) {
			linkBlock := convertLinkToTextBlock(child.(*ast.Link))
			if linkBlock != nil {
				richText = append(richText, linkBlock...)
			}
			continue
		}

		if isStyledText(child) {
			styledBlock := convertStyledTextToBlock(child)
			if styledBlock != nil {
				richText = append(richText, styledBlock...)
			}
			continue
		}

		if child.AsLeaf() != nil {
			content := string(child.AsLeaf().Literal)
			if content != "" {
				richText = append(richText, chunk.RichText(content, nil)...)
			}
		}
	}

	return richText
}
