
from Global import *
from math import *

import numpy as np
import matplotlib.pyplot as plt

Selections = []

def GetColourRepr(pixel):
    #https://en.wikipedia.org/wiki/HSL_and_HSV

    rgb = pixel[:3] #Slice off Alpha channel

    #Convert rgb from 0->255 to 0->1
    if len(rgb[rgb>1]):
        rgb = rgb/255

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

#Numpy array version of GetColourRepr
#Pixels are given as a np.array of each pixel
def GetColourReprNumpy(pixels):
    #https://en.wikipedia.org/wiki/HSL_and_HSV

    rgb = pixels[...,:3] #Slice off Alpha channel

    #Convert rgb from 0->255 to 0->1
    rgb[rgb[...,0]>1] = rgb[rgb[...,0]>1]/255
    rgb[rgb[...,1]>1] = rgb[rgb[...,1]>1]/255
    rgb[rgb[...,2]>1] = rgb[rgb[...,2]>1]/255

    # R,G,B=rgb
    R = rgb[...,0]
    G = rgb[...,1]
    B = rgb[...,2]
    # V=0 #Value
    # C=0 #Chroma
    # L=0 #Lightness
    H=np.zeros(len(R)) #Hue
    Sv=np.zeros(len(R)) #Saturation (Value)
    Sl=np.zeros(len(R)) #Saturation (Lightness)
    A=np.zeros(len(R)) #Alpha

    # A = pixels[3] if len(pixels)==4 else 1
    if pixels.shape[1]==4:
        A = pixels[...,3]

    # maxRGB = max(rgb)
    # minRGB = min(rgb)
    maxRGB = np.array(R)
    maxRGB[G>maxRGB] = G[G>maxRGB]
    maxRGB[B>maxRGB] = B[B>maxRGB]

    minRGB = np.array(R)
    minRGB[G<minRGB] = G[G<minRGB]
    minRGB[B<minRGB] = B[B<minRGB]

    V = np.array(maxRGB)
    C = V - minRGB 
    L = (maxRGB+minRGB)/2 

    # if (C==0): pass
    # elif (R>G and R>B): H = (0 + (G-B)/C)/6
    # elif (G>R and G>B): H = (2 + (B-R)/C)/6
    # elif (B>R and B>G): H = (4 + (R-G)/C)/6
    cond = (C!=0) & (R>G) & (R>B); H[cond] = (0 + (G[cond]-B[cond])/C[cond])/6
    cond = (C!=0) & (G>R) & (G>B); H[cond] = (2 + (G[cond]-B[cond])/C[cond])/6
    cond = (C!=0) & (B>R) & (B>G); H[cond] = (4 + (G[cond]-B[cond])/C[cond])/6

    # Sv = (0 if V==0 else C/V)
    cond = V!=0
    Sv[cond] = C[cond]/V[cond]

    # Sl = 0 if L==0 or L==1 else (V-L)/min(L,1-L)
    minL = np.array(L)
    cond = 1-L < minL; minL[cond] = 1-L[cond]
    cond = (L!=0) & (L!=1); Sl[cond] = (V[cond]-L[cond])/minL[cond]

    return np.array([R,G,B,V,C,L,H,Sv,Sl,A])



#value: float range[0,1]
#target: float range[0,1]
#tolerance: float range[0,1]
def TestTol(value, target, posTol, negTol):
    return value-target <= posTol and target-value <= negTol

def TestTolWrap(value, target, posTol, negTol):
    return (value-target)%1 <= posTol or (target-value)%1 <= negTol

#pixel: float[3] range[0,1]
def SelectPixel(colourRepr):
    R,G,B,V,C,L,H,Sv,Sl,A = colourRepr #GetColourRepr(pixel)

    passSelections = False

    for selection in Selections:
        rules = selection["Rules"]

        #Init to True since all rules must pass
        passRules = True

        for rule in rules:
            fmt = rule["ColourFormat"]
            fmtArgs = rule["Args"]
            tolPos = rule["Tol+"]
            tolNeg = rule["Tol-"]

            #Check the rule
            if "R" in fmt:
                passRules &= TestTol(R, fmtArgs[fmt.find("R")], tolPos, tolNeg)
            if "G" in fmt:
                passRules &= TestTol(G, fmtArgs[fmt.find("G")], tolPos, tolNeg)
            if "B" in fmt:
                passRules &= TestTol(B, fmtArgs[fmt.find("B")], tolPos, tolNeg)
            if "V" in fmt:
                passRules &= TestTol(V, fmtArgs[fmt.find("V")], tolPos, tolNeg)
            if "C" in fmt:
                passRules &= TestTol(C, fmtArgs[fmt.find("C")], tolPos, tolNeg)
            if "L" in fmt:
                passRules &= TestTol(L, fmtArgs[fmt.find("L")], tolPos, tolNeg)
            if "H" in fmt:
                passRules &= TestTolWrap(H, fmtArgs[fmt.find("H")], tolPos, tolNeg)
            if "Sv" in fmt:
                passRules &= TestTol(Sv, fmtArgs[fmt.find("Sv")], tolPos, tolNeg)
            if "Sl" in fmt:
                passRules &= TestTol(Sl, fmtArgs[fmt.find("Sl")], tolPos, tolNeg)
            if "A" in fmt:
                passRules &= TestTol(A, fmtArgs[fmt.find("A")], tolPos, tolNeg)
            if "S" in fmt:
                if "V" in fmt:
                    passRules &= TestTol(Sv, fmtArgs[fmt.find("S")], tolPos, tolNeg)
                if "L" in fmt:
                    passRules &= TestTol(Sl, fmtArgs[fmt.find("S")], tolPos, tolNeg)

        if selection["Negate"]:
            #Selections are and-ed
            passSelections &= not passRules
        else:
            #Selections are added together
            passSelections |= passRules

    if passSelections:
        return [1,1,1,1]
        # return pixel
    else:
        return [0,0,0,1]

#https://www.desmos.com/calculator/mqcfc2lfu0
def Factor(value, posTol, negTol):
    #Small fudge factor for ensureing tolerance of 0 still returns valid pixels
    posTol+=0.01
    negTol+=0.01
    smoothing = 1
    factor = 10/smoothing
    return min(
        Sigmoid(-factor*(value-posTol)),
        Sigmoid( factor*(value+negTol))
    )

def TestTolDist(value, target, posTol, negTol):
    return Factor(value-target, posTol, negTol)

def TestTolWrapDist(value, target, posTol, negTol):
    return Factor(sin(pi*(value-target)), posTol, negTol)


#pixel: float[3] range[0,1]
def SelectPixelDist(colourRepr):
    R,G,B,V,C,L,H,Sv,Sl,A = colourRepr #GetColourRepr(pixel)

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


def SelectImageSections(img, selections, method):
    global Selections
    height, width, depth = img.shape

    Selections = selections

    selPixelFunc = SelectPixel
    if method == "Dist":
        selPixelFunc = SelectPixelDist

    print("> Selecting Colours")
    #Convert to list of pixels, then pass to GetColourReprNumpy
    #Returns lists of each colourRepr
    colourRepr = GetColourReprNumpy(np.reshape(img,(-1,depth)))

    newimg = np.zeros([height, width, 4])
    r=0
    i=0
    ts=time.time()
    for rows in newimg:
        for pixel in rows:
            pixel[...] = selPixelFunc(colourRepr[...,i])
            i+=1
        r+=1
        # print(F"\033[A{(r/height)*100:.2f}%")
        ProgressBar(r, 0, height)
        TimeDisplay(ts)
    print()

    return newimg


