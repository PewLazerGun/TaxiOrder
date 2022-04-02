package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	_ "github.com/mattn/go-sqlite3"

	"go-projects/taxi"
)

type PostmanResp struct {
	Data Order `json:"data"`
}

type Order struct {
	NumberOfCar string `json:"number_of_car"`
	Name        string `json:"name"`
	TypeOfTaxi  string `json:"type_of_taxi"`
	Passnumber  int    `json:"passnumber"`
	Sum         int    `json:"sum"`
}

func sumOfLuggage(passanger []taxi.Passenger) int {
	sum := 0
	for _, v := range passanger {
		sum += v.Luggage
	}
	return sum
}

func main() {
	//Пассажиры
	x := []taxi.Passenger{{10, 20}, {15, 20}, {60, 15}}

	//Подключение к БД
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	Orders := passangersByDistance(x)
	for _, v := range Orders {
		switch { //Условие вызова такси с учетом кол-ва пассажиров и кол-ва багажа:
		case sumOfLuggage(v) > 50 && len(v) <= 2: //Грузовое
			t, ok := taxi.NewCargoTaxi().(taxi.CargoTaxi)
			if !ok {
				fmt.Println(errors.New("ошибка при инициализации грузового такси"))
			}
			result, err := MakeRequest(t.CarNumber, t.DriverName, t.TypeOfTaxi, len(v), t.Price(v)) //post запрос
			if err != nil {
				fmt.Println(err.Error())
			}
			err = insertDB(result)
			if err != nil {
				fmt.Println(err.Error())
			}
		case len(v) <= 4 && sumOfLuggage(v) <= 50: //Пассажирское
			t, ok := taxi.NewPassengerTaxi().(taxi.PassengerTaxi)
			if !ok {
				fmt.Println(errors.New("ошибка при инициализации легкового такси"))
			}
			result, err := MakeRequest(t.CarNumber, t.DriverName, t.TypeOfTaxi, len(v), t.Price(v)) //post запрос
			if err != nil {
				fmt.Println(err.Error())
			}
			err = insertDB(result)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Println("Невозможно заказать машину для такого количества пассажиров и багажа!", v)
		}
	}
	printDB(db) // Вывод базы
	fmt.Println()

}

//Вывод базы в консоль
func printDB(db *sql.DB) {
	rows, err := db.Query("select * from localdb")
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		var id int
		var name string
		var number_of_car string
		var type_of_taxi string
		var number_of_passengers int
		var price int

		err := rows.Scan(&id, &number_of_car, &name, &type_of_taxi, &number_of_passengers, &price)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Водитель: %s, Номер машины: %s, Тип такси: %s, Кол-во пассажиров: %d, Сумма заказа: %d\n", name, number_of_car, type_of_taxi, number_of_passengers, price)

	}

	defer rows.Close()
}

func insertDB(order Order) error {

	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO localdb(number_of_car, name, type, number_of_passengers, price) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(order.NumberOfCar, order.Name, order.TypeOfTaxi, order.Passnumber, order.Sum)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer stmt.Close()
	return nil
}

//получает массив пассажиров, сортирует, возвращая двумерный массив в котором элемент это массив пассажиров с одинаковым расстоянием
func passangersByDistance(array []taxi.Passenger) [][]taxi.Passenger {
	sort.Sort(taxi.PassengerList(array)) //сортировка библиотечная функция
	len := 1
	tmp := array[0].Distance
	//находим кол-во уникальных расстояний это длина результирующего массива
	for _, v := range array {
		if v.Distance != tmp {
			len += 1
			tmp = v.Distance
		}
	}
	Result := make([][]taxi.Passenger, len)
	//разделяем входной массив на подмассивы с одинаковым расстоянием
	i := 0
	tmp = array[0].Distance
	for _, v := range array {
		if v.Distance != tmp {
			i += 1
			tmp = v.Distance
		}
		Result[i] = append(Result[i], v)
	}
	return Result
}

//http-request
func MakeRequest(number_of_car string, name string, type_of_taxi string, passnumber int, sum int) (Order, error) {

	message := map[string]interface{}{
		"number_of_car": number_of_car,
		"name":          name,
		"type_of_taxi":  type_of_taxi,
		"passnumber":    passnumber,
		"sum":           sum,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://postman-echo.com/post", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var result PostmanResp

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err.Error())
	}

	return result.Data, nil
}
