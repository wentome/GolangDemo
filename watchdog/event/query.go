package event

import (
	"time"

	"github.com/astaxie/beego/logs"

	"../../../ctask"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Query(ch chan string) {
	count0 := 0
	count1 := 0
	opts := MQTT.NewClientOptions().AddBroker("tcp://39.99.185.210:1883")
	opts.SetClientID("go-simple")

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(250)
	for {
		ctask.TaskContor(ch)

		if count0 > 3 {
			count0 = 0
			ctask.TaskFeedDog(ch)
		}

		if count1 > 10 {
			count1 = 0
			logs.Warn("star push")
			token := c.Publish("alert/uat", 0, false, "alert")
			token.Wait()
			logs.Warn("end push")
		}

		time.Sleep(time.Second)
		count0++
		count1++
	}
}
