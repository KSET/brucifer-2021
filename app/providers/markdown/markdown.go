package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type markdownProvider struct {
}

var renderer_ markdown.Renderer
var parser_ *parser.Parser

func MarkdownProvider() markdownProvider {
	return markdownProvider{}
}

func (p markdownProvider) Render(text string) string {
	renderer_ = html.NewRenderer(
		html.RendererOptions{
			Flags: html.NofollowLinks |
				html.UseXHTML |
				html.Smartypants |
				html.SmartypantsFractions |
				html.SmartypantsDashes |
				html.SmartypantsLatexDashes |
				html.LazyLoadImages |
				html.NoopenerLinks,
		},
	)

	parser_ = parser.NewWithExtensions(
		parser.NoIntraEmphasis |
			parser.Tables |
			parser.FencedCode |
			parser.Autolink |
			parser.Strikethrough |
			parser.SpaceHeadings |
			parser.HardLineBreak |
			parser.HeadingIDs |
			parser.BackslashLineBreak |
			parser.DefinitionLists |
			parser.MathJax |
			parser.AutoHeadingIDs,
	)

	rendered := markdown.ToHTML([]byte(text), parser_, renderer_)
	// rendered := bluemonday.UGCPolicy().SanitizeBytes(renderedUnsafe)

	return string(rendered)
}
