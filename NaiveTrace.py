
from Global import *

import numpy as np
import matplotlib.pyplot as plt


"""
 +----+   +----+
 |####| 1 |####|
 |##+-------+##|
 +--|-+   +-|--+
  4 |       | 2
 +--|-+   +-|--+
 |##+-------+##|
 |####| 3 |####|
 +----+   +----+
"""

LineData = [
    { 
        "Kernel":[1,1,
                  0,0],
        "Segments":[(2,4)]
    },
    { 
        "Kernel":[0,0,
                  1,1],
        "Segments":[(2,4)]
    },
    { 
        "Kernel":[1,0,
                  1,0],
        "Segments":[(1,3)]
    },
    { 
        "Kernel":[0,1,
                  0,1],
        "Segments":[(1,3)]
    },

    { 
        "Kernel":[0,1,
                  1,1],
        "Segments":[(1,4)]
    },
    { 
        "Kernel":[1,0,
                  1,1],
        "Segments":[(1,2)]
    },
    { 
        "Kernel":[1,1,
                  0,1],
        "Segments":[(3,4)]
    },
    { 
        "Kernel":[1,1,
                  1,0],
        "Segments":[(2,3)]
    },

    { 
        "Kernel":[1,0,
                  0,0],
        "Segments":[(1,4)]
    },
    { 
        "Kernel":[0,1,
                  0,0],
        "Segments":[(1,2)]
    },
    { 
        "Kernel":[0,0,
                  1,0],
        "Segments":[(3,4)]
    },
    { 
        "Kernel":[0,0,
                  0,1],
        "Segments":[(2,3)]
    },

    { 
        "Kernel":[1,0,
                  0,1],
        "Segments":[(1,4),(2,3)]
    },
    { 
        "Kernel":[0,1,
                  1,0],
        "Segments":[(1,2),(3,4)]
    }
]

#Linear interpolation between p1 and p2, based on the relative values of a,b, and 0.5
def InterpolatePos(p1, p2, a,b):
    mu = (0.5-a)/(b-a)
    return p1 + mu*(p2-p1)

#Generate coordinates for a point on a line based on the type of line and the interpolation
def ConvertToCoords(lineIdx, pos, tl,tr,bl,br):
    if lineIdx == 1:
        return InterpolatePos(pos+[0,1], pos+[1,1], tl,tr)
    if lineIdx == 2:
        return InterpolatePos(pos+[1,0], pos+[1,1], br,tr)
    if lineIdx == 3:
        return InterpolatePos(pos+[0,0], pos+[1,0], bl,br)
    if lineIdx == 4:
        return InterpolatePos(pos+[0,0], pos+[0,1], bl,tl)

def GenerateSegments(x,y, tl,tr,bl,br):
    segments = []

    pos = np.array([x,y])
    #Check each to see what line data matches the pixel pattern
    for dat in LineData:
        kernel = dat["Kernel"]

        #Check to see if the pixel pattern matches
        if (kernel[0] == (tl>=0.5) and kernel[1] == (tr>=0.5) and kernel[2] == (bl>=0.5) and kernel[3] == (br>=0.5)):

            #Generate a line for each segment in the pixel pattern
            for seg in dat["Segments"]:
                p1 = ConvertToCoords(seg[0], pos, tl,tr,bl,br)
                p2 = ConvertToCoords(seg[1], pos, tl,tr,bl,br)

                segments += [ [(p1[0], p2[0]), (p1[1], p2[1])] ]
    return segments

def LineDetection(img):
    shape = img.shape
    height = shape[0]
    width = shape[1]
    depth = 1
    if len(shape)==3: depth = shape[2]

    segments = []

    if depth>1:
        grayscaleImg = np.zeros((height, width))
        print("> Creating Grayscale")
        for y in range(height):
            for x in range(width):
                r = img[y][x][0]
                g = img[y][x][1]
                b = img[y][x][2]

                avg = (r+g+b)/3
                grayscaleImg[y][x] = avg #Sigmoid(10*(avg-0.5))
            ProgressBar(y,0,height-1)
        print()

        img = grayscaleImg

    # plt.imsave("SigmoidGrayscale.png", img)

    # img = img > 0.4


    print("> Tracing outline")
    for y in range(height):
        for x in range(width):
            piece = np.resize(img[y:y+2, x:x+2], (2,2))
            # print(piece)

            tl = piece[1][0]
            tr = piece[1][1]
            bl = piece[0][0]
            br = piece[0][1]

            segments += GenerateSegments(x,y, tl,tr,bl,br)

        ProgressBar(y,0,height-1)
    print()

    return segments

def StitchSegments(segments):
    # segments = [
    #     [(0,0), (1,1)],
    #     [(5,20), (10,20)],
    #     [(-1,1), (0,2)],
    #     [(10,0), (20,0)],
    #     [(1,1), (0,2)],
    #     [(15,10), (10,0)],
    #     [(-1,1), (0,0)],
    #     [(20,0), (15,10)],
    # ]
    lineloops = [[]]

    print("> Stiching Segments")
    startNumSeg = len(segments)
    i=0
    noloop=True
    loopStartIdx = -1
    while len(segments):
        lastloop = lineloops[-1]
        seg = segments[i]

        # Loop start
        if len(lastloop)==0:
            print("Loop start")
            lineloops[-1] += [ seg[0] ]
            loopStartIdx = i
        #Loop end
        elif len(lastloop)>1 and lastloop[-1][0] == lastloop[0][0] and lastloop[-1][1] == lastloop[0][1]:
            # print("Loop complete")
            lineloops += [[]]
            noloop = False
        # Connect with 1st part of seg
        elif lastloop[-1][0] == seg[0][0] and lastloop[-1][1] == seg[0][1]:
            lastloop += [ seg[1] ]
            del segments[i]
            noloop = False
        # Connect with 2nd point of seg
        elif lastloop[-1][0] == seg[1][0] and lastloop[-1][1] == seg[1][1]:
            lastloop += [ seg[0] ]
            del segments[i]
            noloop = False
        else:
            i+=1
        # ProgressBar(startNumSeg-len(segments), 0, startNumSeg)
        # print(F" {startNumSeg-len(segments)}/{startNumSeg}", end="")

        if i>=len(segments):
            i=0
            if noloop:
                print("Noloop")
                for seg in segments:
                    print(F"    {lastloop[-1]} :: {seg}")
                lineloops += [[]]
                del segments[loopStartIdx]
            noloop=True
        # print(lineloops)
    print()

    print(F"Generated {len(lineloops)} loops")
    for loop in lineloops:
        print(loop)

    return lineloops


def PlotLines(lines):
    print("> Plotting Segments")
    for i, line in enumerate(lines):
        line[1] = [-line[1][0], -line[1][1]]
        # plt.axline(line[0],line[1])
        plt.plot(line[0], line[1], color="k")
        ProgressBar(i,0,len(lines)-1)
    print()
    # plt.xticks([])
    # plt.yticks([])
    plt.show()
