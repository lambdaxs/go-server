package monitor

import (
    "fmt"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/process"
    "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/load"
    "os"
)

type CpuInfo struct {
    User float64
    System float64
    Idle float64
    Nice float64
    IOWait float64
    IRQ float64
    SoftIRQ float64
    Steal float64
    Guest float64
    GuestNice float64
}

type MemInfo struct {
    Total uint64
    Available uint64 //free + buffer + cache
    Used uint64
    UsedPercent float64
    Free uint64
}

type NetInfo struct {
    BytesRecv uint64
    BytesSent uint64
}

type ProcessInfo struct {
    Fd uint64
}

type LoadInfo struct {
    Load1 float64
    Load5 float64
    Load15 float64
}


type SystemInfo struct {
    CPU CpuInfo
    Mem MemInfo
    NetIO NetInfo
    Process ProcessInfo
    Load LoadInfo
}

func GetSystemInfo() SystemInfo {
    info := SystemInfo{
        CPU:     CpuInfo{},
        Mem:     MemInfo{},
        NetIO:   NetInfo{},
        Process: ProcessInfo{},
        Load:    LoadInfo{},
    }
    cpuInfo,err := getCpu()
    if err == nil {
        info.CPU = CpuInfo{
            User:      cpuInfo.User,
            System:    cpuInfo.System,
            Idle:      cpuInfo.Idle,
            Nice:      cpuInfo.Nice,
            IOWait:    cpuInfo.Iowait,
            IRQ:       cpuInfo.Irq,
            SoftIRQ:   cpuInfo.Softirq,
            Steal:     cpuInfo.Steal,
            Guest:     cpuInfo.Guest,
            GuestNice: cpuInfo.GuestNice,
        }
    }
    memInfo,err := getMem()
    if err == nil {
        info.Mem = MemInfo{
            Total:       memInfo.Total,
            Available:   memInfo.Available,
            Used:        memInfo.Used,
            UsedPercent: memInfo.UsedPercent,
            Free:        memInfo.Free,
        }
    }
    netInfo,err := getNetIO()
    if err == nil {
        info.NetIO = NetInfo{
            BytesRecv: netInfo.BytesRecv,
            BytesSent: netInfo.BytesSent,
        }
    }
    fd,err := getProcess()
    if err == nil {
        info.Process = ProcessInfo{Fd:fd}
    }
    loadInfo,err := getLoad()
    if err == nil {
        info.Load = LoadInfo{
            Load1:  loadInfo.Load1,
            Load5:  loadInfo.Load5,
            Load15: loadInfo.Load15,
        }
    }
    return info
}

func getCpu() (info cpu.TimesStat,err error){
    res,err := cpu.Times(true)
    if err != nil {
        return
    }
    if len(res) == 1 {
        info = res[0]
    }
    err = fmt.Errorf("cpu count error")
    return
}

func getMem() (info *mem.VirtualMemoryStat,err error){
    info,err = mem.VirtualMemory()
    if err != nil {
        return
    }
    return
}

func getNetIO() (info net.IOCountersStat,err error){
    res,err := net.IOCounters(false)
    if err != nil {
        return
    }
    if len(res) == 1 {
        info = res[0]
        return
    }
    return
}

func getLoad() (res *load.AvgStat, err error){
    return load.Avg()
}


func getProcess() (fd uint64,err error){
    pid := os.Getegid()
    pro,err := process.NewProcess(int32(pid))
    if err != nil {
        return
    }
    info,err := pro.OpenFiles()
    if err != nil {
        return
    }
    if len(info) == 1 {
        fd = info[0].Fd
        return
    }
    return
}
