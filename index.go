package main

import (
	"api-calendar/googleAuth"
	"api-calendar/googleCalendarAPI"
	"api-calendar/model"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	r := gin.Default()
	cf := cors.DefaultConfig()
	cf.AllowAllOrigins = true
	cf.AllowCredentials = true
	r.Use(cors.New(cf))

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	ctx := context.Background()
	bCredentialsAdmin, err := ioutil.ReadFile("credentials3.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(bCredentialsAdmin, calendar.CalendarScope)
	// print(conf.Email)
	ts := conf.TokenSource(ctx)
	log.Println(ts)
	conf.Subject = "pin2041to@pintest.page"

	ts = conf.TokenSource(ctx)
	log.Println(ts)
	// conf.Email = "pin2041to@pintest.page"
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// googleAuthAdmin := googleAuth.GoogleAuth{
	// 	Config: configCredentialsAdmin,
	// }

	googleAuth := googleAuth.GoogleAuth{
		Config: config,
	}

	rangeTime := []model.OptionTime{
		model.OptionTime{
			Id:    1,
			Start: "09.00",
			End:   "12.00",
		}, model.OptionTime{
			Id:    2,
			Start: "12.00",
			End:   "15.00",
		}, model.OptionTime{
			Id:    3,
			Start: "15.00",
			End:   "18.00",
		}, model.OptionTime{
			Id:    4,
			Start: "18.00",
			End:   "21.00",
		},
	}

	r.GET("/", func(c *gin.Context) {
		authURL := googleAuth.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		c.JSON(200, gin.H{
			"url": authURL,
		})
	})

	r.POST("/token", func(c *gin.Context) {
		var code model.CodeStruct
		err := c.ShouldBindJSON(&code)
		if err != nil {
			log.Println("error AddUserHandler", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		tok := googleAuth.GetTokenFormCode(code.Code)
		c.JSON(200, gin.H{
			"token": tok,
		})
	})
	r.POST("/sendEventToInstaller", func(c *gin.Context) {
		var event model.Event2
		err = c.ShouldBindJSON(&event)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := conf.Client(oauth2.NoContext)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		newEvent := model.Event{
			Summary: event.Summary,
			Start:   calendar.EventDateTime{DateTime: event.Start.Format(time.RFC3339)},
			End:     calendar.EventDateTime{DateTime: event.End.Format(time.RFC3339)},
		}
		_, err = gService.CreatEvent(newEvent, event.CalendarId)
		if err != nil {
			fmt.Print("gService.CreatEvent", err.Error())
			// c.JSON(500, gin.H{
			// 	"message": err.Error(),
			// })
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	r.POST("/freeDayByDate", func(c *gin.Context) {
		var token model.Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(token.Data)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		res, err := gService.GetDayBusyByDate("primary", token.Date)
		if err != nil {
			log.Fatalf("Freebusy error : %v", err)
		}

		for i, timeOption := range rangeTime {
			disabled := false
			toStart, _ := strconv.Atoi(timeOption.Start[:2])
			toStop, _ := strconv.Atoi(timeOption.End[:2])
			for _, timeEvent := range res.Calendars["pin2041to@gmail.com"].Busy {
				log.Println(timeEvent)
				t, err := time.Parse(time.RFC3339, timeEvent.Start)
				t2, err := time.Parse(time.RFC3339, timeEvent.End)

				if err != nil {
					log.Print(err)
				} else {
					hrStart, _, _ := t.Clock()
					hrStop, _, _ := t2.Clock()
					// log.Println(toStart, toStop, hrStart, hrStop)
					log.Println((toStart >= hrStart))
					if !((toStart <= hrStart && toStop <= hrStart) || (toStart >= hrStop)) {
						disabled = true
					}
					// if (toStop >= hrStart) && (toStop <= hrStop) {
					// 	disabled = true
					// }
				}

			}
			rangeTime[i].Disabled = disabled
			log.Println(rangeTime[i])
		}
		c.JSON(200, gin.H{
			"rangeTime": rangeTime,
		})

	})
	r.POST("/freeDay", func(c *gin.Context) {
		var token model.Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(token.Data)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		res, err := gService.GetDayFee("primary")
		if err != nil {
			log.Fatalf("Freebusy error : %v", err)
		}
		var dateDisables = make(map[string]model.DateDisable)
		var busyDate = make(map[string]model.TimeOfDate)
		for _, timeEvent := range res.Calendars["pin2041to@gmail.com"].Busy {
			dateDisables[timeEvent.Start[:10]] = model.DateDisable{}
			// busyDate[timeEvent.Start[:10]].Times = timeEvent
			newObject := busyDate[timeEvent.Start[:10]]
			newObject.Times = append((newObject.Times), timeEvent)
			busyDate[timeEvent.Start[:10]] = newObject
		}
		var newDateDisables = make(map[string]model.DateDisable)
		for date := range dateDisables {
			var countBusy = 0
			for _, timeOption := range rangeTime {
				toStart, _ := strconv.Atoi(timeOption.Start[:2])
				toStop, _ := strconv.Atoi(timeOption.End[:2])
				var isCount = false
				for _, timePeriod := range busyDate[date].Times {
					t, err := time.Parse(time.RFC3339, timePeriod.Start)
					t2, err := time.Parse(time.RFC3339, timePeriod.End)
					if err != nil {
						log.Print(err)
					} else {
						hrStart, _, _ := t.Clock()
						hrStop, _, _ := t2.Clock()
						log.Println((toStart >= hrStart))
						if !((toStart <= hrStart && toStop <= hrStart) || (toStart >= hrStop)) {
							isCount = true
						}
					}
				}
				if isCount {
					countBusy += 1
				}
			}
			newDD := dateDisables[date]
			if countBusy == len(rangeTime) {
				newDD.FullDayBusy = true
				newDateDisables[date] = newDD
			}

		}
		c.JSON(200, gin.H{
			"dateDisables": newDateDisables,
		})

	})
	r.POST("/eventInstallationFromToken", func(c *gin.Context) {
		var instalationsData [3]model.InstallationData
		instalationsData[0] = model.InstallationData{
			Email: "phutthichod@gmail.com",
			Name:  "ช่าง A",
		}
		instalationsData[1] = model.InstallationData{
			Email: "pin2041to@gmail.com",
			Name:  "ช่าง B",
		}
		instalationsData[2] = model.InstallationData{
			Email: "ritdej.john@gmail.com",
			Name:  "ช่าง C",
		}
		var token model.Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(token.Data)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		for i, item := range instalationsData {
			events, err := gService.ListEvent(item.Email)
			if err != nil {
				log.Print("Unable to retrieve next ten of the user's events: %v", err)
			} else {
				if len(events.Items) == 0 {
					fmt.Println("No upcoming events found.")
				} else {
					var eventList []model.EventList
					for _, item := range events.Items {
						date := item.Start.DateTime
						if date == "" {
							date = item.Start.Date
						}
						eventList = append(eventList, model.EventList{
							Summary: item.Summary,
							Date:    date,
						})
					}
					instalationsData[i].EventList = eventList
				}
			}
		}
		c.JSON(200, gin.H{
			"data": instalationsData,
		})
	})

	r.POST("/eventFromToken", func(c *gin.Context) {

		var token model.Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(token.Data)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		events, err := gService.ListEvent("primary")
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}

		var listEvent []model.EventList
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
				listEvent = append(listEvent, model.EventList{
					Summary: item.Summary,
					Date:    date,
				})
				// print(listEvent)
				// fmt.Printf("%v (%v)\n", item.Summary, date)
			}
		}
		c.JSON(200, gin.H{
			"list": listEvent,
		})
	})

	r.POST("/insertEvent", func(c *gin.Context) {
		var token model.Token
		err = c.ShouldBindJSON(&token)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(token.Data)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		newEvent := model.Event{
			Summary: "test admin",
			Start:   calendar.EventDateTime{DateTime: time.Date(2020, 12, 24, 18, 24, 0, 0, time.UTC).Format(time.RFC3339)},
			End:     calendar.EventDateTime{DateTime: time.Date(2020, 12, 24, 20, 24, 0, 0, time.UTC).Format(time.RFC3339)},
		}
		_, err = gService.CreatEvent(newEvent, "primary")
		if err != nil {
			fmt.Print(err.Error())
			// c.JSON(500, gin.H{
			// 	"message": err.Error(),
			// })
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.POST("/insertEventFromCanlendarID", func(c *gin.Context) {
		var event model.Event2
		err = c.ShouldBindJSON(&event)
		if err != nil {
			log.Fatalf("ShouldBindJSON %v", err)
		}
		client := googleAuth.GetHtpClient(&event.Token)
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
		}
		gService := googleCalendarAPI.Calendar{
			Service: srv,
		}
		newEvent := model.Event{
			Summary: event.Summary,
			Start:   calendar.EventDateTime{DateTime: event.Start.Format(time.RFC3339)},
			End:     calendar.EventDateTime{DateTime: event.End.Format(time.RFC3339)},
		}
		_, err = gService.CreatEvent(newEvent, event.CalendarId)
		if err != nil {
			fmt.Print(err.Error())
			// c.JSON(500, gin.H{
			// 	"message": err.Error(),
			// })
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.Run()

}
