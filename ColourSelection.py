
from math import *

import numpy as np
import matplotlib.pyplot as plt


#hsv: float range[0,1]
def HSVtoRGB(hsv):
    #https://en.wikipedia.org/wiki/HSL_and_HSV
    h = hsv[0]
    s = hsv[1]
    v = hsv[2]

    k_r = (5+6*h)%6
    k_g = (3+6*h)%6
    k_b = (1+6*h)%6
    r = v-v*s*max(0,min([k_r, 4-k_r, 1]))
    g = v-v*s*max(0,min([k_g, 4-k_g, 1]))
    b = v-v*s*max(0,min([k_b, 4-k_b, 1]))

    return [r,g,b]

#rgb: float range[0,1]
def RGBtoHSV(rgb):
    r = rgb[0]
    g = rgb[1]
    b = rgb[2]

    v = max(rgb)
    chroma = v - min(rgb)

    s = (0 if v==0 else chroma/v)

    h = 0
    if (chroma==0): pass
    elif (v==r): h = (0 + (g-b)/chroma)/6
    elif (v==g): h = (2 + (b-r)/chroma)/6
    elif (v==b): h = (4 + (r-g)/chroma)/6

    return [h,s,v]

#value: float range[0,1]
#target: float range[0,1]
#tolerance: float range[0,1]
def TestTol(value, target, tolerance):
    return abs(target-value) <= tolerance/2

def TestTolWrap(value, target, tolerance):
    return (target-value)%1 <= tolerance/2 or (value-target)%1 <= tolerance/2



#pixel: float[3] range[0,1]
def ConvertPixel(pixel):
    h, s, v = RGBtoHSV(pixel)

    if (TestTolWrap(h, 0.1, 0.2) and
        TestTol(s, 0.9, 0.4) and
        TestTol(v, 0.7, 0.4)):
        return pixel
    return [0]


def ConvertImage(imagefile):
    img = plt.imread(imagefile)
    img = img[...,:3]

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


