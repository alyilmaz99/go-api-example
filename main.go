// main paketimizi ekliyoruz go ile gelen ve bir cok ozellige eristigimiz paket
package main

import (
	"fmt"
	"go-api-example/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {
	//print atiyoruz
	fmt.Printf("Starting Server...\n")

	//httprouter olusturuyoruz
	fmt.Printf("Starting httprouter...\n")
	routes := httprouter.New()
	//kontrollerimizdan gelen NewUserController func cagiriyoruz
	fmt.Printf("Starting UserController...\n")
	uc := controllers.NewUserController(getSession())
	fmt.Printf("Routes Created.\n")
	routes.GET("/user/:id", uc.GetUser)
	routes.POST("/user", uc.CreateUser)
	routes.DELETE("/user/:id", uc.DeleteUser)
	
	//dinleyecegimiz portu verdik ve hangi routelari dinleyecegini soyledik
	fmt.Printf("Server Listening on 9000 port...\n")
	http.ListenAndServe("localhost:9000", routes)


}

func getSession() *mgo.Session{
	//burda connect ediyoruz
	session, err := mgo.Dial("mongodb://localhost:27017")
	//her hangi bir sorun varsa err atiyoruz
	if err != nil{
		panic(err)
	}
	return session
} 