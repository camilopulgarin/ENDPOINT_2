package main

//Recomendaciones
//Para el correcto funcionamiento de la api debe ser instalada la dependencia gorilla/mux y mysql respectivamente
// go get -u github.com/gorilla/mux
// github.com/go-sql-driver/mysql"
//La api corre en el puerto 8080
// Puede acceder a la solucion en la siguiente direccion http://localhost:8080/user Para peticiones Post y get sin parametros
// La estructura de las peticiones Post y Put se encuentran el documento .pdf anexado


import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID_user     int    `json:"id_user"`
	Nombre_User string `json:"nombre_user"`
	Ducument     string `json:"ducument"`
	Telefono         string `json:"Telefono"`
	Sexo           int    `json:"sexo"`
	Tipo_Doc_ID  string `json:"tipo_doc_id"`
	Rol_User    string `json:"rol_user"`
}

type Users []User

var usuarios = Users{}

func BD() *sql.DB {
	bd, err := sql.Open("mysql", "usuarios_jikko:kmilo18200@tcp(db4free.net:3306)/usuarios_jikko")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("BD Conectada")
	return bd
}


func GET(w http.ResponseWriter, r *http.Request) {

	bd := BD()
	consulta, err := bd.Query(
		" SELECT ID_user,Nombre_User,Telefono,Tipo_Documento.Tipo_Documento Tipo_Documento,Ducument,Sexo,Tipo_Rol.Tipo_Rol Tipo_Rol "+
			" FROM Usuarios ,Tipo_Rol , Tipo_Documento "+
			" WHERE Tipo_Doc_ID = ID_Documento "+
			" AND Rol_User = ID_Rol ")

	if err != nil {
		panic(err.Error())
	}
	var usuarios User
	var users Users

	for consulta.Next() {
		err = consulta.Scan(&usuarios.ID_user , &usuarios.Nombre_User, &usuarios.Telefono, &usuarios.Tipo_Doc_ID , &usuarios.Ducument, &usuarios.Sexo, &usuarios.Rol_User )

		if err != nil {
			panic(err.Error())
		}
		users = append(users, usuarios)

	}

	w.Header().Set("Contet-Type", "application/json")
	j, err := json.Marshal(users)
	if err != nil {
		panic("error")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	defer consulta.Close()
	defer bd.Close()


}

func GET_ID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	bd := BD()
	consulta, err := bd.Query(
		" SELECT ID_user,Nombre_User,Telefono,Tipo_Documento.Tipo_Documento Tipo_Documento,Ducument,Sexo,Tipo_Rol.Tipo_Rol Tipo_Rol "+
			" FROM Usuarios ,Tipo_Rol , Tipo_Documento "+
			" WHERE Tipo_Doc_ID = ID_Documento "+
			" AND Rol_User = ID_Rol "+
			" AND ID_User = ? ", id)

	if err != nil {
		panic(err.Error())
	}
	var usuario User

	for consulta.Next() {
		err = consulta.Scan(&usuario.ID_user , &usuario.Nombre_User, &usuario.Telefono, &usuario.Tipo_Doc_ID , &usuario.Ducument, &usuario.Sexo, &usuario.Rol_User)

		if err != nil {
			panic(err.Error())
		}

	}

	w.Header().Set("Contet-Type", "application/json")
	j, err := json.Marshal(usuario)
	if err != nil {
		panic("error")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	defer consulta.Close()
	defer bd.Close()

}

func POST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entra al post")
	bd := BD()
	var usuario User


	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &usuario)

	consult_document, err := bd.Query("SELECT ID_Documento From  Tipo_Documento WHERE Tipo_Documento  = ?", usuario.Tipo_Doc_ID)
	consult_document.Next()
	consult_document.Scan(&usuario.Tipo_Doc_ID)
	fmt.Println(usuario.Tipo_Doc_ID)

	consult_rol, err := bd.Query("SELECT ID_Rol From  Tipo_Rol WHERE Tipo_Rol = ?", usuario.Rol_User)
	consult_rol.Next()
	consult_rol.Scan(&usuario.Rol_User)
	fmt.Println(usuario.Rol_User)

	consulta, err := bd.Query(` INSERT INTO Usuarios (Nombre_User , Ducument , Telefono , Sexo , Tipo_Doc_ID , Rol_User) `+
		` VALUES (?,?,?,?,?,?) `,
		usuario.Nombre_User,
		usuario.Ducument,
		usuario.Telefono,
		usuario.Sexo,
		usuario.Tipo_Doc_ID,
		usuario.Rol_User)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err.Error())
	}

	w.Header().Set("Contet-Type", "application/json")
	j, err := json.Marshal(usuario)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		panic("error")

	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	defer consulta.Close()
	defer bd.Close()

}


func PUT(w http.ResponseWriter, r *http.Request) {

	bd := BD()

	var usuario User


	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &usuario)

	consult_document, err := bd.Query("SELECT ID_Documento From  Tipo_Documento WHERE Tipo_Documento = ?", usuario.Tipo_Doc_ID)
	consult_document.Next()
	consult_document.Scan(&usuario.Tipo_Doc_ID)


	consult_rol, err := bd.Query("SELECT ID_Rol From  Tipo_Rol WHERE Tipo_Rol = ?", usuario.Rol_User)
	consult_rol.Next()
	consult_rol.Scan(&usuario.Rol_User)


	vars := mux.Vars(r)
	id := vars["id"]

	consulta, err := bd.Query(` UPDATE Usuarios `+
		` SET Nombre_User = ? , Ducument = ? , Telefono = ? , Sexo = ? , Tipo_Doc_ID = ? , Rol_User = ? `+
		` WHERE ID_User = ? `,
		usuario.Nombre_User,
		usuario.Ducument,
		usuario.Telefono,
		usuario.Sexo,
		usuario.Tipo_Doc_ID,
		usuario.Rol_User,
		id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error")
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Se actualizó el correctamente el usuario: ", id)

	defer consulta.Close()
	defer bd.Close()

}
func DELETE(w http.ResponseWriter, r *http.Request) {

	bd := BD()
	vars := mux.Vars(r)
	id := vars["id"]

	consulta, err := bd.Query(` DELETE FROM Usuarios WHERE ID_User = ? `, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error")
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "El usuario se elimino con exito ")

	defer consulta.Close()
	defer bd.Close()

}
func main() {
	router := mux.NewRouter().StrictSlash(false) // hace las rutas diferentes entre si con slash /
	router.HandleFunc("/user", GET).Methods("GET")
	router.HandleFunc("/user/{id}", GET_ID).Methods("GET")
	router.HandleFunc("/user", POST).Methods("POST")
	router.HandleFunc("/user/{id}", PUT).Methods("PUT")
	router.HandleFunc("/user/{id}", DELETE).Methods("DELETE")

	//Configuracion del servidore
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening....")
	log.Fatal(server.ListenAndServe())
}
