package main

//go build -ldflags "-H=windowsgui"

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugolgst/rich-go/client"
)

var active bool
var activity *client.Activity

const max_users = 7
const app_id = "1236340575800918056"

type SetData struct {
	Wt       string `json:"wt"`
	Usr      string `json:"usr"`
	UsrCount string `json:"usr_count"`
	Text     string `json:"text"`
}

func NewActivity(users string, text string, wt_url string, usr_url string) *client.Activity {
	var user_count int
	user_count, _ = strconv.Atoi(users)

	now := time.Now()
	openBtn := &client.Button{
		Label: "Open anitype",
		Url:   "https://anitype.fun/",
	}

	res := &client.Activity{
		Details: "Anitype - сайт для просмотра аниме",

		LargeImage: "https://cdn.discordapp.com/app-icons/1236340575800918056/acae41c65c1d6977b8ca7529cddc9ecd.png",
		LargeText:  "Anitype.fun",
		SmallImage: "https://cdn0.iconfinder.com/data/icons/font-awesome-solid-vol-3/512/play-circle-1024.png",
		SmallText:  "Watch anime!",

		State: text,
		//party

		Timestamps: &client.Timestamps{
			Start: &now,
		},

		Buttons: []*client.Button{openBtn},
	}

	if user_count > 1 {
		res.Party = &client.Party{
			ID:         "-1",
			Players:    user_count,
			MaxPlayers: max_users,
		}
	}

	if wt_url != "" {
		res.Buttons = append(res.Buttons, &client.Button{
			Label: "Watch together",
			Url:   wt_url,
		})
	}

	if (usr_url != "") && strings.Contains(usr_url, "@") {
		res.Buttons = append(res.Buttons, &client.Button{
			Label: fmt.Sprintf("Open %s's profile", strings.SplitAfter(usr_url, "@")[1]),
			Url:   usr_url,
		})
	}

	return res
}

func main() {
	crit(client.Login(app_id))

	activity = NewActivity("0", "Ожидание...", "", "")
	crit(client.SetActivity(*activity))
	active = true

	// API
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
	}))

	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": active,
		})
	})

	r.GET("/enable", func(c *gin.Context) {
		crit(client.SetActivity(*activity))
		c.Writer.WriteString("enabled")
		client.Login(app_id)
		active = true
	})

	r.GET("/disenabled", func(c *gin.Context) {
		c.Writer.WriteString("disenabled")
		client.Logout()
		active = false
	})

	r.POST("/set", func(c *gin.Context) {
		var sd SetData

		crit(c.ShouldBind(&sd))
		activity := NewActivity(
			sd.UsrCount,
			sd.Text,
			sd.Wt,
			sd.Usr,
		)

		crit(client.SetActivity(*activity))
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	go r.Run("localhost:878")    // Run API in goroutine
	systray.Run(onReady, onExit) // Tray icon
}

func crit(err error) {
	if err != nil {
		panic(err)
	}
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	crit(err)
	return b
}

func onReady() {
	systray.SetTitle("Anitype")
	systray.SetTooltip("Anitype")
	systray.SetIcon(getIcon("icon.ico"))
	mChecked := systray.AddMenuItemCheckbox("Autostart off", "Set autostart", false)
	mQuit := systray.AddMenuItem("Quit", "Close this app")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Autostart off")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Autostart on")
				}
			}
		}
	}()
}

func onExit() {
	// clean up here
}
