package main
//go get -u github.com/gorilla/mux
import (
 "fmt"
 "log"
 "net/http"
 "encoding/json"
 "github.com/gorilla/mux"
 "io/ioutil"
 "os/exec"
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

func leerResumenProcesos(w http.ResponseWriter, r *http.Request){
	ar , err := ioutil.ReadFile("/proc/proc_procesos")
	if err != nil {
		fmt.Println(err)
	}
	entrada := string(ar)
	entradaLimpia := strings.Split(entrada, "-------------------------------------\n")
	resumen := strings.Split(entradaLimpia[1],"\n")
	//valores :=  obtenerResumen(resumen)
	fmt.Println(resumen)
	fmt.Println(entradaLimpia[1])
	/*fmt.Println(valores[0])
	fmt.Println(valores[1])
	fmt.Println(valores[2])
	fmt.Println(valores[3])
	fmt.Println(valores[4])
	fmt.Println(valores[5])*/
}

func obtenerResumen(entrada[] string) []string{
	var valores[]string
	valores[0] = strings.Split(entrada[0],": ")[1]//Total
	valores[1] = strings.Split(entrada[1],": ")[1]//IDLE
	valores[2] = strings.Split(entrada[2],": ")[1]//RUNNING 
	valores[3] = strings.Split(entrada[3],": ")[1]//SLEEP
	valores[4] = strings.Split(entrada[4],": ")[1]//STOPED 
	valores[5] = strings.Split(entrada[5],": ")[1]//ZOMBIE 
	return valores
}

func killProcess(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "POST para eliminar proceso")
	w.Header().Set("Content-Type", "application/json")
	var proce procesoBorrado
	_ = json.NewDecoder(r.Body).Decode(&proce)
	PID := string(proce.Proceso)
	fmt.Println(PID)
	json.NewEncoder(w).Encode(proce)
	//_, err := exec.Command("sh", "-c", "echo '"+ sudopassword +"' | sudo -S pkill -SIGINT my_app_name").Output()
	_, err := exec.Command("sh", "-c", "kill "+PID).Output()
	if err != nil {
		fmt.Println("No se pudo eliminar el proceso")
	} else {
		fmt.Println("Se elimino el proceso con PID: "+PID)
	}
	/*out, err := exec.Command("ls").Output()
	if err != nil{
		fmt.Println("No se pudo eliminar el proceso")
	}
	output := string(out[:])
	fmt.Println(output)*/
	//w.Write([]byte(`{"message": "post called"}`))
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Pagina del endpoint inicial")

}

func handleRequests(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/",homePage)
	myRouter.HandleFunc("/info",memoria_info).Methods("GET")
	myRouter.HandleFunc("/resumen_proc",leerResumenProcesos).Methods("GET")
	myRouter.HandleFunc("/kill",killProcess).Methods("POST")
	log.Fatal(http.ListenAndServe(":5050",myRouter))
}

func main(){
	handleRequests()

}


