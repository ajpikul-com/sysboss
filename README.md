# SYSBOSS 

Need to refactor, it super slppy for moduling. EVerything should be moduled into one global state w/ interface. Modules can provide their own command and command processining. GET + PID. Also, copybuffer on REeadtext has a better option. 


Note:

As of right now, all keys and configs are parsed on every request. This is probably fine given the extremely low volume of requests, and mitigates the need for something that watches and loads file. HOWEVER, it is dumb.


Command processing could be better

It works tho
