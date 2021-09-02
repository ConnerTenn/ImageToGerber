
from Global import *
import numpy as np

def CreateFile(filename):
    filename+=".gbr"
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

def WriteHeader(file):
    file.write("G04 Created by ImageToGerber tool*\n")
    file.write("%FSLAX46Y46*%\n") #Number Format
    file.write("%MOMM*%\n") #Millimeters
    file.write("%LPD*%\n") #Dark polarity


def GeneratePixelated(img, filename):
    file = CreateFile(filename)
    WriteHeader(file)
    file.write("%ADD10R,1X1*%") #Rectangle Object
    file.write("D10*\n") #Use Rectangle Object

    height, width = img.shape

    print("> Writing Gerber")
    for y in range(height):
        for x in range(width):
            if img[y][x]:
                file.write(F"X{NumRepr(x)}Y{NumRepr(height-y)}D03*\n")
        ProgressBar(y, 0, height-1)
    print()

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


