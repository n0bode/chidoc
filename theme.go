package chidoc

type FontSize string

const (
	FontSizeDefault FontSize = "default"
	FontSizeLarge   FontSize = "large"
	FontSizeLargest FontSize = "largest"
)

// FontType const to define font type documentation render
type FontType string

const (
	// FontTypeDefault default font for documentation
	FontTypeDefault FontType = "Nunito"
	// FontTypeSans define open sans as font on documentation render
	FontTypeSans = "Open Sans"
)

// Theme costumizer docs
type Theme struct {
	Schema          string
	PrimaryColor    string
	BackgroundColor string
	TextColor       string
	RenderStyle     string
	HeaderColor     string
	FontType        FontType
	FontSize        FontSize
	FontName        string
	Header          string
	SchemaType      string
}

// DefaultTheme the lighten theme default
var DefaultTheme Theme = Theme{
	Schema:          "light",
	BackgroundColor: "#FFF",
	PrimaryColor:    "#DDD",
	TextColor:       "#444",
	RenderStyle:     "read",
	HeaderColor:     "#FFF",
	FontSize:        FontSizeDefault,
	FontType:        FontTypeDefault,
	FontName:        "Roboto Mono",
	Header:          "false",
	SchemaType:      "table",
}

// DarkTheme the darken theme default
var DarkTheme Theme = Theme{
	Schema:          "dark",
	BackgroundColor: "#2D3133",
	TextColor:       "#CAD9E3",
	PrimaryColor:    "#FF3D00",
	HeaderColor:     "#FFF",
	RenderStyle:     "read",
	FontSize:        FontSizeDefault,
	FontType:        FontTypeDefault,
	FontName:        "Roboto Mono",
	Header:          "false",
	SchemaType:      "table",
}
