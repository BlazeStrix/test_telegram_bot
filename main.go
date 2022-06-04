package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// точка входа
func main() {
	botToken := "5045223984:AAFQ5xhusCiSOBAgu3F42Aflr1lEGqsScus"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	time_format := "2006-01-02 15:04"
	for {
		updates, chatID, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Panicln("Smth went wrong: ", err.Error())
			fmt.Println(chatID)
		}

		fmt.Println(updates)

		recording_time := time.Now()

		DataBase_file_read, err := os.ReadFile("DataBase.txt")

		if err != nil {
			log.Fatal(err)
		}

		DataBase_file_lines := strings.Split(string(DataBase_file_read), "\n")

		/*	первое повторение — сразу по окончании чтения;
			второе повторение — через 20—30 минут после первого повторения;
			третье повторение — через 1 день после второго;
			четвёртое повторение — через 2—3 недели после третьего;
			пятое повторение — через 2—3 месяца после четвёртого повторения.*/

		for i := 0; i < len(DataBase_file_lines); i++ {
			line := strings.Split(DataBase_file_lines[i], "&&&")
			past_time, error := time.Parse(time_format, line[0])
			if error != nil {
				fmt.Println("****ERROR")
				fmt.Println(error)
				fmt.Println("ERROR****")
			}
			buf := recording_time.Format(time_format)
			recording_time, error = time.Parse(time_format, buf)
			if error != nil {
				fmt.Println("!!!ERROR")
				fmt.Println(error)
				fmt.Println("ERROR!!!")
			}
			//fmt.Println("\npast_time = " + past_time.String())
			//fmt.Println("recording_time = " + recording_time.String())
			fmt.Print("\nTIME = ")
			fmt.Println(recording_time.Sub(past_time))

			sub_time := recording_time.Sub(past_time)
			x1, _ := time.ParseDuration("1m")
			x2, _ := time.ParseDuration("2m")
			x3, _ := time.ParseDuration("3m")

			/*
				на данный момент нужно пееписать свич на ифы учитывая ляйн3
				поскольку теперь у меня есть индекс строки файла, то можно будет заметь целиком строку, а потом перезаписать файл
			*/

			if sub_time == x1 && line[3] == "0" {
				fmt.Println("Сработал №1")
				y, err := strconv.Atoi(line[2])

				DataBase_file_lines[i] = line[0] + "&&&" + line[1] + "&&&" + line[2] + "&&&1"

				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}
			}

			if sub_time == x2 && line[3] == "1" {
				fmt.Println("Сработал №2")
				y, err := strconv.Atoi(line[2])

				DataBase_file_lines[i] = line[0] + "&&&" + line[1] + "&&&" + line[2] + "&&&2"

				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}
			}

			if sub_time == x3 && line[3] == "2" {
				fmt.Println("Сработал №3")
				y, err := strconv.Atoi(line[2])

				DataBase_file_lines[i] = line[0] + "&&&" + line[1] + "&&&" + line[2] + "&&&3"

				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}
			}

			/*switch sub_time {
			case x1:
				fmt.Println("Сработал №1")
				y, err := strconv.Atoi(line[2])

				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}

				//offset = updates[0].UpdateId + 1
			case x2:
				fmt.Println("Сработал №2")
				y, err := strconv.Atoi(line[2])
				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}
				//offset = updates[0].UpdateId + 1
			case x3:
				fmt.Println("Сработал №3")
				y, err := strconv.Atoi(line[2])
				err = respond_2(botUrl, y, line[1]) //чтобы эта ебола сработала, надо записывать айдишник юзера, он же чат айди и кидать его сюда, по сути у мея сейчас нет идентификации пользователя, на которого отправлять сообщение
				if err != nil {
					fmt.Println(err)
				}
				//offset = updates[0].UpdateId + 1
			default:
				fmt.Println("Сработал дефолт")
				continue
			}*/
			//fmt.Println(past_time.Add(72 * time.Hour))

		}

		err = os.WriteFile("DataBase.txt", []byte(strings.Join(DataBase_file_lines, "\n")), 0666)
		if err != nil {
			log.Fatal(err)
		}

		if len(updates) != 0 {
			//поправить айдишники, сейчас берется только первый
			DataBase_string := []byte("\n" + recording_time.Format(time_format) + "&&&" + updates[0].Message.Text + "&&&" + strconv.Itoa(updates[0].Message.Chat.ChatId) + "&&&0")

			for _, update := range updates {
				//берем последнюю строку файла, пихаем ее в респонд и вызываем метод с задержкой
				err = respond(botUrl, update)
				offset = update.UpdateId + 1
			}

			DataBase_file_write, err := os.OpenFile("DataBase.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}
			defer DataBase_file_write.Close()

			if _, err = DataBase_file_write.Write(DataBase_string); err != nil {
				panic(err)
			}
		}

	}
}

//запрос обновлений
func getUpdates(botUrl string, offset int) ([]Update, []int, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, nil, err
	}
	return restResponse.Result, restResponse.ChatId, nil
}

//ответ на обновление
func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

//ответ на обновление
func respond_2(botUrl string, update int, resp_message string) error {
	var botMessage BotMessage
	botMessage.ChatId = update
	botMessage.Text = resp_message
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func getUpdates_2(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	var botMessage BotMessage
	err = json.Unmarshal(body, &botMessage)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}
