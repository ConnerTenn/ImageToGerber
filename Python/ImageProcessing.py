
from Global import *

import math
import numpy as np
import matplotlib.pyplot as plt

# typical kernel used for edge detection, sobel
SOBEL_Y = np.array([[-1, -2, -1],
                    [0, 0, 0],
                    [1, 2, 1]])

SOBEL_X = np.array([[-1, 0, 1],
                    [-2, 0, 2],
                    [-1, 0, 1]])

# unnessecary if using openCV
def RGBToGrey(img):
    return np.dot(img[...,:3], [0.2989, 0.5870, 0.1140])

def Convolution2d(img, kernel):
    # check if image is RGB and convert if needed
    if len(img.shape) == 3:
        img = RGBToGrey(img)

    # pull out dimensions for QoL
    img_height, img_width       = img.shape
    kernel_height, kernel_width = kernel.shape

    out = np.zeros((img_height, img_width))

    print("> Convolution")
    # covolution
    for y in range(img_height):
        for x in range(img_width):
            piece     = np.resize(img[y:y+kernel_height, x:x+kernel_width], (3,3)) # pads with copies when needed
            out[y, x] = np.sum(kernel * piece)
        ProgressBar(y, 0, img_height-1)
    print()

    return out

def EdgeDetection(img):
    img_conv_y = Convolution2d(img, SOBEL_Y)
    img_conv_x = Convolution2d(img, SOBEL_X)

    # combine the y and x edges
    img_conv = np.sqrt(np.square(img_conv_x) + np.square(img_conv_y))
    # normalize values
    img_conv *= 255 / img_conv.max()

    return img_conv

# line detection with hough transform
# maths is hard, https://en.wikipedia.org/wiki/Hough_transform
# assuming greyscale and edge detected for now
def LineDetection(img):
    img_height, img_width = img.shape

    # deterministic method of getting local extremes
    # theta = angle between the straight line and the closest point to the origin
    theta_max = 1.0 * math.pi
    theta_min = 0
    theta_dim = 300

    # r = radius from the origin to the straight line
    r_max = math.hypot(img_height, img_width)
    r_min = 0
    r_dim = 200

    # setup hough space
    hough = np.zeros((r_dim, theta_dim))

    print("> Hough Transform")
    for y in range(img_height):
        for x in range(img_width):
            if img[y,x] == 255:
                continue

            for theta_idx in range(theta_dim):
                theta = 1.0 * theta_idx * theta_max / theta_dim
                r     = y * math.cos(theta) + x * math.sin(theta)
                r_idx = r_dim * (1.0 * r) / r_max
                hough[int(r_idx), int(theta_idx)] = hough[int(r_idx), int(theta_idx)] + 1
        ProgressBar(y, 0, img_height-1)
    print()

    # getting extrema
    neighborhood = 20  # adjustable params
    threshold    = 140


    return hough


def GaussianBlur(img):
    kernel = np.array([[1, 2, 1],
                       [2, 4, 2],
                       [1, 2, 1]])

    return Convolution2d(img, kernel)



def GenerateOctree(img, offX=0, offY=0, scale=0):

    height, width = img.shape
    if width*height==0:
        return []

    rects = []

    # print("| "*scale, offX, offY)
    # print(img)
    # print()

    if width != height:
        extra = None
        if width > height:
            extra = GenerateOctree(img[    0:height, height:width], offX+height, offY+0,     scale+1)
            width = height
        if height > width:
            extra = GenerateOctree(img[width:height,      0:width], offX,        offY+width, scale+1)
            height = width
        if extra==True:
            rects+=[[scale, offX, offY]]
        else:
            rects += extra

    elif width==1:
        if img[0,0]:
            return True
        else:
            return []

    tl = GenerateOctree(img[            0:int(height/2),            0:int(width/2)], offX+0,            offY+0,             scale+1)
    tr = GenerateOctree(img[            0:int(height/2), int(width/2):width       ], offX+int(width/2), offY+0,             scale+1)
    bl = GenerateOctree(img[int(height/2):height       ,            0:int(width/2)], offX+0,            offY+int(height/2), scale+1)
    br = GenerateOctree(img[int(height/2):height       , int(width/2):width       ], offX+int(width/2), offY+int(height/2), scale+1)

    if (tl==True and tr==True and bl==True and br==True):
        return True
    else:
        if tl==True:
            rects+=[[scale+1, offX, offY]]
        else:
            rects+=tl

        if tr==True:
            rects+=[[scale+1, offX+int(width/2), offY]]
        else:
            rects+=tr

        if bl==True:
            rects+=[[scale+1, offX, offY+int(height/2)]]
        else:
            rects+=bl

        if br==True:
            rects+=[[scale+1, offX+int(width/2), offY+int(height/2)]]
        else:
            rects+=br

        return rects


def CenterOfMass(img):
    avgx = 0
    avgy = 0
    count = 0
    for y, rows in enumerate(img):
        for x, pixel in enumerate(rows):
            if pixel:
                avgx+=x
                avgy+=y
                count+=1
    if count==0:
        shape = img.shape
        return [shape[1]/2, shape[0]/2]
    avgx = avgx/count
    avgy = avgy/count
    return [avgx, avgy]



