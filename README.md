### apple music playlist manager
Addes/deleted tracks to apple music

### run
```bash
cat config.example.json > config.json # new config
go build -o manager cmd/manager/main.go
./manager -c config.json -f amids.csv -l 100 -o 1000 delete|add
```

### Описание флагов:
 - `-c` - config path
 - `-f` - amids file path (one column csv)
 - `-l` - limit, limit to read
 - `-o` - offset - offset read

### Команды
- add - adds track
- delete - deletes tracks
