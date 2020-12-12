package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/shirou/gopsutil/cpu"
	//"os/exec"
)

var (
	index   = 0
	index2  = 0
	Times   [1000]string
	Values  [1000]string
	Times2  [1000]string
	Values2 [1000]string
)

type Strumemoria struct {
	Total   int    `json:"total"` //poner nombre de atributos en mayuscula
	Libre   int    `json:"libre"`
	Tiempos string `json:"tiempos"`
	Valores string `json:"valores"`
	Index   int    `json:"index"`
}

type Strucpu struct {
	Total   int    `json:"total"` //poner nombre de atributos en mayuscula
	Libre   int    `json:"libre"`
	Tiempos string `json:"tiempos"`
	Valores string `json:"valores"`
	Index   int    `json:"index"`
}

type CPU struct {
	Cores  int32   `json:"cores"`
	Vendor string  `json:"vendor"`
	Family string  `json:"family"`
	Model  string  `json:"model"`
	Speed  string  `json:"speed"`
	Read   float64 `json:"read"`
}

type prueba struct {
	valor int
}
//AGREGANDO LO DE LA PAGINA PRINCIPAL(PROCESOS)
type PROCESOS struct{
	Ejecucion int `json:"ejecucion"`
	Suspendidos int `json:"suspendidos"`
	Detenidos int `json:"detenidos"`
	Zombie int `json:"zombie"`
	Total int `json:"total"`
	Datos string `json:"datos"`
	Arbol string `json:"arbol"`
}

type hijo struct{
	Nombre string
	pid string
}

type Proc struct{
	Nombre string
	Hijos []hijo	
}
var procesos []Proc;

func obtenerram(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENTRE")
	/*	fmt.Fprintf(w, "Welcome to the HomePage!") //IMPRIME EN LA WEB
		fmt.Println("Endpoint Hit: homePage")*/ //IMPRIMIR EN CONSOLA AL
	b, err := ioutil.ReadFile("mockram.json")

	if err != nil { //nil es contrario a null por asi decirlo
		fmt.Println("primer error return")
		return ///si hay error se sale
	}

	str := string(b)
	listainfo := strings.Split(string(str), "\n") //split por salto de linea

	memoriatotal := strings.Replace((listainfo[0])[10:24], " ", "", -1)
	memorialibre := strings.Replace((listainfo[2])[15:24], " ", "", -1)
	fmt.Println("LA DISPONIBLE ES ", memorialibre)
	ramtotal, err1 := strconv.Atoi(memoriatotal)
	ramlibre, err2 := strconv.Atoi(memorialibre)

	if err1 == nil && err2 == nil {
		ramtotalmb := ramtotal / 1024
		ramlibremb := ramlibre / 1024

		porcentajeram := ((ramtotalmb - ramlibremb) * 100) / ramtotalmb
		//obtengo hora
		v1 := time.Now().Format("01-02-2006 15:04:05")
		parte := strings.Split(v1, " ")
		Times[index] = parte[1]
		Values[index] = strconv.Itoa(porcentajeram)

		//fmt.Println("EL VALUES ES ", Values[index])

		/*	if index == 60 {
				var Tiempos2 [1000]string
				var Valores2 [1000]string
				for j := 1; j <= 60; j++ {
					var i = 0
					Tiempos2[i] = Times[j]
					Valores2[i] = Values[j]
					i = i + 1
				}
				//index = index - 1
				Times = Tiempos2
				Values = Valores2
				fmt.Println("la posicion 1 es ", Times[0])
				fmt.Println("La posicion ultima es ", Times[60])
			}
		*/
		respuestamem := Strumemoria{Total: ramtotalmb, Libre: ramlibremb, Tiempos: strings.Join(Times[:], ","), Valores: strings.Join(Values[:], ","), Index: index}

		//	fmt.Println("O SEA SI ENTRO")
		crearj, errorjson := json.Marshal(respuestamem)
		index = index + 1
		/*if index < 60 {
			index = index + 1
		} else {
			index = 60
		}*/

		if errorjson != nil {
			fmt.Println("HAY UN ERROR")
			http.Error(w, errorjson.Error(), http.StatusInternalServerError)
			return
		}

		//c	fmt.Println("la memoria libre es ", respuestamem)
		//conver := string(crearj)
		fmt.Println("EL indice es ", index)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(crearj)
		//	convertir_a_cadena := string(jsonResponse)
		//	fmt.Println(convertir_a_cadena)
		//fmt.Println("LLEGUE AL FINAL A RETORNAR JSON", jsonResponse)
	} else {
		return
	}

}

