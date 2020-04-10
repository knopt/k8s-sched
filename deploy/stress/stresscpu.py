import argparse, subprocess, random

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

def d2s(duration_string):
    last_char = -1
    total_seconds = 0
    for i in range(len(duration_string)):
        if duration_string[i] not in ["h", "H", "s", "S", "d", "D", "m", "M"]:
            continue

        number = int(duration_string[last_char + 1:i])
        if i in ["s", "S"]:
            total_seconds += number
        if i in ["m", "M"]:
            total_seconds += number * 60
        if i in ["h", "H"]:
            total_seconds += number * 60 * 60
        if i in ["d", "D"]:
            total_seconds += number * 60 * 60 * 60
    return total_seconds


def const_cpu_load(cpu_load, timeout):
    subprocess.run("stress-ng", "--cpu", "1", "-l", str(100 * cpu_load), "--timeout", timeout)

def random_cpu_load(cpu_load, timeout, window_size):
    steps = max(1, timeout/window_size)
    for _ in range(steps):
        subprocess.run("stress-ng", "--cpu", "1", "-l", str(random.randint(0, 100)), "--timeout", "{}s".format(window_size))

def uniform_cpu_load(cpu_load, window_size):
    pass

def any_cpu_load(cpu_load, window_size, dist):
    total = 0
    for d in dist:
        total += d
    if total != 100:
        raise Exception("expected sum of dist to be 100")


    pass

parser = argparse.ArgumentParser()
parser.add_argument("--window-size", help="length of const load time window", type=str, default="1s")
parser.add_argument("--cpu-load", help="average CPU load", type=float, default=1)
parser.add_argument("--dist", help="load distribution funcion. One of const, uniform, random, any array of sum 100", type=str, required=True)
parser.add_argument("--timeout", help="total time. 0 if infinite", default="24h", type=str)

args = parser.parse_args()

window_size = d2s(args.window_size)
timeout = args.timeout
cpu_load = args.cpu_load
dist = args.dist

dist = try_parse_num_array(dist)

if dist == "const":
    pass
elif dist in ["uniform", "uni", "u"]:
    pass
elif dist in ["random", "rand", "r"]:
    pass
elif type(dist) == list:
    any_cpu_load(cpu_load, window_size, dist)
else:
    print("incorrect usage. wrong dist argument")
    exit(1)
