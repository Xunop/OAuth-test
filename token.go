package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

func CreateToken(c *gin.Context) {
	pgxConn, _ := pgx.Connect(context.TODO(), os.Getenv("DB_URI"))

	manager := manage.NewDefaultManager()

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	clientStore, _ := pg.NewClientStore(adapter)
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	//	username, ok := c.GetPostForm("username");
	//	if !ok {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"error": "username is empty",
	//		})
	//		return
	//	}
	//	password, ok := c.GetPostForm("password");
	//	if !ok {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"error": "password is empty",
	//		})
	//		return
	//	}
	clientID, ok := c.GetPostForm("client_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client_id is empty",
		})
		return
	}
	clientSecret, ok := c.GetPostForm("client_secret")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client_secret is empty",
		})
		return
	}

	cErr := clientStore.Create(&models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: "http://localhost",
	})
	if cErr != nil {
		log.Println("Internal error: ", cErr)
	}

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	tErr := srv.HandleTokenRequest(c.Writer, c.Request)
	if tErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tErr.Error(),
		})
		return
	}
}

func PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	var user User
	if username == "" || password == "" {
		fmt.Println("username or password is empty")
		return "", nil
	}
	if err := Db.Where("username = ? AND password = ?", username, password).First(user).Error; err != nil {
		fmt.Println("username or password is wrong")
	}
	return fmt.Sprint(user.UID), nil
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func InternalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
