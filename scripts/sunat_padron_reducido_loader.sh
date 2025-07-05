#!/bin/bash
DB_CONTAINER=setup-postgres-1
if [ -f /tmp/padron_reducido_ruc.zip ]; then
    rm /tmp/padron_reducido_ruc.zip
fi
if [ -f /tmp/padron_reducido_ruc.txt ]; then
    rm /tmp/padron_reducido_ruc.txt
fi
if [ -f /tmp/rucs.csv ]; then
    rm /tmp/rucs.csv
    rm /tmp/rucsutf8.csv
fi
# Descargar padron reducido
wget -O /tmp/padron_reducido_ruc.zip http://www2.sunat.gob.pe/padron_reducido_ruc.zip
echo -e "\nUnzip padron_reducido_ruc.zip ..."
unzip /tmp/padron_reducido_ruc.zip -d /tmp/
echo -e "\nFormatear txt ..."
# formatear
awk -F '|' 'FNR > 1 {
  gsub(/\"/,"",$2); \
  print $1 "|" $2";;"\
  ($3=="ACTIVO" ? "1":$3)";;" \
  ($4=="HABIDO" ? "1":($4=="NO HABIDO" ? "0":$4))\
  ($5=="-" ? "": ";;" $5 ";;" $6 ";;" $7 ";;" $8 ";;" $9 ";;" $10 ";;" $11 ";;" $12 ";;" $13 ";;" $14 ";;" $15)\
  }' /tmp/padron_reducido_ruc.txt > /tmp/padron_reducido_ruc.csv
echo -e "\nCambiar encode a utf-8 ..."
iconv -f iso-8859-1 -t utf-8 /tmp/padron_reducido_ruc.csv -o /tmp/padron_reducido_ruc_utf8.csv
echo -e "\nCopiar padron_reducido_ruc_utf8.csv a db:/tmp ...\n"
# copiar rucs a /tmp docker
sudo docker cp /tmp/padron_reducido_ruc_utf8.csv "$DB_CONTAINER":/tmp
# sudo docker cp /home/robot/deco/decolecta-ruc/scripts/sunat_padron_reducido_clear_table.sh "$DB_CONTAINER":/tmp/
# exec copy
echo -e "\nexec sunat_padron_reducido_sql in docker ..."
sudo docker exec -i "$DB_CONTAINER" /bin/bash  /usr/local/bin/sunat_padron_reducido_clear_table.sh