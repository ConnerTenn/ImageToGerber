
from Global import *

ValidSelectionTokens = (
    "R","G","B","V","C","L","H","Sv","Sl","A",
    "RGB", "RGBA", "HSV", "HSVA", "HSL", "HSLA"
)


def ParseFilename(line):
    tokens = line.split("\"")

    #This line must only consist of a single name between quotes
    if len(tokens) != 3 or len(tokens[0]) or len(tokens[2]):
        Error("Malformed filepath specifier")

    path = tokens[1]
    print(F"{TERM_MAGENTA}>> {path} <<{TERM_RESET}")


def ParseRule(rule):
    rule = rule.strip()
    print(F"    {rule}")

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

    colourSel = tokens[0]
    params = tokens[1]
    tolerance = tokens[2]

    for sel in ValidSelectionTokens:
        if colourSel == sel:
            # ParseParams(colourSel, params)
            # ParseTolerance(tolerance)
            pass


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
