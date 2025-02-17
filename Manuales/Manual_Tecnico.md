# Manual Tecnico practica 1 grupo 5


### Integrantes
| Carnet | Nombre |
| ------ | ------ |
| 201404423 | Jairo Pablo Hernández Guzmán |
| 201504042 | Julio Estuardo Gómez Alonzo  |
| 201503750 | Carlos Eduardo Carías Salan |

## Servidor Web
Este servidor contedra los modulos que leeran los procesos del sistema y el modulo de que leera los datos de la RAM utilizada. Tambien en este servidor tendra un servidor go el cual tendra la funcion de enviar la informacion de los modulos a traves de API's.

### Requisitos
Herramientas necesarias para el correcto funcionamiento de la aplicación.
  - Gcc
  - Make
  - Go
  - Git

### API's
Estas son las encargadas de comunicar la informacion extraida de los modulos con un servicio cliente.

El programa escucha en el puerto 5050
```
log.Fatal(http.ListenAndServe(":5050",myRouter))
```
Para poder crear las API's se utilizo la libreria gorilla/mux la cual facilita la creacion de las mismas. Para usar se debe instalar ingresando el siguinte comando en la terminal.
```
go get -u github.com/gorilla/mux
```
Las librerias que se deben importar para el funcionamiento de la aplicacion en go son las siguientes:
```
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
```
La api que extrae la informacion que se encarga de enviar toda la informacion relacionada con los procesos del sistema es **/todo**
```
myRouter.HandleFunc("/todo",toditotodito).Methods("GET")
```
Esta api hace llama a la funcion **toditotodito** la cual se encarga de leer el archivo /proc/proc_procesos la cual contiene la informacion relacionada con los proceso del sistema.
```
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


}
```
Esta funcion leer el archivo /proc/proc_procesos, lo convierte en json y retorna el json con la informacion.

La API que extrae la inforamcion de la RAM es **/info** 
```
myRouter.HandleFunc("/info",memoria_info).Methods("GET")
```
Esta api hace llama a la funcion **memoria_info** la cual se encarga de leer el archivo /proc/m_grupo5 la cual contiene la informacion relacionada con los RAM del sistema.
```
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
```
Para poder eliminar procesos del sistema se hace a traves de una peticion POST, la cual se encuentra en la ruta **/kill**. Esta peticion recibe un json el cual contiene el PID del proceso a eliminar.
```
myRouter.HandleFunc("/kill",killProcess).Methods("POST")
```
formato del json (el valor ingresado es solo un ejemplo)
```
{
    "proceso":"10785"
}
```
Esta api hace llama a la funcion **killProcess** la cual se encarga de procesar la peticion POST y de eliminar el proceso.
```
func killProcess(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var proce procesoBorrado
	_ = json.NewDecoder(r.Body).Decode(&proce)
	PID := string(proce.Proceso)
	fmt.Println(PID)
	json.NewEncoder(w).Encode(proce)
	_, err := exec.Command("sh", "-c", "kill "+PID).Output()
	if err != nil {
		fmt.Println("No se pudo eliminar el proceso")
	} else {
		fmt.Println("Se elimino el proceso con PID: "+PID)
	}

}
```
## MODULOS
Estos son los modulos creados para la lectura de la informacion de la RAM y los procesos del sistema.
### Modulo de RAM
Este modulo se llama m_grupo5, se encarga de guardar la informacion del total de RAM del sistema, Cantidad de RAM utilizada y % de utilizacion del sistema. Esta información la guarda en el archivo /proc/m_grupo5 a la cual se encarga de acceder el programa en Go para exportarlos.

El archivo que contiene el codigo del modulo es **m_grupo5.c** en este se obtiene la informacion correspondiente.

Funcion que se encarga de extraer la informacion
```
static int my_proc_show(struct seq_file *m, void *v){

	unsigned long ram_total;
	unsigned long ram_libre;
	unsigned long ram_consumida;
	unsigned long unidad;
	unsigned long porcentaje;
	struct sysinfo info_sistema;
	si_meminfo(&info_sistema);
	unidad = (unsigned long long)info_sistema.mem_unit;
	ram_total = (info_sistema.totalram*unidad)/(1024*1024);
	ram_libre = (info_sistema.freeram*unidad)/(1024*1024);
	ram_consumida = ram_total - ram_libre;
	porcentaje = ((ram_total - ram_libre)*100)/ram_total;
	seq_printf(m, "%ld %ld %ld",ram_total,ram_consumida,porcentaje);
	//seq_printf(m, "porcentaje ocupado %ld ram total: %ld MB ram libre %ld MB\n",porcentaje,ram_total,ram_libre);
	return 0;
}
```
funcion que crea el archivo **m_grupo5** en la carpeta /proc
```
static int __init test_init(void){
	struct proc_dir_entry *entry;
	printk(KERN_INFO "Buenas, att: 5, monitor de memoria\n");
	entry = proc_create("m_grupo5",0777,NULL, &my_fops);
	if(!entry){
		return -1;
	} else {
		printk(KERN_INFO "Inicio modulo ram \n");	
	}
	return 0;
}
```
El archivo Makefile es el que se encarga de crear el modulo que esta escrito en el archivo c, para ello se crea un archivo llamado **Makefile** y en la terminal se ingresa **make** para ejecutarlo.

