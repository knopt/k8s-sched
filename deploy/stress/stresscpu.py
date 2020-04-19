import argparse, subprocess, random
import scipy.stats as stats

def try_parse_num_array(string_arg):
    if string_arg[0] in ["[", "("]:
        string_arg = string_arg[1:]
    if string_arg[-1] in ["]", ")"]:
        string_arg = string_arg[:-1]

    nums = str.split(string_arg, ",")
    res = []
    for n in nums:
        try:
            res.append(float(n))
        except:
            return string_arg

    return res

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

def uniform_numbers(min, max, mu, sigma, cnt):
    while True:
        dist = stats.truncnorm((min - mu) / sigma, (max - mu) / sigma, loc=mu, scale=sigma)
        values = dist.rvs(cnt)

        if abs(values.mean() - mu) <= 0.05 * mu:
            print("returning with mean {:.2f}. wanted {}".format(values.mean(), mu))
            return values

def shuffle_array(arr):
    i = 0
    while i < len(arr) - 1:
        j = random.randint(i + 1, len(arr) - 1)
        arr[i], arr[j] = arr[j], arr[i]
        i += 1

def const_cpu_load(cpu_load, timeout):
    subprocess.run(["stress-ng", "--cpu", "1", "-l", str(int(100 * cpu_load)), "--timeout", "{}s".format(timeout)])

def random_cpu_load(cpu_load, timeout, window_size):
    steps = max(1, int(timeout/window_size))
    for _ in range(steps):
        load = random.randint(0, 100)
        print("running for {} with load {}".format(window_size, load))
        subprocess.run(["stress-ng", "--cpu", "1", "-l", str(load), "--timeout", "{}s".format(window_size)])

def uniform_cpu_load(cpu_load, min_load, max_load, sigma, timeout, window_size):
    steps = max(1, int(timeout/window_size))
    steps_done = 0

    if sigma == 0:
        granularity = 1
        loads = [min_load + float(i * granularity) for i in range(int((max_load - min_load) / granularity))]
    else:
        loads = uniform_numbers(min_load, max_load, cpu_load, sigma, min(1000, steps))

    while True:
        shuffle_array(loads)

        for load in loads:
            print("running for {} with load {:.2f}".format(window_size, load))
            subprocess.run(["stress-ng", "--cpu", "1", "-l", "{:.2f}".format(load), "--timeout", "{}s".format(window_size)])
            steps_done += 1
            if steps_done >= steps:
                return


def any_cpu_load(cpu_load, timeout, window_size, dist):
    total = 0
    for d in dist:
        total += d
    if total != 100:
        raise Exception("expected sum of dist to be 100")

    steps, steps_done = max(1, timeout/window_size), 0

    while True:
        for load in dist:
            subprocess.run(["stress-ng", "--cpu", "1", "-l", str(load), "--timeout", "{}s".format(window_size)])
            steps_done += 1
            if steps_done >= steps:
                return


parser = argparse.ArgumentParser()
parser.add_argument("--window-size", help="length of const load time window", type=str, default="1s")
parser.add_argument("--cpu-load", help="average CPU load", type=float, default=1)
parser.add_argument("--dist", help="load distribution funcion. One of const, uniform, random, any array of sum 100", type=str, required=True)
parser.add_argument("--timeout", help="total time. 0s if infinite", default="24h", type=str)
parser.add_argument("--sigma", help="standard dev, only with dist uniform", default=0, type=float)
parser.add_argument("--min", help="min cpu load", default=0, type=int)
parser.add_argument("--max", help="max cpu load", default=100, type=int)

args = parser.parse_args()

window_size = s2d(args.window_size)
timeout = s2d(args.timeout)
cpu_load = args.cpu_load
dist = args.dist
sigma = args.sigma

dist = try_parse_num_array(dist)

if dist == "const":
    const_cpu_load(cpu_load, timeout)
elif dist in ["uniform", "uni", "u"]:
    uniform_cpu_load(cpu_load, args.min, args.max, sigma, timeout, window_size)
elif dist in ["random", "rand", "r"]:
    random_cpu_load(cpu_load, timeout, window_size)
elif type(dist) == list:
    any_cpu_load(cpu_load, timeout, window_size, dist)
else:
    print("incorrect usage. wrong dist argument")
    exit(1)
