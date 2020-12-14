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
	Idle string `json:"idle"`
	Running string `json:"running"`
	Sleep string `json:"sleep"`
	Stoped string `json:"stoped"`
	Zombie string `json:"zombie"`
}

type procesoBorrado struct {
	Proceso string `json:"proceso"`
}

type todo struct {
	Procesos []todos_procesos `json:"procesos"`
}

type todos_procesos struct {
	Pid string `json:"pid"`
	Process string `json:"process"`
	Uid string `json:"uid"`
	State string `json:"state"`
	Childs []child `json:"childs"`

}

type child struct {
	Pid string `json:"pid"`
	Process string `json:"process"`
	Uid string `json:"uid"`
	State string `json:"state"`
}

type todoUnificado struct {
	Procesos []todo `json:"procesos"`
	Total_proc string `json:"total_procesos"`
	Idle string `json:"idle"`
	Running string `json:"running"`
	Sleep string `json:"sleep"`
	Stoped string `json:"stoped"`
	Zombie string `json:"zombie"`
}
func toditotodito(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	proc_data , err2 := ioutil.ReadFile("/proc/proc_procesos")
	if err2 != nil {
		fmt.Println(err2)
	}
	entradaProc := string(proc_data)
	entradaLimpiap := strings.Split(entradaProc, "-------------------------------------\n")
	resumen := strings.Split(entradaLimpiap[1],"\n")
	byteValue := []byte(entradaLimpiap[0])
	valores :=  obtenerResumen(resumen)
	//var procesoss todo
	//var completo todoUnificado
	out := map[string]interface{}{}
	json.Unmarshal(byteValue, &out)
	out["total_procesos"] = valores[0]
	out["idle"] = valores[1]
	out["running"] = valores[2]
	out["sleep"] = valores[3]
	out["stoped"] = valores[4]
	out["zombie"] = valores[5]
	json.NewEncoder(w).Encode(out)
	
	//completo = todoUnificado{Procesos:procesoss,Total_proc: valores[0], Idle: valores[1], Running: valores[2], Sleep: valores[3], Stoped: valores[4],Zombie:valores[5]}


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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(xd)
	fmt.Println(entrada)
	fmt.Printf("%q\n", strings.Split(entrada, " "))
}

func leerResumenProcesos(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ar , err := ioutil.ReadFile("/proc/proc_procesos")
	if err != nil {
		fmt.Println(err)
	}
	entrada := string(ar)
	entradaLimpia := strings.Split(entrada, "-------------------------------------\n")
	resumen := strings.Split(entradaLimpia[1],"\n")
	valores :=  obtenerResumen(resumen)
	var procs procesos
	procs = procesos{Total_proc: valores[0], Idle: valores[1], Running: valores[2], Sleep: valores[3], Stoped: valores[4],Zombie:valores[5]}
	fmt.Println(resumen)
	json.NewEncoder(w).Encode(procs)
}

func leerTodosLosProcesos(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ar , err := ioutil.ReadFile("/proc/proc_procesos")
	if err != nil {
		fmt.Println(err)
	}
	entrada := string(ar)
	entradaLimpia := strings.Split(entrada, "-------------------------------------\n")
	//procs := strings.Split(entradaLimpia[0],"\n")
	var procesoss todo
	byteValue := []byte(entradaLimpia[0])
	json.Unmarshal(byteValue, &procesoss)
	json.NewEncoder(w).Encode(procesoss)
	//fmt.Println(procs)
}

func obtenerResumen(entrada[] string) []string{
	var valores[]string
	//valores = make([]string, 0,10)
	valores = append(valores,strings.Split(entrada[0],": ")[1])//Total
	valores = append(valores,strings.Split(entrada[1],": ")[1])//IDLE
	valores = append(valores,strings.Split(entrada[2],": ")[1])//RUNNING 
	valores = append(valores,strings.Split(entrada[3],": ")[1])//SLEEP
	valores = append(valores,strings.Split(entrada[4],": ")[1])//STOPED 
	valores = append(valores,strings.Split(entrada[5],": ")[1])//ZOMBIE 
	return valores
}

func killProcess(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "POST para eliminar proceso")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
    	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
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
	myRouter.HandleFunc("/todo",toditotodito).Methods("GET")
	myRouter.HandleFunc("/resumen_proc",leerResumenProcesos).Methods("GET")
	myRouter.HandleFunc("/procesos",leerTodosLosProcesos).Methods("GET")
	myRouter.HandleFunc("/kill",killProcess).Methods("POST")
	log.Fatal(http.ListenAndServe(":5050",myRouter))
}

func main(){
	handleRequests()

}


