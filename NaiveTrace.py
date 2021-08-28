
from Global import *

import numpy as np
import matplotlib.pyplot as plt



def LineDetection(img):
    height, width, depth = img.shape

    lines = []

    img = np.dot(img[...,:3], [0.2989, 0.5870, 0.1140])

    img = img > 0.4


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
                lines += [ [[x,x+1], [y+1,y+1]] ]
            # ..
            # ++
            if ((not tl) and (not tr) and (    bl) and (    br)):
                lines += [ [[x,x+1], [y,y]] ]
            # +.
            # +.
            if ((    tl) and (not tr) and (    bl) and (not br)):
                lines += [ [[x,x], [y+1,y]] ]
            # .+
            # .+
            if ((not tl) and (    tr) and (not bl) and (    br)):
                lines += [ [[x+1,x+1], [y+1,y]] ]

            # .+
            # ++
            if ((not tl) and (    tr) and (    bl) and (    br)):
                lines += [ [[x,x+1], [y,y+1]] ]
            # +.
            # ++
            if ((    tl) and (not tr) and (    bl) and (    br)):
                lines += [ [[x,x+1], [y+1,y]] ]
            # ++
            # .+
            if ((    tl) and (    tr) and (not bl) and (    br)):
                lines += [ [[x,x+1], [y+1,y]] ]
            # ++
            # +.
            if ((    tl) and (    tr) and (    bl) and (not br)):
                lines += [ [[x,x+1], [y,y+1]] ]

            # .+
            # +.
            if ((not tl) and (    tr) and (    bl) and (not br)):
                lines += [ [[x,x+1], [y,y+1]] ]
            # +.
            # .+
            if ((    tl) and (not tr) and (not bl) and (    br)):
                lines += [ [[x,x+1], [y+1,y]] ]

                # exit()
        ProgressBar(y,0,height-1)
    print()

    return lines

img = plt.imread("Test3.png")
height, width, depth = img.shape
lines = LineDetection(img)

for line in lines:
    line[1] = [height-line[1][0], height-line[1][1]]
    # plt.axline(line[0],line[1])
    plt.plot(line[0], line[1], color="k")
# plt.xticks([])
# plt.yticks([])
plt.show()
