package controlepeso

type Entry struct {
	UserId int     `json:"userId"`
	Weight float64 `json:"weight"`
	Date   string  `json:"date"`
}
