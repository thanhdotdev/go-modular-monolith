package customerapplication

type CustomerDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Tier  string `json:"tier"`
}
