#!/bin/python3

import sys
import matplotlib.pyplot as plt
import ImageProcessing

TERM_RESET =   "\033[m"
TERM_RED =     "\033[1;31m"
TERM_GREEN =   "\033[1;32m"
TERM_YELLOW =  "\033[1;33m"
TERM_BLUE =    "\033[1;34m"
TERM_MAGENTA = "\033[1;35m"
TERM_CYAN =    "\033[1;36m"
TERM_WHITE =   "\033[1;37m"


def GetOptions():
    options = {}
    argv = iter(sys.argv[1:])
    for arg in argv:
        #Config file
        if arg=="-c":
            try:
                options["ConfigFilename"] = next(argv)
            except:
                print(TERM_RED+"Error: -c must be followed by a filename"+TERM_RESET)
        #Default (Last) argument
        else:
            options["ImageFilename"] = arg

    return options

options = GetOptions()
if not "ImageFilename" in options:
    print(TERM_RED+"Error: An image file must be specified"+TERM_RESET)

img = plt.imread(options["ImageFilename"])

img = ImageProcessing.EdgeDetection(img)

plt.imsave("TestOutput.jpg", img)
