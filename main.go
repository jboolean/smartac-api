package main

import (
	"fmt"
	// "github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/api/power", handlePower)
	http.ListenAndServe(":"+port, nil)
}

func handlePower(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	powerOn, err := strconv.ParseBool(r.FormValue("power"))

	if err != nil {
		panic(err)
	}

	switchPower(email, password, powerOn)

	w.WriteHeader(200)
}

func switchPower(email string, password string, powerOn bool) {
	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open("https://mymodlet.com/SmartAC")
	if err != nil {
		panic(err)
	}

	fm, err := bow.Form("#header-login-form")

	fm.Input("loginForm.Email", email)
	fm.Input("loginForm.Password", password)

	err = fm.Submit()

	if err != nil {
		panic(err)
	}

	bow.Click("a.acOnOff.Off")

	fmt.Println(bow.Dom().Find("a.acOnOff.Off").Text())

	applianceId, hasAppliance := bow.Dom().Find("[name=ApplianceId]").Attr("value")

	if !hasAppliance {
		panic("There are no devices")
	}

	currentTemp, _ := bow.Dom().Find(".temperatureDropDown option[selected]").Attr("value")

	body := fmt.Sprintf("{\"applianceId\":\"%s\",\"targetTemperature\":\"%s\",\"thermostated\":%t}", applianceId, currentTemp, powerOn)

	err = bow.Post("https://mymodlet.com/SmartAC/UserSettings", "application/json", strings.NewReader(body))

	if err != nil {
		panic(err)
	}

	bow.Open("https://mymodlet.com/Account/Logout")
}
