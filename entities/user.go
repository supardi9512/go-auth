package entities

type User struct {
	Id        int64
	Name      string `validate:"required"`
	Email     string `validate:"required,email,isunique=users-email"`
	Username  string `validate:"required,gte=3,isunique=users-username"`
	Password  string `validate:"required,gte=6"`
	Cpassword string `validate:"required,eqfield=Password" label:"Confirm Password"`
}
