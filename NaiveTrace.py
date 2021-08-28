
from Global import *

import numpy as np
import matplotlib.pyplot as plt


def Sigmoid(x):
    return 1 / (1 + 2**(-x))


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

def InterpolatePos(a,b):
    return 1.5-(a+b)

def ConvertToCoords(lineIdx, tl,tr,bl,br):
    if lineIdx == 1:
        return [0.5,1]#[InterpolatePos(tl,tr), 1]
    if lineIdx == 2:
        return [1,0.5]#[1, InterpolatePos(br,tr)]
    if lineIdx == 3:
        return [0.5,0]#[InterpolatePos(bl,br), 0]
    if lineIdx == 4:
        return [0,0.5]#[0, InterpolatePos(bl,tl)]

def GenerateLines(x,y, tl,tr,bl,br):
    lines = []
    for dat in LineData:
        kernel = dat["Kernel"]

        if (kernel[0] == (tl>0.5) and
            kernel[1] == (tr>0.5) and
            kernel[2] == (bl>0.5) and
            kernel[3] == (br>0.5)):

            for seg in dat["Segments"]:
                p1 = ConvertToCoords(seg[0], tl,tr,bl,br)
                p2 = ConvertToCoords(seg[1], tl,tr,bl,br)

                p1 = (p1[0]+x, p1[1]+y)
                p2 = (p2[0]+x, p2[1]+y)
                # print()
                # print(p1)
                # print(p2)
                # exit()

                lines += [ [(p1[0], p2[0]), (p1[1], p2[1])] ]
    return lines

def LineDetection(img):
    height, width, depth = img.shape

    lines = []

    # img = np.dot(img[...,:3], [0.2989, 0.5870, 0.1140])

    grayscaleImg = np.zeros((height, width))
    print("> Creating Grayscale")
    for y in range(height):
        for x in range(width):
            r = img[y][x][0]
            g = img[y][x][1]
            b = img[y][x][2]

            avg = (r+g+b)/3
            grayscaleImg[y][x] = Sigmoid(20*(avg-0.5))
        ProgressBar(y,0,height-1)
    print()

    img = grayscaleImg

    plt.imsave("grayscale.png", img)

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

            lines += GenerateLines(x,y, tl,tr,bl,br)

        ProgressBar(y,0,height-1)
    print()

    return lines

img = plt.imread("Test3.png")
height, width, depth = img.shape

lines = LineDetection(img)

print("> Plotting Segments")
for i, line in enumerate(lines):
    line[1] = [height-line[1][0], height-line[1][1]]
    # plt.axline(line[0],line[1])
    plt.plot(line[0], line[1], color="k")
    ProgressBar(i,0,len(lines)-1)
print()
# plt.xticks([])
# plt.yticks([])
plt.show()
