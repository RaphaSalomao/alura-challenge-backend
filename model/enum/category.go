package enum

type Category string

const (
	CategoryUndefined  Category = ""
	CategoryFood       Category = "FOOD"
	CategoryHealth     Category = "HEALTH"
	CategoryHome       Category = "HOME"
	CategoryTransport  Category = "TRANSPORT"
	CategoryEducation  Category = "EDUCATION"
	CategoryLeisure    Category = "LEISURE"
	CategoryUnforeseen Category = "UNFORESEEN"
	CategoryOthers     Category = "OTHERS"
)
