
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
    return value-target <= posTol and target-value <= negTol

def TestTolWrap(value, target, posTol, negTol):
    return (value-target)%1 <= posTol or (target-value)%1 <= negTol

#https://www.desmos.com/calculator/mqcfc2lfu0
def Factor(value, posTol, negTol):
    smoothing = 1
    factor = 30/smoothing
    return min(
        Sigmoid(-factor*(value-posTol)),
        Sigmoid( factor*(value+negTol))
    )

def TestTolDist(value, target, posTol, negTol):
    return Factor(value-target, posTol, negTol)

def TestTolWrapDist(value, target, posTol, negTol):
    return Factor(sin(pi*(value-target)), posTol, negTol)


#pixel: float[3] range[0,1]
def SelectPixel(pixel):
    R,G,B,V,C,L,H,Sv,Sl,A = GetColourRepr(pixel)

    selectionDist = 10**100

    for selection in Selections:
        rules = selection["Rules"]

        #Init to True since all rules must pass
        dist = 10**100

        for rule in rules:
            fmt = rule["ColourFormat"]
            fmtArgs = rule["Args"]
            tolPos = rule["Tol+"]
            tolNeg = rule["Tol-"]

            #Check the rule
            if "R" in fmt:
                dist = min(dist, TestTolDist(R, fmtArgs[fmt.find("R")], tolPos, tolNeg))
            if "G" in fmt:
                dist = min(dist, TestTolDist(G, fmtArgs[fmt.find("G")], tolPos, tolNeg))
            if "B" in fmt:
                dist = min(dist, TestTolDist(B, fmtArgs[fmt.find("B")], tolPos, tolNeg))
            if "V" in fmt:
                dist = min(dist, TestTolDist(V, fmtArgs[fmt.find("V")], tolPos, tolNeg))
            if "C" in fmt:
                dist = min(dist, TestTolDist(C, fmtArgs[fmt.find("C")], tolPos, tolNeg))
            if "L" in fmt:
                dist = min(dist, TestTolDist(L, fmtArgs[fmt.find("L")], tolPos, tolNeg))
            if "H" in fmt:
                dist = min(dist, TestTolWrapDist(H, fmtArgs[fmt.find("H")], tolPos, tolNeg))
            if "Sv" in fmt:
                dist = min(dist, TestTolDist(Sv, fmtArgs[fmt.find("Sv")], tolPos, tolNeg))
            if "Sl" in fmt:
                dist = min(dist, TestTolDist(Sl, fmtArgs[fmt.find("Sl")], tolPos, tolNeg))
            if "A" in fmt:
                dist = min(dist, TestTolDist(A, fmtArgs[fmt.find("A")], tolPos, tolNeg))
            if "S" in fmt:
                if "V" in fmt:
                    dist = min(dist, TestTolDist(Sv, fmtArgs[fmt.find("S")], tolPos, tolNeg))
                if "L" in fmt:
                    dist = min(dist, TestTolDist(Sl, fmtArgs[fmt.find("S")], tolPos, tolNeg))

        if selection["Negate"]:
            #Selections are and-ed
            # passSelections &= not passRules
            selectionDist = min(1-dist, selectionDist)
        else:
            #Selections are added together
            # passSelections |= passRules
            selectionDist = min(dist, selectionDist)

    return [selectionDist]*3 +[1]


def SelectImageSections(imagefile, selections):
    global Selections
    img = plt.imread(imagefile)
    img = img

    height, width, depth = img.shape

    Selections = selections

    print(F"Image Resolution: {width}x{height} [{depth}]")
    print("> Selecting Colours")
    r=0
    for rows in img:
        for pixel in rows: pixel[...] = SelectPixel(list(pixel))
        r+=1
        # print(F"\033[A{(r/height)*100:.2f}%")
        ProgressBar(r, 0, height)
    print()

    return img


