package compat

import (
	"github.com/evanw/esbuild/internal/css_ast"
)

type CSSFeature uint8

const (
	HexRGBA CSSFeature = 1 << iota
	InlineStyle
	InsetProperty
	IsPseudoClass
	Modern_RGB_HSL
	Nesting
	RebeccaPurple
)

var StringToCSSFeature = map[string]CSSFeature{
	"hex-rgba":        HexRGBA,
	"inline-style":    InlineStyle,
	"inset-property":  InsetProperty,
	"is-pseudo-class": IsPseudoClass,
	"modern-rgb-hsl":  Modern_RGB_HSL,
	"nesting":         Nesting,
	"rebecca-purple":  RebeccaPurple,
}

func (features CSSFeature) Has(feature CSSFeature) bool {
	return (features & feature) != 0
}

func (features CSSFeature) ApplyOverrides(overrides CSSFeature, mask CSSFeature) CSSFeature {
	return (features & ^mask) | (overrides & mask)
}

var cssTable = map[CSSFeature]map[Engine][]versionRange{
	HexRGBA: {
		Chrome:  {{start: v{62, 0, 0}}},
		Edge:    {{start: v{79, 0, 0}}},
		Firefox: {{start: v{49, 0, 0}}},
		IOS:     {{start: v{9, 3, 0}}},
		Opera:   {{start: v{49, 0, 0}}},
		Safari:  {{start: v{9, 1, 0}}},
	},
	InsetProperty: {
		Chrome:  {{start: v{87, 0, 0}}},
		Edge:    {{start: v{87, 0, 0}}},
		Firefox: {{start: v{66, 0, 0}}},
		IOS:     {{start: v{14, 5, 0}}},
		Opera:   {{start: v{73, 0, 0}}},
		Safari:  {{start: v{14, 1, 0}}},
	},
	IsPseudoClass: {
		Chrome:  {{start: v{88, 0, 0}}},
		Edge:    {{start: v{88, 0, 0}}},
		Firefox: {{start: v{78, 0, 0}}},
		IOS:     {{start: v{14, 0, 0}}},
		Opera:   {{start: v{75, 0, 0}}},
		Safari:  {{start: v{14, 0, 0}}},
	},
	Modern_RGB_HSL: {
		Chrome:  {{start: v{66, 0, 0}}},
		Edge:    {{start: v{79, 0, 0}}},
		Firefox: {{start: v{52, 0, 0}}},
		IOS:     {{start: v{12, 2, 0}}},
		Opera:   {{start: v{53, 0, 0}}},
		Safari:  {{start: v{12, 1, 0}}},
	},
	Nesting: {
		Chrome: {{start: v{112, 0, 0}}},
	},
	RebeccaPurple: {
		Chrome:  {{start: v{38, 0, 0}}},
		Edge:    {{start: v{12, 0, 0}}},
		Firefox: {{start: v{33, 0, 0}}},
		IE:      {{start: v{11, 0, 0}}},
		IOS:     {{start: v{8, 0, 0}}},
		Opera:   {{start: v{25, 0, 0}}},
		Safari:  {{start: v{9, 0, 0}}},
	},
}

// Return all features that are not available in at least one environment
func UnsupportedCSSFeatures(constraints map[Engine][]int) (unsupported CSSFeature) {
	for feature, engines := range cssTable {
		if feature == InlineStyle {
			continue // This is purely user-specified
		}
		for engine, version := range constraints {
			if !engine.IsBrowser() {
				// Specifying "--target=es2020" shouldn't affect CSS
				continue
			}
			if versionRanges, ok := engines[engine]; !ok || !isVersionSupported(versionRanges, version) {
				unsupported |= feature
			}
		}
	}
	return
}

type CSSPrefix uint8

const (
	WebkitPrefix CSSPrefix = 1 << iota
	MozPrefix
	MsPrefix
	OPrefix

	NoPrefix CSSPrefix = 0
)

type prefixData struct {
	// Note: In some cases, earlier versions did not require a prefix but later
	// ones do. This is the case for Microsoft Edge for example, which switched
	// the underlying browser engine from a custom one to the one from Chrome.
	// However, we assume that users specifying a browser version for CSS mean
	// "works in this version or newer", so we still add a prefix when a target
	// is an old Edge version.
	engine        Engine
	withoutPrefix v
	prefix        CSSPrefix
}

