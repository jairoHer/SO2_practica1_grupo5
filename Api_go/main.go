package main
//go get -u github.com/gorilla/mux
import (
 "fmt"
 "log"
 "net/http"
 "encoding/json"
 "github.com/gorilla/mux"
 "io/ioutil"
 "strings"
)

type Memoria struct {
  Total_RAM string `json:"total"`
  RAM_consumida string `json:"consumida"`
  Porcentaje string `json:"porcentaje"`
}

type procesos struct {
	Total_proc string `json:"total"`
	Ejecucion string `json:"ejecucion"`
	Suspendidos string `json:"suspendidos"`
	Detenidos string `json:"detenidos"`
	Zombies string `json:"zombies"`
}

type procesoBorrado struct {
	Proceso string `json:"proceso"`
}

func memoria_info(w http.ResponseWriter, r *http.Request){
	ar , err := ioutil.ReadFile("/proc/m_grupo5")
	if err != nil {
		fmt.Println(err)
	}
	entrada := string(ar)
	entradaLimpia := strings.Split(entrada, " ")
	Ram_tot := entradaLimpia[0]
	ram_lib := entradaLimpia[1]
	porcen := strings.TrimSuffix(entradaLimpia[2], "\n")
	var xd Memoria
	xd = Memoria{Total_RAM: Ram_tot, RAM_consumida: ram_lib, Porcentaje: porcen}
	fmt.Println("Endpoint de los valores del modulo")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(xd)
	fmt.Println(entrada)
	fmt.Printf("%q\n", strings.Split(entrada, " "))
}

func killProcess(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "POST para eliminar proceso")
	w.Header().Set("Content-Type", "application/json")
	var proce procesoBorrado
	_ = json.NewDecoder(r.Body).Decode(&proce)
	fmt.Println(proce.Proceso)
	json.NewEncoder(w).Encode(proce)
	//w.Write([]byte(`{"message": "post called"}`))
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Pagina del endpoint inicial")

}

func handleRequests(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/",homePage)
	myRouter.HandleFunc("/info",memoria_info).Methods("GET")
	myRouter.HandleFunc("/kill",killProcess).Methods("POST")
	log.Fatal(http.ListenAndServe(":5050",myRouter))
}

func main(){
	handleRequests()

}


