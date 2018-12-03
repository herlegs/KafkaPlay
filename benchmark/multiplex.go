package benchmark

import (
	"fmt"
	"time"

	"github.com/herlegs/KafkaPlay/benchmark/rate"
)

type Event struct {
	ID          string
	ProduceTime time.Time
}

func StartProducer(id string, qps float64, dataCh chan interface{}) {
	rateLimiter := rate.NewLimiter(rate.Limit(qps), 1)
	for {
		if rateLimiter.Allow() {
			//fmt.Printf("insert id:%v\n",id)
			dataCh <- &Event{
				ID:          id,
				ProduceTime: time.Now(),
			}
		}
	}
}

func StartConsumer(id string, dataCh chan interface{}, consumeTime time.Duration, statsEnabled bool) {
	fmt.Printf("consumer %v started\n", id)
	count := 0
	lastTime := time.Now()
	eventDelays := time.Duration(0)
	for data := range dataCh {
		event := data.(*Event)
		eventDelays += time.Now().Sub(event.ProduceTime)
		count++
		timeDiff := time.Now().Sub(lastTime)
		if statsEnabled && timeDiff.Nanoseconds() > (time.Second).Nanoseconds() {
			fmt.Printf("[%v]average event delay in ns %v\n", id, float64(eventDelays.Nanoseconds())/float64(count))
			fmt.Printf("[%v]average qps %v\n", event.ID, float64(count)/timeDiff.Seconds())
			lastTime = time.Now()
			eventDelays = time.Duration(0)
			count = 0
		}
		time.Sleep(consumeTime)
	}
}

func WithoutMultiplex() {
	dataCh := make(chan interface{})
	go StartConsumer("t1-p1", dataCh, 0, true)
	go StartProducer("t1-p1", 100, dataCh)
	time.Sleep(time.Second * 30)
	//dataCh, t1_p1_ch := make(chan interface{}, 1), make(chan interface{}, 1)
	//go Multiplexer(dataCh, map[string]chan interface{}{
	//	"t1-p1": t1_p1_ch,
	//})
	//go StartProducer("t1-p1", 100, dataCh)
	//
	//go StartConsumer("t1-p1", t1_p1_ch, 0, true)
	//
	//time.Sleep(time.Second * 30)
}

func Multiplexer(dataCh chan interface{}, config map[string]chan interface{}) {
	for data := range dataCh {
		id := data.(*Event).ID
		if ch, ok := config[id]; ok {
			ch <- data
		}
	}
}

func WithMultiplex() {
	dataCh, t1_p1_ch, t2_p1_ch := make(chan interface{}), make(chan interface{}), make(chan interface{}, 1000)
	go Multiplexer(dataCh, map[string]chan interface{}{
		"light": t1_p1_ch,
		"heavy": t2_p1_ch,
	})
	go StartProducer("light", 100, dataCh)
	go StartProducer("heavy", 1000000, dataCh)

	go StartConsumer("light", t1_p1_ch, 0, true)
	go StartConsumer("heavy", t2_p1_ch, time.Microsecond, false)

	time.Sleep(time.Second * 30)
}
