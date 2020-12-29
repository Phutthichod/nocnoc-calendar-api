package mainOld

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var srv *calendar.Service
var code = ""

func test() {

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	r := gin.Default()
	cf := cors.DefaultConfig()
	cf.AllowAllOrigins = true
	cf.AllowCredentials = true
	r.Use(cors.New(cf))
	r.GET("/", func(c *gin.Context) {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		c.JSON(200, gin.H{
			"url": authURL,
		})
	})
	r.POST("/listening", func(c *gin.Context) {

		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		type Token struct {
			Data *oauth2.Token `json:"token"`
		}
		var token Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		print("refresh :", token.Data.RefreshToken)
		client := getClientByToken(config, token.Data)
		srv, err = calendar.New(client)
		// newEvent := calendar.Event{
		// 	Summary: "Testevent",
		// 	Start:   &calendar.EventDateTime{DateTime: time.Date(2020, 12, 24, 18, 24, 0, 0, time.UTC).Format(time.RFC3339)},
		// 	End:     &calendar.EventDateTime{DateTime: time.Date(2020, 12, 24, 20, 24, 0, 0, time.UTC).Format(time.RFC3339)},
		// }
		// _, err = srv.Events.Watch("primary").Do()
		// if err != nil {
		// 	fmt.Print(err)
		// 	// return
		// }
		// c.JSON(200, gin.H{
		// 	"status": "ok",
		// })
	})
	r.POST("/insertEvent", func(c *gin.Context) {

		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		type Token struct {
			Data *oauth2.Token `json:"token"`
		}
		var token Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		print("refresh :", token.Data.RefreshToken)
		client := getClientByToken(config, token.Data)
		srv, err = calendar.New(client)
		newEvent := calendar.Event{
			Summary: "Testevent",
			Start:   &calendar.EventDateTime{DateTime: time.Date(2020, 12, 22, 18, 24, 0, 0, time.UTC).Format(time.RFC3339)},
			End:     &calendar.EventDateTime{DateTime: time.Date(2020, 12, 22, 20, 24, 0, 0, time.UTC).Format(time.RFC3339)},
		}
		_, err = srv.Events.Insert("primary", &newEvent).Do()
		if err != nil {
			fmt.Print(err)
			// return
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.POST("/token", func(c *gin.Context) {
		type CodeStruct struct {
			Code string `json:"code"`
		}
		var code CodeStruct
		err := c.ShouldBindJSON(&code)
		if err != nil {
			log.Println("error AddUserHandler", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		tok := getTokenFormCode(config, code.Code)
		c.JSON(200, gin.H{
			"token": tok,
		})
	})

	r.POST("/eventFromToken", func(c *gin.Context) {
		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		type Token struct {
			Data *oauth2.Token `json:"token"`
		}
		var token Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		print("refresh :", token.Data.RefreshToken)
		client := getClientByToken(config, token.Data)
		srv, err = calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		fmt.Println("Upcoming events:")
		type EventList struct {
			Name interface{} `json:"name"`
			Date interface{} `json:"date"`
		}

		var listEvent []EventList
		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found.")
		} else {
			print("loop")
			for _, item := range events.Items {
				date := item.Start.DateTime
				// print(item.Summary)
				if date == "" {
					date = item.Start.Date
				}
				listEvent = append(listEvent, EventList{
					Name: item.Summary,
					Date: date,
				})
				// print(listEvent)
				// fmt.Printf("%v (%v)\n", item.Summary, date)
			}
		}
		print(listEvent)
		c.JSON(200, gin.H{
			"list": listEvent,
		})
	})
	r.POST("/event", func(c *gin.Context) {
		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		type CodeStruct struct {
			Code string `json:"code"`
		}
		var code CodeStruct
		err = c.ShouldBindJSON(&code)
		print(code.Code)
		client := getClient(config, code.Code)
		srv, err = calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		fmt.Println("Upcoming events:")
		type EventList struct {
			Name interface{} `json:"name"`
			Date interface{} `json:"date"`
		}

		var listEvent []EventList
		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found.")
		} else {
			print("loop")
			for _, item := range events.Items {
				date := item.Start.DateTime
				// print(item.Summary)
				if date == "" {
					date = item.Start.Date
				}
				listEvent = append(listEvent, EventList{
					Name: item.Summary,
					Date: date,
				})
				// print(listEvent)
				// fmt.Printf("%v (%v)\n", item.Summary, date)
			}
		}
		print(listEvent)
		c.JSON(200, gin.H{
			"list": listEvent,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getClient(config *oauth2.Config, code string) *http.Client {

	// // tokFile := "token.json"
	// tok := getTokenFormCode(config, code)
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFormCode(config, code)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}
func getClientByToken(config *oauth2.Config, token *oauth2.Token) *http.Client {
	return config.Client(context.Background(), token)
}
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFormCode(config *oauth2.Config, code string) *oauth2.Token {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
