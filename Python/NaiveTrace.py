
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

    #Check if any work has to be done for this selection
    if not ((tl>=0.5) == (tr>=0.5) == (bl>=0.5) == (br>=0.5)):
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

                    segments += [ [(p1[0], p1[1]), (p2[0], p2[1])] ]
                return segments
    return []

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
    ts=time.time()
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
        TimeDisplay(ts)
    print()

    return segments

def StitchSegments(segments):
    lineloops = [[]]

    print("> Stiching Segments")
    startNumSeg = len(segments)
    i=0
    ts=time.time()
    while len(segments):
        lastloop = lineloops[-1]
        seg = segments[i]

        # Loop start
        if len(lastloop)==0:
            # print("Loop start")
            lineloops[-1] += seg #Add both points to the loop
            del segments[i]
        #Loop end
        elif len(lastloop)>2 and lastloop[-1][0] == lastloop[0][0] and lastloop[-1][1] == lastloop[0][1]:
            # print("Loop complete")
            lineloops += [[]]
        # Connect with 1st part of seg
        elif lastloop[-1][0] == seg[0][0] and lastloop[-1][1] == seg[0][1]:
            lastloop += [ seg[1] ]
            del segments[i]
        # Connect with 2nd point of seg
        elif lastloop[-1][0] == seg[1][0] and lastloop[-1][1] == seg[1][1]:
            lastloop += [ seg[0] ]
            del segments[i]
        else:
            i+=1
        ProgressBar(startNumSeg-len(segments), 0, startNumSeg)
        TimeDisplay(ts)

        #Loop around when reached the end of the segment list
        if i>=len(segments):
            i=0
    print()

    print(F"Generated {len(lineloops)} loops")

    return lineloops

def PlotLineLoops(lineloops):
    lines = []
    for loop in lineloops:
        if len(loop):
            lastpt = loop[0]
            for point in loop[1:]:
                lines += [ [lastpt, point] ]
                lastpt = point
    PlotLines(lines)

def PlotLines(lines):
    print("> Plotting Segments")
    ts=time.time()
    for i, line in enumerate(lines):
        #Convert to [ (x1, x2), (y1, y2) ] format
        line = [ [line[0][0], line[1][0]], [-line[0][1], -line[1][1]] ]
        # plt.axline(line[0],line[1])
        plt.plot(line[0], line[1], color="k")
        ProgressBar(i,0,len(lines)-1)
        TimeDisplay(ts)
    print()
    # plt.xticks([])
    # plt.yticks([])
    plt.show()
