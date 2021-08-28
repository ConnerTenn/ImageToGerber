#!/bin/python3

from Global import *

import sys
import matplotlib.pyplot as plt
import ImageProcessing

from ColourSelection import *
from ConfigParser import *


def ShowHelp():
    print( \
"""\
Usage:
./ImageToGerber.py [-h] [-c <config file>] <image file>
    -h --help       Show the help menu. 
                    This argument is optional and will cause the program to 
                    exit immediately.
    -c --config     Specify the config file to be used for processing the image.
                    If none is specified, a default config will be used.
                    This argument is optional.
    <image file>    The image file to process. This argument is reqired.
""")
    exit()

def GetOptions():
    options = {}
    argv = iter(sys.argv[1:])
    for arg in argv:
        #Help menu
        if arg=="-h" or arg=="--helph":
            ShowHelp()
        #Config file
        if arg=="-c" or arg=="--config":
            try:
                options["ConfigFilename"] = next(argv)
            except:
                Error("-c must be followed by a filename")
                ShowHelp()

        #Default (Last) argument
        else:
            options["ImageFilename"] = arg

    return options


options = GetOptions()
if not "ImageFilename" in options:
    Error("An image file must be specified")
    ShowHelp()

config = None

if "ConfigFilename" in options:
    config = ParseConfig(options["ConfigFilename"])

try:
    img = plt.imread(options["ImageFilename"])
except:
    Error("Failed to open Image")

img = ConvertImage(options["ImageFilename"], config)
plt.imsave("SelectedRegions.png", img)

img_edge = ImageProcessing.EdgeDetection(img)
plt.imsave("EdgeDetection.png", img_edge)

img_hough = ImageProcessing.LineDetection(img_edge)
plt.imsave("LineDetection.png", img_hough)
