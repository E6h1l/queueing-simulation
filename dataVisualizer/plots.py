import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


def readCSV(path : str) -> pd.DataFrame:
    return pd.read_csv(path)


def drawAvgSystemQueueMMs(path : str, lamda : float, mu : float, servers : int) -> None:
    data = readCSV(path)

    args = list(data.Time)
    values = list(data.AvgQueue)

    waitingMean = plt.figure()
    ax3 = waitingMean.add_subplot()
    ax3.plot(args, values)
    ax3.axhline(y=values[-1:], color='r', linestyle='-')
    ax3.set_xlabel('Час симуляції')
    ax3.set_ylabel('Середня к-ть заявок в системі')

    ax3.set_title(f"Середня к-ть заявок в СМО з втратами M/M/{servers}")
    plt.figtext(.75, .75, f"\u03BB = {lamda}\n\u03BC = {mu}", {'size'   : 12})

    plt.savefig('/home/eshil/Programming/queueing_theory_golang/images/system_queue_mean_MMs.png')


def drawAvgSystemQueueMMn(path : str, lamda : float, mu : float, servers : int) -> None:
    data = readCSV(path)

    args = list(data.Time)
    values = list(data.AvgQueue)

    waitingMean = plt.figure()
    ax2 = waitingMean.add_subplot()
    ax2.plot(args, values)
    ax2.axhline(y=values[-1:], color='r', linestyle='-')
    ax2.set_xlabel('Час симуляції')
    ax2.set_ylabel('Середня к-ть заявок в системі')

    ax2.set_title(f"Середня к-ть заявок в СМО M/M/{servers}")
    plt.figtext(.75, .75, f"\u03BB = {lamda}\n\u03BC = {mu}", {'size'   : 12})

    plt.savefig('/home/eshil/Programming/queueing_theory_golang/images/system_queue_mean_MMn.png')


def drawWaitingMeanPlot(path : str, lamda : float, mu : float, servers : int) -> None:
    data = readCSV(path)

    args = list(data.Time)
    values = list(data.AvgWaitingTime)

    waitingMean = plt.figure()
    ax1 = waitingMean.add_subplot()
    ax1.plot(args, values)
    ax1.axhline(y=values[-1:], color='r', linestyle='-')
    ax1.set_xlabel('Час симуляції')
    ax1.set_ylabel('Середній час очікування в черзі')

    ax1.set_title(f"Середній час очікування в черзі СМО M/M/{servers}")
    plt.figtext(.75, .75, f"\u03BB = {lamda}\n\u03BC = {mu}", {'size'   : 12})

    plt.savefig('/home/eshil/Programming/queueing_theory_golang/images/waiting_mean.png')

def drawStatesProbsMMn(path : str, lamda : float, mu : float, servers : int) -> None:
    data = readCSV(path)

    customersCount = list(data.index)
    stateProbabilities = list(data.Probabilities)

    customers = [f"{i}" for i in customersCount]

    _, ax = plt.subplots()

    ax.bar(customers, stateProbabilities)

    ax.set_ylabel('Ймовірність стану')
    ax.set_xlabel('К-ть заявок в системі')
    ax.set_title(f'Ймовірності станів СМО M/M/{servers}')
    plt.figtext(.75, .75, f"\u03BB = {lamda}\n\u03BC = {mu}", {'size'   : 12})

    plt.savefig('/home/eshil/Programming/queueing_theory_golang/images/states_probs_MMn.png')


def drawStatesProbsMMs(path : str, lamda : float, mu : float, servers : int) -> None:
    data = readCSV(path)

    customersCount = list(data.index)
    stateProbabilities = list(data.Probabilities)

    customers = [f"{i}" for i in customersCount]

    _, ax4 = plt.subplots()

    ax4.bar(customers, stateProbabilities)

    ax4.set_ylabel('Ймовірність стану')
    ax4.set_xlabel('К-ть заявок в системі')
    ax4.set_title(f'Ймовірності станів СМО з втратами M/M/{servers}')
    plt.figtext(.75, .75, f"\u03BB = {lamda}\n\u03BC = {mu}", {'size'   : 12})

    plt.savefig('/home/eshil/Programming/queueing_theory_golang/images/states_probs_MMs.png')