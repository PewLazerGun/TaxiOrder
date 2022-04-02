package taxi

//Определяем интерфейс "Такси" с методом рассчета цены
type Taxi interface {
	Price(passengers []Passenger) int
}

//Наследуемая структура легковое такси
type PassengerTaxi struct {
	CarNumber  string
	DriverName string
	TypeOfTaxi string
}

//Наследуемая структура грузовое такси
type CargoTaxi struct {
	CarNumber  string
	DriverName string
	TypeOfTaxi string
}

//Метод для рассчета стоимости поездки на легковом такси
func (p PassengerTaxi) Price(passengers []Passenger) int {
	sum := passengers[0].Distance * 2
	return sum
}

//Метод для рассчета стоимости поездки на грузовом такси
func (c CargoTaxi) Price(passengers []Passenger) int {
	sum := passengers[0].Distance * 3
	return sum
}

func NewCargoTaxi() Taxi {
	var result = CargoTaxi{"a111aa", "Name", "Cargo"}
	return result
}

func NewPassengerTaxi() Taxi {
	var result = PassengerTaxi{"a111aa", "Name", "Passenger"}
	return result
}
