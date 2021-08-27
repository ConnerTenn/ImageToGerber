
TERM_RESET =   "\033[m"
TERM_RED =     "\033[1;31m"
TERM_GREEN =   "\033[1;32m"
TERM_YELLOW =  "\033[1;33m"
TERM_BLUE =    "\033[1;34m"
TERM_MAGENTA = "\033[1;35m"
TERM_CYAN =    "\033[1;36m"
TERM_WHITE =   "\033[1;37m"

def Error(msg):
    print(F"{TERM_RED}Error: {msg}{TERM_RESET}")
    exit(1)

def Assert(value, expected, msg=""):
    if (value==expected):
        print(F"{TERM_GREEN} PASS: {msg}{TERM_RESET}")
    else:
        print(F"{TERM_RED} FAIL: {msg}")
        print(F" - - - ({value}) != Expected ({expected}){TERM_RESET}")

