#!/bin/bash

destination_directory="../gosql/migrations"

# Loop de 1 a 15
for i in {1..15}; do
  # Use printf para formatar o número com zeros à esquerda
  num=$(printf "%04d" $i)
  # Crie o arquivo com o nome desejado
  touch "${destination_directory}/${num}.teste.sql"
done
