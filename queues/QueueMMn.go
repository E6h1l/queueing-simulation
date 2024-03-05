package queues

import (
	"fmt"
	"math"
	"math/rand"
)

type customerDataKey struct {
	id        int
	eventType string
}

type QueueMMn struct {
	customerData     map[customerDataKey]float64
	serversData      []float64
	thetaData        []float64
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
	qLen             int
}

func NewQueueMMn(simulationTime float64, numServers int, arrivalRate float64, serviceRate float64) *QueueMMn {
	q := QueueMMn{
		customerData:     make(map[customerDataKey]float64),
		serversData:      make([]float64, numServers+1),
		thetaData:        make([]float64, 0),
		statesData:       make([]float64, 0),
		statesQLen:       make([]int, 0),
		arrivalRate:      arrivalRate,
		serviceRate:      serviceRate,
		simulationTime:   simulationTime,
		numServers:       numServers,
		totalWaitTime:    0.0,
		totalServiceTime: 0.0,
		lastTime:         0.0,
		customersCount:   0,
		qLen:             0,
	}

	q.serversData[0] = 0.0

	for i := 1; i <= numServers; i++ {
		q.serversData[i] = math.Inf(1)
	}

	return &q
}

func (q *QueueMMn) getArrivalTime() float64 {
	return rand.ExpFloat64() / q.arrivalRate
}

func (q *QueueMMn) getServiceTime() float64 {
	return rand.ExpFloat64() / q.serviceRate
}

func (q *QueueMMn) RunSimulation() {
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

			q.thetaData = append(q.thetaData, nextArrivalTime)

			if q.qLen+int(q.serversData[0])+1 > len(q.statesData) {
				q.addState()
			}
			q.statesData[q.qLen+int(q.serversData[0])] += nextArrivalTime

			//customerID := len(q.customerData)/2 + 1

			if q.serversData[0] == float64(q.numServers) {
				//q.serversData[0] = float64(q.numServers)
				//q.totalBusyServers += q.numServers

				q.customerData[customerDataKey{q.customersCount, "arrival"}] = actualTime
				q.customerData[customerDataKey{q.customersCount, "startService"}] = 0.0

				for i := 1; i <= q.numServers; i++ {
					q.serversData[i] -= nextArrivalTime
				}

				q.qLen += 1

			} else if q.serversData[0] < float64(q.numServers) {
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

			q.thetaData = append(q.thetaData, nextfreeServerTime)

			if q.qLen+int(q.serversData[0])+1 > len(q.statesData) {
				q.addState()
			}

			q.statesData[q.qLen+int(q.serversData[0])] += nextfreeServerTime

			if q.qLen > 0 {
				q.serversData[0] = float64(q.numServers)
				//q.totalBusyServers += q.numServers

				firstQueueID := q.findFirstInQueue()
				//fmt.Println("ID", firstQueueID)
				q.customerData[customerDataKey{firstQueueID, "startService"}] = actualTime
				//fmt.Println("ID", firstQueueID)
				//fmt.Println("actual time:", actualTime)
				//fmt.Println("Wait time:", q.customerData[customerDataKey{firstQueueID, "arrival"}])
				q.totalWaitTime += q.customerData[customerDataKey{firstQueueID, "startService"}] - q.customerData[customerDataKey{firstQueueID, "arrival"}]
				//fmt.Println("Total wait time:", q.totalWaitTime)
				//fmt.Println("Before DELETE:", q.customerData)
				delete(q.customerData, customerDataKey{firstQueueID, "startService"})
				delete(q.customerData, customerDataKey{firstQueueID, "arrival"})
				//fmt.Println("After DELETE:", q.customerData)

				for i := 1; i <= q.numServers; i++ {
					if i == nextfreeServerIndex {
						continue
					}
					q.serversData[i] -= nextfreeServerTime
				}

				tmpServiceTime := q.getServiceTime()
				q.totalServiceTime += tmpServiceTime
				q.serversData[nextfreeServerIndex] = tmpServiceTime
				q.qLen -= 1

			} else {
				q.serversData[0] -= 1.0

				for i := 1; i <= q.numServers; i++ {
					if i == nextfreeServerIndex {
						continue
					}
					q.serversData[i] -= nextfreeServerTime
				}

				q.serversData[nextfreeServerIndex] = math.Inf(1)
			}
		}

		q.statesQLen = append(q.statesQLen, q.qLen)
		//fmt.Println("States:", q.statesData)

		/*
			//iter++
			fmt.Println(q.customerData)
			//fmt.Println("Iteration:", iter)
			fmt.Println("Total wait time:", q.totalWaitTime)
			fmt.Println(q.serversData)
			fmt.Println("Наступна заявка надійде через:", nextArrivalTime)
			fmt.Println("Розмір черги:", q.qLen)
			fmt.Println("Час, який пройшов:", actualTime)
			fmt.Println()
		*/
	}

}

