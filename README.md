# go-musthave-shortener-tpl
Шаблон репозитория для практического трека «Go в веб-разработке».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

git add --all
git commit -m "My first commit"
git push --set-upstream origin inc12


go get github.com/caarlos0/env
go get github.com/go-chi/chi
go get github.com/go-chi/chi/v5/middleware
go get github.com/gofrs/uuid
go get github.com/google/uuid
go get github.com/jackc/pgerrcode
go get github.com/lib/pq
go get github.com/mihailcoc/shortener/internal/app/workers
go get golang.org/x/sync/errgroup
