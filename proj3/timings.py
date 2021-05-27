import subprocess
from statistics import median
import pandas as pd
import matplotlib.pyplot as plt

#Create the parameter tuples that we want to run the tests on
tests = [400] * 5 + [800] * 5 +[1000] * 5
threads = [1,2,4,8,10] * 3
# image_set = ['small','mixture','big']
# image_set = [[x] * 6 for x in image_set] * 2
# image_set = [item for sublist in image_set for item in sublist] #flatten it
combos = list(zip(tests,threads))
print(len(combos))
print(combos)

#Collect our timings for each tuple of parameters
times = []
for parameters in combos:
  # print(parameters)
  thread_times = [] #collect the runtimes for this set of parameters
  for i in range(5): #Take the median of five runs
    p = subprocess.run(f"time go run main.go {parameters[1]} {parameters[0]} 0",shell=True,capture_output=True)
    #'\nreal\t0m0.003s\nuser\t0m0.001s\nsys\t0m0.001s\n'
    # real = p.stderr.decode()
    # print(real.split("\t"))
    real = p.stderr.decode().split("\n")[1].split("\t")[1]
    minutes = float(real.split("m")[0])
    seconds = float(real.split("m")[1].replace("s",""))
    wallTime = 60*minutes+seconds
    thread_times.append(wallTime)
    print(parameters,wallTime)
  times.append(median(thread_times))

# Split our final times into chunks of 6 - these 6 points will represent a line on a graph
def chunks(lst, n):
  for i in range(0, len(lst), n):
    yield lst[i:i + n]
groupings = chunks(times,5)
speedups = [group[0]/x for group in groupings for x in group]

# Make a df with our data
df = pd.DataFrame({'size': tests , 'threads': threads , 'speedup': speedups})

#Suppress annoying matplotlib warnings on linux cluster
import warnings
warnings.filterwarnings("ignore")

#For each test size
# for this_size in [400,800,1000]:
#create a local df with this size only
# local_df = df[(df['size'] == this_size) & (df['threads'] > 1)]
#create a plt
fig,ax = plt.subplots()
#create a line on the plt for each image_set; x=threads, y=speedup
for size in set(df['size']):
  ax.plot('threads','speedup', data = df[df['size']==size], label=f"size:{size} Documents")
#format and save the graph
ax.set_xlabel("threads")
ax.set_ylabel("speedup")
ax.legend(loc='best')
ax.set_title(f"Speedup")
plt.savefig(f"Speedup.png")
# plt.show()