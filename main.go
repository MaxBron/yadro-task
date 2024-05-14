package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func remove(list []string, i int) []string {
	list[i] = list[0]
	return list[1:]
}

func main() {
	file, _ := os.Open(os.Args[1])
	defer file.Close()
	sc := bufio.NewScanner(file)
	sc.Scan()
	countTables, _ := strconv.Atoi(sc.Text()) // количество столов
	var tables []Table
	var clients, queue []string
	for i := 0; i < countTables; i++ {
		table := Table{}
		tables = append(tables, table)
	}

	sc.Scan()
	start, _ := time.Parse("15:04", strings.Split(sc.Text(), " ")[0]) // открытие
	end, _ := time.Parse("15:04", strings.Split(sc.Text(), " ")[1])   // закрытие
	sc.Scan()
	cost, _ := strconv.Atoi(sc.Text()) // цена за час
	fmt.Println(start.Format("15:04"))
	for sc.Scan() {
		fmt.Println(sc.Text())
		str := strings.Split(sc.Text(), " ")
		clientStartTime, _ := time.Parse("15:04", str[0])
		// клиент пришел
		if str[1] == "1" {
			flag := true
			// проверка на то, находится ли клиент уже в списке
			for _, client := range clients {
				if client == str[2] {
					flag = false
					fmt.Println(str[0], 13, "YouShallNotPass")
				}

			}

			//проверка на то, пришёл ли клиент после открытия
			if start.After(clientStartTime) && flag {
				fmt.Println(str[0], 13, "NotOpenYet")
			} else {
				clients = append(clients, str[2])
			}

			//клиент сел за стол
		} else if str[1] == "2" {
			flag := false
			// есть ли клиент в списке
			for _, client := range clients {
				if client == str[2] {
					flag = true
					break
				}

			}

			if !flag {
				fmt.Println(str[0], 13, "ClientUnknown")
			} else {
				number, _ := strconv.Atoi(str[3])
				// место занято
				if tables[number-1].busy {
					fmt.Println(str[0], 13, "PlaceIsBusy")
				} else {
					tables[number-1].start, _ = time.Parse("15:04", str[0])
					tables[number-1].busy = true
					tables[number-1].clientName = str[2]
					countTables -= 1
				}

				// удаление из очереди
				if len(queue) > 0 {
					if queue[0] == str[2] {
						queue = queue[1:]
					}

				}

			}

			// клиент ожидает
		} else if str[1] == "3" {
			if countTables < 0 {
				for i := 0; i < len(clients); i++ {
					if clients[i] == str[2] {
						clients[i] = ""
						break
					}

				}

				fmt.Println(str[0], 11, str[2])
			} else {
				queue = append(queue, str[2])
				// если есть свободные столы, клиент не хочет ждать
				if countTables > 0 {
					fmt.Println(str[0], 13, "ICanWaitNoLonger!")
				}

			}

			// клиент ушёл
		} else if str[1] == "4" {
			flag := false
			var clientsId int
			// есть ли клиент в списке
			for i, client := range clients {
				if client == str[2] {
					flag = true
					clientsId = i
					break
				}

			}

			if !flag {
				fmt.Println(str[0], 13, "ClientUnknown")
			} else {
				clients = remove(clients, clientsId)
				for i := 0; i < len(tables); i++ {
					// освобождение стола и подсчёт времени
					if tables[i].clientName == str[2] {
						tables[i].clientName = ""
						tables[i].busy = false
						tables[i].end, _ = time.Parse("15:04", str[0])
						tables[i].hours += int(tables[i].end.Hour()) - int(tables[i].start.Hour())
						tables[i].duration = tables[i].duration.Add(tables[i].end.Sub(tables[i].start))
						if len(queue) > 0 {
							tables[i].clientName = queue[0]
							tables[i].start, _ = time.Parse("15:04", str[0])
							tables[i].busy = true
							fmt.Println(str[0], 12, queue[0], i+1)
							queue = queue[1:]
						} else {
							if int(tables[i].end.Minute()) > 0 {
								tables[i].hours += 1
							}

							// подсчёт вырчучки
							tables[i].revenue += (cost * tables[i].hours)
						}

						break
					}

				}

				countTables += 1
			}

		}

	}

	// оставшиеся клиенты к концу дня
	if len(clients) > 0 {
		sort.Strings(clients)
		for _, client := range clients {
			for i := 0; i < len(tables); i++ {
				if client == tables[i].clientName {
					tables[i].end = end
					tables[i].hours = int(tables[i].end.Hour()) - int(tables[i].start.Hour())
					if int(tables[i].end.Minute()) > 0 {
						tables[i].hours += 1
					}

					tables[i].revenue += tables[i].hours * cost
					tables[i].duration = tables[i].duration.Add(tables[i].end.Sub(tables[i].start))
					break
				}

			}

			fmt.Println(end.Format("15:04"), 11, client)
		}

	}

	fmt.Println(end.Format("15:04"))
	// вывод выручки с каждого стола
	for i, table := range tables {
		fmt.Println(i+1, table.revenue, table.duration.Format("15:04"))
	}

}
