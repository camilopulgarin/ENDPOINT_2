package main

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

	fmt.Println("Comienzo")
	bd, err := sql.Open("mysql", "usuarios_jikko:kmilo18200@tcp(db4free.net:3306)/usuarios_jikko")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Conectado a BD")
	return bd
}


func f_get(w http.ResponseWriter, r *http.Request) {

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
		//fmt.Println(usuarios)
	}
	//fmt.Println(users)
	//escribir
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

func f_get_id(w http.ResponseWriter, r *http.Request) {

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
	//fmt.Println(users)
	//escribir
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

func f_post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entra al post")
	bd := BD()
	var usuario User

	//err := json.NewDecoder(r.Body).Decode(&l_in)
	reqBody, err := ioutil.ReadAll(r.Body) //leemos tooodo
	json.Unmarshal(reqBody, &usuario)

	consulta_documento, err := bd.Query("SELECT ID_Documento From  Tipo_Documento WHERE Tipo_Documento  = ?", usuario.Tipo_Doc_ID)
	consulta_documento.Next()
	consulta_documento.Scan(&usuario.Tipo_Doc_ID)
	fmt.Println(usuario.Tipo_Doc_ID)

	consulta_nacionalidad, err := bd.Query("SELECT ID_Rol From  Tipo_Rol WHERE Tipo_Rol = ?", usuario.Rol_User)
	consulta_nacionalidad.Next()
	consulta_nacionalidad.Scan(&usuario.Rol_User)
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
		fmt.Fprintf(w, "ERRORRRR")
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


func f_put(w http.ResponseWriter, r *http.Request) {

	bd := BD()

	var usuario User

	//err := json.NewDecoder(r.Body).Decode(&l_in)
	reqBody, err := ioutil.ReadAll(r.Body) //leemos tooodo
	json.Unmarshal(reqBody, &usuario)

	consulta_documento, err := bd.Query("SELECT ID_Documento From  Tipo_Documento WHERE Tipo_Documento = ?", usuario.Tipo_Doc_ID)
	consulta_documento.Next()
	consulta_documento.Scan(&usuario.Tipo_Doc_ID)
	//fmt.Println(usuario.Tipo_Doc_ID)

	consulta_nacionalidad, err := bd.Query("SELECT ID_Rol From  Tipo_Rol WHERE Tipo_Rol = ?", usuario.Rol_User)
	consulta_nacionalidad.Next()
	consulta_nacionalidad.Scan(&usuario.Rol_User)
	//fmt.Println(usuario.Rol_User)

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
		fmt.Fprintf(w, "ERRORRRR")
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ACTUALIZACION REALIZADA CORRECTAMENTE A ID ", id)

	defer consulta.Close()
	defer bd.Close()

}
func f_delete(w http.ResponseWriter, r *http.Request) {

	bd := BD()
	fmt.Println("DELETE")
	vars := mux.Vars(r)
	id := vars["id"]

	consulta, err := bd.Query(` DELETE FROM Usuarios WHERE ID_User = ? `, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERRORRRR")
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "EL REGISTRO SE HA BORRADO CORRECTAMENTE ")

	defer consulta.Close()
	defer bd.Close()

}
func main() {
	router := mux.NewRouter().StrictSlash(false) // hace las rutas diferentes entre si con slash /
	router.HandleFunc("/user", f_get).Methods("GET")
	router.HandleFunc("/user/{id}", f_get_id).Methods("GET")
	router.HandleFunc("/user", f_post).Methods("POST")
	router.HandleFunc("/user/{id}", f_put).Methods("PUT")
	router.HandleFunc("/user/{id}", f_delete).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",          // puerto
		Handler:        router,           //
		ReadTimeout:    20 * time.Second, // tiempo de lectura
		WriteTimeout:   20 * time.Second, // tiempo de escritura
		MaxHeaderBytes: 1 << 20,          // 1mega en bits
	}
	log.Println("Listening....")
	log.Fatal(server.ListenAndServe())
}
