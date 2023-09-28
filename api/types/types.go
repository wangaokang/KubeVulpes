package types

type User struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Status      int8   `json:"status"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	Description string `json:"description"`

	TimeOption `json:",inline"`
}
