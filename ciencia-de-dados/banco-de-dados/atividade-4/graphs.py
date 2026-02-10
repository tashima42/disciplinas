import matplotlib.pyplot as plt
import pandas as pd

acoes = pd.read_csv("./acoes.csv")

fig, ax1 = plt.subplots(figsize=(20, 10))
color = 'tab:blue'
ax1.set_xlabel('Numero de ocorrencias')
ax1.set_ylabel('Acao')
ax1.barh(acoes["acao"], acoes['total'], label='Ocorrencias')
ax1.tick_params(axis='y')
ax1.tick_params(axis='x')

# Title and legend
plt.title('Ocorrencias de acoes')

plt.savefig("acoes.eps", format='eps')
