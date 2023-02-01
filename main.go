package main

import (
	"fmt"
	"strconv"
	"time"
)

func MeasureFastPredict() {
	fpin := make(chan fpInput)
	defer close(fpin)
	go fastPredictProxyHolder(fpin)
	var totalFunc time.Duration
	var totalProxy time.Duration
	for i := 0; i < 1000; i++ {
		data := strconv.FormatInt(int64(i), 10)
		start := time.Now()
		FastPredict(data)
		totalFunc += time.Since(start)
		start = time.Now()
		fastPredictProxy(data, fpin)
		totalProxy += time.Since(start)
	}
	fmt.Printf("FastPredict: %s, FastPredictProxy: %s\n",(totalFunc / time.Millisecond * 1000), (totalProxy / time.Millisecond * 1000))

}

func main() {
	var testRes string
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for i := 0; i < 64; i++ {
				out <- i * 1000
			}
		}),
		job(GetProfile),
		job(GetGroup),
		job(ConcatProfiles),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				panic("cant convert result data to string")
			}
			testRes = data
		}),
	}

	start := time.Now()

	Advertise(hashSignJobs...)

	end := time.Since(start)

	expectedTime := 3 * time.Second

	testRes += "asd"
	fmt.Println(end, expectedTime)
}
