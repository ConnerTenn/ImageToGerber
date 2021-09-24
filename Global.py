
import sys
import time

TERM_RESET =   "\033[m"
TERM_RED =     "\033[1;31m"
TERM_GREEN =   "\033[1;32m"
TERM_YELLOW =  "\033[1;33m"
TERM_BLUE =    "\033[1;34m"
TERM_MAGENTA = "\033[1;35m"
TERM_CYAN =    "\033[1;36m"
TERM_WHITE =   "\033[1;37m"

def Error(msg):
    print(F"{TERM_RED}Error: {msg}{TERM_RESET}")
    exit(1)

def Assert(value, expected, msg=""):
    if (value==expected):
        print(F"{TERM_GREEN} PASS: {msg}{TERM_RESET}")
    else:
        print(F"{TERM_RED} FAIL: {msg}")
        print(F" - - - ({value}) != Expected ({expected}){TERM_RESET}")

def PrintList(obj, level=""):
    print(level + "[")
    for item in obj:
        if type(item)==dict:
            PrintDict(item, level+"| ")
        elif type(item)==list:
            PrintList(item, level+"| ")
        else:
            print(level + F"{item}")
    print(level + "]")

def PrintDict(obj, level=""):
    print(level + "{")
    for key, value in obj.items():
        # print(key, value)
        print(level + F"{key} : ", end="")
        if type(value)==dict:
            print()
            PrintDict(value, level+"! ")
        elif type(value)==list:
            print()
            PrintList(value, level+"! ")
        else:
            print(F"{value}")
    print(level + "}")


#val in range[minimum,maximum]
#Must be followed by an empty print before printing anything after the progress bar
def ProgressBar(val, minimum, maximum):
    barSize = 50
    progress = (val-minimum)/maximum
    numbars = int(barSize*progress)
    bars = "="*numbars + " "*(barSize-numbars)
    print(F"\033[1G[{bars}] {100*progress:5.1f}%", end="", flush=True)

def TimeDisplay(ts):
    elapsed = time.time()-ts
    seconds = elapsed%60
    minutes = int(elapsed/60)
    print(F"   Elapsed time: {minutes:02d}m {seconds:05.2f}s", end="", flush=True)


def Sigmoid(x):
    return 1 / (1 + 2**(-x))

def Norm(x, factor):
    return 2**(-factor*(x**2))
