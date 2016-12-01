package main

import (
	"encoding/json"
	"fmt"
	"net/http"
//	"time"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	. "github.com/sideshow/apns2/payload"
	"log"
	"github.com/172478394/InnosmartAPNS/conf"
)

type mqtt_msg struct {
    Alert string
    Sound string
    Badge int
    Devicetoken string
    Did string
	Project string
	Topic string
}

var clients map[string]*apns.Client

func pushHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    var t mqtt_msg
    err := decoder.Decode(&t)
    if err != nil {
        //panic(0)
	fmt.Fprintf(w,"{\"code\":-1}")
	return
    }
    notification := &apns.Notification{}
    notification.DeviceToken = t.Devicetoken
    notification.Topic = t.Topic//"cn.innosmart.heimdallrrn"
    notification.Payload = NewPayload().Alert(t.Alert).Sound(t.Sound).Badge(t.Badge)
	client,ok := clients[t.Project]
	if !ok {
		fmt.Fprintf(w,"{\"code\":-2}")
		return
	}
  res, err := client.Push(notification)

  if err != nil {
    log.Println("Error:", err)
    fmt.Fprintf(w,"{\"code\":-3}")
    return
  }
log.Println("APNs ID:", res.ApnsID)
	//time.Sleep(3*time.Second)
	//fmt.Printf("request end")
	fmt.Fprintf(w,"{\"code\":0}")
}

func main() {
	clients = make(map[string]*apns.Client)
	myConfig := new(conf.Config)
	myConfig.InitConfig("client.conf")
	for k, v := range myConfig.Mymap {
	    fmt.Printf("key[%s] value[%s]\n", k, v)
		cert, pemErr := certificate.FromP12File(v["path"], v["password"])
		if pemErr != nil {
			log.Println("Cert Error:", pemErr)
			continue
		}
		client := apns.NewClient(cert).Production()
		clients[k] = client
	}
	http.HandleFunc("/doPush",pushHandler)
	http.ListenAndServe(":9999", nil)
}
