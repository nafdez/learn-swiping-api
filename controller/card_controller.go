package controller

type CardController interface {
	Create()       // POST
	Card()         // GET
	Update()       // PUT
	Delete()       // DELETE
	CreateAnswer() // POST
	Answers()      // GET
	UpdateAnswer() // PUT
	DeleteAnswer() // DELETE
}
