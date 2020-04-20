import argparse, subprocess, random
import scipy.stats as stats

def s2d(duration_string):
    last_char = -1
    total_seconds = 0
    for i in range(len(duration_string)):
        if duration_string[i] not in ["h", "H", "s", "S", "d", "D", "m", "M"]:
            continue

        number = int(duration_string[last_char + 1:i])
        if duration_string[i] in ["s", "S"]:
            total_seconds += number
        if duration_string[i] in ["m", "M"]:
            total_seconds += number * 60
        if duration_string[i] in ["h", "H"]:
            total_seconds += number * 60 * 60
        if duration_string[i] in ["d", "D"]:
            total_seconds += number * 60 * 60 * 60
            
    return total_seconds


def shuffle_array(arr):
    i = 0
    while i < len(arr) - 1:
        j = random.randint(i + 1, len(arr) - 1)
        arr[i], arr[j] = arr[j], arr[i]
        i += 1

def uniform_numbers(min, max, mu, sigma, cnt):
    while True:
        dist = stats.truncnorm((min - mu) / sigma, (max - mu) / sigma, loc=mu, scale=sigma)
        values = dist.rvs(cnt)

        if abs(values.mean() - mu) <= 0.03 * mu:
            return values

def const_cpu_load(cpu_load, timeout):
    return [cpu_load]

def random_cpu_load(min_load, max_load, num):
    return [round(random.uniform(min_load, max_load), 3) for i in range(num)]

def uniform_cpu_load(cpu_load, min_load, max_load, sigma, num):
    if sigma == 0:
        loads = []
        for i in range(num):
            delta = (max_load - min_load) * i / float(num)
            loads.append(min_load + delta)
            
        shuffle_array(loads)
    else:
        loads = uniform_numbers(min_load, max_load, cpu_load, sigma, num)
    
    loads = [round(f, 3) for f in loads]

    return loads

parser = argparse.ArgumentParser()
parser.add_argument("--cpu-load", help="average CPU load", type=float, default=1)
parser.add_argument("--dist", help="load distribution funcion. One of const, uniform, random", type=str, required=True)
parser.add_argument("--sigma", help="standard dev, only with dist uniform", default=0, type=float)
parser.add_argument("--min", help="min cpu load", default=0, type=int)
parser.add_argument("--max", help="max cpu load", default=100, type=int)
parser.add_argument("--num", help="total num", default=500, type=int)

args = parser.parse_args()

cpu_load = args.cpu_load
dist = args.dist
sigma = args.sigma
num = args.num

if dist == "const":
    print(const_cpu_load(cpu_load))
elif dist in ["uniform", "uni", "u"]:
    print(uniform_cpu_load(cpu_load, args.min, args.max, sigma, num))
elif dist in ["random", "rand", "r"]:
    print(random_cpu_load(cpu_load, timeout, window_size))
else:
    print("incorrect usage. wrong dist argument")
    exit(1)