func obtenercpu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENTRE")
	/*	fmt.Fprintf(w, "Welcome to the HomePage!") //IMPRIME EN LA WEB
		fmt.Println("Endpoint Hit: homePage")*/ //IMPRIMIR EN CONSOLA AL
	b, err := ioutil.ReadFile("/proc/meminfo")

	if err != nil { //nil es contrario a null por asi decirlo
		fmt.Println("primer error return")
		return ///si hay error se sale
	}

	str := string(b)
	listainfo := strings.Split(string(str), "\n") //split por salto de linea

	//para el cpu
	//cpuStat, err := cpu.Info()
	percentage, err := cpu.Percent(0, true)
	//fmt.Println("la info del cpu es ", cpuStat)
	fmt.Println("EL PORCENTAJE DE USO ES ", percentage[0])
	memoriatotal := strings.Replace((listainfo[0])[10:24], " ", "", -1)
	memorialibre := strings.Replace((listainfo[2])[15:24], " ", "", -1)
	//	fmt.Println("LA DISPONIBLE ES ", memorialibre)
	ramtotal, err1 := strconv.Atoi(memoriatotal)
	ramlibre, err2 := strconv.Atoi(memorialibre)

	if err1 == nil && err2 == nil {
		ramtotalmb := ramtotal / 1024
		ramlibremb := ramlibre / 1024

		//	porcentajeram := ((ramtotalmb - ramlibremb) * 100) / ramtotalmb
		//obtengo hora
		v1 := time.Now().Format("01-02-2006 15:04:05")
		parte := strings.Split(v1, " ")
		Times2[index2] = parte[1]
		Values2[index2] = strconv.FormatFloat(percentage[0], 'f', 5, 64)

		respuestamem := Strumemoria{Total: ramtotalmb, Libre: ramlibremb, Tiempos: strings.Join(Times2[:], ","), Valores: strings.Join(Values2[:], ","), Index: index2}

		//	fmt.Println("O SEA SI ENTRO")
		crearj, errorjson := json.Marshal(respuestamem)
		index2 = index2 + 1
		/*if index < 60 {
			index = index + 1
		} else {
			index = 60
		}*/

		if errorjson != nil {
			fmt.Println("HAY UN ERROR")
			http.Error(w, errorjson.Error(), http.StatusInternalServerError)
			return
		}

		//c	fmt.Println("la memoria libre es ", respuestamem)
		//conver := string(crearj)
		//	fmt.Println("EL indice es ", index)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(crearj)
		//	convertir_a_cadena := string(jsonResponse)
		//	fmt.Println(convertir_a_cadena)
		//fmt.Println("LLEGUE AL FINAL A RETORNAR JSON", jsonResponse)
	} else {
		return
	}

}

