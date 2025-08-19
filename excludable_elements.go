package unstructured

// ExcludeableElement represents elements that can be excluded during document processing.
type ExcludeableElement string

// Excludeable element constants for document processing.
const (
	ExcludableElementFigureCaption     ExcludeableElement = "FigureCaption"
	ExcludableElementNarrativeText     ExcludeableElement = "NarrativeText"
	ExcludableElementListItem          ExcludeableElement = "ListItem"
	ExcludableElementTitle             ExcludeableElement = "Title"
	ExcludableElementAddress           ExcludeableElement = "Address"
	ExcludableElementTable             ExcludeableElement = "Table"
	ExcludableElementPageBreak         ExcludeableElement = "PageBreak"
	ExcludableElementHeader            ExcludeableElement = "Header"
	ExcludableElementFooter            ExcludeableElement = "Footer"
	ExcludableElementUncategorizedText ExcludeableElement = "UncategorizedText"
	ExcludableElementImage             ExcludeableElement = "Image"
	ExcludableElementFormula           ExcludeableElement = "Formula"
	ExcludableElementEmailAddress      ExcludeableElement = "EmailAddress"
)
