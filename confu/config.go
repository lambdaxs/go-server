package confu

import (
    "errors"
    "fmt"
    "github.com/hashicorp/consul/api"
    "github.com/jinzhu/configor"
    "io/ioutil"
    "os"
    "time"
)

//通过本地文件初始化配置
func InitWithFilePath(path string, data interface{}) error {
    if err := configor.Load(data, path); err != nil {
        return err
    }
    return nil
}

//通过consul远端文件初始化配置
func InitWithRemotePath(path string, data interface{}, remoteAddr string) error {
    if remoteAddr == "" {
        remoteAddr = os.Getenv("CONSUL_ADDR")
    }
    if remoteAddr == "" {
        return errors.New("env var consul_addr is empty")
    }
    config := api.DefaultConfig()
    config.Address = remoteAddr
    client, err := api.NewClient(config)
    if err != nil {
        return fmt.Errorf("new consul client err:%s", err.Error())
    }
    config.HttpClient.Timeout = time.Second*5

    //存储在本地文件
    localPath := fmt.Sprintf("./consul-%s",path)
    kv, _, err := client.KV().Get(path, nil)
    if err != nil {//远端数据查询失败,容错从本地文件获取配置数据
        if localErr := configor.Load(data, localPath); localErr != nil {
            return fmt.Errorf("load remote config error:%s local config error:%s", err.Error(), localErr.Error())
        }else {
            return nil
        }
    }
    if err := ioutil.WriteFile(localPath, kv.Value, os.ModePerm);err != nil {
        return fmt.Errorf("write file error"+err.Error())
    }
    if err := configor.Load(data, localPath); err != nil {
        return err
    }
    return nil
}
