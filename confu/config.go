package confu

import (
    "errors"
    "flag"
    "fmt"
    "github.com/hashicorp/consul/api"
    "github.com/jinzhu/configor"
    "io/ioutil"
    "os"
    "time"
)

//解析本地和远端配置参数
func ParseFlag() (string, string){
    var configPath string
    var remoteConfigPath string
    flag.StringVar(&configPath, "config", "", "config file path")
    flag.StringVar(&remoteConfigPath, "remote-config", "", "remote config file path")
    flag.Parse()
    return configPath, remoteConfigPath
}

//通过本地文件初始化配置
func InitWithFilePath(path string, data interface{}) error {
    if err := configor.Load(data, path); err != nil {
        return err
    }
    return nil
}

//通过consul远端文件初始化配置
func InitWithRemotePath(path string, data interface{}, remoteAddr string) (content []byte,err error) {
    if remoteAddr == "" {
        remoteAddr = os.Getenv("CONSUL_ADDR")
    }
    if remoteAddr == "" {
        err = errors.New("env var consul_addr is empty")
        return
    }
    config := api.DefaultConfig()
    config.Address = remoteAddr
    client, err := api.NewClient(config)
    if err != nil {
        err = fmt.Errorf("new consul client err:%s", err.Error())
        return
    }
    config.HttpClient.Timeout = time.Second*5

    //存储在本地文件
    localPath := fmt.Sprintf("./consul-%s",path)
    kv, _, err := client.KV().Get(path, nil)
    if err != nil {//远端数据查询失败,容错从本地文件获取配置数据
        if localErr := configor.Load(data, localPath); localErr != nil {
            err = fmt.Errorf("load remote config error:%s local config error:%s", err.Error(), localErr.Error())
            return
        }else {
            buf,IOErr := ioutil.ReadFile(localPath)
            if IOErr != nil {
                err = errors.New("local config file load err:"+IOErr.Error())
                return
            }
            content = buf
            return
        }
    }

    if err = ioutil.WriteFile(localPath, kv.Value, os.ModePerm);err != nil {
        err = fmt.Errorf("write file error"+err.Error())
        return
    }
    if err = configor.Load(data, localPath); err != nil {
        return
    }
    content = kv.Value
    return
}
