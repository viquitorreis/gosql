#!/bin/bash

# Diretório onde estão os arquivos
diretorio="/home/victor/Projetos/golang/go-migrations/gosql/migrations"

# Número inicial para renomear
numero_inicial=10

# Loop para renomear os arquivos
for arquivo in "$diretorio"/*; do
    if [ -f "$arquivo" ]; then
        novo_numero=$(printf "%04d" "$numero_inicial")
        nome_arquivo=$(basename "$arquivo")
        novo_nome="${novo_numero}${nome_arquivo#????}"
        mv "$arquivo" "${diretorio}/${novo_nome}"
        ((numero_inicial+=10))
    fi
done