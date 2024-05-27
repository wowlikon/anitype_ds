package main

//go build -ldflags "-H=windowsgui" && ds_rich.exe

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugolgst/rich-go/client"
)

var active, hidden bool
var activity *client.Activity
var block_genres = []string{"?"}
var mActive, mAuto *systray.MenuItem

const max_users = 10
const app_id = "1236340575800918056"

type BlockData struct {
	Genre string `json:"genre"`
}

type SetData struct {
	Wt       string `json:"wt"`
	Usr      string `json:"usr"`
	Text     string `json:"text"`
	Genres   string `json:"genres"`
	UsrCount int    `json:"usr_count"`
}

func NewActivity(user_count int, text string, wt_url string, usr_url string) *client.Activity {
	var openBtn *client.Button
	now := time.Now()

	if wt_url != "" {
		openBtn = &client.Button{
			Label: "Присоединиться к просмотру",
			Url:   wt_url,
		}
	} else if (usr_url != "") && strings.Contains(usr_url, "@") {
		openBtn = &client.Button{
			Label: fmt.Sprintf("Открыть профиль @%s", strings.SplitAfter(usr_url, "@")[1]),
			Url:   usr_url,
		}
	} else {
		openBtn = &client.Button{
			Label: "Перейти на сайт",
			Url:   "https://anitype.fun/",
		}
	}

	srcBtn := &client.Button{
		Label: "Исходный код Anitype DS RPC",
		Url:   "https://github.com/wowlikon/anitype_ds",
	}

	res := client.Activity{
		Details: "Сайт для просмотра аниме",

		LargeImage: "https://cdn.discordapp.com/app-icons/1236340575800918056/acae41c65c1d6977b8ca7529cddc9ecd.png",
		LargeText:  "Anitype.fun",
		SmallImage: "https://cdn0.iconfinder.com/data/icons/font-awesome-solid-vol-3/512/play-circle-1024.png",
		SmallText:  "Watch anime!",

		State: text,
		//party

		Timestamps: &client.Timestamps{
			Start: &now,
		},

		Buttons: []*client.Button{openBtn, srcBtn},
	}

	if user_count > 1 {
		res.Party = &client.Party{
			ID:         "-1",
			Players:    user_count,
			MaxPlayers: max_users,
		}
	}

	return &res
}

func main() {
	// API
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
	}))

	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"active": active,
			"hidden": hidden,
		})
	})

	r.POST("/set", func(c *gin.Context) {
		var sd SetData

		crit(c.ShouldBind(&sd))
		if sd.Text != activity.State {
			fmt.Println("Updated!")
			activity = NewActivity(sd.UsrCount, sd.Text, sd.Wt, sd.Usr)
			crit(client.SetActivity(*activity))
		}

		// Скрытие статуса
		hidden = contain(sd.Genres, block_genres)

		if active && hidden {
			disenable()
		}

		if !active && !hidden {
			enable()
		}

		fmt.Printf("Text: %s\n", sd.Text)
		fmt.Printf("WT: %s\n", sd.Wt)
		fmt.Printf("Usr: %s\n", sd.Usr)
		fmt.Printf("UsrCount: %d\n", sd.UsrCount)
		fmt.Printf("Genres: %s\n", sd.Genres)
		fmt.Printf("Hiden: %t\n", hidden)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	r.GET("/enable", func(c *gin.Context) {
		c.Writer.WriteString("enabled")
		enable()
	})

	r.GET("/disenabled", func(c *gin.Context) {
		c.Writer.WriteString("disenabled")
		disenable()
	})

	r.POST("/add_block", func(c *gin.Context) {
		var bd BlockData

		crit(c.ShouldBind(&bd))
		block_genres = append(block_genres, bd.Genre)
		c.JSON(http.StatusOK, gin.H{"genres": block_genres})
	})

	r.POST("/del_block", func(c *gin.Context) {
		var bd BlockData

		crit(c.ShouldBind(&bd))
		block_genres = remove(block_genres, bd.Genre)
		c.JSON(http.StatusOK, gin.H{"genres": block_genres})
	})

	r.GET("/get_block", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"genres": block_genres})
	})

	activity = NewActivity(0, "Ожидание...", "", "")
	go r.Run("localhost:878")    // Запуск API
	systray.Run(onReady, onExit) // Добавление в трэй
}

func crit(err error) {
	if err != nil {
		panic(err)
	}
}

func contain(genres string, arr []string) bool {
	if len(genres) == 0 {
		return false
	}

	arr1 := strings.Split(genres, ", ")
	set := make(map[string]bool)
	for _, elem := range arr1 {
		set[elem] = true
	}

	for _, elem := range arr {
		if set[elem] {
			return true
		}
	}

	return false
}

func enable() {
	crit(client.SetActivity(*activity))
	crit(client.Login(app_id))
	mActive.SetTitle("Discord RPC включен")
	mActive.Check()
	active = true
}

func disenable() {
	client.Logout()
	mActive.SetTitle("Discord RPC выключен")
	mActive.Uncheck()
	active = false
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	crit(err)
	return b
}

func remove[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func onReady() {
	systray.SetTitle("Anitype")
	systray.SetTooltip("Anitype")
	systray.SetIcon(getIcon("icon.ico"))
	mAuto = systray.AddMenuItemCheckbox("Автозапуск выключен", "Автозапуск приложения", false)
	mActive = systray.AddMenuItemCheckbox("Discord startus включен", "Turn on/off", active)
	mQuit := systray.AddMenuItem("Выйти", "Закрыть приложение")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mAuto.ClickedCh:
				if mAuto.Checked() {
					mAuto.Uncheck()
					mAuto.SetTitle("Автозапуск выключен")
				} else {
					mAuto.Check()
					mAuto.SetTitle("Автозапуск включен")
				}
			case <-mActive.ClickedCh:
				if mActive.Checked() {
					disenable()
				} else {
					enable()
				}
			}
		}
	}()

	enable()
}

func onExit() {
	// clean up here
}
