import plots

STATES_PROBABILITIES_MMn_PATH = "/home/eshil/Programming/queueing_theory_golang/data/statesProbsMMn.csv"
STATES_PROBABILITIES_MMs_PATH = "/home/eshil/Programming/queueing_theory_golang/data/statesProbsMMs.csv"
WAITING_MEAN_PATH = "/home/eshil/Programming/queueing_theory_golang/data/AvgWaiting.csv"
SYSTEM_QUEUE_MEAN_MMn_PATH = "/home/eshil/Programming/queueing_theory_golang/data/AvgSystemQueueMMn.csv"
SYSTEM_QUEUE_MEAN_MMs_PATH = "/home/eshil/Programming/queueing_theory_golang/data/AvgSystemQueueMMs.csv"


def run(lamda: float, mu: float, servers: int) -> None:
    plots.drawAvgSystemQueueMMn(SYSTEM_QUEUE_MEAN_MMn_PATH, lamda, mu, servers)
    plots.drawWaitingMeanPlot(WAITING_MEAN_PATH, lamda, mu, servers)
    plots.drawStatesProbsMMn(STATES_PROBABILITIES_MMn_PATH, lamda, mu, servers)
    plots.drawStatesProbsMMs(STATES_PROBABILITIES_MMs_PATH, lamda, mu, servers)
    plots.drawAvgSystemQueueMMs(SYSTEM_QUEUE_MEAN_MMs_PATH, lamda, mu, servers)
    

if __name__ == "__main__":
    run(lamda=1.0, mu=20.0, servers=1)