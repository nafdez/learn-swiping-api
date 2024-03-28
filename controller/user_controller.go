package controller

type UserController interface {
	Register() // POST
	Login()    // POST
	Logout()   // POST
	Account()  // POST
	User()     // GET
	Update()   // PUT
	Delete()   // DELETE
}
