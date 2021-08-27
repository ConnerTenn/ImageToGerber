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
    print()
    # covolution
    for y in range(img_height):
        for x in range(img_width):
            piece     = np.resize(img[y:y+kernel_height, x:x+kernel_width], (3,3)) # pads with copies when needed
            out[y, x] = np.sum(kernel * piece)
        print(F"\033[A{((y+1)/img_height)*100:.2f}%")

    return out

def EdgeDetection(img):
    img_conv_y = Convolution2d(img, SOBEL_Y)
    img_conv_x = Convolution2d(img, SOBEL_X)

    # combine the y and x edges
    img_conv = np.sqrt(np.square(img_conv_x) + np.square(img_conv_y))
    # normalize values
    img_conv *= 255 / img_conv.max()

    return img_conv
