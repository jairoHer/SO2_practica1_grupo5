#include <linux/module.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/uaccess.h>
#include <linux/fs.h>
#include <linux/sysinfo.h>
#include <linux/seq_file.h>
#include <linux/slab.h>
#include <linux/mm.h>
#include <linux/swap.h>
#include <linux/timekeeping.h>
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



static ssize_t my_proc_write(struct file* file, const char __user *buffer, size_t count, loff_t *f_pos){
	return 0;
}

static int my_proc_open(struct inode *inode, struct file *file){
	return single_open(file,my_proc_show,NULL);

}

static struct file_operations my_fops={
	.owner = THIS_MODULE,
	.open = my_proc_open,
	.release = single_release,
	.read = seq_read,
	.llseek = seq_lseek,
	.write = my_proc_write
};


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

static void __exit test_exit(void){
	remove_proc_entry("m_grupo5",NULL);
	printk(KERN_INFO "Bai, att: 5, y este fue el monitor de memoria\n");
}

module_init(test_init);
module_exit(test_exit);
MODULE_LICENSE("GPL");
