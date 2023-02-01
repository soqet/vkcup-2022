package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type fpInput struct {
	data string
	out chan string
}


func Advertise(freeFlowJobs ...job) {
	out := make(chan any)
	for _, j := range freeFlowJobs {
		in := out
		out = make(chan any)
		go func(j job, in, out chan any) {
			j(in, out)
			close(out)
		}(j, in, out)
	}
	<-out
}

// Результатом выполнения этой функции будет строка
// SlowPredict (data)+"-"+SlowPredict (FastPredict (data))
// где data — то, что пришло на вход
func GetProfile(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	fpin := make(chan fpInput)
	defer close(fpin)
	go fastPredictProxyHolder(fpin)
	for dataRaw := range in {
		data := fmt.Sprint(dataRaw)
		wg.Add(1)
		go func() {
			firstRes := make(chan string)
			go func () {
				firstRes <- SlowPredict(data)
			}()
			secondRes := make(chan string)
			go func () {
				secondRes <- SlowPredict(fastPredictProxy(data, fpin))
			}()
			out <- <-firstRes + "-" + <-secondRes
			wg.Done()
		}()
	}
	wg.Wait()
}


// Результатом выполнения этой функции будет строка
// SlowPredict (th+data) 
// (конкатенация цифры, приведённой к строке, и строки), где th=0..5
// потом берёт конкатенацию результатов в порядке расчёта (0..5)
func GetGroup(in, out chan interface{}) {
	size := 6 // [0, 6)
	wg := sync.WaitGroup{}
	for dataRaw := range in {
		data := fmt.Sprint(dataRaw)
		wg.Add(1)
		go func() {
			res := make([]string, size)
			allGroups := sync.WaitGroup{}
			allGroups.Add(size)
			for i := 0; i < size; i++ {
				i := i
				go func() {
					res[i] = SlowPredict(strconv.FormatInt(int64(i), 10) + data)
					allGroups.Done()
				}()
			}
			allGroups.Wait()
			out <- strings.Join(res, "")
			wg.Done()
		}()
	}
	wg.Wait()
}


// получает все результаты для заданного набора пользователей, 
// сортирует, объединяет отсортированный результат через _
func ConcatProfiles(in, out chan interface{}) {
	res := []string{}
	for dataRaw := range in {
		data := fmt.Sprint(dataRaw)
		i := sort.SearchStrings(res, data)
		res = append(res, "")
		copy(res[i+1:], res[i:])
		res[i] = data
	}
	out <- strings.Join(res, "_")
}

func fastPredictProxyHolder(in chan fpInput) {
	for input := range in {
		input.out <- FastPredict(input.data)
	}
}

func fastPredictProxy(data string, fpin chan fpInput) string {
	out := make(chan string)
	fpin <- fpInput{data: data, out: out}
	return <-out
}