Codigo del Makefile:
```
obj-m +=  m_grupo5.o
all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules
clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```
Despues de correr el comando make, se debe insertar el modulo. Para esto se usa los siguiente comandos en la terminal, ingresando el archivo .ko:
```
sudo insmod m_grupo5.ko
```
### Modulo de procesos
Este modulo se encarga de obtener la informacion de todos los procesos que se estan en el sistema operativo. Tambien extrae la informacion de los procesos hijos. Despues de recolectar toda la informcionde los procesos, este obtiene cuantos procesos se encuentran en los diferentes estados (Running, Sleep, Zombie, Stoped, Idle)

El archivo que contiene el codigo del modulo es **m_grupo5.c** en este se obtiene la informacion correspondiente

Funcion que se encarga de extraer la informacion:
```
static int myread (struct seq_file *buff, void *v){
    int proc_idle = 0;
    int proc_running = 0;
    int proc_sleep = 0;
    int proc_stoped = 0;
    int proc_zombie = 0;
    int num_proc = 0;
    int num_proceso = 0;
    int num_proceso2= 0;
    int num_hijo = 0;
    int num_hijo2 = 0;
    uid_t uid;
    printk(KERN_INFO "EMPEZANDO MODULO PROCESOS\n");
    seq_printf(buff, "{%s\"procesos\":[","\n");
    for_each_process( task ){ 
	num_proceso2++;
    }
    for_each_process( task ){            /*    for_each_process() es un MACRO para iterar ubicado en linux\sched\signal.h    */
        uid = __kuid_val(task_uid(task));
        num_proc++;
	    num_proceso++;
        num_hijo = 0;
        num_hijo2 = 0;
        seq_printf(buff, "%s","{\n");
        if(task->state == 1026){
            seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\", \n\t\"CHILDS\":[\n", task->pid, task->comm, uid, "I(idle)");
            proc_idle++;
        }
        if(task->state == 0){
            seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\", \n\t\"CHILDS\":[\n", task->pid, task->comm, uid, "R(running)");
            proc_running++;
        }
        if(task->state == 1){
            seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\", \n\t\"CHILDS\":[\n", task->pid, task->comm, uid, "S(sleep)");
            proc_sleep++;
        }
        list_for_each(list, &task->children){
            num_hijo ++;
        }
        list_for_each(list, &task->children){  
            num_hijo2++;                      /*    list_for_each MACRO para iterar task->children    */
            task_child = list_entry( list, struct task_struct, sibling );    /*    using list_entry to declare all vars in task_child struct    */
            seq_printf(buff, "%s","{\n");
            uid = __kuid_val(task_uid(task_child));
            if(task_child->state == 1026){
                seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\" ", task_child->pid, task_child->comm, uid, "I(idle)");
                proc_idle++;
            }
            if(task_child->state == 0){
                seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\"", task_child->pid, task_child->comm, uid, "R(running)");
                proc_running++;
            }
            if(task_child->state == 1){
                seq_printf(buff, "\t\t\"PID\": \"%d\",\n\t\t\"PROCESS\": \"%s\", \n\t\t\"UID\": \"%d\", \n\t\t\"STATE\": \"%s\"", task_child->pid, task_child->comm, uid, "S(sleep)");
                proc_sleep++;
            }
            if(num_hijo==num_hijo2){
                seq_printf(buff, "%s","\n}\n");
            }else{
                seq_printf(buff, "%s","\n},\n");
            }
            //seq_printf(buff, "%s","\n},\n");
        }
        if (num_proceso==num_proceso2){
            seq_printf(buff, "%s","\t]\n\n}\n");
        }else{
            seq_printf(buff, "%s","\t]\n\n},\n");
        }
        //seq_printf(buff, "%s","\t]\n\n},\n");
    }    
    seq_printf(buff, "%s","]}\n");
    seq_printf(buff, "%s","-------------------------------------\n");
    seq_printf(buff, "TOTAL: %d\n",proc_idle+proc_running+proc_sleep);
    seq_printf(buff, "IDLE: %d\n",proc_idle);
    seq_printf(buff, "RUNNING: %d\n",proc_running);
    seq_printf(buff, "SLEEP: %d\n",proc_sleep);
    seq_printf(buff, "STOPED: %d\n",proc_stoped);
    seq_printf(buff, "ZOMBIE: %d\n",proc_zombie);
    
    return 0;
}
```
funcion que crea el archivo **proc_procesos** en la carpeta /proc
```
static int simple_init(void){

    printk(KERN_INFO "EMPEZANDO MODULO PROCESOS\n");
    ent=proc_create("proc_procesos",0,NULL,&myops);
    return 0;
}
```
El archivo Makefile es el que se encarga de crear el modulo que esta escrito en el archivo c, para ello se crea un archivo llamado **Makefile** y en la terminal se ingresa **make** para ejecutarlo.

Codigo del Makefile:
```
obj-m +=  procesos.o
all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules
clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```
Despues de correr el comando make, se debe insertar el modulo. Para esto se usa los siguiente comandos en la terminal, ingresando el archivo .ko:
```
sudo insmod procesos.ko
```

