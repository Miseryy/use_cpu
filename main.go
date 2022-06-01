package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"os"
	"strconv"
	"strings"
	// "time"
)

const file_name string = "last"

type CPU_Status struct {
    user     int
    nice     int
    sys      int
    idle     int
    iowait   int
    irq      int
    softirq  int
}

func set_status(source []byte, status *CPU_Status) {
    str := string(source)
    str_arr := strings.Split(str, " ")

    status.user, _ = strconv.Atoi(str_arr[2])
    status.nice, _ = strconv.Atoi(str_arr[3])
    status.sys, _ = strconv.Atoi(str_arr[4])
    status.idle, _ = strconv.Atoi(str_arr[5])
    status.iowait, _ = strconv.Atoi(str_arr[6])
    status.irq, _ = strconv.Atoi(str_arr[7])
    status.softirq, _ = strconv.Atoi(str_arr[8])

}

func get_status() CPU_Status {
    current := CPU_Status {
        user:    0,
        nice:    0,
        sys:     0,
        idle:    0,
        iowait:  0,
        irq:     0,
        softirq: 0,
    }

    out, _ := exec.Command("head", "-1", "/proc/stat").Output()

    set_status(out, &current)

    cd, _:= os.Getwd()
    ioutil.WriteFile(cd+"/"+file_name, out, 0666)

    return current
}

func read_last() CPU_Status {
    last := CPU_Status {
        user:    0,
        nice:    0,
        sys:     0,
        idle:    0,
        iowait:  0,
        irq:     0,
        softirq: 0,
    }

    cd, _:= os.Getwd()
    bytes, err := ioutil.ReadFile(cd+"/"+file_name)

    if err != nil {
        return last
    }

    set_status(bytes, &last)

    return last

}

func culc(current int, last int, total int) (int, int){
    num   := ((current - last)*10000/total)/100
    comma := ((current - last)*10000/total)%100

    return num, comma

}

func show() {
    last := read_last()

    current := get_status()

    total := current.user - last.user +
             current.nice - last.nice +
             current.sys - last.sys +
             current.idle - last.idle +
             current.iowait - last.iowait +
             current.irq - last.irq +
             current.softirq - last.softirq

    user_num, user_comma := culc(current.user, last.user, total)
    nice_num, nice_comma := culc(current.nice, last.nice, total)
    sys_num, sys_comma := culc(current.sys, last.sys, total)
    idle_num, idle_comma := culc(current.idle, last.idle, total)

    us := fmt.Sprintf("us:%d.%d", user_num, user_comma)
    ni := fmt.Sprintf("ni:%d.%d", nice_num, nice_comma)
    sy := fmt.Sprintf("sy:%d.%d", sys_num, sys_comma)
    id := fmt.Sprintf("id:%d.%d", idle_num, idle_comma)

    fmt.Printf("%s %s %s %s\n", us, ni, sy, id)

    last = current
}

func main() {
    show()
}
