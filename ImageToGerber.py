#!/bin/python3

from Global import *

import sys
import os
import numpy as np
import matplotlib.pyplot as plt
import ImageProcessing

import ColourSelection
import ConfigParser
import NaiveTrace
import GerberWriter


def ShowHelp():
    print( \
"""\
Usage:
./ImageToGerber.py [-h] [-c <config file>] [-m <method>] <image file>
    -h --help       Show the help menu. 
                    This argument is optional and will cause the program to 
                    exit immediately.
    -c --config     Specify the config file to be used for processing the image.
                    If none is specified, a default config will be used.
                    This argument is optional.
    -m --method     <Dist/Blur/Blocky>
                    Select which pixel selection and processing method should be used.
                    Defaults to Dist
    <image file>    The image file to process. This argument is reqired.
""")
    exit()

def GetOptions():
    options = {"ConfigFilename":None, "Method":"Dist"}
    argv = iter(sys.argv[1:])
    for arg in argv:
        #Help menu
        if arg=="-h" or arg=="--help":
            ShowHelp()
        #Config file
        if arg=="-c" or arg=="--config":
            try:
                options["ConfigFilename"] = next(argv)
            except:
                Error("-c must be followed by a filename")
                ShowHelp()
        #Processing method
        if arg=="-m" or arg=="--method":
            try:
                options["Method"] = next(argv)
            except:
                Error("-m must be followed by a method selection")
                ShowHelp()

            if not options["Method"] in ["Dist", "Blur", "Blocky"]:
                Error("Invalid method selection")
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
    config = ConfigParser.ParseConfig(options["ConfigFilename"])


srcimg=None
try:
    srcimg = plt.imread(options["ImageFilename"])
except:
    Error("Failed to open Image")

dim = srcimg.shape
imgDim = Object()
imgDim.Width=dim[1]
imgDim.Height=dim[0]
imgDim.Depth=dim[2]
print(F"Image Resolution: {imgDim.Width}x{imgDim.Height} [{imgDim.Depth}]")

gerberDim = Object()
gerberDim.Width=None
gerberDim.Height=None
if "Width" in config["Options"]:
    gerberDim.Width = config["Options"]["Width"]
if "Height" in config["Options"]:
    gerberDim.Height = config["Options"]["Height"]

if gerberDim.Width==None:
    gerberDim.Width = imgDim.Width*gerberDim.Height/imgDim.Height
elif gerberDim.Height==None:
    gerberDim.Height = imgDim.Height*gerberDim.Width/imgDim.Width
else:
    if abs(gerberDim.Width/imgDim.Width - gerberDim.Height/imgDim.Height) > 0.01:
        Error("Specified gerber aspect ratio does not match image aspect ratio")



for i, proc in enumerate(config["Processes"]):
    print()
    print(F"{TERM_BLUE}=== Processing {i} ==={TERM_RESET}")
    outPath = proc["Path"]
    print(F"{TERM_WHITE}>> {outPath}-{proc['Type']}{TERM_RESET}")
    print()

    splitPath = outPath.rpartition("/")
    # filename = splitPath[-1]
    
    if splitPath[1] == "/":
        directory = splitPath[0]
        if not os.path.isdir(directory):
            os.mkdir(directory)

    img = ColourSelection.SelectImageSections(srcimg, proc["Selections"], options["Method"])
    if options["Method"] == "Blur":
        img = ImageProcessing.GaussianBlur(img)
    plt.imsave(F"{outPath}_SelectedRegions-{proc['Type']}.png", img)

    #Convert to scalar array (taking red channel only)
    img = np.dot(img[...,:1], [1])

    if proc["Method"] == "Fill":
        GerberWriter.GeneratePixelatedFillLines(img, gerberDim, outPath, proc["Type"])
        # octree = ImageProcessing.GenerateOctree(img)
        # GerberWriter.GeneratePixelatedOctree(octree, outPath, img.shape)
    elif proc["Method"] == "Trace":
        segments = NaiveTrace.LineDetection(img)
        GerberWriter.GenerateTraceFromSegments(segments, imgDim, gerberDim, outPath, proc["Type"])
        # lineloops = NaiveTrace.StitchSegments(segments)
        # NaiveTrace.PlotLineLoops(lineloops)
        # GerberWriter.GenerateTrace(lineloops, img.shape, outPath, proc["Type"])

        # img_edge = ImageProcessing.EdgeDetection(img)
        # plt.imsave(outPath+"_EdgeDetection.png", img_edge)

        # img_hough = ImageProcessing.LineDetection(img_edge)
        # plt.imsave(outPath+"_LineDetection.png", img_hough)
    else:
        Error(F"Invalid process selected: {proc['Method']}")
