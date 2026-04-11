# Читаем .env вручную
$envContent = Get-Content .env

foreach ($line in $envContent) { #чтение каждого значения из файла env
    if ($line -and $line[0] -ne '#') {
        $parts = $line -split '=', 2
        if ($parts.Count -eq 2) { #если значение = 2
            Set-Item -Path "env:$($parts[0])" -Value $parts[1] #создает переменные окружения (только в текущем сеансе PowerShell)
        }
    }
}

# Формируем строку подключения
$dbUrl = "postgres://$($env:DB_USER):$($env:DB_PASSWORD)@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)?sslmode=disable"

Write-Host "Connecting to: postgres://$($env:DB_USER):****@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)" -ForegroundColor Green

# Запускаем миграцию
& migrate -path ./migrations -database $dbUrl up

if ($LASTEXITCODE -eq 0) { #проверка кода возврата последней выполненной программы
    Write-Host "Migrations applied successfully!" -ForegroundColor Green
} else {
    Write-Host "Migration failed" -ForegroundColor Red
}

Read-Host "Press Enter to exit"