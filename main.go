package main

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

func main() {
	client := resty.New()
	defer client.Close()

	res, err := client.R().Get("https://aniwatch-api-rosy-one.vercel.app/api/v2/hianime/home")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(res)
}
