#!/bin/bash
DB_CONTAINER=api-postgres-1
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
# formatear and clear
tr -d '<>\\' < /tmp/padron_reducido_ruc.txt > /tmp/remove_padron_reducido_ruc.txt
awk -F '|' 'FNR > 1 {
  gsub(/\"/,"",$2); \
  print $1 "|" $2">"\
  ($3=="ACTIVO" ? "1":$3)">" \
  ($4=="HABIDO" ? "1":($4=="NO HABIDO" ? "0":$4))\
  ($5=="-" ? "": ">" $5 ">" $6 ">" $7 ">" $8 ">" $9 ">" $10 ">" $11 ">" $12 ">" $13 ">" $14 ">" $15)\
  }' /tmp/remove_padron_reducido_ruc.txt > /tmp/remove_padron_reducido_ruc.csv
echo -e "\nCambiar encode a utf-8 ..."
iconv -f iso-8859-1 -t utf-8//TRANSLIT /tmp/remove_padron_reducido_ruc.csv -o /tmp/padron_reducido_ruc_utf8.csv
# remove ;
echo -e "\nCopiar padron_reducido_ruc_utf8.csv a db:/tmp ...\n"
sudo docker cp /tmp/padron_reducido_ruc_utf8.csv "$DB_CONTAINER":/tmp
echo -e "\nCopy sunat_padron_reducido_sql"
sudo docker cp ~/cronjobs/sunat_padron_reducido_sql.sh "$DB_CONTAINER":/tmp/
# exec copy
echo -e "\nexec sunat_padron_reducido_sql in docker ..."
sudo docker exec -i "$DB_CONTAINER" /bin/bash /tmp/sunat_padron_reducido_sql.sh
