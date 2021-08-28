
from Global import *

ColourFormats = {
    "R":1,"G":1,"B":1,"V":1,"C":1,"L":1,"H":1,"Sv":1,"Sl":1,"A":1,
    "RGB":3, "RGBA":4, "HSV":3, "HSVA":4, "HSL":3, "HSLA":3
}


def ParseFilename(line):
    tokens = line.split("\"")

    #This line must only consist of a single name between quotes
    if len(tokens) != 3 or len(tokens[0]) or len(tokens[2]):
        Error("Malformed filepath specifier")

    path = tokens[1]
    print(F"{TERM_MAGENTA}>> {path} <<{TERM_RESET}")


def ParseParams(colourFormat, params):
    params = params.split(",")

    if len(params) != ColourFormats[colourFormat]:
        Error("Incorrect number of arguments for colour format")

    args = []
    for arg in params:
        args += [float(arg)]
    
    print(F"{F'[{colourFormat}]':6s}: {args[0]}")
    for arg in args[1:]:
        print(F"      : {arg}")

def ParseTolerance(tolerance):
    tolPos = 0
    tolNeg = 0

    if tolerance.count("+-")==1:
        if tolerance.count("+")>1 or tolerance.count("-")>1:
            Error("A +- tolerance must not be followed by additional tolerances")

        plusminusIdx = tolerance.find("+-")
        percentIdx = tolerance.find("%",plusminusIdx+1)
        tolPos = tolNeg = int(tolerance[plusminusIdx+2:percentIdx])


    else:
        try:
            if tolerance.find("+")!=-1:
                plusIdx = tolerance.find("+")
                percentIdx = tolerance.find("%",plusIdx+1)
                if percentIdx==-1:
                    Error("Malformed tolerance specifier")
                tolPos = int(tolerance[plusIdx+1:percentIdx])

            if tolerance.find("-")!=-1:
                minusIdx = tolerance.find("+")
                percentIdx = tolerance.find("%",minusIdx+1)
                if percentIdx==-1:
                    Error("Malformed tolerance specifier")
                tolPos = int(tolerance[minusIdx+1:percentIdx])
        except:
            Error("Failed to parse tolerances")
        
    print(F"    Tolerance +{tolPos}% -{tolNeg}%")


def ParseRule(rule):
    rule = rule.strip()

    #Syntax checking
    if " " in rule:
        Error("There must not be space inside rules")
    if "#" in rule:
        Error("Invalid character in rule")
    if rule.count("(")!=1 and rule.count(")")!=1:
        Error("Malformed rule")
    if rule.find("(") > rule.find(")"):
        Error("Malformed rule")


    #Delimit into sections
    rule = rule.replace("(","#")
    rule = rule.replace(")","#")

    tokens = rule.split("#")
    if len(tokens) != 3:
        Error("Malformed rule")

    colourFormat = tokens[0]
    params = tokens[1]
    tolerance = tokens[2]

    for fmt in ColourFormats:
        if colourFormat == fmt:
            ParseParams(colourFormat, params)
            ParseTolerance(tolerance)


def ParseSelection(line):

    #syntax checking
    negate = line.count("!")
    if negate>1:
        Error("The negation operator must only be specified once")

    if negate:
        if line[0]!="!":
            Error("Negation operator must appear first in the string")
        line=line[1:]


    for rule in line.split("&"):
        ParseRule(rule)



def ParseConfig(filename):
    config=[]

    cfgfile=None
    try:
        cfgfile=open(filename)
    except:
        print(F"{TERM_RED}Failed to open file {filename}{TERM_RESET}")


    for line in cfgfile:
        line = line.split("#")[0].rstrip()

        if len(line):

            if line[0].isspace():
                line = line.strip()

                ParseSelection(line)
            else:
                ParseFilename(line)
