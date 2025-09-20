import numpy as np
import pandas as pd


dfX_testToronto = pd.read_csv('data/X_testToronto.csv')

random_predictions = np.random.choice([0, 1], size=dfX_testToronto.shape[0])

dfX_testToronto['destaque'] = random_predictions

# Please pay attention: Kaggle demands that you add "business_id" and "destaque" column headers
dfX_testToronto.to_csv("data/randomTeste.csv", columns=['business_id','destaque'],index=False)
