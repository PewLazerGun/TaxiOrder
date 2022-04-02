package taxi

//Структура "Пассажир" с весом багажа и расстоянием поездки
type Passenger struct {
	Luggage  int
	Distance int
}

//Двумерка для работы функции passengersByDistance
type PassengerList []Passenger

func (e PassengerList) Len() int {
	return len(e)
}

func (e PassengerList) Less(i, j int) bool {
	return e[i].Distance > e[j].Distance
}

func (e PassengerList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
