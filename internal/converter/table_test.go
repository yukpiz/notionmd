package converter

import (
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/stretchr/testify/assert"
)

func TestConvertTable(t *testing.T) {
	t.Run("can convert simple table", func(t *testing.T) {
		markdownText := `| Name | Age |
| --- | --- |
| Alice | 30 |
| Bob | 25 |
`

		expected := []notion.Block{
			notion.TableBlock{
				TableWidth:      2,
				HasColumnHeader: true,
				HasRowHeader:    false,
				Children: []notion.Block{
					notion.TableRowBlock{
						Cells: [][]notion.RichText{
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "Name"},
									PlainText: "Name",
								},
							},
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "Age"},
									PlainText: "Age",
								},
							},
						},
					},
					notion.TableRowBlock{
						Cells: [][]notion.RichText{
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "Alice"},
									PlainText: "Alice",
								},
							},
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "30"},
									PlainText: "30",
								},
							},
						},
					},
					notion.TableRowBlock{
						Cells: [][]notion.RichText{
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "Bob"},
									PlainText: "Bob",
								},
							},
							{
								{
									Type:      notion.RichTextTypeText,
									Text:      &notion.Text{Content: "25"},
									PlainText: "25",
								},
							},
						},
					},
				},
			},
		}

		result, err := Convert(markdownText)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("can convert table with styled text", func(t *testing.T) {
		markdownText := `| Feature | Status |
| --- | --- |
| **Auth** | Done |
`

		result, err := Convert(markdownText)
		assert.NoError(t, err)
		assert.Len(t, result, 1)

		table, ok := result[0].(notion.TableBlock)
		assert.True(t, ok)
		assert.Equal(t, 2, table.TableWidth)
		assert.True(t, table.HasColumnHeader)
		assert.Len(t, table.Children, 2)

		// Check bold in body row
		bodyRow, ok := table.Children[1].(notion.TableRowBlock)
		assert.True(t, ok)
		assert.Len(t, bodyRow.Cells, 2)
		assert.Equal(t, "Auth", bodyRow.Cells[0][0].PlainText)
		assert.True(t, bodyRow.Cells[0][0].Annotations.Bold)
	})

	t.Run("can convert table with heading", func(t *testing.T) {
		markdownText := `## Summary

| Key | Value |
| --- | --- |
| Name | Test |
`

		result, err := Convert(markdownText)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		_, isHeading := result[0].(notion.Heading2Block)
		assert.True(t, isHeading)

		_, isTable := result[1].(notion.TableBlock)
		assert.True(t, isTable)
	})
}
