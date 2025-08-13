# Поиск строк с подстрокой "ошибка"
go run grep.go "ошибка" sample.txt

# С контекстом после (2 строки после)
go run grep.go -A 2 "ошибка" sample.txt

# С контекстом до (1 строка до)
go run grep.go -B 1 "ошибка" sample.txt

# С контекстом вокруг (до и после)
go run grep.go -C 1 "ошибка" sample.txt

# Только количество совпадений
go run grep.go -c "ошибка" sample.txt

# Инвертировать: всё, где НЕТ "ошибка"
go run grep.go -v "ошибка" sample.txt

# Игнорировать регистр (найдёт и "Ошибка", и "ОШИБКА")
go run grep.go -i "Ошибка" sample.txt

# Как фиксированную строку (без регэкспа)
go run grep.go -F "ошибка: файл не найден" sample.txt

# С номерами строк
go run grep.go -n "ошибка" sample.txt
