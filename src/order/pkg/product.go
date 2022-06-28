package pkg

type Product struct {
	ID             string   `json:"id" yaml:"id"`
	URL            string   `json:"url,omitempty" yaml:"url,omitempty"`
	Name           string   `json:"name" yaml:"name"`
	Category       string   `json:"category" yaml:"category"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	Aliases        []string `json:"aliases,omitempty" yaml:"aliases,omitempty"` // keywords for use with e.g. Alexa
	Price          float32  `json:"price" yaml:"price"`
	Image          string   `json:"image,omitempty" yaml:"image,omitempty"`
	Featured       string   `json:"featured,omitempty" yaml:"featured,omitempty"`
	GenderAffinity string   `json:"gender,omitempty" yaml:"gender,omitempty"`
	Affinity       []string `json:"affinity" yaml:"affinity"`
	Stock          int      `json:"stock" yaml:"stock"`
}

// Initialized - indicates if instance has been initialized or not
func (p *Product) Initialized() bool { return p != nil && len(p.ID) > 0 }

// Products Array
type Products []Product

// Inventory Struct
type Inventory struct {
	StockDelta int `json:"stock_delta" yaml:"stock_delta"`
}
