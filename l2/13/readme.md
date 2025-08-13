# Вывести 1 и 3 поле (нумерация с 1)
cat data.txt | go run cut.go -f 1,3 -d ":"

# Диапазон: со 2 по 4 включительно
cat data.txt | go run cut.go -f 2-4 -d ":"

# Только строки, содержащие разделитель (исключит "hello world")
cat data.txt | go run cut.go -f 1 -d ":" -s