func getprocesos(w http.ResponseWriter, r *http.Request){
	fmt.Println("ENTRE A LEER PROCESOS");
	b, err := ioutil.ReadFile("procesos")
	if err != nil { //nil es contrario a null por asi decirlo
		fmt.Println("Error abriendo procesos")
		return ///si hay error se sale
	}
	
}
func obtenerprincipal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENTRE a PRINCIPAL")
	archivos, err := ioutil.ReadDir("/proc")
    if err != nil {
        log.Fatal(err)
	}
	procejecucion:=0;
	procsuspendidos:=0;
	procdetenidos:=0;
	proczombies:=0;
	contador:=0;
	textocompleto:="";
	//totalram:=0.0;
	//totalrammegas:=0.0;
	//pruebaram:=0.0;
	procesos=nil;
    for _, archivo := range archivos {
		if(archivo.Name()[0]>48 && archivo.Name()[0]<58){
			//TODOS ESTOS SON PROCESOS
			nombreArchivo:= "/proc/";
			nombreArchivo+= archivo.Name();
			nombreArchivo+="/status";
			bytesLeidos, err := ioutil.ReadFile(nombreArchivo)
			if err != nil {
				fmt.Printf("Error leyendo archivo: %v", err)
			}
			contenido := string(bytesLeidos)
			splitsalto:=strings.Split(contenido,"\n");
			hayram:=true;
			nombre:="";
			for _,salto:= range splitsalto{
				splitpuntos:=strings.Split(salto,":");
				if(splitpuntos[0]=="Name"){
					//fmt.Printf("El nombre del proceso con id: %s es: %s\n",nombreArchivo,splitpuntos[1])
					textocompleto+=archivo.Name()+",";
					aux:=strings.ReplaceAll(splitpuntos[1],"\t","");
					aux=strings.ReplaceAll(aux," ","");
					nombre=aux;
					textocompleto+=aux+",";
				}else if(splitpuntos[0]=="Uid"){
					//USUARIO, falta cambiar por nombre del usuario
					aux:=strings.ReplaceAll(splitpuntos[1],"\t","");
					aux=strings.ReplaceAll(aux," ","");
					textocompleto+=aux+",";					
				}else if(splitpuntos[0]=="State"){
					aux:=strings.ReplaceAll(splitpuntos[1],"\t","");
					aux=strings.ReplaceAll(aux," ","");
					textocompleto+=aux+",";
					if(aux=="R(running)"){
						procejecucion+=1;
					}else if(aux=="S(sleeping)"){
						procsuspendidos+=1;
					}else if(aux=="I(idle)"){
						procdetenidos+=1;
					}else if(aux=="Z(zombie)"){
						proczombies+=1;
					}else{
						fmt.Println("Proceso en estado: %s",aux)
					}
					//PPID es el padre
				}else if(splitpuntos[0]=="PPid"){
					var hj hijo;
					hj.Nombre=nombre;
					hj.pid=archivo.Name();
					aux:=strings.ReplaceAll(splitpuntos[1],"\t","");
					aux=strings.ReplaceAll(aux," ","");
					//fmt.Printf("Llamando metodo de meter hijo")
					meterhijo(hj,aux);
				}else if(splitpuntos[0]=="VmRSS"){
					aux:=strings.ReplaceAll(splitpuntos[1],"\t","");
					aux=strings.ReplaceAll(aux," ","");
					aux=strings.ReplaceAll(aux,"kB","");
					numero, err:= strconv.Atoi(aux);
					if err != nil {
						fmt.Printf("Error convirtiendo a numero la ram: %v", err)
						numero=0;
					}
					//fmt.Printf("Numero convertido: %d",numero);
					dec:=1.565;
					dec=float64(numero)/80611.80;
					//totalram+=dec;
					//dec=dec/8000;
					//aux=strconv.Itoa(dec);
					aux= fmt.Sprintf("%f", dec)
					textocompleto+=aux;
					textocompleto+="%\n";	
					hayram=false;	
				}
			}
			if(hayram){
				//SI NO ENCONTRO EL APARTADO DE RAM
				textocompleto+="0.0%\n";	
			}
		contador++;
    	}
	}
	//TERMINO DE LEER TODOS LOS PROCESOS
	respuestamem := PROCESOS{Ejecucion: procejecucion, Suspendidos: procsuspendidos, Detenidos: procdetenidos, Zombie: proczombies, Total: contador,Datos:textocompleto, Arbol: getprocesos()}
	crearj, errorjson := json.Marshal(respuestamem)
	if errorjson != nil {
		fmt.Println("HAY UN ERROR")
		http.Error(w, errorjson.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println("Mandando json: %s",string(crearj))
	//c	fmt.Println("la memoria libre es ", respuestamem)
	//conver := string(crearj)
	//	fmt.Println("EL indice es ", index)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(crearj)
	fmt.Println("TERMINO CORRECTAMENTE EL PROCESO")
	//imprimirprocesos();
	
}

func imprimirprocesos(){
	for i:=0;i<len(procesos);i++{
		fmt.Printf("Padre: %s{\n",procesos[i].Nombre);
		for j:=0;j<len(procesos[i].Hijos);j++{
			fmt.Printf("\tid: %s,",procesos[i].Hijos[j].pid);
			fmt.Printf("\tNombre: %s\n",procesos[i].Hijos[j].Nombre);
		}
		fmt.Printf("}\n");
	}
}
func meterhijo(proceso hijo,nombre string){
	//proceso es el proceso a meter en hijos
	//procesos la lista de procesos que voy manejando(los padres)
	//nombre el id del padre a buscar
	if(nombre=="0"){
		return
	}
	for i:=0;i<len(procesos);i++ {
		if(nombre== procesos[i].Nombre){
			//SI  ES EL PADRE
			procesos[i].Hijos=append(procesos[i].Hijos,proceso);
			return;
		}
	}
	var aux Proc;
	aux.Nombre=nombre;
	aux.Hijos=append(aux.Hijos,proceso);
	procesos=append(procesos,aux);
}
func getprocesos() string {
	texto:="";
	for i:=0;i<len(procesos);i++{
		texto+=procesos[i].Nombre;
		for j:=0;j<len(procesos[i].Hijos);j++{
			texto+=","+procesos[i].Hijos[j].pid;
			texto+=":"+procesos[i].Hijos[j].Nombre;
		}
		texto+="\n"

	}
	return texto;
}
func reply(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==============================")
	fmt.Println("ENTRE A REPLY")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")


    if err := r.ParseForm(); err != nil {
        fmt.Println(err);
    }

    //this is my first impulse.  It makes the most sense to me.
    //fmt.Println(r.PostForm);          //out -> `map[]`  would be `map[string]string` I think
    //fmt.Println(r.PostForm["hat"]);   //out -> `[]`  would be `fez` or `["fez"]`
   // fmt.Println(r.Body);              //out -> `&{0xc82000e780 <nil> <nil> false true {0 0} false false false}`


    type Hat struct {
        Hat string
	}
	
	type Datos struct{
		Nombre string `json:"Nombre"`
		Departamento string `json:"Departamento"`
		Edad   int    `json:"Edad"`
		Forma string `json:"Forma de contagio"`
		Estado string `json:"Estado"`
	}

    //this is the way the linked SO post above said should work.  I don't see how the r.Body could be decoded.
    decoder := json.NewDecoder(r.Body)
    var t Datos  
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println(err);
	}
	if(t.Hat !=""){
		fmt.Println("Se recibio nombre: ",t.Nombre);
		//app := "kill"
		//cmd:=exec.Command(app,t.Hat);
		//stdout, err := cmd.Output()
		if err != nil{
			fmt.Println("Hubo un error!!!------------");
			fmt.Println(err);
		}else{
			fmt.Println("TODO BINE!!!!");
			fmt.Println(stdout);
		}
	}else{
		fmt.Println("Vino vacio");
	}
}
func main() {
	//obtener ram
	http.HandleFunc("/", obtenerram)
	//obtener info del cpu
	http.HandleFunc("/cpu", obtenercpu)

	http.HandleFunc("/principal", obtenerprincipal)

	http.HandleFunc("/post", reply)

	log.Fatal(http.ListenAndServe(":3030", nil)) //log.fatal como que lo mantiene a la escucha y permite pararlo
}
