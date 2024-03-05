package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"queueing_theory_golang/queues"

	"github.com/schollz/progressbar/v3"
)

func exportCSV(trialsSimulationTime *[]float64, avgsWaitingTime *[]float64,
	avgsSystemQueue *[]float64, statesProbs *[]float64, path string, header []string) {

	csvFile, err := os.Create(path)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()

	w.Write(header)

	if statesProbs != nil {
		for _, record := range *statesProbs {
			row := []string{fmt.Sprintf("%0.20f", record)}
			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	} else if avgsWaitingTime != nil && trialsSimulationTime != nil {
		args := *trialsSimulationTime
		values := *avgsWaitingTime

		for i := 0; i < len(args); i++ {
			row := []string{fmt.Sprintf("%0.20f", args[i]), fmt.Sprintf("%0.20f", values[i])}
			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	} else if trialsSimulationTime != nil && avgsSystemQueue != nil {
		args := *trialsSimulationTime
		values := *avgsSystemQueue

		for i := 0; i < len(args); i++ {
			row := []string{fmt.Sprintf("%0.20f", args[i]), fmt.Sprintf("%0.20f", values[i])}
			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	}
}

func repeatSimulation(q queues.QueueMMn, simulationCount int) ([]float64, []float64, []float64, []float64) {
	var trialsSimulationTime []float64 = make([]float64, 0)
	var avgsSystemQueueMMn []float64 = make([]float64, 0)
	var avgsSystemQueueMMs []float64 = make([]float64, 0)
	var avgsWaitingTime []float64 = make([]float64, 0)

	bar := progressbar.Default(int64(simulationCount))

	simulationTime := q.GetSimulationTime()
	step := simulationTime / float64(simulationCount)

	for i := step; i <= simulationTime; i += step {
		q := queues.NewQueueMMn(i, 1, 1.0, 20.0)
		qs := queues.NewQueueMMs(i, 1, 1.0, 20.0)
		q.RunSimulation()
		qs.RunSimulation()

		trialsSimulationTime = append(trialsSimulationTime, i)
		avgsSystemQueueMMn = append(avgsSystemQueueMMn, q.FindAvarageSystemQueue())
		avgsSystemQueueMMs = append(avgsSystemQueueMMs, qs.FindAvarageSystemQueue())
		avgsWaitingTime = append(avgsWaitingTime, q.FindWaitingMean())

		bar.Add(1)
	}

	return trialsSimulationTime, avgsSystemQueueMMn, avgsWaitingTime, avgsSystemQueueMMs

}

func main() {
	q := queues.NewQueueMMn(1.0e5, 1, 1.0, 20.0)
	qs := queues.NewQueueMMs(1.0e5, 1, 1.0, 20.0)
	q.RunSimulation()
	qs.RunSimulation()

	fmt.Println("Теоретичні ймовірності станів", qs.FindTheoreticalStateProbability())
	fmt.Println("Пораховані ймовірності станів", qs.FindStateProbability())

	fmt.Println("Теоретична середня к-ть заявок в системі:", qs.FindTheoreticalSystemQueue())
	fmt.Println("Порахована середня к-ть заявок в системі:", qs.FindAvarageSystemQueue())

	fmt.Println("Пораховане:", q.FindWaitingMean())
	rho := 2.0 / 20.0
	mw := rho / (20.0 - 2.0)
	fmt.Println("Теоретичне:", mw)

	fmt.Println("Порахована середня к-ть заявок в системі:", q.FindAvarageSystemQueue())
	fmt.Println("Теоретична середня к-ть заявок в системі:", rho/(1-rho))

	customersCount := 2

	fmt.Println("Порахована ймовірність стану:", q.FindStateProbability())
	fmt.Println("Теоретична ймовірність стану:", math.Pow(rho, float64(customersCount))*(1-rho))
	//drawStatesHistogram(q.FindStateProbability(customersCount))

	probsMMn := q.FindStateProbability()
	probsMMs := qs.FindStateProbability()
	argsTime, avgSystemQueueMMn, avgWaitingTime, avgSystemQueueMMs := repeatSimulation(*q, 10000)

	//fmt.Println("Час:", argsTime)
	//fmt.Println("К-ть заявок:", avgSystemQueue)
	//fmt.Println("Очікування:", avgWaitingTime)

	exportCSV(nil, nil, nil, &probsMMn, "data/statesProbsMMn.csv", []string{"Probabilities"})
	exportCSV(&argsTime, nil, &avgSystemQueueMMn, nil, "data/AvgSystemQueueMMn.csv", []string{"Time", "AvgQueue"})
	exportCSV(&argsTime, &avgWaitingTime, nil, nil, "data/AvgWaiting.csv", []string{"Time", "AvgWaitingTime"})
	exportCSV(nil, nil, nil, &probsMMs, "data/statesProbsMMs.csv", []string{"Probabilities"})
	exportCSV(&argsTime, nil, &avgSystemQueueMMs, nil, "data/AvgSystemQueueMMs.csv", []string{"Time", "AvgQueue"})

}
