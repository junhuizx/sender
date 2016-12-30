package cron

import (
	"github.com/open-falcon/sender/g"
	"github.com/open-falcon/sender/model"
	"github.com/open-falcon/sender/proc"
	"github.com/open-falcon/sender/redis"
	"github.com/toolkits/net/httplib"
	"log"
	"time"
)

func ConsumeSms() {
	queue := g.Config().Queue.Sms
	for {
		L := redis.PopAllSms(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*model.Sms) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

func SendSms(sms *model.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config().Api.Sms
	chatUrl := g.Config().Api.Chat
	resp, err := MsgPost(url, sms.Tos,sms.Content)
	resp2, err := MsgPost(chatUrl, sms.Tos,sms.Content)
	if err != nil {
		log.Println(err)
	}

	proc.IncreSmsCount()

	if g.Config().Debug {
		log.Println("==sms==>>>>", sms)
		log.Println("<<<<==sms==", resp)
		log.Println("<<<<==chat==", resp2)
	}

}

func  MsgPost(url,tos,content string)(string, error){
	r := httplib.Post(url).SetTimeout(5*time.Second, 2*time.Minute)
	r.Param("tos", tos)
	r.Param("content", content)
	return r.String()
}
