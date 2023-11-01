package router

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func Routes(app *gin.Engine) {
	app.GET("/", homePage)

	app.GET("/subscribe", handlerSubscribe)

	app.GET("/subscription", subscriptionPage)

	app.GET("/send-notification", sendNotification)
}

func homePage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Main website",
	})
}

func handlerSubscribe(c *gin.Context) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	randStr := hex.EncodeToString(b)
	c.Set("rand", randStr)

	successURL := "http://localhost:8080/subscription?rand=" + randStr
	failureURL := "http://localhost:8080/subscription?failure=1"

	subscriptionToken := "subscription-token"

	c.Redirect(http.StatusSeeOther, "https://pushover.net/subscribe/"+subscriptionToken+
		"?success="+successURL+
		"&failure="+failureURL)
}

func subscriptionPage(c *gin.Context) {
	success := c.Query("rand")
	failure := c.Query("failure")
	var message string

	if success != "" {
		message = "Sucesso ao se inscrever"
	}
	if failure != "" {
		message = "Erro ao se inscrever"
	}

	c.HTML(200, "subscription.html", gin.H{
		"title":   "Subscription page",
		"message": message,
	})
}

func sendNotification(c *gin.Context) {
	data := url.Values{
		"token":   {"api token"},
		"user":    {"user token"},
		"title":   {"Teste envio de notificação"},
		"message": {"Esta é uma mensagem de teste"},
		"url":     {"https://www.google.com"},
	}

	req, err := http.NewRequest("POST", "https://api.pushover.net/1/messages.json", strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	defer resp.Body.Close()

	c.HTML(http.StatusOK, "subscription.html", gin.H{
		"title": "Subscription page",
		"resp":  fmt.Sprintf("Notificação enviada com sucesso! %v", resp),
	})
}
