package api

//################file for endpoints form tag messages##########
import (
	"net/http"
	//"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
	"github.com/julienschmidt/httprouter"
)

type SendMessageRequest struct {
	Content              string //`json:"content"`
	RespondedToMessageID int    //`json:"respondedToMessageId"`
	ConversationID       int    //`json:"conversationId"`/// trzeba tez ustalic jak wygladaje te id
	SenderUsername       string //`json:"senderUsername"`
	ReciverUsername      string //`json:"receiverUsername"`
	Attachmemnts         attachments.AttachmentsPack
}

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var rqst SendMessageRequest

	err := r.ParseMultipartForm(10 << 22)
	if err != nil {
		rt.baseLogger.Printf("Error parsing multipart form")
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rqst.SenderUsername = r.FormValue("senderUsername")
	if rqst.SenderUsername == "" {
		rt.baseLogger.Printf("Username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, rqst.SenderUsername)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rqst.Content = r.FormValue("content")
	rqst.SenderUsername = r.FormValue("senderUsername")

	_, _ = rt.db.DoesUsersOwnConversation(rqst.SenderUsername, rqst.ReciverUsername) /// tutaj trzeba edytowac ta metode bazy tak aby zwracalo mi id/-1 a nie boola

	// + sprawdzic autoryzacje
	// pobrac wiadomosc z body
	// jesli to pierwsza wiadomosc w konwersacji, stworzyc konwersacje
	// wyslac dane do db

	// zwrocic odpowiedz
}
