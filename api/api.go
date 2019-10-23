// Package api handles routing to the tweerchivist functionalities.
package api

import (
	"net/http"

	"github.com/ghostbar/tweerchivist/archiver"
	"github.com/ghostbar/tweerchivist/retriever"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// API represents the API main object from which to manage the API main
// routines.
type API struct {
	arc    *archiver.Archiver
	ret    *retriever.Retriever
	router *gin.Engine
}

// New returns an initialized API pointer, ready to be runned or served.
func New(tcli *twitter.Client, mcli *mongo.Client) *API {
	arc := archiver.New(mcli)
	ret := retriever.New(tcli)
	a := &API{arc: arc, ret: ret}
	a.router = gin.Default()

	a.declareRoutes()

	return a
}

func (a *API) declareRoutes() {
	a.router.GET("/archive/:username", a.getArchive)
	a.router.GET("/archive", a.getListOfUsers)
	a.router.POST("/archive", a.addUserToArchive)
}

// Run is a proxy for running Run against the GinGonic router. This is the
// start button. You should send here the address argument, like ":8080" for
// localhost:8080.
func (a *API) Run(port string) {
	a.router.Run(port)
}

// ServeHTTP is helpful while testing the API and is a proxy for
// gin.Engine.ServeHTTP.
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *API) getArchive(c *gin.Context) {
	status := http.StatusOK
	c.JSON(status, gin.H{"status": "success"})
}

func (a *API) getListOfUsers(c *gin.Context) {
}

func (a *API) addUserToArchive(c *gin.Context) {
}
