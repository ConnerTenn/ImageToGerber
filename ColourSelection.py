
from Global import *
from math import *

import numpy as np
import matplotlib.pyplot as plt

Selections = []

def GetColourRepr(pixel):
    #https://en.wikipedia.org/wiki/HSL_and_HSV

    rgb = pixel[:3] #Slice off Alpha channel

    R,G,B=rgb
    V=0 #Value
    C=0 #Chroma
    L=0 #Lightness
    H=0 #Hue
    Sv=0 #Saturation (Value)
    Sl=0 #Saturation (Lightness)
    A=0 #Alpha

    A = pixel[3] if len(pixel)==4 else 1

    maxRGB = max(rgb)
    minRGB = min(rgb)

    V = maxRGB 
    C = V - minRGB 
    L = (maxRGB+minRGB)/2 

    if (C==0): pass
    elif (R>G and R>B): H = (0 + (G-B)/C)/6
    elif (G>R and G>B): H = (2 + (B-R)/C)/6
    elif (B>R and B>G): H = (4 + (R-G)/C)/6

    Sv = (0 if V==0 else C/V)

    Sl = 0 if L==0 or L==1 else (V-L)/min(L,1-L)

    return R,G,B,V,C,L,H,Sv,Sl,A




#value: float range[0,1]
#target: float range[0,1]
#tolerance: float range[0,1]
def TestTol(value, target, posTol, negTol):
    return (value-target) <= posTol/2 and (target-value) <= negTol/2

def TestTolWrap(value, target, posTol, negTol):
    return (value-target)%1 <= posTol/2 or (target-value)%1 <= negTol/2



#pixel: float[3] range[0,1]
def ConvertPixel(pixel):
    R,G,B,V,C,L,H,Sv,Sl,A = GetColourRepr(pixel)

    # if (TestTolWrap(H, 0.1, 0.2, 0.2) and
    #     TestTol(Sv, 0.9, 0.4, 0.4) and
    #     TestTol(V, 0.7, 0.4, 0.4)):
    #     # return [1,1,1,1]
    #     return pixel
    # return [0,0,0,1]

    for selection in Selections:
        rules = selection["Rules"]

        for rule in rules:
            fmt = rule["ColourFormat"]
            fmtArgs = rule["Args"]
            tolPos = rule["Tol+"]
            tolNeg = rule["Tol-"]

            res = True

            if "R" in fmt:
                res &= TestTol(R, fmtArgs[fmt.find("R")], tolPos, tolNeg)
            if "G" in fmt:
                res &= TestTol(G, fmtArgs[fmt.find("G")], tolPos, tolNeg)
            if "B" in fmt:
                res &= TestTol(B, fmtArgs[fmt.find("B")], tolPos, tolNeg)
            if "V" in fmt:
                res &= TestTol(V, fmtArgs[fmt.find("V")], tolPos, tolNeg)
            if "C" in fmt:
                res &= TestTol(C, fmtArgs[fmt.find("C")], tolPos, tolNeg)
            if "L" in fmt:
                res &= TestTol(L, fmtArgs[fmt.find("L")], tolPos, tolNeg)
            if "H" in fmt:
                res &= TestTolWrap(H, fmtArgs[fmt.find("H")], tolPos, tolNeg)
            if "Sv" in fmt:
                res &= TestTol(Sv, fmtArgs[fmt.find("Sv")], tolPos, tolNeg)
            if "Sl" in fmt:
                res &= TestTol(Sl, fmtArgs[fmt.find("Sl")], tolPos, tolNeg)
            if "A" in fmt:
                res &= TestTol(A, fmtArgs[fmt.find("A")], tolPos, tolNeg)

            if not res:
                return [0,0,0,1]

    # return [1,1,1,1]
    return pixel


def ConvertImage(imagefile, config):
    global Selections
    img = plt.imread(imagefile)
    img = img

    height, width, depth = img.shape

    Selections = config["Processes"][0]["Selections"]

    print(F"Image Resolution: {width}x{height} [{depth}]")
    print("> Selecting Colours")
    print()
    r=0
    for rows in img:
        for pixel in rows: pixel[...] = ConvertPixel(list(pixel))
        r+=1
        print(F"\033[A{(r/height)*100:.2f}%")

    return img


