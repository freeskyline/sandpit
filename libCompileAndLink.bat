go tool compile -I libDir lib.go
go tool link -o lib.exe -L libDir lib.o
pause
