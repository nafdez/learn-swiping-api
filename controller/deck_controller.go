package controller

type DeckController interface {
	Create()     // POST
	Deck()       // GET
	OwnedDecks() // GET
	Decks()      // GET (from other users
	Update()     // PUT
	Delete()     // DELETE
}