var cssPrefixTable = map[css_ast.D][]prefixData{
	css_ast.DAppearance: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{84, 0, 0}},
		{engine: Edge, prefix: WebkitPrefix, withoutPrefix: v{84, 0, 0}},
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{80, 4, 0}},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{73, 4, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DBackdropFilter: {
		{engine: IOS, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix},
	},
	css_ast.DBackgroundClip: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{14, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{14, 0, 0}},
	},
	css_ast.DBoxDecorationBreak: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix},
	},
	css_ast.DClipPath: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{55, 0, 0}},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{13, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{42, 0, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{13, 1, 0}},
	},
	css_ast.DFontKerning: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{33, 0, 0}},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{12, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{20, 0, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{9, 1, 0}},
	},
	css_ast.DHyphens: {
		{engine: Edge, prefix: MsPrefix, withoutPrefix: v{79, 0, 0}},
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{43, 0, 0}},
		{engine: IE, prefix: MsPrefix},
		{engine: IOS, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix},
	},
	css_ast.DInitialLetter: {
		{engine: IOS, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix},
	},
	css_ast.DMaskImage: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DMaskOrigin: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DMaskPosition: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DMaskRepeat: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DMaskSize: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DPosition: {
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{13, 0, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{13, 0, 0}},
	},
	css_ast.DPrintColorAdjust: {
		{engine: Chrome, prefix: WebkitPrefix},
		{engine: Edge, prefix: WebkitPrefix},
		{engine: Opera, prefix: WebkitPrefix},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{15, 4, 0}},
	},
	css_ast.DTabSize: {
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{91, 0, 0}},
		{engine: Opera, prefix: OPrefix, withoutPrefix: v{15, 0, 0}},
	},
	css_ast.DTextDecorationColor: {
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{36, 0, 0}},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{12, 2, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{12, 1, 0}},
	},
	css_ast.DTextDecorationLine: {
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{36, 0, 0}},
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{12, 2, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{12, 1, 0}},
	},
	css_ast.DTextDecorationSkip: {
		{engine: IOS, prefix: WebkitPrefix, withoutPrefix: v{12, 2, 0}},
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{12, 1, 0}},
	},
	css_ast.DTextEmphasisColor: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Edge, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{85, 0, 0}},
	},
	css_ast.DTextEmphasisPosition: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Edge, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{85, 0, 0}},
	},
	css_ast.DTextEmphasisStyle: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Edge, prefix: WebkitPrefix, withoutPrefix: v{99, 0, 0}},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{85, 0, 0}},
	},
	css_ast.DTextOrientation: {
		{engine: Safari, prefix: WebkitPrefix, withoutPrefix: v{14, 0, 0}},
	},
	css_ast.DTextSizeAdjust: {
		{engine: Edge, prefix: MsPrefix, withoutPrefix: v{79, 0, 0}},
		{engine: IOS, prefix: WebkitPrefix},
	},
	css_ast.DUserSelect: {
		{engine: Chrome, prefix: WebkitPrefix, withoutPrefix: v{54, 0, 0}},
		{engine: Edge, prefix: MsPrefix, withoutPrefix: v{79, 0, 0}},
		{engine: Firefox, prefix: MozPrefix, withoutPrefix: v{69, 0, 0}},
		{engine: IE, prefix: MsPrefix},
		{engine: IOS, prefix: WebkitPrefix},
		{engine: Opera, prefix: WebkitPrefix, withoutPrefix: v{41, 0, 0}},
		{engine: Safari, prefix: WebkitPrefix},
	},
}

func CSSPrefixData(constraints map[Engine][]int) (entries map[css_ast.D]CSSPrefix) {
	for property, items := range cssPrefixTable {
		prefixes := NoPrefix
		for engine, version := range constraints {
			if !engine.IsBrowser() {
				// Specifying "--target=es2020" shouldn't affect CSS
				continue
			}
			for _, item := range items {
				if item.engine == engine && (item.withoutPrefix == v{} || compareVersions(item.withoutPrefix, version) > 0) {
					prefixes |= item.prefix
				}
			}
		}
		if prefixes != NoPrefix {
			if entries == nil {
				entries = make(map[css_ast.D]CSSPrefix)
			}
			entries[property] = prefixes
		}
	}
	return
}
