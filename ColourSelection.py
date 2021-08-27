
from math import *


#hsv in range [0,1]
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

#rgb in range [0,1]
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



