package pkg

type Category struct {
	ID             string `json:"id" yaml:"id"`
	URL            string `json:"url" yaml:"url"`
	Name           string `json:"name" yaml:"name"`
	Image          string `json:"image" yaml:"image"`
	GenderAffinity string `json:"gender_affinity" yaml:"gender_affinity"`
}

// Initialized - indicates if instance has been initialized or not
func (c *Category) Initialized() bool { return c != nil && len(c.ID) > 0 }

// Categories Array
type Categories []Category
