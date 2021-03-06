
from Global import *

ColourFormats = {
    "R":1,"G":1,"B":1,"V":1,"C":1,"L":1,"H":1,"Sv":1,"Sl":1,"A":1,
    "RGB":3, "RGBA":4, "HSV":3, "HSVA":4, "HSL":3, "HSLA":3
}

"""
== Call Stack ==
ParseConfig
    ParsePath
    ParseSelection
        ParseRule
            ParseParams
            ParseTolerance

"""


#Parse the config file
def ParseConfig(filename):
    print(F"{TERM_BLUE}=== Parsing Config ==={TERM_RESET}")
    config={"Processes":[], "Options":{}}

    cfgfile=None
    try:
        cfgfile=open(filename)
    except:
        Error(F"Failed to open file {filename}")


    for line in cfgfile:
        line = line.split("#")[0].rstrip()

        if len(line):

            if line[0].isspace():
                line = line.strip()
                config["Processes"][-1]["Selections"] += [ParseSelection(line)]
            elif line[0] == "\"":
                config["Processes"] += [{}]
                config["Processes"][-1] |= ParseProcess(line)
                config["Processes"][-1]["Selections"] = []
            elif line.count("=") == 1:
                key, value = ParseOption(line)
                config["Options"][key] = value
            else:
                Error("Unrecognized config line")

    # print(F"{TERM_GREEN}Config Dump{TERM_RESET}")

    # print(config)
    # PrintDict(config)

    if not ("Width" in config["Options"] or "Height" in config["Options"]):
        Error("At least 1 board dimension must be specified")

    return config


#Parse an option line
def ParseOption(line):
    key, value = line.split("=")
    key = key.strip()
    value = value.strip()

    if key in ["Width", "Height"]:
        mm = value.find("mm")
        if mm==-1:
            Error("Incorrect dimension format")

        try:
            value = float(value[:mm])
        except:
            Error("Could not parse dimension value")
    else:
        Error(F"Unrecognized config option \"{key}\"")

    return key, value


#Parse a new process
def ParseProcess(line):
    tokens = line.split(":")

    #This line must contain 2 parts
    if len(tokens) != 2:
        Error("Malformed process line")

    proc = {}
    proc |= ParsePath(tokens[0].strip())

    if tokens[1].count("[") != 1 or tokens[1].count("]") != 1 or "#" in tokens[1]:
        Error("Invalid process config syntax")

    config = tokens[1].strip()
    config = config.replace("[", "#")
    config = config.replace("]", "#")
    tokens = config.split("#")

    if len(tokens) != 3 or len(tokens[2]):
        Error("Invalid process config syntax")

    proc["Type"] = tokens[0].strip()
    proc["Method"] = tokens[1].strip()

    if not proc["Type"] in ["F_Cu", "B_Cu", "F_Mask", "B_Mask", "F_SilkS", "B_SilkS", "Edge_Cuts"]:
        Error(F"Invalid gerber type: {proc['Type']}")
    return proc


#Parse the file names at the beginning of a new section
def ParsePath(line):
    tokens = line.split("\"")

    #This line must only consist of a single name between quotes
    if len(tokens) != 3 or len(tokens[0]) or len(tokens[2]):
        Error("Malformed filepath specifier")

    path = tokens[1]
    print(F"{TERM_GREEN}>> {path}{TERM_RESET}")

    return {"Path":path}



#Parse a line of the selection specifier
def ParseSelection(line):

    config = {"Rules":[], "Negate":False}

    #syntax checking
    negate = line.count("!")
    if negate>1:
        Error("The negation operator must only be specified once")

    if negate:
        if line[0]!="!":
            Error("Negation operator must appear first in the string")
        line=line[1:]

        config["Negate"] = True

    for rule in line.split("&"):
        config["Rules"] += [ParseRule(rule)]

    return config



#Parse each individual rule
def ParseRule(rule):
    rule = rule.strip()

    config = {}

    #Syntax checking
    if " " in rule:
        Error("There must not be space inside rules")
    if rule.count("(")!=1 and rule.count(")")!=1:
        Error("Malformed rule")
    if rule.find("(") > rule.find(")"):
        Error("Malformed rule")


    #Delimit into sections, seperated by the brackets '(' and ')'
    #Replace with temporary character and then split by that character
    rule = rule.replace("(","\x80")
    rule = rule.replace(")","\x80")

    tokens = rule.split("\x80")
    if len(tokens) != 3:
        Error("Malformed rule")

    colourFormat = tokens[0]
    params = tokens[1]
    tolerance = tokens[2]

    #Find the matching colour format
    if not colourFormat in ColourFormats:
        Error(F"Invalid colour format specified \"{colourFormat}\"")
    for fmt in ColourFormats:
        if colourFormat == fmt:
            config |= ParseParams(colourFormat, params)
            config |= ParseTolerance(tolerance)

    return config


#Parse the parameters in a colour format
def ParseParams(colourFormat, params):
    params = params.split(",")

    #Make sure that the correct number of args is specified for the format
    if len(params) != ColourFormats[colourFormat]:
        Error("Incorrect number of arguments for colour format")

    args = []
    for arg in params:
        try:
            args += [float(arg)]
        except:
            Error("Failed to parse colour format argument")
    
    print(F"{F'[{colourFormat}]':6s}: {args[0]}")
    for arg in args[1:]:
        print(F"      : {arg}")

    return {"ColourFormat":colourFormat, "Args":args}

#Parse the tolerance specifiers
def ParseTolerance(tolerance):
    tolPos = 0
    tolNeg = 0

    #Check if the tolerance is specified as a +-
    if tolerance.count("+-")==1:
        if tolerance.count("+")>1 or tolerance.count("-")>1:
            Error("A +- tolerance must not be followed by additional tolerances")

        plusminusIdx = tolerance.find("+-")
        percentIdx = tolerance.find("%",plusminusIdx+1) #Make sure to find the '%' immediately following
        if percentIdx==-1:
            Error("Malformed tolerance specifier")
        #Convert to float
        tolPos = tolNeg = float(tolerance[plusminusIdx+2:percentIdx])/100
    #Tolerance must have separate pos/neg specifiers
    else:
        try:
            #Find the pos specifier
            if tolerance.find("+")!=-1:
                plusIdx = tolerance.find("+")
                percentIdx = tolerance.find("%",plusIdx+1) #Make sure to find the '%' immediately following
                if percentIdx==-1:
                    Error("Malformed tolerance specifier")
                #Convert to float
                tolPos = float(tolerance[plusIdx+1:percentIdx])/100

            #Find the neg specifier
            if tolerance.find("-")!=-1:
                minusIdx = tolerance.find("-")
                percentIdx = tolerance.find("%",minusIdx+1) #Make sure to find the '%' immediately following
                if percentIdx==-1:
                    Error("Malformed tolerance specifier")
                #Convert to float
                tolNeg = float(tolerance[minusIdx+1:percentIdx])/100
        except:
            Error("Failed to parse tolerances")
        
    print(F"    Tolerance +{tolPos*100}% -{tolNeg*100}%")

    return {"Tol+":tolPos, "Tol-":tolNeg}

