package controllers

import (
	"encoding/json"
	"fmt"
	"go-api-example/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
type UserController struct{
	//struct yapisi tanimliyoruz
	session *mgo.Session
}

//yeni bir userController olusturma fonksiyonu *UserController pointer tipinde return etmek zorunda
func NewUserController(s *mgo.Session) *UserController{
return &UserController{s}
}

//istenilen id ye gore user getirme fonksiyonu
//go da functionlarin stuctur i olabiliyor burda UserControlleri ekleyip function icerisinde kullandik

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	//httprouter un params inda yer alan idmizi degiskene atadik
	id := p.ByName("id")

	//eger bson dan yani mongo dan gelen data da bir eksiklik var ise ekrana hata yazdir diyoruz
	if !bson.IsObjectIdHex(id){
		w.WriteHeader(http.StatusNotFound)
	}
	//bson dan gelen datadaki istenilen id ye gore ObjectId dondertiyoruz
	oid := bson.ObjectIdHex(id)

	//models ten gelen user modelimizi u degiskenine atiyoruz
	u := models.User{}

	//eger benim usercontrollerimin session inin DB ye baglarsam ve bu benim mongodaki db adim ile ayniysa ve collection adim users ise burdaki objectid querryden bir tanesini execute et diyoruz.
	// eger tersiyse 404 hatasi atiyoruz
  if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u); err != nil{
	w.WriteHeader(404)
	return
   }
   //burda model userden gelen ve icerisine verilen datayi json a convert ediyoruz
   uj, err :=json.Marshal(u)
   if err!= nil{
	   fmt.Println(err)
   }
//response writer in headerina  content type olarak application json u set ediyorum
w.Header().Set("Content-Type", "application/json")

//header a status ok bastiriyorum
w.WriteHeader(http.StatusOK)
//cikti verdirtiyorum
fmt.Printf("UserGet function worked.\n")
fmt.Fprintf(w, "%s\n", uj)

}

//user olusturma functionu
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	u := models.User{}
	//request ile gelen data nin body sini u ya gore decode ediyoruz
	json.NewDecoder(r.Body).Decode(&u)

	//user imin id sine bson dan gelen object id yi atadim
	u.Id = bson.NewObjectId()

	//yine ayni sekilde usercontroller ile yeni user u insert ediyorum
	uc.session.DB("mongo-golang").C("users").Insert(u)

	//user u convert ediyorum
	uj, err := json.Marshal(u)

	//hata alirsam convertte bastiriyorum
	if err != nil{
		fmt.Println(err)
	}

	//ayni sekilde set edip satatus bastiriyorum
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("CreateUser function worked.\n")
	fmt.Fprintf(w, "%s\n", uj)
}

//kullanici silme function u
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params){

	//atama yapildi
	id := p.ByName("id")

	//id yoksa hata bas
	if !bson.IsObjectIdHex(id){
		w.WriteHeader(404)
		return
	}
	//object id ata
	oid := bson.ObjectIdHex(id)

	//idyi sil tablodan ve collectiondan eger varsa hata bas
	if err := uc.session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
	}
	//sorunsuz biterse statusok bas silinen user bilgisi bas
	w.WriteHeader(http.StatusOK)
	fmt.Printf("DeleteUser function worked.\n")
	fmt.Fprint(w, "Deleted user", oid, "\n")
}