func (q *QueueMMn) UpdateSimulationTime(newTime float64) {
	q.simulationTime = newTime
}

func (q *QueueMMn) GetSimulationTime() float64 {
	return q.simulationTime
}

func (q *QueueMMn) DataReset() {
	q.customerData = make(map[customerDataKey]float64)
	q.serversData = make([]float64, q.numServers+1)
	q.thetaData = make([]float64, 0)
	q.statesData = make([]float64, 0)
	q.statesQLen = make([]int, 0)
	q.totalWaitTime = 0.0
	q.totalServiceTime = 0.0
	q.lastTime = 0.0
	q.customersCount = 0
	q.qLen = 0
}

func (q *QueueMMn) addState() {
	q.statesData = append(q.statesData, 0.0)
}

func (q *QueueMMn) FindStateProbability() []float64 {

	//fmt.Println("States count", len(q.statesCountData))
	//fmt.Println("States count my", q.statesCount)
	statesProbabilities := make([]float64, len(q.statesData))

	for i := 0; i < len(q.statesData); i++ {
		statesProbabilities[i] = q.statesData[i] / q.simulationTime
	}

	return statesProbabilities
}

func (q *QueueMMn) FindAvarageSystemQueue() float64 {
	l := len(q.thetaData)
	sum := 0.0

	//fmt.Println("Theta", q.thetaData)
	//fmt.Println("States", q.statesCountData)
	//fmt.Println("Len Theta", len(q.thetaData))
	//fmt.Println("Len States", len(q.thetaData))

	for i := 0; i < l-1; i++ {
		sum += float64(q.statesQLen[i]) * q.thetaData[i]
	}

	//sum += float64(q.statesQLen[l-1]) * (q.simulationTime - q.thetaData[l-1])

	res := (sum / q.simulationTime) + (q.totalServiceTime / q.simulationTime)

	if res > 1 {
		fmt.Println("simulation time:", q.simulationTime)
		fmt.Println("total service time:", q.totalServiceTime)
		fmt.Println("total time in queue:", sum)
	}

	return res
}

func (q *QueueMMn) FindWaitingMean() float64 {
	//fmt.Println("Total wait time:", q.totalWaitTime)
	//fmt.Println("Customer count:", q.customersCount)
	//fmt.Println(q.customerData)

	return q.totalWaitTime / float64(q.customersCount)
}

func (q *QueueMMn) findFirstInQueue() int {
	var minValueKey int
	var minValue float64

	for i := range q.customerData {
		if q.customerData[customerDataKey{i.id, "startService"}] == 0.0 {
			minValueKey = i.id
			minValue = q.customerData[customerDataKey{i.id, "arrival"}]
			break
		}
	}

	for i := range q.customerData {
		if q.customerData[customerDataKey{i.id, "startService"}] == 0.0 && q.customerData[customerDataKey{i.id, "arrival"}] < minValue {
			minValueKey = i.id
			minValue = q.customerData[customerDataKey{i.id, "arrival"}]
		}
	}

	return minValueKey

}

func (q *QueueMMn) findNextFreeServerIndex() (int, float64) {
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

func (q *QueueMMn) findFreeServerIndex() int {
	index := 0

	for i := 1; i <= q.numServers; i++ {
		if q.serversData[i] == math.Inf(1) {
			index = i
			break
		}
	}

	return index
}
