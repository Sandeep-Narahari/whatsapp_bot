package main

import (
	"github.com/AutonomyNetwork/whatsapp_bot/modules/whatsapp"
)

func main() {
	// whatsAppClient := make(chan *whatsmeow.Client)

	whatsapp.NewWhatsAppModule()
	// fmt.Println("Server started successfully")
	// whatsapp.ConnectWhatsAppClient(whatsAppClient)
}
