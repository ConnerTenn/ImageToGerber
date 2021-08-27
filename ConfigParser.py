
from Global import *


def ParseFilename(line):
    tokens = line.split("\"")

    if len(tokens) != 3:
        print("{TERM_RED}Malformed filepath specifier")

    path = tokens[1]
    print(F">> {path} <<")

def ParseSelection(line):
    print(line)




def ParseConfig(filename):
    config=[]

    cfgfile=None
    try:
        cfgfile=open(filename)
    except:
        print(F"{TERM_RED}Failed to open file {filename}{TERM_RESET}")


    for line in cfgfile:
        line = line.split("#")[0].rstrip()

        if len(line):

            if line[0].isspace():
                line = line.strip()

                ParseSelection(line)
            else:
                ParseFilename(line)
