package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/sorcererxw/jikeview-bot/bot"
)

//Handler handles Vercel serverless request
func Handler(w http.ResponseWriter, r *http.Request) {
	stream, err := r.GetBody()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := ioutil.ReadAll(stream)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var update telebot.Update
	if err := json.Unmarshal(bytes, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bot.Bot.ProcessUpdate(update)
}
