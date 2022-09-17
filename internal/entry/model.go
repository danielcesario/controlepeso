package entry

type Entry struct {
	ID     int     `json:"id"`
	UserId int     `json:"user_id"`
	Weight float64 `json:"weight"`
	Date   string  `json:"date"`
}
