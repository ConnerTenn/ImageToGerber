
#https://www.ucamco.com/en/guest/downloads/gerber-format

from Global import *
import numpy as np

def CreateFile(filename, gerberType):
    filename+=F"-{gerberType}.gbr"
    try:
        return open(filename, "w")
    except:
        Error(F"Failed to create file \"{filename}\"")
        return None

def FinishFile(file):
    file.write("M02*")
    file.close()

def NumRepr(num):
    integer = int(num)
    fraction = num
    return F"{num:011.6f}".replace(".","")

def WriteHeader(file, gerberType):
    file.write("%TF.GenerationSoftware,ImageToGerber,Design,1.1.0*%\n")
    file.write("%TF.SameCoordinates,Original*%\n")
    #Copper
    if gerberType=="F_Cu":
        file.write("%TF.FileFunction,Copper,L1,Top*%\n")
        file.write("%TF.FilePolarity,Positive*%\n")
    elif gerberType=="B_Cu":
        file.write("%TF.FileFunction,Copper,L2,Bot*%\n")
        file.write("%TF.FilePolarity,Positive*%\n")
    #Solder Mask
    elif gerberType=="F_Mask":
        file.write("%TF.FileFunction,Soldermask,Top*%\n")
        file.write("%TF.FilePolarity,Negative*%\n")
    elif gerberType=="B_Mask":
        file.write("%TF.FileFunction,Soldermask,Bot*%\n")
        file.write("%TF.FilePolarity,Negative*%\n")
    #Silk screen
    elif gerberType=="F_SilkS":
        file.write("%TF.FileFunction,Legend,Top*%\n")
        file.write("%TF.FilePolarity,Positive*%\n")
    elif gerberType=="B_SilkS":
        file.write("%TF.FileFunction,Legend,Bot*%\n")
        file.write("%TF.FilePolarity,Positive*%\n")
    #Edge Cuts
    elif gerberType=="Edge_Cuts":
        file.write("%TF.FileFunction,Profile,NP*%\n")

    file.write("%FSLAX46Y46*%\n") #Number Format
    file.write("%MOMM*%\n") #Millimeters
    file.write("%LPD*%\n") #Dark polarity
    file.write("G04 == End Header ==*\n")


def GeneratePixelated(img, filename):
    file = CreateFile(filename)
    WriteHeader(file)
    file.write("%ADD10R,1X1*%\n") #Rectangle Object
    file.write("D10*\n") #Use Rectangle Object

    height, width = img.shape

    print("> Writing Gerber")
    for y in range(height):
        for x in range(width):
            if img[y][x]:
                file.write(F"X{NumRepr(x)}Y{NumRepr(height-y)}D03*\n")
        ProgressBar(y, 0, height-1)
    print()

def GeneratePixelatedFillLines(img, gerberDim, filename, gerberType):
    file = CreateFile(filename, gerberType)
    WriteHeader(file, gerberType)
    height, width = img.shape
    scaleWidth = gerberDim.Width/width
    scaleHeight = gerberDim.Height/height

    lines = 0
    rectsize = 0

    def DrawLine(x, rectsize, line):
        file.write(F"%ADD{line+10}R,{scaleWidth*rectsize:.6f}X{scaleHeight:.6f}*%\n") #Rectangle Object
        file.write(F"D{line+10}*\n") #Use Rectangle Object
        file.write(F"X{NumRepr(scaleWidth*(x-rectsize/2))}Y{NumRepr(scaleHeight*(height-y+1/2))}D03*\n") #Place rectangle

    print("> Writing Gerber")
    ts=time.time()
    for y in range(height):
        for x in range(width):
            if img[y][x]:
                rectsize += 1
            #Handle line
            elif rectsize>0:
                DrawLine(x, rectsize, lines)
                rectsize = 0
                lines+=1
        #Handle last line
        if rectsize>0:
            DrawLine(width-1, rectsize, lines)
            rectsize = 0
            lines+=1
        ProgressBar(y, 0, height-1)
        TimeDisplay(ts)

    print()
    FinishFile(file)

def GeneratePixelatedOctree(octree, filename, dim):
    file = CreateFile(filename)
    WriteHeader(file)

    sizes = 0
    depths = {}

    for item in octree:
        if not item[0] in depths:
            depths |= {item[0]:10+sizes}
            scale = min(dim)/(2**item[0])
            file.write(F"%ADD{10+sizes}R,{scale:.6f}X{scale:.6f}*%\n") #Rectangle Object
            sizes+=1

    for item in octree:
        file.write(F"D{depths[item[0]]}*\n") #Use Rectangle Object
        s=min(dim)/(2**item[0])
        x=item[1] + s/2
        y=item[2] + s/2
        file.write(F"X{NumRepr(x)}Y{NumRepr(dim[1]-y)}D03*\n")


def GenerateTraceFromLoops(lineloops, dim, filename, gerberType):
    file = CreateFile(filename, gerberType)
    WriteHeader(file, gerberType)

    file.write("G01*\n") #Linear interpolation mode
    # file.write("%TA.AperFunction,Profile*%")
    file.write("%ADD10C,0.200000*%") #Circle aperture
    # file.write("%TD*%")
    file.write("D10*") #Use aperture

    height, width = dim

    for loop in lineloops:
        file.write(F"X{NumRepr(loop[0][0])}Y{NumRepr(height-loop[0][1])}D02*\n")
        for point in loop[1:]:
            file.write(F"X{NumRepr(point[0])}Y{NumRepr(height-point[1])}D01*\n")
    print()
    FinishFile(file)

def GenerateTraceFromSegments(segments, imgDim, gerberDim, filename, gerberType):
    file = CreateFile(filename, gerberType)
    WriteHeader(file, gerberType)

    scaleWidth = gerberDim.Width/imgDim.Width
    scaleHeight = gerberDim.Height/imgDim.Height

    file.write("G01*\n") #Linear interpolation mode
    # file.write("%TA.AperFunction,Profile*%")
    file.write("%ADD10C,0.200000*%") #Circle aperture
    # file.write("%TD*%")
    file.write("D10*") #Use aperture

    print("> Writing Gerber")
    ts=time.time()
    for i, line in enumerate(segments):
        file.write(F"X{NumRepr(scaleWidth*line[0][0])}Y{NumRepr(scaleHeight*(imgDim.Height-line[0][1]))}D02*\n")
        file.write(F"X{NumRepr(scaleWidth*line[1][0])}Y{NumRepr(scaleHeight*(imgDim.Height-line[1][1]))}D01*\n")
        ProgressBar(i, 0, len(segments)-1)
        TimeDisplay(ts)
    print()

    FinishFile(file)


