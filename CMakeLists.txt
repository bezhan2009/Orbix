cmake_minimum_required(VERSION 3.10)
project(CompileEachCppFileSeparately)

# Установите стандарт C++
set(CMAKE_CXX_STANDARD 17)

# Найдите все .cpp файлы в текущей директории
file(GLOB CPP_FILES "${CMAKE_SOURCE_DIR}/*.cpp")

# Для каждого .cpp файла создаем отдельный исполняемый файл
foreach(CPP_FILE ${CPP_FILES})
    # Извлекаем имя файла без расширения, чтобы использовать его как имя исполняемого файла
    get_filename_component(EXE_NAME ${CPP_FILE} NAME_WE)

    # Создаем исполняемый файл с именем, совпадающим с именем исходного файла
    add_executable(${EXE_NAME} ${CPP_FILE})
endforeach()
