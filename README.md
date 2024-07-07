# RaGodahn
Go PE Injector 

To build on Linux:

GOOS=windows GOARCH=amd64 go build main.go

To build and execute on windows:

go build main.go && ./main.exe

OR 

go run main.go

(in case of "go run" a binary file will be built but stored in %temp% )



![Alt text](https://github.com/0xDagobert/RaGodahn/blob/main/images/MasterRaGodahn.png)


Usage:

There is 3 mode for the injector:
- PE
- Process-auto
- Process-manuel

PE mode will inject the shellcode in the binary file that you will need to input 
ex: main.exe -mode PE --path "C:\Users\test\Desktop\test.exe"

Process-auto will spawn a notepad process and inject the shellcode in the process
ex: main.exe -mode Process-auto

Process-manuel will display all the running processes on the computer and let you choose which one to infect (still in beta might not work for every process, but will work for notepad, calculator and msedge processes)
ex: main.exe 

