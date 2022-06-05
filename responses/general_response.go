package responses

type GeneralResponse struct {
	Status  int
	Message string
	Data    interface{}
}
