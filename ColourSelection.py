
from math import *

import numpy as np
import matplotlib.pyplot as plt


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
    elif (maxRGB==R): H = (0 + (G-B)/C)/6
    elif (maxRGB==G): H = (2 + (B-R)/C)/6
    elif (maxRGB==B): H = (4 + (R-G)/C)/6

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

    if (TestTolWrap(H, 0.1, 0.2, 0.2) and
        TestTol(Sv, 0.9, 0.4, 0.4) and
        TestTol(V, 0.7, 0.4, 0.4)):
        # return [1,1,1,1]
        return pixel
    return [0,0,0,1]


def ConvertImage(imagefile):
    img = plt.imread(imagefile)
    img = img

    height, width, depth = img.shape

    print(F"Image Resolution: {width}x{height} [{depth}]")
    print("> Selecting Colours")
    print()
    r=0
    for rows in img:
        for pixel in rows: pixel[...] = ConvertPixel(list(pixel))
        r+=1
        print(F"\033[A{(r/height)*100:.2f}%")

    return img


