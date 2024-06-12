## ORPXI - Расширяемая командная строка на Go

![GitHub top language](https://img.shields.io/github/languages/top/bezhan2009/ORPXI) 
![GitHub language count](https://img.shields.io/github/languages/count/bezhan2009/ORPXI)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/bezhan2009/ORPXI)
![GitHub repo size](https://img.shields.io/github/repo-size/bezhan2009/ORPXI) 
![GitHub](https://img.shields.io/github/license/bezhan2009/ORPXI) 
![GitHub last commit](https://img.shields.io/github/last-commit/bezhan2009/ORPXI)
![GitHub User's stars](https://img.shields.io/github/stars/bezhan2009?style=social)

<p align="left">
    <img src="https://visitor-badge.laobi.icu/badge?page_id=bezhan2009.ORPXI" alt="visitors"/>
</p>

Читать на [English](README.md)

**ORPXI** - это альтернативная командная строка, написанная на языке программирования Go. Она предоставляет те же основные команды, что и стандартная CMD, но также включает в себя дополнительные команды, специфичные для этого инструмента.

### Особенности

- **Совместимость с CMD**: ORPXI поддерживает большинство команд CMD, обеспечивая плавный переход для пользователей.
- **Дополнительные команды**: ORPXI включает в себя ряд дополнительных команд, разработанных для упрощения повседневных задач системных администраторов и разработчиков.
- **Расширяемость**: Поскольку ORPXI написан на Go, его функциональность легко расширяется за счет разработки пользовательских команд.
- **Производительность**: Go известен своей высокой производительностью, что делает ORPXI быстрым и отзывчивым.
- **Поддержка и обратная связь**: Есть возможность обратиться к разработчику напрямую через Telegram.

### Установка и Запуск

```bash
git clone https://github.com/bezhan2009/ORPXI.git
cd ORPXI
go run .
```

### Список доступных команд

- **CREATE**: Создает новый файл.
- **CLEAN**: Очищает экран.
- **CD**: Смена текущего каталога.
- **COPUSOURCE**: копирует исходный код файла
- **LS**: Выводит содержимое каталога.
- **NEWSHABLON**: Создает новый шаблон команд для выполнения.
- **REMOVE**: Удаляет файл.
- **READ**: Выводит на экран содержимое файла.
- **PROMPT**: Изменяет ORPXI.
- **PINGVIEW**: Показывает пинг.
- **NEWUSER**: Новый пользователь для ORPXI.
- **ORPXI**: Запускает еще одну ORPXI.
- **SHABLON**: Выполняет определенный шаблон команд.
- **SYSTEMGOCMD**: Вывод информации о ORPXI.
- **SYSTEMINFO**: Вывод информации о системе.
- **SIGNOUT**: Пользователь выходит из ORPXI.
- **TREE**: Графически отображает структуру каталогов диска или пути.
- **WRITE**: Записывает данные в файл.
- **EDIT**: Редактирует файл.
- **EXTRACTZIP**: Распаковывает архивы .zip.
- **SCANPORT**: Сканирование портов.
- **WHOIS**: Информация о домене.
- **DNSLOOKUP**: DNS-запросы.
- **WIFIUTILS**: Запускает утилиту для работы с WiFi.
- **IPINFO**: Информация об IP-адресе.
- **GEOIP**: Геолокация IP-адреса.
- **EXIT**: Выход.

### Руководство по разработке

Если вы хотите расширить функциональность ORPXI, вы можете написать разработчику [здесь](https://t.me/Rust_Bezhan).

### Обратная связь и поддержка

Если у вас возникли вопросы, предложения или проблемы, пожалуйста, создайте новый [issue](https://github.com/bezhan2009/ORPXI/issues/new) в нашем репозитории на GitHub.
