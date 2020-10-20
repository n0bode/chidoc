package chidoc

type FontSize string

const (
	FontSizeDefault FontSize = "default"
	FontSizeLarge            = "large"
	FontSizeLargest          = "largest"
)

type FontType string

const (
	FontTypeDefault FontType = "Nunito"
	FontTypeSans             = "Open Sans"
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
