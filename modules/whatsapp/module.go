package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"go.mau.fi/whatsmeow"

	"github.com/AutonomyNetwork/whatsapp_bot/database"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Module struct {
	sheetsService  *sheets.Service
	whatsAppClient *whatsmeow.Client
	db             *sql.DB
}

func NewWhatsAppModule() (*Module, error) {
	whatsAppClient := make(chan *whatsmeow.Client)

	// go connectSheets()
	sheetsService, err := connectSheets()
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	go ConnectWhatsAppClient(whatsAppClient)
	database, err := database.Connect()
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	router := mux.NewRouter()
	module := &Module{
		sheetsService:  sheetsService,
		whatsAppClient: <-whatsAppClient,
		db:             database,
	}
	go handleRequests(router, module)
	select {}
	// fmt.Println(module)
	// fmt.Println("AFTERRRRRRRRRRRRRR")
}

func ConnectWhatsAppClient(ch chan *whatsmeow.Client) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("postgres", "postgres://sandeep:admin123@localhost/whatsapp?sslmode=disable", dbLog)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	// client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
	}
	ch <- client
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()

}
func connectSheets() (*sheets.Service, error) {
	ctx := context.Background()
	credsFile := "/home/sandeep/go/src/github.com/AutonomyNetwork/whatsapp_bot/newCred.json"
	creds, err := ioutil.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
		return nil, err
	}
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}
	return srv, nil
}

func handleRequests(router *mux.Router, module *Module) {
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

	router.HandleFunc("/message", module.SendMessageHandler).Methods("POST")
	go func() {
		if err := http.ListenAndServe(":9996", router); err != nil {
			log.Fatal(err)
			return
		}
	}()
}
