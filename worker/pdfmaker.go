package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

const pathToPdfs = "../store/pdfs"

func getImg(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}
	defer response.Body.Close()

	dotIdx := strings.LastIndex(url, ".")
	if dotIdx == -1 {
		return "", errors.New("bad url")
	}

	path := "tmp/" + strconv.Itoa(time.Now().Nanosecond()) + "_" + strconv.Itoa(len(url)) + "." + url[dotIdx+1:]
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func createPdf(msg *QueueMsg) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 25)
	pdf.MoveTo(75, 10)
	pdf.Cell(50, 15, "S O A - M A F I A")

	pdf.SetFont("Arial", "B", 20)
	pdf.MoveTo(20, 30)
	pdf.Cell(40, 15, "Player Info:")

	pdf.SetFont("Arial", "", 15)
	pdf.MoveTo(20, 45)
	pdf.Cell(40, 10, "Username: ")
	pdf.MoveTo(20, 55)
	pdf.Cell(40, 10, "Sex: ")
	pdf.MoveTo(20, 65)
	pdf.Cell(40, 10, "Email: ")
	pdf.MoveTo(20, 75)

	pdf.SetFont("Arial", "B", 20)
	pdf.MoveTo(20, 95)
	pdf.Cell(40, 15, "Player Stats:")

	pdf.SetFont("Arial", "", 15)
	pdf.MoveTo(20, 110)
	pdf.Cell(40, 10, "Session count:")
	pdf.MoveTo(20, 120)
	pdf.Cell(40, 10, "Win count:")
	pdf.MoveTo(20, 130)
	pdf.Cell(40, 10, "Loss count:")
	pdf.MoveTo(20, 140)
	pdf.Cell(40, 10, "Duration:")

	pdf.MoveTo(60, 45)
	pdf.Cell(40, 10, msg.Username)
	pdf.MoveTo(60, 55)
	pdf.Cell(40, 10, msg.Sex)
	pdf.MoveTo(60, 65)
	pdf.Cell(40, 10, msg.Email)

	pdf.MoveTo(60, 110)
	pdf.Cell(40, 10, strconv.Itoa(msg.LossCount+msg.WinCount))
	pdf.MoveTo(60, 120)
	pdf.Cell(40, 10, strconv.Itoa(msg.WinCount))
	pdf.MoveTo(60, 130)
	pdf.Cell(40, 10, strconv.Itoa(msg.LossCount))
	pdf.MoveTo(60, 140)
	pdf.Cell(40, 10, strconv.Itoa(msg.Duration)+" min")

	pdf.ImageOptions("mafia-sign.jpg",
		8, 180,
		0, 0,
		false,
		fpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
		0,
		"",
	)

	pdf.SetFont("Arial", "", 10)
	path, imgErr := getImg(msg.Avatar)
	if imgErr == nil {
		pdf.ImageOptions(path,
			128, 32,
			60, 60,
			false,
			fpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
			0,
			"",
		)
		pdf.MoveTo(153, 92)
		pdf.Cell(40, 10, "Avatar")

		os.Remove(path)
	} else {
		pdf.MoveTo(140, 60)
		pdf.Cell(40, 10, "Avatar URL is Incorrect!")
		log.Print(imgErr.Error())
	}

	path = pathToPdfs + "/" + msg.ID

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err := pdf.Output(writer)
	if err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}
