package main

type GoogleDoc struct {
	Body struct {
		Content []struct {
			EndIndex     float64 `json:"endIndex,omitempty"`
			SectionBreak struct {
				SectionStyle struct {
					ColumnSeparatorStyle string `json:"columnSeparatorStyle,omitempty"`
					ContentDirection     string `json:"contentDirection,omitempty"`
					SectionType          string `json:"sectionType,omitempty"`
				} `json:"sectionStyle,omitempty"`
			} `json:"sectionBreak,omitempty"`
		} `json:"content,omitempty"`
	} `json:"body,omitempty"`
	DocumentId    string `json:"documentId,omitempty"`
	DocumentStyle struct {
		Background struct {
			Color struct {
			} `json:"color,omitempty"`
		} `json:"background,omitempty"`
		DefaultHeaderId string `json:"defaultHeaderId,omitempty"`
		MarginBottom    struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginBottom,omitempty"`
		MarginFooter struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginFooter,omitempty"`
		MarginHeader struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginHeader,omitempty"`
		MarginLeft struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginLeft,omitempty"`
		MarginRight struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginRight,omitempty"`
		MarginTop struct {
			Magnitude float64 `json:"magnitude,omitempty"`
			Unit      string  `json:"unit,omitempty"`
		} `json:"marginTop,omitempty"`
		PageNumberStart float64 `json:"pageNumberStart,omitempty"`
		PageSize        struct {
			Height struct {
				Magnitude float64 `json:"magnitude,omitempty"`
				Unit      string  `json:"unit,omitempty"`
			} `json:"height,omitempty"`
			Width struct {
				Magnitude float64 `json:"magnitude,omitempty"`
				Unit      string  `json:"unit,omitempty"`
			} `json:"width,omitempty"`
		} `json:"pageSize,omitempty"`
		UseCustomHeaderFooterMargins bool `json:"useCustomHeaderFooterMargins,omitempty"`
	} `json:"documentStyle,omitempty"`
	Headers struct {
		Kix_Dt94tg7t3bsh struct {
			Content []struct {
				EndIndex  float64 `json:"endIndex,omitempty"`
				Paragraph struct {
					Elements []struct {
						EndIndex float64 `json:"endIndex,omitempty"`
						TextRun  struct {
							Content   string `json:"content,omitempty"`
							TextStyle struct {
							} `json:"textStyle,omitempty"`
						} `json:"textRun,omitempty"`
					} `json:"elements,omitempty"`
					ParagraphStyle struct {
						Direction      string `json:"direction,omitempty"`
						NamedStyleType string `json:"namedStyleType,omitempty"`
					} `json:"paragraphStyle,omitempty"`
				} `json:"paragraph,omitempty"`
			} `json:"content,omitempty"`
			HeaderId string `json:"headerId,omitempty"`
		} `json:"kix.dt94tg7t3bsh,omitempty"`
	} `json:"headers,omitempty"`
	InlineObjects struct {
		Kix_I6pahq99cw6k struct {
			InlineObjectProperties struct {
				EmbeddedObject struct {
					EmbeddedObjectBorder struct {
						Color struct {
							Color struct {
								RgbColor struct {
								} `json:"rgbColor,omitempty"`
							} `json:"color,omitempty"`
						} `json:"color,omitempty"`
						DashStyle     string `json:"dashStyle,omitempty"`
						PropertyState string `json:"propertyState,omitempty"`
						Width         struct {
							Unit string `json:"unit,omitempty"`
						} `json:"width,omitempty"`
					} `json:"embeddedObjectBorder,omitempty"`
					ImageProperties struct {
						ContentUri     string `json:"contentUri,omitempty"`
						CropProperties struct {
							OffsetLeft  float64 `json:"offsetLeft,omitempty"`
							OffsetRight float64 `json:"offsetRight,omitempty"`
						} `json:"cropProperties,omitempty"`
					} `json:"imageProperties,omitempty"`
					MarginBottom struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"marginBottom,omitempty"`
					MarginLeft struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"marginLeft,omitempty"`
					MarginRight struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"marginRight,omitempty"`
					MarginTop struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"marginTop,omitempty"`
					Size struct {
						Height struct {
							Magnitude float64 `json:"magnitude,omitempty"`
							Unit      string  `json:"unit,omitempty"`
						} `json:"height,omitempty"`
						Width struct {
							Magnitude float64 `json:"magnitude,omitempty"`
							Unit      string  `json:"unit,omitempty"`
						} `json:"width,omitempty"`
					} `json:"size,omitempty"`
				} `json:"embeddedObject,omitempty"`
			} `json:"inlineObjectProperties,omitempty"`
			ObjectId string `json:"objectId,omitempty"`
		} `json:"kix.i6pahq99cw6k,omitempty"`
	} `json:"inlineObjects,omitempty"`
	Lists struct {
		Kix_C9u0abyvo3qw struct {
			ListProperties struct {
				NestingLevels []struct {
					BulletAlignment string `json:"bulletAlignment,omitempty"`
					GlyphFormat     string `json:"glyphFormat,omitempty"`
					GlyphType       string `json:"glyphType,omitempty"`
					IndentFirstLine struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"indentFirstLine,omitempty"`
					IndentStart struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"indentStart,omitempty"`
					StartNumber float64 `json:"startNumber,omitempty"`
					TextStyle   struct {
					} `json:"textStyle,omitempty"`
				} `json:"nestingLevels,omitempty"`
			} `json:"listProperties,omitempty"`
		} `json:"kix.c9u0abyvo3qw,omitempty"`
		Kix_Dnigbincn0qe struct {
			ListProperties struct {
				NestingLevels []struct {
					BulletAlignment string `json:"bulletAlignment,omitempty"`
					GlyphFormat     string `json:"glyphFormat,omitempty"`
					GlyphSymbol     string `json:"glyphSymbol,omitempty"`
					IndentFirstLine struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"indentFirstLine,omitempty"`
					IndentStart struct {
						Magnitude float64 `json:"magnitude,omitempty"`
						Unit      string  `json:"unit,omitempty"`
					} `json:"indentStart,omitempty"`
					StartNumber float64 `json:"startNumber,omitempty"`
					TextStyle   struct {
					} `json:"textStyle,omitempty"`
				} `json:"nestingLevels,omitempty"`
			} `json:"listProperties,omitempty"`
		} `json:"kix.dnigbincn0qe,omitempty"`
	} `json:"lists,omitempty"`
	NamedStyles struct {
		Styles []struct {
			NamedStyleType string `json:"namedStyleType,omitempty"`
			ParagraphStyle struct {
				Alignment           string `json:"alignment,omitempty"`
				AvoidWidowAndOrphan bool   `json:"avoidWidowAndOrphan,omitempty"`
				BorderBetween       struct {
					Color struct {
					} `json:"color,omitempty"`
					DashStyle string `json:"dashStyle,omitempty"`
					Padding   struct {
						Unit string `json:"unit,omitempty"`
					} `json:"padding,omitempty"`
					Width struct {
						Unit string `json:"unit,omitempty"`
					} `json:"width,omitempty"`
				} `json:"borderBetween,omitempty"`
				BorderBottom struct {
					Color struct {
					} `json:"color,omitempty"`
					DashStyle string `json:"dashStyle,omitempty"`
					Padding   struct {
						Unit string `json:"unit,omitempty"`
					} `json:"padding,omitempty"`
					Width struct {
						Unit string `json:"unit,omitempty"`
					} `json:"width,omitempty"`
				} `json:"borderBottom,omitempty"`
				BorderLeft struct {
					Color struct {
					} `json:"color,omitempty"`
					DashStyle string `json:"dashStyle,omitempty"`
					Padding   struct {
						Unit string `json:"unit,omitempty"`
					} `json:"padding,omitempty"`
					Width struct {
						Unit string `json:"unit,omitempty"`
					} `json:"width,omitempty"`
				} `json:"borderLeft,omitempty"`
				BorderRight struct {
					Color struct {
					} `json:"color,omitempty"`
					DashStyle string `json:"dashStyle,omitempty"`
					Padding   struct {
						Unit string `json:"unit,omitempty"`
					} `json:"padding,omitempty"`
					Width struct {
						Unit string `json:"unit,omitempty"`
					} `json:"width,omitempty"`
				} `json:"borderRight,omitempty"`
				BorderTop struct {
					Color struct {
					} `json:"color,omitempty"`
					DashStyle string `json:"dashStyle,omitempty"`
					Padding   struct {
						Unit string `json:"unit,omitempty"`
					} `json:"padding,omitempty"`
					Width struct {
						Unit string `json:"unit,omitempty"`
					} `json:"width,omitempty"`
				} `json:"borderTop,omitempty"`
				Direction string `json:"direction,omitempty"`
				IndentEnd struct {
					Unit string `json:"unit,omitempty"`
				} `json:"indentEnd,omitempty"`
				IndentFirstLine struct {
					Unit string `json:"unit,omitempty"`
				} `json:"indentFirstLine,omitempty"`
				IndentStart struct {
					Unit string `json:"unit,omitempty"`
				} `json:"indentStart,omitempty"`
				LineSpacing    float64 `json:"lineSpacing,omitempty"`
				NamedStyleType string  `json:"namedStyleType,omitempty"`
				Shading        struct {
					BackgroundColor struct {
					} `json:"backgroundColor,omitempty"`
				} `json:"shading,omitempty"`
				SpaceAbove struct {
					Unit string `json:"unit,omitempty"`
				} `json:"spaceAbove,omitempty"`
				SpaceBelow struct {
					Magnitude float64 `json:"magnitude,omitempty"`
					Unit      string  `json:"unit,omitempty"`
				} `json:"spaceBelow,omitempty"`
				SpacingMode string `json:"spacingMode,omitempty"`
			} `json:"paragraphStyle,omitempty"`
			TextStyle struct {
				BackgroundColor struct {
				} `json:"backgroundColor,omitempty"`
				BaselineOffset string `json:"baselineOffset,omitempty"`
				FontSize       struct {
					Magnitude float64 `json:"magnitude,omitempty"`
					Unit      string  `json:"unit,omitempty"`
				} `json:"fontSize,omitempty"`
				ForegroundColor struct {
					Color struct {
						RgbColor struct {
							Blue  float64 `json:"blue,omitempty"`
							Green float64 `json:"green,omitempty"`
							Red   float64 `json:"red,omitempty"`
						} `json:"rgbColor,omitempty"`
					} `json:"color,omitempty"`
				} `json:"foregroundColor,omitempty"`
				WeightedFontFamily struct {
					FontFamily string  `json:"fontFamily,omitempty"`
					Weight     float64 `json:"weight,omitempty"`
				} `json:"weightedFontFamily,omitempty"`
			} `json:"textStyle,omitempty"`
		} `json:"styles,omitempty"`
	} `json:"namedStyles,omitempty"`
	RevisionId          string `json:"revisionId,omitempty"`
	SuggestionsViewMode string `json:"suggestionsViewMode,omitempty"`
	Title               string `json:"title,omitempty"`
}
