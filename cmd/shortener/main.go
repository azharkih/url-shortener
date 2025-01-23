package main

import "url-shortener/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
