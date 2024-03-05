package queues

import (
	"fmt"
	"math"
	"math/rand"
)

type QueueMMs struct {
	serversData      []float64
	arrivalRate      float64
	serviceRate      float64
	simulationTime   float64
	totalWaitTime    float64
	lastTime         float64
	totalServiceTime float64
	statesData       []float64
	statesQLen       []int
	customersCount   int
	numServers       int
}

func NewQueueMMs(simulationTime float64, numServers int, arrivalRate float64, serviceRate float64) *QueueMMs {
	q := QueueMMs{
		serversData:      make([]float64, numServers+1),
		statesData:       make([]float64, 0),
		statesQLen:       make([]int, 0),
		arrivalRate:      arrivalRate,
		serviceRate:      serviceRate,
		simulationTime:   simulationTime,
		numServers:       numServers,
		totalServiceTime: 0.0,
		lastTime:         0.0,
		customersCount:   0,
	}

	q.serversData[0] = 0.0

	for i := 1; i <= numServers; i++ {
		q.serversData[i] = math.Inf(1)
	}

	return &q
}

func (q *QueueMMs) getArrivalTime() float64 {
	return rand.ExpFloat64() / q.arrivalRate
}

func (q *QueueMMs) getServiceTime() float64 {
	return rand.ExpFloat64() / q.serviceRate
}

func (q *QueueMMs) RunSimulation() {
	var nextArrivalTime float64 = q.getArrivalTime()
	//iter := 0

	actualTime := 0.0
	/*
		q.queueData = append(q.queueData, q.qLen)
	*/

	//fmt.Println(q.serversData)
	//fmt.Println("Наступна заявка надійде через:", nextArrivalTime)
	//fmt.Println("Розмір черги:", q.qLen)
	//fmt.Println("Час, який пройшов:", actualTime)
	//fmt.Println()

	for actualTime < q.simulationTime {
		nextfreeServerIndex, nextfreeServerTime := q.findNextFreeServerIndex()
		//fmt.Println(q.customerData)
		//fmt.Println("Індекс наступного вільного сервера:", nextfreeServerIndex)
		//fmt.Println("Час до звільнення сервера:", nextfreeServerTime)

		if nextArrivalTime < nextfreeServerTime {
			actualTime += nextArrivalTime
			q.customersCount += 1

			if int(q.serversData[0])+1 > len(q.statesData) {
				q.addState()
			}
			q.statesData[int(q.serversData[0])] += nextArrivalTime

			//customerID := len(q.customerData)/2 + 1

			if q.serversData[0] < float64(q.numServers) {
				q.serversData[0] += 1.0

				//q.customerData[customerDataKey{customerID, "arrival"}] = actualTime
				//q.customerData[customerDataKey{customerID, "startService"}] = actualTime

				freeServerIndex := q.findFreeServerIndex()

				for i := 1; i <= q.numServers; i++ {
					q.serversData[i] -= nextArrivalTime
				}

				tmpServiceTime := q.getServiceTime()
				q.totalServiceTime += tmpServiceTime
				q.serversData[freeServerIndex] = tmpServiceTime
			}

			nextArrivalTime = q.getArrivalTime()

		} else if nextArrivalTime > nextfreeServerTime {
			actualTime += nextfreeServerTime
			nextArrivalTime -= nextfreeServerTime

			if int(q.serversData[0])+1 > len(q.statesData) {
				q.addState()
			}

			q.statesData[int(q.serversData[0])] += nextfreeServerTime

			q.serversData[0] -= 1.0

			for i := 1; i <= q.numServers; i++ {
				if i == nextfreeServerIndex {
					continue
				}
				q.serversData[i] -= nextfreeServerTime
			}

			q.serversData[nextfreeServerIndex] = math.Inf(1)
		}
		//fmt.Println("States:", q.statesData)
		/*
			iter++
			fmt.Println("Iteration:", iter)
			fmt.Println("Total wait time:", q.totalWaitTime)
			fmt.Println(q.serversData)
			fmt.Println("Наступна заявка надійде через:", nextArrivalTime)
			fmt.Println("Час, який пройшов:", actualTime)
			fmt.Println()
		*/
	}
}

func (q *QueueMMs) FindTheoreticalStateProbability() []float64 {
	res := make([]float64, q.numServers+1)

	P0 := 0.0
	rho := q.arrivalRate / q.serviceRate

	for i := 0; i <= q.numServers; i++ {

		P0 += math.Pow(rho, float64(i)) / float64(factorial(i))
	}

	res[0] = 1.0 / P0

	for j := 1; j < len(res); j++ {
		res[j] = res[j-1] * math.Pow(rho, float64(j)) / float64(factorial(j))
	}

	return res
}

func (q *QueueMMs) FindTheoreticalSystemQueue() float64 {
	statesProbabilities := q.FindTheoreticalStateProbability()
	res := 0.0
	for i := 0; i < len(statesProbabilities); i++ {
		res += float64(i) * statesProbabilities[i]
	}

	return res
}

func (q *QueueMMs) UpdateSimulationTime(newTime float64) {
	q.simulationTime = newTime
}

func (q *QueueMMs) GetSimulationTime() float64 {
	return q.simulationTime
}

func (q *QueueMMs) DataReset() {
	q.serversData = make([]float64, q.numServers+1)
	q.statesData = make([]float64, 0)
	q.statesQLen = make([]int, 0)
	q.totalWaitTime = 0.0
	q.totalServiceTime = 0.0
	q.lastTime = 0.0
	q.customersCount = 0
}

func (q *QueueMMs) addState() {
	q.statesData = append(q.statesData, 0.0)
}

func (q *QueueMMs) FindStateProbability() []float64 {

	//fmt.Println("States count", len(q.statesCountData))
	//fmt.Println("States count my", q.statesCount)
	statesProbabilities := make([]float64, len(q.statesData))

	for i := 0; i < len(q.statesData); i++ {
		statesProbabilities[i] = q.statesData[i] / q.simulationTime
	}

	return statesProbabilities
}

func (q *QueueMMs) FindAvarageSystemQueue() float64 {

	//fmt.Println("Theta", q.thetaData)
	//fmt.Println("States", q.statesCountData)
	//fmt.Println("Len Theta", len(q.thetaData))
	//fmt.Println("Len States", len(q.thetaData))

	//sum += float64(q.statesQLen[l-1]) * (q.simulationTime - q.thetaData[l-1])

	res := q.totalServiceTime / q.simulationTime

	if res > 1 {
		fmt.Println("simulation time:", q.simulationTime)
		fmt.Println("total service time:", q.totalServiceTime)
	}

	return res
}

func (q *QueueMMs) findNextFreeServerIndex() (int, float64) {
	index := 1
	smallest := q.serversData[index]

	for i := 2; i <= q.numServers; i++ {
		if q.serversData[i] < smallest {
			index = i
			smallest = q.serversData[i]
		}
	}

	return index, smallest
}

func (q *QueueMMs) findFreeServerIndex() int {
	index := 0

	for i := 1; i <= q.numServers; i++ {
		if q.serversData[i] == math.Inf(1) {
			index = i
			break
		}
	}

	return index
}

func factorial(n int) int {
	factVal := 1

	if n < 0 {
		fmt.Print("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			factVal *= i // mismatched types int64 and int
		}

	}
	return factVal
}
