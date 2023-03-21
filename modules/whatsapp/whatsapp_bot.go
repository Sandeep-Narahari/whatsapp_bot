package whatsapp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AutonomyNetwork/whatsapp_bot/logger"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (m Module) sendMessage(contactNumber string, message string) error {
	_, err := m.whatsAppClient.SendMessage(context.Background(), types.JID{
		User:   contactNumber,
		Server: types.DefaultUserServer,
	}, &waProto.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		return err
	}

	return nil
}

//	func (m Module) retrieveData() {
//		// Retrieve the values in the second column of the sheet
//		spreadsheetId := "11cs1NLrBdUq3Vamg_ORza-Zt2fGiXV7oZkF79fyuAE4"
//		rangeName := "Sheet1!A:B"
//		resp, err := m.sheetsService.Spreadsheets.Values.Get(spreadsheetId, rangeName).Do()
//		if err != nil {
//			log.Fatalf("Unable to retrieve data from sheet: %v", err)
//		}
//		if len(resp.Values) > 1 {
//			resp.Values = resp.Values[1:]
//		} else {
//			fmt.Println("No data found")
//			return
//		}
//		// Extract the values from the response
//		for _, row := range resp.Values {
//			contactNumber := row[0].(string)
//			message := row[1].(string)
//			err = m.sendMessage(contactNumber, message)
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//		// Print the retrieved column
//	}
func (m Module) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	status := logger.StatusWriter{
		ResponseWriter: w,
		Status:         0,
	}
	spreadsheetId := "11cs1NLrBdUq3Vamg_ORza-Zt2fGiXV7oZkF79fyuAE4"
	rangeName := "Sheet1!A:B"

	resp, err := m.sheetsService.Spreadsheets.Values.Get(spreadsheetId, rangeName).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		bz := logger.ProcessResponseBody(false, err.Error())
		status.WriteHeader(http.StatusInternalServerError)
		_, _ = status.Write(bz)
		return
	}
	if len(resp.Values) > 1 {
		resp.Values = resp.Values[1:]
	} else {
		fmt.Println("No data found")
		msg := "No data found in spread sheets"
		bz := logger.ProcessResponseBody(false, msg)
		status.WriteHeader(http.StatusInternalServerError)
		_, _ = status.Write(bz)
		return
	}
	// Extract the values from the response
	for _, row := range resp.Values {
		contactNumber := row[0].(string)
		message := row[1].(string)
		_, err := m.whatsAppClient.SendMessage(context.Background(), types.JID{
			User:   contactNumber,
			Server: types.DefaultUserServer,
		}, &waProto.Message{
			Conversation: proto.String(message),
		})
		if err != nil {
			bz := logger.ProcessResponseBody(false, err.Error())
			status.WriteHeader(http.StatusInternalServerError)
			_, _ = status.Write(bz)
			return
		}
	}
}
