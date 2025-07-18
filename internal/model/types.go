package model

type Optimized struct {
	OptimizedMasterList []string `json:"optimizedMasterList"` // Optimized master list
	PreferenceTags      []string `json:"preferenceTags"`      // Optimized user preference
}

type Transformed struct {
	Title               string   `json:"title"`               // !Model's title based on text content
	Alignment           int      `json:"alignment"`           // !Model's framework alignment score (0–10)
	Tags                []string `json:"tags"`                // !Model's categories/topics, including framework themes
	PoliticalBiases     []string `json:"politicalBiases"`     // !Model's examination and tagging for political biases eg: leftist, right, progressive, conservative
	Summary             string   `json:"summary"`             // !Model's thoughtful summary with framework reflection
	InSufficientContent bool     `json:"InSufficientContent"` //!Model's flagging
}
