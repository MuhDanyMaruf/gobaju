package users

type Users struct {
	UsersId  int    `json:"users_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
