package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Bike struct {
	ID       int
	Title    string
	Price    int
	ImageURL string
}

func main() {
	db, err := sql.Open("sqlite", "shop.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initSQL := `
	CREATE TABLE IF NOT EXISTS bikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		price INTEGER,
		image_url TEXT
	);`
	_, err = db.Exec(initSQL)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM bikes").Scan(&count)
	if count == 0 {
		insertSQL := `INSERT INTO bikes (title, price, image_url) VALUES 
		('Горный Велосипед', 25000, 'https://githubusercontent.com'),
		('Городской Круизер', 18000, 'https://githubusercontent.com'),
		('BMX Экстрим', 32000, 'https://githubusercontent.com');`
		db.Exec(insertSQL)
	}

	rows, err := db.Query("SELECT id, title, price, image_url FROM bikes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bikesList []Bike

	for rows.Next() {
		var b Bike
		err := rows.Scan(&b.ID, &b.Title, &b.Price, &b.ImageURL)
		if err != nil {
			log.Fatal(err)
		}
		bikesList = append(bikesList, b)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Магазин Велосипедов из БД SQL")
	myWindow.Resize(fyne.NewSize(750, 450))

	mainRow := container.NewHBox()

	for _, bike := range bikesList {
		currentBike := bike

		txtTitle := widget.NewLabel(currentBike.Title)
		txtPrice := widget.NewLabel(fmt.Sprintf("Цена: %d руб.", currentBike.Price))

		uri := storage.NewURI(currentBike.ImageURL)
		img := canvas.NewImageFromURI(uri)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(150, 150))

		wasPressed := false
		var btn *widget.Button
		btn = widget.NewButton("Купить", func() {
			if wasPressed {
				btn.SetText("Уже в корзине!")
				return
			}
			btn.SetText("Успешно куплено!")
			wasPressed = true
			time.AfterFunc(2*time.Second, func() {
				btn.SetText("Купить")
			})
		})

		card := container.NewVBox(
			container.NewCenter(txtTitle),
			container.NewCenter(txtPrice),
			container.NewCenter(img),
			btn,
		)

		mainRow.Add(card)
	}

	scrollableRow := container.NewHScroll(mainRow)
	myWindow.SetContent(container.NewCenter(scrollableRow))
	myWindow.ShowAndRun()
}
