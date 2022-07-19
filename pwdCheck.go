package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
)

func main() {

	scoringSystem := map[int]string{
		1: "Poor",
		2: "Unsatisfactory",
		3: "Satisfactory",
		4: "Good",
		5: "Very Good",
		6: "Excellent",
	}

	var pass string
	var isLeaked bool
	var count int64

	fmt.Println("ENTER YOUR PASSWORD: ")

	fmt.Scanln(&pass)

	res, err := checkpass(pass)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Your score is [%d]. It is considered [%s] \n", res, scoringSystem[res])
	}

	isLeaked, count = requestApi(pass)

	if isLeaked {
		fmt.Println("This password has been leaked. Number of occurrences in leaks:", count)
	}
}

func checkpass(pass string) (int, error) {

	score := 0

	if len(pass) > 0 {

		capitalLetters := regexp.MustCompile(`[A-Z]`)
		capitalLettersFound := capitalLetters.MatchString(pass)

		lowerLetters := regexp.MustCompile(`[a-z]`)
		lowerLettersFound := lowerLetters.MatchString(pass)

		moreThanTwelve := regexp.MustCompile(`[a-zA-Z0-9$&+,:;=?@#|'<>.^*()%!-]{12,}$`)
		moreThanTwelveFound := moreThanTwelve.MatchString(pass)

		anyNumbers := regexp.MustCompile(`[0-9]`)
		anyNumbersFound := anyNumbers.MatchString(pass)

		specialCharacter := regexp.MustCompile(`^.*[$&+,:;=?@#|'<>.^*()%!-]+.*$`)
		specialCharacterFound := specialCharacter.MatchString(pass)

		if capitalLettersFound {
			score++
		}
		if lowerLettersFound {
			score++
		}
		if moreThanTwelveFound {
			score += 2
		}
		if anyNumbersFound {
			score++
		}
		if specialCharacterFound {
			score++
		}
		return score, nil
	}
	err := errors.New("PASSWORD HAS NOT BEEN SET")
	return score, err
}

func requestApi(pass string) (bool, int64) {
	url := "https://check.cybernews.com/chk-pw/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("p", pass)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	p, _ := jsonparser.GetBoolean(body, "p")
	o, _ := jsonparser.GetInt(body, "o")

	return p, o
}