
from Global import *

import numpy as np
import matplotlib.pyplot as plt


def Sigmoid(x):
    return 1 / (1 + 2**(-x))


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
            grayscaleImg[y][x] = Sigmoid(50*(avg-0.5))
        ProgressBar(y,0,height-1)
    print()

    img = grayscaleImg

    plt.imsave("grayscale.png", img)

    img = img > 0.4


    print("> Tracing outline")
    for y in range(height):
        for x in range(width):
            piece = np.resize(img[y:y+2, x:x+2], (2,2))
            # print(piece)

            tl = piece[1][0]
            tr = piece[1][1]
            bl = piece[0][0]
            br = piece[0][1]


            # ++
            # ..
            if ((    tl) and (    tr) and (not bl) and (not br)):
                lines += [ [[x+0.5,x+1.5], [y+1,y+1]] ]
            # ..
            # ++
            if ((not tl) and (not tr) and (    bl) and (    br)):
                lines += [ [[x+0.5,x+1.5], [y+1,y+1]] ]
            # +.
            # +.
            if ((    tl) and (not tr) and (    bl) and (not br)):
                lines += [ [[x+1,x+1], [y+0.5,y+1.5]] ]
            # .+
            # .+
            if ((not tl) and (    tr) and (not bl) and (    br)):
                lines += [ [[x+1,x+1], [y+0.5,y+1.5]] ]

            # .+
            # ++
            if ((not tl) and (    tr) and (    bl) and (    br)):
                lines += [ [[x+0.5,x+1], [y+1,y+1.5]] ]
            # +.
            # ++
            if ((    tl) and (not tr) and (    bl) and (    br)):
                lines += [ [[x+1,x+1.5], [y+1.5,y+1]] ]
            # ++
            # .+
            if ((    tl) and (    tr) and (not bl) and (    br)):
                lines += [ [[x+0.5,x+1], [y+1,y+0.5]] ]
            # ++
            # +.
            if ((    tl) and (    tr) and (    bl) and (not br)):
                lines += [ [[x+1,x+1.5], [y+0.5,y+1]] ]

            # +.
            # ..
            if ((    tl) and (not tr) and (not bl) and (not br)):
                lines += [ [[x+0.5,x+1], [y+1,y+1.5]] ]
            # .+
            # ..
            if ((not tl) and (    tr) and (not bl) and (not br)):
                lines += [ [[x+1,x+1.5], [y+1.5,y+1]] ]
            # ..
            # +.
            if ((not tl) and (not tr) and (    bl) and (not br)):
                lines += [ [[x+0.5,x+1], [y+1,y+0.5]] ]
            # ..
            # .+
            if ((not tl) and (not tr) and (not bl) and (    br)):
                lines += [ [[x+1,x+1.5], [y+0.5,y+1]] ]

            # +.
            # .+
            if ((    tl) and (not tr) and (not bl) and (    br)):
                lines += [ [[x+0.5,x+1], [y+1,y+1.5]] ]
                lines += [ [[x+1,x+1.5], [y+0.5,y+1]] ]
            # .+
            # +.
            if ((not tl) and (    tr) and (    bl) and (not br)):
                lines += [ [[x+1,x+1.5], [y+1.5,y+1]] ]
                lines += [ [[x+0.5,x+1], [y+1,y+0.5]] ]

                # exit()
        ProgressBar(y,0,height-1)
    print()

    return lines

img = plt.imread("Test1.png")
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
